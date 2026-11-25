// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"bytes"
	"context"
	"crypto"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"runtime/debug"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	internalattestation "github.com/Azure/azure-sdk-for-go/sdk/azidentity/internal/attestation"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
)

const vaultScope = "https://vault.azure.net/.default"

type configuration struct {
	vaultURL, secretName           string
	clientID, objectID, resourceID string
	timeout                        time.Duration
	concurrency                    int
	isolationClientID              string
	dac                            bool
	disableMTLSPoP                 bool
}

func main() {
	var cfg configuration
	flag.StringVar(&cfg.vaultURL, "vault-url", "", "HTTPS URL of an mTLS PoP-enabled public Azure Key Vault (required)")
	flag.StringVar(&cfg.secretName, "secret-name", "", "Existing secret the managed identity can read (required)")
	flag.StringVar(&cfg.clientID, "client-id", "", "User-assigned identity client ID; omit all identity flags for system-assigned")
	flag.StringVar(&cfg.objectID, "object-id", "", "User-assigned identity object ID (mutually exclusive with other identity flags)")
	flag.StringVar(&cfg.resourceID, "resource-id", "", "User-assigned identity resource ID (mutually exclusive with other identity flags)")
	flag.BoolVar(&cfg.dac, "dac", false, "Use DefaultAzureCredential; requires AZURE_TOKEN_CREDENTIALS=prod and matching AZURE_CLIENT_ID for UAMI")
	flag.BoolVar(&cfg.disableMTLSPoP, "disable-mtls-pop", false, "Run the explicit MIC bearer opt-out control instead of the positive PoP checks")
	flag.IntVar(&cfg.concurrency, "concurrency", 0, "Opt in to shared/multiple-client concurrent reads with 2-8 workers; zero disables")
	flag.StringVar(&cfg.isolationClientID, "isolation-client-id", "", "Opt in to alternating system-assigned and this UAMI client ID; both must read the secret")
	flag.DurationVar(&cfg.timeout, "timeout", 5*time.Minute, "Overall test timeout")
	flag.Parse()
	if err := run(context.Background(), cfg); err != nil {
		fmt.Fprintln(os.Stderr, "FAIL:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, cfg configuration) error {
	if err := cfg.validate(); err != nil {
		return err
	}
	if internalattestation.Registered() != attestationImported {
		return errors.New("attestation activation doesn't match this smoke build's import configuration")
	}
	bearerControl := cfg.disableMTLSPoP || !attestationImported
	if bearerControl && (cfg.concurrency != 0 || cfg.isolationClientID != "") {
		return errors.New("concurrency and identity isolation are positive PoP checks; run bearer controls separately")
	}
	if cfg.disableMTLSPoP && !attestationImported {
		return errors.New("run disable-mtls-pop with the default imported build, and the no-import build without that flag")
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return errors.New("Go build information is unavailable")
	}
	if err := verifyBuild(info); err != nil {
		return err
	}
	// Refuse source overrides instead of changing the caller's environment or testing another service.
	for _, name := range []string{
		"IDENTITY_ENDPOINT", "IDENTITY_HEADER", "IDENTITY_SERVER_THUMBPRINT",
		"IMDS_ENDPOINT", "MSI_ENDPOINT", "MSI_SECRET", "AZURE_POD_IDENTITY_AUTHORITY_HOST",
		"MSAL_MI_DISABLE_IMDS_V2",
	} {
		if os.Getenv(name) != "" {
			return fmt.Errorf("%s is set; run from a clean Azure VM environment using native IMDS", name)
		}
	}
	ctx, cancel := context.WithTimeout(ctx, cfg.timeout)
	defer cancel()
	credential, err := newSmokeCredential(cfg)
	if err != nil {
		return fmt.Errorf("creating credential: %w", err)
	}
	if bearerControl {
		if err := verifyBearerControl(ctx, cfg, credential, nil); err != nil {
			return fmt.Errorf("bearer control: %w", err)
		}
		if cfg.disableMTLSPoP {
			fmt.Println("PASS: imported integration with explicit MIC opt-out returns bearer tokens and SDK reads remain bearer")
		} else {
			fmt.Println("PASS: no-import control returns bearer tokens and SDK reads remain bearer")
		}
		return nil
	}
	if err := verifySDK(ctx, cfg, credential); err != nil {
		return fmt.Errorf("automatic SDK mTLS PoP: %w", err)
	}
	fmt.Println("PASS: cold Key Vault SDK reads secret twice with automatic mTLS PoP signaling")
	bearer, err := credential.GetToken(ctx, policy.TokenRequestOptions{Scopes: []string{vaultScope}})
	if err != nil {
		return fmt.Errorf("initial bearer acquisition: %w", err)
	}
	if err := validateBearer(bearer); err != nil {
		return err
	}
	fmt.Println("PASS: ordinary bearer acquisition returns no binding certificate")

	bound, err := credential.GetToken(ctx, policy.TokenRequestOptions{Scopes: []string{vaultScope}, EnableMTLSPoP: true})
	if err != nil {
		return fmt.Errorf("mTLS PoP acquisition (requires Windows, Credential Guard/KeyGuard, and IMDSv2): %w", err)
	}
	thumbprint, err := validateBound(bound, time.Now())
	if err != nil {
		return err
	}
	fmt.Printf("PASS: bound token cnf.x5t#S256 matches certificate %s; token expires %s\n",
		thumbprint, bound.ExpiresOn.UTC().Format(time.RFC3339))

	target := strings.TrimRight(cfg.vaultURL, "/") + "/secrets/" + url.PathEscape(cfg.secretName) + "?api-version=7.4"
	if err := verifyResource(ctx, target, bound); err != nil {
		return err
	}
	fmt.Println("PASS: direct mTLS resource request succeeds; same token without certificate returns MissingClientCertificate")

	bearer, err = credential.GetToken(ctx, policy.TokenRequestOptions{Scopes: []string{vaultScope}})
	if err != nil {
		return fmt.Errorf("bearer acquisition after mTLS PoP: %w", err)
	}
	if err := validateBearer(bearer); err != nil {
		return err
	}
	fmt.Println("PASS: bearer acquisition remains unbound after mTLS PoP acquisition")

	if cfg.concurrency > 0 {
		if err := verifySDKConcurrency(ctx, cfg, credential, nil); err != nil {
			return fmt.Errorf("concurrent SDK mTLS PoP: %w", err)
		}
		fmt.Printf("PASS: %d concurrent reads with one client and with two clients sharing a credential\n", cfg.concurrency)
	}
	if cfg.isolationClientID != "" {
		var credentials [2]azcore.TokenCredential
		for i, id := range []azidentity.ManagedIDKind{nil, azidentity.ClientID(cfg.isolationClientID)} {
			credentials[i], err = azidentity.NewManagedIdentityCredential(&azidentity.ManagedIdentityCredentialOptions{
				ID: id,
			})
			if err != nil {
				return fmt.Errorf("creating identity-isolation credential #%d: %w", i+1, err)
			}
		}
		if err := verifyIdentityIsolation(ctx, cfg, credentials, nil); err != nil {
			return fmt.Errorf("identity isolation: %w", err)
		}
		fmt.Println("PASS: alternating system-assigned/UAMI SDK reads retain distinct token identities")
	}
	fmt.Println("PASS: all mTLS PoP smoke checks")
	return nil
}

func (c configuration) validate() error {
	u, err := url.Parse(c.vaultURL)
	if err != nil || u.Scheme != "https" || u.Hostname() == "" || u.User != nil ||
		u.RawQuery != "" || u.ForceQuery || u.Fragment != "" || (u.Path != "" && u.Path != "/") {
		return errors.New("vault-url must be an HTTPS origin without credentials, path, query, or fragment")
	}
	if c.secretName == "" || strings.ContainsAny(c.secretName, "/\\?#") || c.timeout <= 0 {
		return errors.New("secret-name must be a nonempty secret name and timeout must be positive")
	}
	selectors := 0
	for _, id := range []string{c.clientID, c.objectID, c.resourceID} {
		if id != "" {
			selectors++
		}
	}
	if selectors > 1 {
		return errors.New("specify at most one of client-id, object-id, and resource-id")
	}
	if c.concurrency != 0 && (c.concurrency < 2 || c.concurrency > 8) {
		return errors.New("concurrency must be zero (disabled) or between 2 and 8")
	}
	if c.dac && (c.objectID != "" || c.resourceID != "" || c.isolationClientID != "") {
		return errors.New("dac supports only client-id selection; run identity isolation separately without dac")
	}
	if c.dac && c.disableMTLSPoP {
		return errors.New("disable-mtls-pop is a MIC option; DAC has no additional public opt-out")
	}
	return nil
}

func newSmokeCredential(cfg configuration) (azcore.TokenCredential, error) {
	if !cfg.dac {
		return azidentity.NewManagedIdentityCredential(&azidentity.ManagedIdentityCredentialOptions{
			ID: cfg.identity(), DisableMTLSPoP: cfg.disableMTLSPoP,
		})
	}
	if err := validateDACEnvironment(cfg); err != nil {
		return nil, err
	}
	return azidentity.NewDefaultAzureCredential(nil)
}

func validateDACEnvironment(cfg configuration) error {
	if os.Getenv("AZURE_TOKEN_CREDENTIALS") != "prod" {
		return errors.New("dac requires AZURE_TOKEN_CREDENTIALS=prod to exercise probing without developer-credential fallback")
	}
	for _, name := range []string{
		"AZURE_TENANT_ID", "AZURE_CLIENT_SECRET", "AZURE_CLIENT_CERTIFICATE_PATH",
		"AZURE_USERNAME", "AZURE_PASSWORD", "AZURE_FEDERATED_TOKEN_FILE",
	} {
		if os.Getenv(name) != "" {
			return fmt.Errorf("%s is set; dac smoke requires EnvironmentCredential and WorkloadIdentityCredential to be unconfigured", name)
		}
	}
	id, set := os.LookupEnv("AZURE_CLIENT_ID")
	if cfg.clientID == "" && set {
		return errors.New("unset AZURE_CLIENT_ID for system-assigned dac smoke")
	}
	if cfg.clientID != "" && (!set || !strings.EqualFold(id, cfg.clientID)) {
		return errors.New("AZURE_CLIENT_ID must match client-id for UAMI dac smoke")
	}
	return nil
}

func (c configuration) identity() azidentity.ManagedIDKind {
	switch {
	case c.clientID != "":
		return azidentity.ClientID(c.clientID)
	case c.objectID != "":
		return azidentity.ObjectID(c.objectID)
	case c.resourceID != "":
		return azidentity.ResourceID(c.resourceID)
	default:
		return nil
	}
}

func verifyBuild(info *debug.BuildInfo) error {
	// azsecrets enables PoP by installing keyvault/internal's token binding policy; a released
	// version without it stays bearer-only against the local keyvault/internal
	for _, name := range []string{"azcore", "azidentity", "internal", "security/keyvault/azsecrets", "security/keyvault/internal"} {
		path := "github.com/Azure/azure-sdk-for-go/sdk/" + name
		local := false
		for _, dep := range info.Deps {
			if dep.Path == path && dep.Replace != nil && dep.Replace.Path != "" &&
				(dep.Replace.Version == "" || dep.Replace.Version == "(devel)") {
				fmt.Printf("Using local %s: %s\n", name, dep.Replace.Path)
				local = true
				break
			}
		}
		if !local {
			return fmt.Errorf("%s must be replaced with this checkout; dependency replace directives aren't transitive", name)
		}
	}
	return nil
}

func validateBearer(token azcore.AccessToken) error {
	if token.Token == "" || !token.ExpiresOn.After(time.Now()) || token.BindingCertificate != nil {
		return errors.New("bearer acquisition must return a nonempty, unexpired token without a binding certificate")
	}
	return nil
}

func validateBound(token azcore.AccessToken, now time.Time) (string, error) {
	cert := token.BindingCertificate
	if token.Token == "" || !token.ExpiresOn.After(now) || cert == nil || len(cert.Certificate) == 0 || cert.PrivateKey == nil {
		return "", errors.New("mTLS PoP must return a nonempty, unexpired token and a binding certificate with a private key")
	}
	leaf, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return "", fmt.Errorf("parsing binding certificate: %w", err)
	}
	if now.Before(leaf.NotBefore) || !now.Before(leaf.NotAfter) {
		return "", errors.New("binding certificate isn't currently valid")
	}
	signer, ok := cert.PrivateKey.(crypto.Signer)
	if !ok {
		return "", errors.New("binding private key doesn't implement crypto.Signer")
	}
	publicKey, err := x509.MarshalPKIXPublicKey(signer.Public())
	if err != nil {
		return "", fmt.Errorf("marshaling binding public key: %w", err)
	}
	if !bytes.Equal(publicKey, leaf.RawSubjectPublicKeyInfo) {
		return "", errors.New("binding private key doesn't match the leaf certificate")
	}
	payload, err := tokenPayload(token.Token)
	if err != nil {
		return "", err
	}
	var claims struct {
		Confirmation struct {
			Thumbprint string `json:"x5t#S256"`
		} `json:"cnf"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return "", fmt.Errorf("decoding bound token claims: %w", err)
	}
	hash := sha256.Sum256(leaf.Raw)
	thumbprint := base64.RawURLEncoding.EncodeToString(hash[:])
	if claims.Confirmation.Thumbprint != thumbprint {
		return "", errors.New("bound token cnf.x5t#S256 is missing or doesn't match the binding certificate")
	}
	return thumbprint, nil
}

func tokenPayload(token string) ([]byte, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("access token isn't a three-part JWT")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("decoding token payload: %w", err)
	}
	return payload, nil
}

type tokenIdentity struct {
	clientID, objectID, tenantID string
}

func readTokenIdentity(token string) (tokenIdentity, error) {
	payload, err := tokenPayload(token)
	if err != nil {
		return tokenIdentity{}, err
	}
	var claims struct {
		AppID    string `json:"appid"`
		AZP      string `json:"azp"`
		ObjectID string `json:"oid"`
		TenantID string `json:"tid"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return tokenIdentity{}, errors.New("decoding token identity claims failed")
	}
	if claims.AppID == "" {
		claims.AppID = claims.AZP
	}
	if claims.AppID == "" || claims.ObjectID == "" || claims.TenantID == "" {
		return tokenIdentity{}, errors.New("token identity requires appid/azp, oid, and tid claims")
	}
	if claims.AZP != "" && !strings.EqualFold(claims.AppID, claims.AZP) {
		return tokenIdentity{}, errors.New("token appid and azp claims disagree")
	}
	return tokenIdentity{strings.ToLower(claims.AppID), strings.ToLower(claims.ObjectID), strings.ToLower(claims.TenantID)}, nil
}

func newResourceClient(cert *tls.Certificate, roots *x509.CertPool) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{
				RootCAs:       roots,
				MinVersion:    tls.VersionTLS12,
				MaxVersion:    tls.VersionTLS12,
				Renegotiation: tls.RenegotiateOnceAsClient,
				// Key Vault asks for the certificate after reading the authorization header.
				GetClientCertificate: func(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
					if cert == nil {
						return &tls.Certificate{}, nil
					}
					return cert, nil
				},
			},
		},
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return errors.New("redirects aren't permitted for mTLS PoP")
		},
	}
}

