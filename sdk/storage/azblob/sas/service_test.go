// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package sas

import (
	"errors"
	"net/url"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/stretchr/testify/require"
)

func TestContainerPermissions_String(t *testing.T) {
	testdata := []struct {
		input    ContainerPermissions
		expected string
	}{
		{input: ContainerPermissions{Read: true}, expected: "r"},
		{input: ContainerPermissions{Add: true}, expected: "a"},
		{input: ContainerPermissions{Create: true}, expected: "c"},
		{input: ContainerPermissions{Write: true}, expected: "w"},
		{input: ContainerPermissions{Delete: true}, expected: "d"},
		{input: ContainerPermissions{DeletePreviousVersion: true}, expected: "x"},
		{input: ContainerPermissions{List: true}, expected: "l"},
		{input: ContainerPermissions{Tag: true}, expected: "t"},
		{input: ContainerPermissions{FilterByTags: true}, expected: "f"},
		{input: ContainerPermissions{Move: true}, expected: "m"},
		{input: ContainerPermissions{Execute: true}, expected: "e"},
		{input: ContainerPermissions{ModifyOwnership: true}, expected: "o"},
		{input: ContainerPermissions{ModifyPermissions: true}, expected: "p"},
		{input: ContainerPermissions{SetImmutabilityPolicy: true}, expected: "i"},
		{input: ContainerPermissions{
			Read:                  true,
			Add:                   true,
			Create:                true,
			Write:                 true,
			Delete:                true,
			DeletePreviousVersion: true,
			List:                  true,
			Tag:                   true,
			FilterByTags:          true,
			Move:                  true,
			Execute:               true,
			ModifyOwnership:       true,
			ModifyPermissions:     true,
			SetImmutabilityPolicy: true,
		}, expected: "racwdxltfmeopi"},
	}
	for _, c := range testdata {
		require.Equal(t, c.expected, c.input.String())
	}
}

func TestContainerPermissions_Parse(t *testing.T) {
	testdata := []struct {
		input    string
		expected ContainerPermissions
	}{
		{expected: ContainerPermissions{Read: true}, input: "r"},
		{expected: ContainerPermissions{Add: true}, input: "a"},
		{expected: ContainerPermissions{Create: true}, input: "c"},
		{expected: ContainerPermissions{Write: true}, input: "w"},
		{expected: ContainerPermissions{Delete: true}, input: "d"},
		{expected: ContainerPermissions{DeletePreviousVersion: true}, input: "x"},
		{expected: ContainerPermissions{List: true}, input: "l"},
		{expected: ContainerPermissions{Tag: true}, input: "t"},
		{expected: ContainerPermissions{FilterByTags: true}, input: "f"},
		{expected: ContainerPermissions{Move: true}, input: "m"},
		{expected: ContainerPermissions{Execute: true}, input: "e"},
		{expected: ContainerPermissions{ModifyOwnership: true}, input: "o"},
		{expected: ContainerPermissions{ModifyPermissions: true}, input: "p"},
		{expected: ContainerPermissions{SetImmutabilityPolicy: true}, input: "i"},
		{expected: ContainerPermissions{
			Read:                  true,
			Add:                   true,
			Create:                true,
			Write:                 true,
			Delete:                true,
			DeletePreviousVersion: true,
			List:                  true,
			Tag:                   true,
			FilterByTags:          true,
			Move:                  true,
			Execute:               true,
			ModifyOwnership:       true,
			ModifyPermissions:     true,
			SetImmutabilityPolicy: true,
		}, input: "racwdxltfmeopi"},
		{expected: ContainerPermissions{
			Read:                  true,
			Add:                   true,
			Create:                true,
			Write:                 true,
			Delete:                true,
			DeletePreviousVersion: true,
			List:                  true,
			Tag:                   true,
			FilterByTags:          true,
			Move:                  true,
			Execute:               true,
			ModifyOwnership:       true,
			ModifyPermissions:     true,
			SetImmutabilityPolicy: true,
		}, input: "ctpwxfmreodail"}, // Wrong order parses correctly
	}
	for _, c := range testdata {
		permissions, err := parseContainerPermissions(c.input)
		require.Nil(t, err)
		require.Equal(t, c.expected, permissions)
	}
}

func TestContainerPermissions_ParseNegative(t *testing.T) {
	_, err := parseContainerPermissions("cpwxtfmreodailz") // Here 'z' is invalid
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "122")
}

