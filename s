[33mcommit 7b0aa258822ec949db14c8e858f3ee0aa7255e4f[m
Author: Ken Faulkner <ken.faulkner@gmail.com>
Date:   Wed Nov 9 10:21:39 2016 +1100

    shared key lite for blob and queue initial changes

[1mdiff --git a/storage/blob.go b/storage/blob.go[m
[1mindex 696ccca..fa750be 100644[m
[1m--- a/storage/blob.go[m
[1m+++ b/storage/blob.go[m
[36m@@ -596,14 +596,16 @@[m [mfunc (b BlobStorageClient) leaseCommonPut(container string, name string, headers[m
 }[m
 [m
 // AcquireLease creates a lease for a blob as per https://msdn.microsoft.com/en-us/library/azure/ee691972.aspx[m
[31m-[m
 // returns leaseID acquired[m
 func (b BlobStorageClient) AcquireLease(container string, name string, leaseTimeInSeconds int, proposedLeaseID string) (returnedLeaseID string, err error) {[m
 	headers := b.client.getStandardHeaders()[m
 	headers[leaseAction] = acquireLease[m
[31m-	headers[leaseProposedID] = proposedLeaseID[m
 	headers[leaseDuration] = strconv.Itoa(leaseTimeInSeconds)[m
 [m
[32m+[m	[32mif proposedLeaseID != "" {[m
[32m+[m		[32mheaders[leaseProposedID] = proposedLeaseID[m
[32m+[m	[32m}[m
[32m+[m
 	respHeaders, err := b.leaseCommonPut(container, name, headers, http.StatusCreated)[m
 	if err != nil {[m
 		return "", err[m
[1mdiff --git a/storage/client.go b/storage/client.go[m
[1mindex a41a586..816c37c 100644[m
[1m--- a/storage/client.go[m
[1m+++ b/storage/client.go[m
[36m@@ -10,6 +10,7 @@[m [mimport ([m
 	"fmt"[m
 	"io"[m
 	"io/ioutil"[m
[32m+[m	[32m"log"[m
 	"net/http"[m
 	"net/url"[m
 	"regexp"[m
[36m@@ -127,7 +128,8 @@[m [mfunc NewBasicClient(accountName, accountKey string) (Client, error) {[m
 	if accountName == StorageEmulatorAccountName {[m
 		return NewEmulatorClient()[m
 	}[m
[31m-	return NewClient(accountName, accountKey, DefaultBaseURL, DefaultAPIVersion, defaultUseHTTPS)[m
[32m+[m	[32mreturn NewEmulatorClient()[m
[32m+[m	[32m//return NewClient(accountName, accountKey, DefaultBaseURL, DefaultAPIVersion, defaultUseHTTPS)[m
 }[m
 [m
 //NewEmulatorClient contructs a Client intended to only work with Azure[m
[36m@@ -357,6 +359,26 @@[m [mfunc (c Client) buildCanonicalizedResource(uri string) (string, error) {[m
 	return cr, nil[m
 }[m
 [m
[32m+[m[32mfunc (c Client) buildCanonicalizedLiteString(verb string, headers map[string]string, canonicalizedResource string) string {[m
[32m+[m
[32m+[m	[32mlog.Printf("verb %s", verb)[m
[32m+[m	[32mlog.Printf("md5 %s", headers["Content-MD5"])[m
[32m+[m	[32mlog.Printf("t %s", headers["Content-Type"])[m
[32m+[m	[32mlog.Printf("d %s", headers["Date"])[m
[32m+[m
[32m+[m	[32mcanonicalizedString := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",[m
[32m+[m		[32mverb,[m
[32m+[m		[32mheaders["Content-MD5"],[m
[32m+[m		[32mheaders["Content-Type"],[m
[32m+[m		[32mheaders["Date"],[m
[32m+[m		[32mc.buildCanonicalizedHeader(headers),[m
[32m+[m		[32mcanonicalizedResource)[m
[32m+[m
[32m+[m	[32mlog.Printf("str %s", canonicalizedString)[m
[32m+[m	[32mhmac := c.computeHmac256(canonicalizedString)[m
[32m+[m	[32mreturn fmt.Sprintf("SharedKeyLite %s:%s", c.accountName, hmac)[m
[32m+[m[32m}[m
[32m+[m
 func (c Client) buildCanonicalizedString(verb string, headers map[string]string, canonicalizedResource string) string {[m
 	contentLength := headers["Content-Length"][m
 	if contentLength == "0" {[m
[36m@@ -382,15 +404,14 @@[m [mfunc (c Client) buildCanonicalizedString(verb string, headers map[string]string,[m
 }[m
 [m
 func (c Client) exec(verb, url string, headers map[string]string, body io.Reader) (*storageResponse, error) {[m
[31m-	authHeader, err := c.getAuthorizationHeader(verb, url, headers)[m
[31m-	if err != nil {[m
[31m-		return nil, err[m
[31m-	}[m
[31m-	headers["Authorization"] = authHeader[m
[32m+[m	[32mcanonicalizedResource, err := c.buildCanonicalizedResource(url)[m
 	if err != nil {[m
 		return nil, err[m
 	}[m
 [m
[32m+[m	[32mheaders["Authorization"] = c.buildCanonicalizedLiteString(verb, headers, canonicalizedResource)[m
[32m+[m	[32mlog.Printf("auth is %s", headers["Authorization"])[m
[32m+[m
 	req, err := http.NewRequest(verb, url, body)[m
 	if err != nil {[m
 		return nil, errors.New("azure/storage: error creating request: " + err.Error())[m
[36m@@ -494,6 +515,8 @@[m [mfunc (c Client) execInternalJSON(verb, url string, headers map[string]string, bo[m
 }[m
 [m
 func (c Client) createSharedKeyLite(url string, headers map[string]string) (string, error) {[m
[32m+[m
[32m+[m	[32mlog.Printf("createSharedKeyLite start: %s", headers)[m
 	can, err := c.buildCanonicalizedResourceTable(url)[m
 [m
 	if err != nil {[m
