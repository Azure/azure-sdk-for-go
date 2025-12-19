// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package cache

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity/cache/internal/aescbc"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity/cache/internal/jwe"
	"github.com/AzureAD/microsoft-authentication-extensions-for-go/cache/accessor"
	"golang.org/x/sys/unix"
)

const (
	keySize = 32
	userKey = "user"
)

var (
	cacheDir = func() (d string, err error) {
		if d = os.Getenv("XDG_CACHE_HOME"); d == "" {
			if d, err = os.UserHomeDir(); err == nil {
				d = filepath.Join(d, ".cache")
			}
		}
		return
	}
	storage = func(name string) (accessor.Accessor, error) {
		return newKeyring(name)
	}
)

// keyring encrypts cache data with a key stored on the user keyring and writes the encrypted
// data to a file. The encryption key, and thus the data, is lost when the system shuts down.
type keyring struct {
	description, file string
	key               []byte
	keyID, ringID     int
}

func newKeyring(name string) (*keyring, error) {
	p, err := cacheFilePath(name)
	if err != nil {
		return nil, err
	}
	// the user keyring is available to all processes owned by the user whereas the user
	// *session* keyring is available only to processes in the current session i.e. shell
	ringID, err := unix.KeyctlGetKeyringID(unix.KEY_SPEC_USER_KEYRING, true)
	if err != nil {
		return nil, fmt.Errorf("couldn't get the user keyring due to error %q", err)
	}
	// Link the session keyring to the user keyring so the process possesses any key[ring] it links
	// to the user keyring and thereby has permission to read/write/search them (see the "Possession"
	// section of the keyrings man page). This step isn't always necessary but in some cases prevents
	// weirdness such as a process adding keys it can't read. Ignore errors because failure here
	// doesn't guarantee this process can't perform all required operations on the user keyring.
	if sessionID, err := unix.KeyctlGetKeyringID(unix.KEY_SPEC_SESSION_KEYRING, true); err == nil {
		_, _ = unix.KeyctlInt(unix.KEYCTL_LINK, ringID, sessionID, 0, 0)
	}
	// Attempt to link a persistent keyring to the user keyring. This keyring is persistent in that
	// its linked keys survive all the user's login sessions being deleted but like all user keys,
	// they exist only in memory and are therefore lost on system shutdown. If the attempt fails
	// (some systems don't support persistent keyrings) continue with the user keyring.
	if persistentRing, err := unix.KeyctlInt(unix.KEYCTL_GET_PERSISTENT, -1, ringID, 0, 0); err == nil {
		ringID = persistentRing
	}
	return &keyring{description: name, file: p, ringID: ringID}, nil
}

func (k *keyring) Delete(context.Context) error {
	if k.keyID != 0 && k.ringID != 0 {
		_, err := unix.KeyctlInt(unix.KEYCTL_UNLINK, k.keyID, k.ringID, 0, 0)
		if err != nil && !isKeyInvalidOrNotFound(err) {
			return fmt.Errorf("failed to delete cache data due to error %q", err)
		}
	}
	err := os.Remove(k.file)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func (k *keyring) Read(context.Context) ([]byte, error) {
	b, err := os.ReadFile(k.file)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read cache data due to error %q", err)
	}
	if len(b) == 0 {
		return nil, nil
	}
	j, err := jwe.ParseCompactFormat(b)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse cache data due to error %q", err)
	}
	plaintext, err := k.decrypt(j)
	return plaintext, err
}

func (k *keyring) Write(_ context.Context, data []byte) error {
	if len(data) == 0 {
		return nil
	}
	j, err := k.encrypt(data)
	if err != nil {
		return err
	}
	content, err := j.Serialize()
	if err != nil {
		return fmt.Errorf("couldn't serialize cache data due to error %q", err)
	}
	err = os.WriteFile(k.file, []byte(content), 0600)
	if errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(filepath.Dir(k.file), 0700)
		if err == nil {
			err = os.WriteFile(k.file, []byte(content), 0600)
		}
	}
	return err
}

func (k *keyring) createKey() ([]byte, error) {
	// allocate an extra byte because keyring payloads must have a null terminator
	key := make([]byte, keySize+1)
	_, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("couldn't create cache encryption key due to error %q", err)
	}
	key[keySize] = 0
	id, err := unix.AddKey(userKey, k.description, key, k.ringID)
	if err != nil {
		return nil, fmt.Errorf("couldn't store cache encryption key due to error %q", err)
	}
	k.key = key[:keySize]
	k.keyID = id
	return k.key, nil
}

func (k *keyring) decrypt(j jwe.JWE) ([]byte, error) {
	for tries := 0; tries < 2; tries++ {
		key, err := k.getKey()
		if err != nil {
			if err == unix.ENOKEY {
				return nil, nil
			}
			return nil, err
		}
		plaintext, err := j.Decrypt(key)
		if err == nil {
			return plaintext, nil
		}
		// try again, getting the key from the keyring first in case it was overwritten
		// by the user (with keyctl) or another process (in a Write() race)
		k.key = nil
		k.keyID = 0
	}
	// data is unreadable; the next Write will overwrite the file
	return nil, nil
}

func (k *keyring) encrypt(data []byte) (jwe.JWE, error) {
	key, err := k.getKey()
	if isKeyInvalidOrNotFound(err) {
		key, err = k.createKey()
	}
	if err != nil {
		return jwe.JWE{}, fmt.Errorf("couldn't get cache encryption key due to error %q", err)
	}
	alg, err := aescbc.NewAES128CBCHMACSHA256(key)
	if err != nil {
		return jwe.JWE{}, err
	}
	return jwe.Encrypt(data, fmt.Sprint(k.keyID), alg)
}

func (k *keyring) getKey() ([]byte, error) {
	if k.key != nil {
		// we created, or got, the key earlier
		return k.key, nil
	}
	if k.keyID == 0 {
		// search for a key matching the description i.e. the cache name
		keyID, err := unix.KeyctlSearch(k.ringID, userKey, k.description, 0)
		if err != nil {
			return nil, err
		}
		k.keyID = keyID
	}
	pl := make([]byte, keySize+1) // extra byte for the payload's null terminator
	_, err := unix.KeyctlBuffer(unix.KEYCTL_READ, k.keyID, pl, 0)
	if err != nil {
		return nil, err
	}
	k.key = pl[:keySize]
	return k.key, nil
}

func isKeyInvalidOrNotFound(err error) bool {
	return errors.Is(err, unix.EKEYEXPIRED) || errors.Is(err, unix.EKEYREVOKED) || errors.Is(err, unix.ENOENT) || errors.Is(err, unix.ENOKEY)
}

var _ accessor.Accessor = (*keyring)(nil)