func TestBlobPermissions_String(t *testing.T) {
	testdata := []struct {
		input    BlobPermissions
		expected string
	}{
		{input: BlobPermissions{Read: true}, expected: "r"},
		{input: BlobPermissions{Add: true}, expected: "a"},
		{input: BlobPermissions{Create: true}, expected: "c"},
		{input: BlobPermissions{Write: true}, expected: "w"},
		{input: BlobPermissions{Delete: true}, expected: "d"},
		{input: BlobPermissions{DeletePreviousVersion: true}, expected: "x"},
		{input: BlobPermissions{PermanentDelete: true}, expected: "y"},
		{input: BlobPermissions{List: true}, expected: "l"},
		{input: BlobPermissions{Tag: true}, expected: "t"},
		{input: BlobPermissions{Move: true}, expected: "m"},
		{input: BlobPermissions{Execute: true}, expected: "e"},
		{input: BlobPermissions{Ownership: true}, expected: "o"},
		{input: BlobPermissions{Permissions: true}, expected: "p"},
		{input: BlobPermissions{SetImmutabilityPolicy: true}, expected: "i"},
		{input: BlobPermissions{
			Read:                  true,
			Add:                   true,
			Create:                true,
			Write:                 true,
			Delete:                true,
			DeletePreviousVersion: true,
			PermanentDelete:       true,
			List:                  true,
			Tag:                   true,
			Move:                  true,
			Execute:               true,
			Ownership:             true,
			Permissions:           true,
			SetImmutabilityPolicy: true,
		}, expected: "racwdxyltmeopi"},
	}
	for _, c := range testdata {
		require.Equal(t, c.expected, c.input.String())
	}
}

func TestBlobPermissions_Parse(t *testing.T) {
	testdata := []struct {
		expected BlobPermissions
		input    string
	}{
		{expected: BlobPermissions{Read: true}, input: "r"},
		{expected: BlobPermissions{Add: true}, input: "a"},
		{expected: BlobPermissions{Create: true}, input: "c"},
		{expected: BlobPermissions{Write: true}, input: "w"},
		{expected: BlobPermissions{Delete: true}, input: "d"},
		{expected: BlobPermissions{DeletePreviousVersion: true}, input: "x"},
		{expected: BlobPermissions{PermanentDelete: true}, input: "y"},
		{expected: BlobPermissions{List: true}, input: "l"},
		{expected: BlobPermissions{Tag: true}, input: "t"},
		{expected: BlobPermissions{Move: true}, input: "m"},
		{expected: BlobPermissions{Execute: true}, input: "e"},
		{expected: BlobPermissions{Ownership: true}, input: "o"},
		{expected: BlobPermissions{Permissions: true}, input: "p"},
		{expected: BlobPermissions{SetImmutabilityPolicy: true}, input: "i"},
		{expected: BlobPermissions{
			Read:                  true,
			Add:                   true,
			Create:                true,
			Write:                 true,
			Delete:                true,
			DeletePreviousVersion: true,
			PermanentDelete:       true,
			List:                  true,
			Tag:                   true,
			Move:                  true,
			Execute:               true,
			Ownership:             true,
			Permissions:           true,
			SetImmutabilityPolicy: true,
		}, input: "racwdxyltmeopi"},
		{expected: BlobPermissions{
			Read:                  true,
			Add:                   true,
			Create:                true,
			Write:                 true,
			Delete:                true,
			DeletePreviousVersion: true,
			PermanentDelete:       true,
			List:                  true,
			Tag:                   true,
			Move:                  true,
			Execute:               true,
			Ownership:             true,
			Permissions:           true,
			SetImmutabilityPolicy: true,
		}, input: "apwecxrdlmyiot"}, // Wrong order parses correctly
	}
	for _, c := range testdata {
		permissions, err := parseBlobPermissions(c.input)
		require.Nil(t, err)
		require.Equal(t, c.expected, permissions)
	}
}

func TestBlobPermissions_ParseNegative(t *testing.T) {
	_, err := parseBlobPermissions("apwecxrdlfmyiot") // Here 'f' is invalid
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "102")
}

