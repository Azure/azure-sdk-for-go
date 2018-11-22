package easykey

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.0/keyvault"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/golang/protobuf/proto"
	"github.com/kylelemons/godebug/pretty"
)

func TestSecretIDName(t *testing.T) {
	tests := []struct {
		desc string
		url  string
		err  bool
		want ID
	}{
		{
			desc: "Too many parameters in URL",
			err:  true,
			url:  "https://myvault.vault.azure.net/secrets/mysecret",
		},
		{
			desc: "success",
			err:  false,
			want: ID{Name: "mysecret", Version: "2", url: "https://myvault.vault.azure.net/secrets/mysecret/2"},
			url:  "https://myvault.vault.azure.net/secrets/mysecret/2",
		},
	}

	for _, test := range tests {
		u, err := url.Parse(test.url)
		if err != nil {
			panic(err)
		}
		got, err := urlToID(u)
		if test.err && err == nil {
			t.Errorf("TestSecretIDName(%s): got %s, want error", test.desc, got)
			continue
		}
		if test.err {
			continue
		}

		if diff := pretty.Compare(test.want, got); diff != "" {
			t.Errorf("TestURLToID(%s): -want/+got:\n%s", test.desc, diff)
		}
	}
}

func mustURLToID(u *url.URL) ID {
	id, err := urlToID(u)
	if err != nil {
		panic(err)
	}
	return id
}

func TestSecretToBundle(t *testing.T) {
	id := "https://myvault.vault.azure.net/secrets/mysecret/1"

	u, err := url.Parse(id)
	if err != nil {
		panic(err)
	}
	created := time.Now()
	notBefore := created.Add(10 * time.Second)

	s := Secret{
		ID:    mustURLToID(u),
		Value: "secretvalue",
		Attr: SecretAttr{
			RecoveryLevel: Purgeable,
			Enabled:       true,
			Created:       created,
			NotBefore:     notBefore,
			Updated:       created,
		},
	}

	dNotBefore := date.UnixTime(notBefore)
	dCreated := date.UnixTime(created)

	want := &keyvault.SecretBundle{
		Value: proto.String("secretvalue"),
		ID:    proto.String(id),
		Attributes: &keyvault.SecretAttributes{
			RecoveryLevel: Purgeable,
			Enabled:       proto.Bool(true),
			NotBefore:     &dNotBefore,
			Created:       &dCreated,
			Updated:       &dCreated,
		},
	}

	got := s.toBundle()

	if diff := pretty.Compare(want, got); diff != "" {
		t.Fatalf("TestSecretToBundle: -want/+got:\n%s", diff)
	}
}