func resourceGET(ctx context.Context, client *http.Client, target, token string) (int, []byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "mtls_pop "+token)
	req.Header.Set("x-ms-tokenboundauth", "true")
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	const maxBody = 1 << 20
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBody+1))
	if err != nil {
		return resp.StatusCode, nil, fmt.Errorf("reading resource response: %w", err)
	}
	if len(body) > maxBody {
		return resp.StatusCode, nil, errors.New("resource response exceeds 1 MiB")
	}
	return resp.StatusCode, body, nil
}

func verifyResource(ctx context.Context, target string, token azcore.AccessToken) error {
	client := newResourceClient(token.BindingCertificate, nil)
	defer client.CloseIdleConnections()
	status, body, err := resourceGET(ctx, client, target, token.Token)
	if err != nil {
		return fmt.Errorf("direct bound resource request: %w", err)
	}
	if err := validateResourceResponse(status, body); err != nil {
		return fmt.Errorf("direct bound resource request: %w", err)
	}
	// A separate transport prevents the negative request from reusing an authenticated connection.
	withoutCert := newResourceClient(nil, nil)
	defer withoutCert.CloseIdleConnections()
	status, body, err = resourceGET(ctx, withoutCert, target, token.Token)
	if err != nil {
		return fmt.Errorf("missing-certificate request failed without an HTTP rejection; this isn't proof of binding: %w", err)
	}
	return validateRejection(status, body)
}