func TestGetCanonicalName(t *testing.T) {
	testdata := []struct {
		inputAccount   string
		inputContainer string
		inputBlob      string
		inputDirectory string
		expected       string
	}{
		{inputAccount: "fakestorageaccount", inputContainer: "fakestoragecontainer", expected: "/blob/fakestorageaccount/fakestoragecontainer"},
		{inputAccount: "fakestorageaccount", inputContainer: "fakestoragecontainer", inputBlob: "fakestorageblob", expected: "/blob/fakestorageaccount/fakestoragecontainer/fakestorageblob"},
		{inputAccount: "fakestorageaccount", inputContainer: "fakestoragecontainer", inputBlob: "fakestoragevirtualdir/fakestorageblob", expected: "/blob/fakestorageaccount/fakestoragecontainer/fakestoragevirtualdir/fakestorageblob"},
		{inputAccount: "fakestorageaccount", inputContainer: "fakestoragecontainer", inputBlob: "fakestoragevirtualdir\\fakestorageblob", expected: "/blob/fakestorageaccount/fakestoragecontainer/fakestoragevirtualdir/fakestorageblob"},
		{inputAccount: "fakestorageaccount", inputContainer: "fakestoragecontainer", inputBlob: "fakestoragedirectory", expected: "/blob/fakestorageaccount/fakestoragecontainer/fakestoragedirectory"},
		// Directory SAS canonical names (blobName is empty, directoryName is set)
		{inputAccount: "fakestorageaccount", inputContainer: "fakestoragecontainer", inputDirectory: "foo", expected: "/blob/fakestorageaccount/fakestoragecontainer/foo"},
		{inputAccount: "fakestorageaccount", inputContainer: "fakestoragecontainer", inputDirectory: "foo/bar", expected: "/blob/fakestorageaccount/fakestoragecontainer/foo/bar"},
		{inputAccount: "fakestorageaccount", inputContainer: "fakestoragecontainer", inputDirectory: "foo/bar/hello", expected: "/blob/fakestorageaccount/fakestoragecontainer/foo/bar/hello"},
	}
	for _, c := range testdata {
		require.Equal(t, c.expected, getCanonicalName(c.inputAccount, c.inputContainer, c.inputBlob, c.inputDirectory))
	}
}

func TestGetDirectoryDepth(t *testing.T) {
	testdata := []struct {
		input    string
		expected string
	}{
		{input: "", expected: ""},
		{input: "myfile", expected: "1"},
		{input: "mydirectory", expected: "1"},
		{input: "mydirectory/myfile", expected: "2"},
		{input: "mydirectory/mysubdirectory", expected: "2"},
		{input: "foo/bar/hello", expected: "3"},
		{input: "a/b/c/d", expected: "4"},
		{input: "/", expected: "0"},
		{input: "foo/", expected: "1"},
		{input: "/foo", expected: "1"},
		{input: "/foo/bar/", expected: "2"},
	}
	for _, c := range testdata {
		require.Equal(t, c.expected, getDirectoryDepth(c.input))
	}
}

func TestBlobSignatureValues_SignWithSharedKey(t *testing.T) {
	cred, err := exported.NewSharedKeyCredential("fakeaccountname", "AKIAIOSFODNN7EXAMPLE")
	require.Nil(t, err, "error creating valid shared key credentials.")

	expiryDate, err := time.Parse("2006-01-02", "2023-07-20")
	require.Nil(t, err, "error creating valid expiry date.")

	testdata := []struct {
		object        BlobSignatureValues
		expected      QueryParameters
		expectedError error
	}{
		{
			object:        BlobSignatureValues{ContainerName: "fakestoragecontainer", Permissions: "a", ExpiryTime: expiryDate},
			expected:      QueryParameters{version: Version, permissions: "a", expiryTime: expiryDate, resource: "c"},
			expectedError: nil,
		},
		{
			object:        BlobSignatureValues{ContainerName: "fakestoragecontainer", Permissions: "", ExpiryTime: expiryDate},
			expected:      QueryParameters{},
			expectedError: errors.New("service SAS is missing at least one of these: ExpiryTime or Permissions"),
		},
		{
			object:        BlobSignatureValues{ContainerName: "fakestoragecontainer", Permissions: "a", ExpiryTime: *new(time.Time)}, //nolint
			expected:      QueryParameters{},
			expectedError: errors.New("service SAS is missing at least one of these: ExpiryTime or Permissions"),
		},
		{
			object:        BlobSignatureValues{ContainerName: "fakestoragecontainer", Permissions: "", ExpiryTime: *new(time.Time), Identifier: "fakepolicyname"}, //nolint
			expected:      QueryParameters{version: Version, resource: "c", identifier: "fakepolicyname"},
			expectedError: nil,
		},
	}
	for _, c := range testdata {
		act, err := c.object.SignWithSharedKey(cred)
		// ignore signature value
		act.signature = ""
		require.Equal(t, c.expected, act)
		require.Equal(t, c.expectedError, err)
	}
}