var sampleCert = []byte{48, 130, 3, 52, 48, 130, 2, 28, 160, 3, 2, 1, 2, 2, 16, 41, 95, 2, 251, 140,
	236, 73, 70, 129, 144, 238, 1, 248, 8, 230, 217, 48, 13, 6, 9, 42, 134,
	72, 134, 247, 13, 1, 1, 11, 5, 0, 48, 23, 49, 21, 48, 19, 6, 3, 85, 4, 3, 19,
	12, 109, 121, 100, 111, 109, 97, 105, 110, 46, 99, 111, 109, 48, 30, 23,
	13, 49, 56, 49, 49, 48, 55, 50, 49, 49, 54, 53, 57, 90, 23, 13, 49, 57, 49,
	49, 48, 55, 50, 49, 50, 54, 53, 57, 90, 48, 23, 49, 21, 48, 19, 6, 3, 85, 4,
	3, 19, 12, 109, 121, 100, 111, 109, 97, 105, 110, 46, 99, 111, 109, 48,
	130, 1, 34, 48, 13, 6, 9, 42, 134, 72, 134, 247, 13, 1, 1, 1, 5, 0, 3, 130, 1,
	15, 0, 48, 130, 1, 10, 2, 130, 1, 1, 0, 173, 90, 215, 70, 143, 222, 126, 2,
	133, 57, 90, 191, 190, 31, 221, 165, 108, 143, 141, 213, 10, 139, 225,
	231, 160, 66, 219, 99, 150, 253, 194, 28, 187, 117, 95, 121, 52, 108,
	211, 13, 218, 198, 214, 108, 24, 162, 156, 216, 66, 243, 196, 210, 133,
	65, 193, 126, 90, 141, 177, 3, 203, 231, 146, 103, 153, 14, 242, 42, 232,
	173, 229, 246, 241, 179, 140, 37, 125, 16, 221, 242, 108, 238, 29, 219,
	64, 163, 149, 201, 41, 235, 8, 142, 44, 206, 158, 76, 223, 182, 3, 83,
	211, 27, 37, 145, 47, 31, 210, 102, 49, 240, 23, 29, 65, 159, 73, 114, 77,
	73, 88, 239, 81, 166, 24, 67, 44, 251, 20, 36, 247, 252, 234, 155, 231,
	253, 69, 33, 44, 122, 191, 128, 180, 186, 104, 229, 20, 22, 14, 220, 177,
	96, 229, 21, 34, 36, 184, 254, 30, 66, 227, 120, 50, 200, 244, 70, 224,
	75, 209, 200, 243, 85, 199, 18, 75, 185, 163, 51, 9, 105, 231, 78, 60, 80,
	198, 22, 102, 103, 169, 102, 60, 179, 99, 239, 61, 152, 219, 188, 97,
	217, 208, 108, 4, 101, 192, 160, 167, 120, 94, 44, 8, 218, 41, 125, 216,
	93, 166, 105, 112, 98, 56, 185, 160, 50, 150, 181, 31, 98, 30, 116, 253,
	77, 149, 13, 7, 48, 201, 215, 182, 186, 2, 60, 46, 165, 245, 171, 21, 28,
	101, 10, 175, 92, 167, 152, 37, 221, 115, 173, 2, 3, 1, 0, 1, 163, 124, 48,
	122, 48, 14, 6, 3, 85, 29, 15, 1, 1, 255, 4, 4, 3, 2, 5, 160, 48, 9, 6, 3, 85, 29,
	19, 4, 2, 48, 0, 48, 29, 6, 3, 85, 29, 37, 4, 22, 48, 20, 6, 8, 43, 6, 1, 5, 5, 7,
	3, 1, 6, 8, 43, 6, 1, 5, 5, 7, 3, 2, 48, 31, 6, 3, 85, 29, 35, 4, 24, 48, 22, 128,
	20, 152, 38, 160, 228, 232, 224, 228, 4, 237, 73, 102, 128, 12, 204, 67, 186,
	48, 111, 145, 103, 48, 29, 6, 3, 85, 29, 14, 4, 22, 4, 20, 152, 38, 160, 228,
	232, 224, 228, 4, 237, 73, 102, 128, 12, 204, 67, 186, 48, 111, 145, 103, 48,
	13, 6, 9, 42, 134, 72, 134, 247, 13, 1, 1, 11, 5, 0, 3, 130, 1, 1, 0, 160, 193,
	50, 93, 238, 248, 208, 48, 169, 253, 24, 230, 97, 28, 184, 169, 175, 242, 165,
	21, 118, 224, 145, 235, 190, 66, 251, 12, 159, 125, 99, 118, 107, 115, 61,
	234, 253, 170, 205, 203, 202, 210, 189, 197, 149, 110, 98, 152, 144, 101, 30,
	93, 95, 20, 69, 190, 201, 13, 70, 31, 30, 137, 167, 31, 177, 17, 111, 170, 195,
	204, 201, 194, 221, 160, 4, 156, 127, 122, 233, 197, 198, 243, 31, 152, 165,
	221, 185, 180, 96, 182, 58, 68, 206, 55, 10, 111, 203, 195, 203, 162, 2, 4, 3,
	230, 67, 255, 70, 24, 9, 25, 188, 36, 238, 81, 187, 202, 52, 70, 161, 130, 70,
	173, 70, 35, 59, 247, 135, 230, 96, 191, 76, 39, 168, 185, 216, 218, 156, 67,
	180, 145, 198, 10, 26, 251, 185, 106, 64, 150, 105, 55, 51, 93, 246, 223, 149,
	163, 182, 169, 136, 88, 51, 199, 48, 216, 91, 72, 92, 177, 31, 11, 30, 24, 194,
	16, 129, 24, 39, 132, 229, 41, 108, 130, 241, 12, 12, 204, 65, 75, 97, 124, 49,
	186, 87, 37, 91, 213, 8, 255, 86, 166, 232, 190, 45, 85, 175, 151, 30, 199, 122,
	59, 9, 49, 117, 165, 219, 38, 25, 9, 144, 126, 52, 178, 162, 69, 194, 176, 184,
	205, 184, 49, 117, 33, 130, 58, 127, 204, 196, 89, 83, 121, 152, 145, 201, 13,
	195, 125, 194, 222, 211, 59, 177, 173, 197, 49, 102, 39}

func TestCertificateToFromBundle(t *testing.T) {
	id := "https://myvault.vault.azure.net/certificates/mysecret/1"

	u, err := url.Parse(id)
	if err != nil {
		panic(err)
	}
	created := time.Now()
	notBefore := created.Add(10 * time.Second)
	expires := notBefore.Add(10 * time.Second)

	c := Certificate{
		ID:         mustURLToID(u),
		Thumbprint: "thumbprint",
		Content:    sampleCert,
		Attr: CertAttr{
			RecoveryLevel: Purgeable,
			Enabled:       true,
			Expires:       expires,
			Created:       created,
			NotBefore:     notBefore,
			Updated:       created,
		},
	}

	dNotBefore := date.UnixTime(notBefore)
	dCreated := date.UnixTime(created)
	dExpires := date.UnixTime(expires)

	want := &keyvault.CertificateBundle{
		Cer:            &sampleCert,
		X509Thumbprint: proto.String("thumbprint"),
		ID:             proto.String(id),
		Attributes: &keyvault.CertificateAttributes{
			RecoveryLevel: Purgeable,
			Enabled:       proto.Bool(true),
			NotBefore:     &dNotBefore,
			Created:       &dCreated,
			Updated:       &dCreated,
			Expires:       &dExpires,
		},
	}

	got := c.toBundle()

	if diff := pretty.Compare(want, got); diff != "" {
		t.Fatalf("TestCertificateToFromBundle: toBundle problem:  -want/+got:\n%s", diff)
	}

	fromGot := Certificate{}
	if err := fromGot.fromBundle(got); err != nil {
		t.Fatalf("TestCertificateToFromBundle: problem with fromBundle(): %s", err)
	}

	if diff := pretty.Compare(c, fromGot); diff != "" {
		t.Fatalf("TestCertificateToFromBundle: fromBundle issue: -want/+got:\n%s", diff)

	}
}