func validateResourceResponse(status int, body []byte) error {
	if status == http.StatusOK {
		return nil
	}
	// Only report error codes, never a secret value, error message, or raw response body.
	var response struct {
		Error struct {
			Code       string `json:"code"`
			InnerError struct {
				Code string `json:"code"`
			} `json:"innererror"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("HTTP %d, want 200 (couldn't decode Key Vault error response; body omitted)", status)
	}
	if response.Error.Code == "" {
		return fmt.Errorf("HTTP %d, want 200 (no Key Vault error code in response; body omitted)", status)
	}
	return fmt.Errorf("HTTP %d, want 200 (Key Vault code %q, inner code %q; check vault opt-in, RBAC, network access, and secret name)",
		status, response.Error.Code, response.Error.InnerError.Code)
}

func validateRejection(status int, body []byte) error {
	if status < 400 || status >= 500 || !bytes.Contains(body, []byte("MissingClientCertificate")) {
		return fmt.Errorf("missing-certificate request: HTTP %d without the required 4xx/MissingClientCertificate rejection", status)
	}
	return nil
}

type checkedCredential struct {
	credential azcore.TokenCredential
	calls      atomic.Int32
	identity   *tokenIdentity
	clientID   string
	bearer     bool
}

func (c *checkedCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	if !options.EnableMTLSPoP {
		return azcore.AccessToken{}, errors.New("SDK didn't propagate EnableMTLSPoP to the credential")
	}
	token, err := c.credential.GetToken(ctx, options)
	if err != nil {
		return azcore.AccessToken{}, err
	}
	if c.bearer {
		if err := validateBearer(token); err != nil {
			return azcore.AccessToken{}, err
		}
	} else {
		if _, err := validateBound(token, time.Now()); err != nil {
			return azcore.AccessToken{}, err
		}
	}
	if c.identity != nil || c.clientID != "" {
		got, err := readTokenIdentity(token.Token)
		if err != nil {
			return azcore.AccessToken{}, err
		}
		if c.identity != nil && got != *c.identity {
			return azcore.AccessToken{}, errors.New("credential returned a token for a different identity")
		}
		if c.clientID != "" && !strings.EqualFold(got.clientID, c.clientID) {
			return azcore.AccessToken{}, errors.New("credential returned a token for a different managed identity client ID")
		}
	}
	c.calls.Add(1)
	return token, nil
}

func verifySDK(ctx context.Context, cfg configuration, credential azcore.TokenCredential) error {
	checked := &checkedCredential{credential: credential, clientID: cfg.clientID}
	client, err := newSecretClient(cfg, checked, nil)
	if err != nil {
		return err
	}
	for i := 1; i <= 2; i++ {
		if err := readSecret(ctx, client, cfg.secretName); err != nil {
			return fmt.Errorf("GetSecret #%d: %w", i, err)
		}
	}
	if checked.calls.Load() == 0 {
		return errors.New("SDK resource calls never requested a bound token")
	}
	return nil
}

func verifyBearerControl(ctx context.Context, cfg configuration, credential azcore.TokenCredential, options *azsecrets.ClientOptions) error {
	for _, requestPoP := range []bool{false, true} {
		token, err := credential.GetToken(ctx, policy.TokenRequestOptions{Scopes: []string{vaultScope}, EnableMTLSPoP: requestPoP})
		if err != nil {
			return fmt.Errorf("direct control acquisition: %w", err)
		}
		if err := validateBearer(token); err != nil {
			return err
		}
	}
	cp := azsecrets.ClientOptions{}
	if options != nil {
		cp = *options
	}
	observer := &bearerControlPolicy{}
	cp.PerRetryPolicies = append(slices.Clone(cp.PerRetryPolicies), observer)
	checked := &checkedCredential{credential: credential, clientID: cfg.clientID, bearer: true}
	client, err := newSecretClient(cfg, checked, &cp)
	if err != nil {
		return err
	}
	for i := 1; i <= 2; i++ {
		if err := readSecret(ctx, client, cfg.secretName); err != nil {
			return fmt.Errorf("control GetSecret #%d: %w", i, err)
		}
	}
	if checked.calls.Load() == 0 || observer.anonymous != 1 || observer.authenticated != 2 {
		return errors.New("bearer control didn't observe one cold challenge and two authenticated SDK reads")
	}
	return nil
}

type bearerControlPolicy struct {
	anonymous, authenticated int
}

func (p *bearerControlPolicy) Do(req *policy.Request) (*http.Response, error) {
	if req.Raw().Header.Get("x-ms-tokenboundauth") != "" {
		return nil, errors.New("bearer control unexpectedly sent a token-binding header")
	}
	if authorization := req.Raw().Header.Get("Authorization"); authorization == "" {
		p.anonymous++
	} else {
		if !strings.HasPrefix(authorization, "Bearer ") {
			return nil, errors.New("bearer control unexpectedly sent non-bearer authorization")
		}
		p.authenticated++
	}
	return req.Next()
}

func newSecretClient(cfg configuration, credential azcore.TokenCredential, options *azsecrets.ClientOptions) (*azsecrets.Client, error) {
	cp := azsecrets.ClientOptions{}
	if options != nil {
		cp = *options
	}
	cp.Retry.MaxRetries = -1
	return azsecrets.NewClient(cfg.vaultURL, credential, &cp)
}

func readSecret(ctx context.Context, client *azsecrets.Client, name string) error {
	_, err := client.GetSecret(ctx, name, "", nil)
	// Don't print the response body: it may contain a secret value.
	var response *azcore.ResponseError
	if errors.As(err, &response) {
		return fmt.Errorf("HTTP %d, code %s", response.StatusCode, response.ErrorCode)
	}
	return err
}

func verifySDKConcurrency(ctx context.Context, cfg configuration, credential azcore.TokenCredential, options *azsecrets.ClientOptions) error {
	checked := &checkedCredential{credential: credential}
	for _, clientCount := range []int{1, 2} {
		clients := make([]*azsecrets.Client, clientCount)
		for i := range clients {
			client, err := newSecretClient(cfg, checked, options)
			if err != nil {
				return err
			}
			clients[i] = client
		}
		if err := readSecret(ctx, clients[0], cfg.secretName); err != nil {
			return fmt.Errorf("%d-client baseline: %w", clientCount, err)
		}
		cx, cancel := context.WithCancel(ctx)
		var wg sync.WaitGroup
		results := make([]error, cfg.concurrency)
		start := make(chan struct{})
		for i := range cfg.concurrency {
			wg.Go(func() {
				<-start
				if err := readSecret(cx, clients[i%len(clients)], cfg.secretName); err != nil {
					results[i] = fmt.Errorf("%d-client read #%d: %w", clientCount, i+1, err)
					cancel()
				}
			})
		}
		close(start)
		wg.Wait()
		cancel()
		if err := errors.Join(results...); err != nil {
			return err
		}
	}
	if checked.calls.Load() == 0 {
		return errors.New("concurrent SDK reads never requested a bound token")
	}
	return nil
}

func verifyIdentityIsolation(ctx context.Context, cfg configuration, credentials [2]azcore.TokenCredential, options *azsecrets.ClientOptions) error {
	var identities [2]tokenIdentity
	var checked [2]*checkedCredential
	for i, credential := range credentials {
		checked[i] = &checkedCredential{credential: credential}
		token, err := checked[i].GetToken(ctx, policy.TokenRequestOptions{Scopes: []string{vaultScope}, EnableMTLSPoP: true})
		if err != nil {
			return fmt.Errorf("identity #%d baseline: %w", i+1, err)
		}
		identities[i], err = readTokenIdentity(token.Token)
		if err != nil {
			return err
		}
		checked[i].identity = &identities[i]
	}
	if !strings.EqualFold(identities[1].clientID, cfg.isolationClientID) {
		return errors.New("UAMI token doesn't match isolation-client-id")
	}
	if identities[0].clientID == identities[1].clientID || identities[0].objectID == identities[1].objectID {
		return errors.New("system-assigned and UAMI tokens must identify different principals")
	}
	if identities[0].tenantID != identities[1].tenantID {
		return errors.New("system-assigned and UAMI tokens must belong to the same tenant")
	}
	var clients [2]*azsecrets.Client
	for i := range clients {
		client, err := newSecretClient(cfg, checked[i], options)
		if err != nil {
			return err
		}
		clients[i] = client
	}
	for round := range 2 {
		for i := range clients {
			// Reenter each credential even when the SDK's own token cache is warm.
			if _, err := checked[i].GetToken(ctx, policy.TokenRequestOptions{Scopes: []string{vaultScope}, EnableMTLSPoP: true}); err != nil {
				return fmt.Errorf("identity #%d round #%d: %w", i+1, round+1, err)
			}
			if err := readSecret(ctx, clients[i], cfg.secretName); err != nil {
				return fmt.Errorf("identity #%d SDK round #%d: %w", i+1, round+1, err)
			}
		}
	}
	return nil
}