func TestFormatSignedRequestHeaders(t *testing.T) {
	testdata := []struct {
		desc                  string
		input                 map[string]string
		expectedNames         string
		expectedCanonicalized string
	}{
		{
			desc:                  "nil map",
			input:                 nil,
			expectedNames:         "",
			expectedCanonicalized: "",
		},
		{
			desc:                  "empty map",
			input:                 map[string]string{},
			expectedNames:         "",
			expectedCanonicalized: "",
		},
		{
			desc:                  "single header",
			input:                 map[string]string{"abra": "cadabra"},
			expectedNames:         "abra",
			expectedCanonicalized: "abra:cadabra\n",
		},
	}
	for _, c := range testdata {
		t.Run(c.desc, func(t *testing.T) {
			names, canonicalized := formatSignedRequestHeaders(c.input)
			require.Equal(t, c.expectedNames, names)
			require.Equal(t, c.expectedCanonicalized, canonicalized)
		})
	}

	// Multi-key test: we can't predict map iteration order, so validate both parts independently
	t.Run("multiple headers", func(t *testing.T) {
		input := map[string]string{"foo": "123", "bar": "456"}
		names, canonicalized := formatSignedRequestHeaders(input)

		// names should contain both keys comma-separated
		require.Contains(t, names, "foo")
		require.Contains(t, names, "bar")
		require.Contains(t, names, ",")

		// canonicalized should contain both key:value pairs each ending with \n
		require.Contains(t, canonicalized, "foo:123\n")
		require.Contains(t, canonicalized, "bar:456\n")
	})
}

func TestFormatSignedRequestQueryParameters(t *testing.T) {
	testdata := []struct {
		desc                  string
		input                 map[string]string
		expectedNames         string
		expectedCanonicalized string
	}{
		{
			desc:                  "nil map",
			input:                 nil,
			expectedNames:         "",
			expectedCanonicalized: "",
		},
		{
			desc:                  "empty map",
			input:                 map[string]string{},
			expectedNames:         "",
			expectedCanonicalized: "",
		},
		{
			desc:                  "single param",
			input:                 map[string]string{"foo": "123"},
			expectedNames:         "foo",
			expectedCanonicalized: "\nfoo:123",
		},
	}
	for _, c := range testdata {
		t.Run(c.desc, func(t *testing.T) {
			names, canonicalized := formatSignedRequestQueryParameters(c.input)
			require.Equal(t, c.expectedNames, names)
			require.Equal(t, c.expectedCanonicalized, canonicalized)
		})
	}

	// Multi-key test: validate both parts independently due to map iteration order
	t.Run("multiple params", func(t *testing.T) {
		input := map[string]string{"hello": "world", "abra": "cadabra"}
		names, canonicalized := formatSignedRequestQueryParameters(input)

		// names should contain both keys comma-separated
		require.Contains(t, names, "hello")
		require.Contains(t, names, "abra")
		require.Contains(t, names, ",")

		// canonicalized should contain both key:value pairs each prefixed with \n
		require.Contains(t, canonicalized, "\nhello:world")
		require.Contains(t, canonicalized, "\nabra:cadabra")
	})
}