func TestConfigValidate(t *testing.T) {
	const (
		clientKey = 0
		tenantKey = 1
		subKey    = 2
		secretKey = 3
		clientVal = "client"
		tenantVal = "tenant"
		subVal    = "subscription"
		secretVal = "secret"
		badVal    = ""
	)
	tests := [][]string{
		{badVal, tenantVal, subVal, secretVal},
		{clientVal, badVal, subVal, secretVal},
		{clientVal, tenantVal, badVal, secretVal},
		{clientVal, tenantVal, subVal, badVal},
		{clientVal, tenantVal, subVal, secretVal},
	}

	var errCount int
	for _, test := range tests {
		c := Config{
			ClientID:       test[clientKey],
			TenantID:       test[tenantKey],
			SubscriptionID: test[subKey],
			ClientSecret:   test[secretKey],
		}

		err := c.validate()
		if err != nil {
			errCount++
		}
	}
	if errCount != 4 {
		t.Fatalf("TestConfigValidate: got %d validate() calls to fail, want 4 to fail", errCount)
	}
}

func TestSecret(t *testing.T) {
	ctx := context.Background()
	m, err := NewMock("myVault")
	if err != nil {
		panic(err)
	}

	m.AddSecret("cosmos", "0", Secret{Value: "secret0"})
	m.AddSecret("cosmos", "1", Secret{Value: "secret1"})
	m.AddSecret("empty", "1", Secret{})

	kv, err := New(ctx, "myVault", Config{}, MockService(m))
	if err != nil {
		panic(err)
	}

	tests := []struct {
		desc string
		key  string
		ver  string
		want string
		err  bool
	}{
		{
			desc: "Secret doesn't exist",
			key:  "notExist",
			ver:  LatestVersion,
			err:  true,
		},
		{
			desc: "Secret was empty string",
			key:  "empty",
			ver:  LatestVersion,
			err:  true,
		},
		{
			desc: "Success latest secret",
			key:  "cosmos",
			ver:  LatestVersion,
			want: "secret1",
		},
		{
			desc: "Success version 0",
			key:  "cosmos",
			ver:  "0",
			want: "secret0",
		},
	}

	for _, test := range tests {
		got, err := kv.Secret(ctx, test.key, test.ver)

		switch {
		case err == nil && test.err:
			t.Errorf("TestSecret(%s): got err == nil, want err != nil", test.desc)
		case err != nil && !test.err:
			t.Errorf("TestSecret(%s): got err == %s, want err == nil", test.desc, err)
		case err != nil:
			continue
		}

		if test.want != got.Value {
			t.Errorf("TestSecret(%s): got %s, want %s", test.desc, got.Value, test.want)
		}
	}
}

func TestCertificate(t *testing.T) {
	ctx := context.Background()
	m, err := NewMock("myVault")
	if err != nil {
		panic(err)
	}

	cert0 := []byte("cert0")
	cert1 := []byte("cert1")

	m.AddCertificate("cosmos", "0", Certificate{Content: cert0})
	m.AddCertificate("cosmos", "1", Certificate{Content: cert1})
	m.AddCertificate("empty", "1", Certificate{})

	kv, err := New(ctx, "myVault", Config{}, MockService(m))
	if err != nil {
		panic(err)
	}

	tests := []struct {
		desc string
		key  string
		ver  string
		want []byte
		err  bool
	}{
		{
			desc: "Certificate doesn't exist",
			key:  "notExist",
			ver:  LatestVersion,
			err:  true,
		},
		{
			desc: "Certificate was empty",
			key:  "empty",
			ver:  LatestVersion,
			err:  true,
		},
		{
			desc: "Success latest certificate",
			key:  "cosmos",
			ver:  LatestVersion,
			want: cert1,
		},
		{
			desc: "Success version 0",
			key:  "cosmos",
			ver:  "0",
			want: cert0,
		},
	}

	for _, test := range tests {
		got, err := kv.Certificate(ctx, test.key, test.ver)

		switch {
		case err == nil && test.err:
			t.Errorf("TestCertificate(%s): got err == nil, want err != nil", test.desc)
		case err != nil && !test.err:
			t.Errorf("TestCertificate(%s): got err == %s, want err == nil", test.desc, err)
		case err != nil:
			continue
		}

		if string(test.want) != string(got.Content) {
			t.Errorf("TestCertificate(%s): got %s, want %s", test.desc, string(got.Content), string(test.want))
		}
	}
}