func TestBlobSignatureValues_SignWithSharedKey_DirectorySAS(t *testing.T) {
	cred, err := exported.NewSharedKeyCredential("fakeaccountname", "AKIAIOSFODNN7EXAMPLE")
	require.Nil(t, err, "error creating valid shared key credentials.")

	expiryDate, err := time.Parse("2006-01-02", "2023-07-20")
	require.Nil(t, err, "error creating valid expiry date.")

	testdata := []struct {
		desc              string
		object            BlobSignatureValues
		expectedResource  string
		expectedDepth     string
		expectedError     error
		shouldBlobBeEmpty bool
	}{
		{
			desc: "single-level directory",
			object: BlobSignatureValues{
				ContainerName: "fakestoragecontainer",
				Directory:     "foo",
				Permissions:   "rl",
				ExpiryTime:    expiryDate,
			},
			expectedResource:  "d",
			expectedDepth:     "1",
			expectedError:     nil,
			shouldBlobBeEmpty: true,
		},
		{
			desc: "two-level directory",
			object: BlobSignatureValues{
				ContainerName: "fakestoragecontainer",
				Directory:     "foo/bar",
				Permissions:   "rl",
				ExpiryTime:    expiryDate,
			},
			expectedResource:  "d",
			expectedDepth:     "2",
			expectedError:     nil,
			shouldBlobBeEmpty: true,
		},
		{
			desc: "three-level directory",
			object: BlobSignatureValues{
				ContainerName: "fakestoragecontainer",
				Directory:     "foo/bar/hello",
				Permissions:   "racwdxl",
				ExpiryTime:    expiryDate,
			},
			expectedResource:  "d",
			expectedDepth:     "3",
			expectedError:     nil,
			shouldBlobBeEmpty: true,
		},
		{
			desc: "directory SAS clears BlobName",
			object: BlobSignatureValues{
				ContainerName: "fakestoragecontainer",
				BlobName:      "shouldbeignored",
				Directory:     "mydir",
				Permissions:   "rl",
				ExpiryTime:    expiryDate,
			},
			expectedResource:  "d",
			expectedDepth:     "1",
			expectedError:     nil,
			shouldBlobBeEmpty: true,
		},
		{
			desc: "directory SAS with identifier only",
			object: BlobSignatureValues{
				ContainerName: "fakestoragecontainer",
				Directory:     "dir1/dir2",
				Identifier:    "mypolicy",
			},
			expectedResource:  "d",
			expectedDepth:     "2",
			expectedError:     nil,
			shouldBlobBeEmpty: true,
		},
		{
			desc: "directory SAS missing permissions and identifier",
			object: BlobSignatureValues{
				ContainerName: "fakestoragecontainer",
				Directory:     "mydir",
				ExpiryTime:    expiryDate,
			},
			expectedError: errors.New("service SAS is missing at least one of these: ExpiryTime or Permissions"),
		},
	}
	for _, c := range testdata {
		t.Run(c.desc, func(t *testing.T) {
			act, err := c.object.SignWithSharedKey(cred)
			if c.expectedError != nil {
				require.Equal(t, c.expectedError, err)
				return
			}
			require.Nil(t, err)
			require.Equal(t, c.expectedResource, act.Resource())
			require.Equal(t, c.expectedDepth, act.SignedDirectoryDepth())
			require.NotEmpty(t, act.Signature())
		})
	}
}

func TestBlobSignatureValues_SignWithUserDelegation_DirectorySAS(t *testing.T) {
	oid := "fakeoid"
	tid := "faketid"
	svc := "b"
	ver := "2020-02-10"
	val := "AKIAIOSFODNN7EXAMPLE"
	start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	expiry := time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC)

	udc := exported.NewUserDelegationCredential("fakeaccountname", exported.UserDelegationKey{
		SignedOID:     &oid,
		SignedTID:     &tid,
		SignedStart:   &start,
		SignedExpiry:  &expiry,
		SignedService: &svc,
		SignedVersion: &ver,
		Value:         &val,
	})

	expiryDate, err := time.Parse("2006-01-02", "2023-07-20")
	require.Nil(t, err, "error creating valid expiry date.")

	testdata := []struct {
		desc             string
		object           BlobSignatureValues
		expectedResource string
		expectedDepth    string
		expectedError    error
	}{
		{
			desc: "single-level directory UDK",
			object: BlobSignatureValues{
				ContainerName: "fakestoragecontainer",
				Directory:     "foo",
				Permissions:   "rl",
				ExpiryTime:    expiryDate,
			},
			expectedResource: "d",
			expectedDepth:    "1",
			expectedError:    nil,
		},
		{
			desc: "three-level directory UDK",
			object: BlobSignatureValues{
				ContainerName: "fakestoragecontainer",
				Directory:     "foo/bar/hello",
				Permissions:   "racwdxl",
				ExpiryTime:    expiryDate,
			},
			expectedResource: "d",
			expectedDepth:    "3",
			expectedError:    nil,
		},
		{
			desc: "directory UDK clears BlobName",
			object: BlobSignatureValues{
				ContainerName: "fakestoragecontainer",
				BlobName:      "shouldbeignored",
				Directory:     "mydir/subdir",
				Permissions:   "rl",
				ExpiryTime:    expiryDate,
			},
			expectedResource: "d",
			expectedDepth:    "2",
			expectedError:    nil,
		},
		{
			desc: "directory UDK missing permissions",
			object: BlobSignatureValues{
				ContainerName: "fakestoragecontainer",
				Directory:     "mydir",
				ExpiryTime:    expiryDate,
			},
			expectedError: errors.New("user delegation SAS is missing at least one of these: ExpiryTime or Permissions"),
		},
		{
			desc: "directory UDK nil credential",
			object: BlobSignatureValues{
				ContainerName: "fakestoragecontainer",
				Directory:     "mydir",
				Permissions:   "rl",
				ExpiryTime:    expiryDate,
			},
		},
	}
	for _, c := range testdata {
		t.Run(c.desc, func(t *testing.T) {
			var act QueryParameters
			var signErr error
			if c.desc == "directory UDK nil credential" {
				act, signErr = c.object.SignWithUserDelegation(nil)
				require.Error(t, signErr)
				require.Contains(t, signErr.Error(), "cannot sign SAS query without User Delegation Key")
				return
			}
			act, signErr = c.object.SignWithUserDelegation(udc)
			if c.expectedError != nil {
				require.Equal(t, c.expectedError, signErr)
				return
			}
			require.Nil(t, signErr)
			require.Equal(t, c.expectedResource, act.Resource())
			require.Equal(t, c.expectedDepth, act.SignedDirectoryDepth())
			require.NotEmpty(t, act.Signature())
			require.Equal(t, oid, act.SignedOID())
			require.Equal(t, tid, act.SignedTID())
		})
	}
}

func TestDirectorySAS_StringToSign(t *testing.T) {
	// Verify that the string-to-sign for a directory SAS uses the directory path in the canonical name
	// and that the resource is "d"
	cred, err := exported.NewSharedKeyCredential("fakeaccountname", "AKIAIOSFODNN7EXAMPLE")
	require.Nil(t, err)

	expiryDate := time.Date(2023, 7, 20, 0, 0, 0, 0, time.UTC)

	// Sign two directory SAS tokens with different depths and verify they produce different signatures
	qp1, err := BlobSignatureValues{
		ContainerName: "container",
		Directory:     "foo/bar/hello",
		Permissions:   "rl",
		ExpiryTime:    expiryDate,
	}.SignWithSharedKey(cred)
	require.Nil(t, err)
	require.Equal(t, "d", qp1.Resource())
	require.Equal(t, "3", qp1.SignedDirectoryDepth())

	qp2, err := BlobSignatureValues{
		ContainerName: "container",
		Directory:     "foo/bar",
		Permissions:   "rl",
		ExpiryTime:    expiryDate,
	}.SignWithSharedKey(cred)
	require.Nil(t, err)
	require.Equal(t, "d", qp2.Resource())
	require.Equal(t, "2", qp2.SignedDirectoryDepth())

	// Different directories should produce different signatures
	require.NotEqual(t, qp1.Signature(), qp2.Signature())
}

func TestDirectorySAS_EncodeRoundTrip(t *testing.T) {
	// Test that a directory SAS token round-trips correctly through Encode/NewQueryParameters
	cred, err := exported.NewSharedKeyCredential("fakeaccountname", "AKIAIOSFODNN7EXAMPLE")
	require.Nil(t, err)

	expiryDate := time.Date(2023, 7, 20, 0, 0, 0, 0, time.UTC)

	original, err := BlobSignatureValues{
		ContainerName: "container",
		Directory:     "foo/bar/hello",
		Permissions:   "rl",
		ExpiryTime:    expiryDate,
	}.SignWithSharedKey(cred)
	require.Nil(t, err)

	// Encode to query string
	encoded := original.Encode()
	require.Contains(t, encoded, "sr=d")
	require.Contains(t, encoded, "sdd=3")

	// Parse back from query string
	parsedURL := "https://fakeaccount.blob.core.windows.net/container/foo/bar/hello?" + encoded
	u, err := url.Parse(parsedURL)
	require.Nil(t, err)

	parsed := NewQueryParameters(u.Query(), false)
	require.Equal(t, "d", parsed.Resource())
	require.Equal(t, "3", parsed.SignedDirectoryDepth())
	require.Equal(t, original.Signature(), parsed.Signature())
	require.Equal(t, original.Permissions(), parsed.Permissions())
	require.Equal(t, original.Version(), parsed.Version())
}
