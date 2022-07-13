package azkeys

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"errors"
	"math/big"
)

// Public returns the public key included the object, as a crypto.PublicKey object.
// This method returns an error if it's invoked on a JSONWebKey object representing a symmetric key.
func (key JSONWebKey) Public() (crypto.PublicKey, error) {
	if key.Kty == nil {
		return nil, errors.New("property Kty is nil")
	}

	switch {
	case key.Kty.IsRSAKey():
		return key.publicRSA()
	case key.Kty.IsECKey():
		return key.publicEC()
	}

	return nil, errors.New("unsupported key type")
}

func (key JSONWebKey) publicRSA() (*rsa.PublicKey, error) {
	res := &rsa.PublicKey{}

	// N = modulus
	if len(key.N) == 0 {
		return nil, errors.New("property N is empty")
	}
	res.N = &big.Int{}
	res.N.SetBytes(key.N)

	// e = public exponent
	if len(key.E) == 0 {
		return nil, errors.New("property e is empty")
	}
	res.E = int(big.NewInt(0).SetBytes(key.E).Uint64())

	return res, nil
}

func (key JSONWebKey) publicEC() (*ecdsa.PublicKey, error) {
	res := &ecdsa.PublicKey{}

	if key.Crv == nil {
		return nil, errors.New("property Crv is nil")
	}
	switch *key.Crv {
	case JSONWebKeyCurveNameP256:
		res.Curve = elliptic.P256()
	case JSONWebKeyCurveNameP384:
		res.Curve = elliptic.P384()
	case JSONWebKeyCurveNameP521:
		res.Curve = elliptic.P521()
	case JSONWebKeyCurveNameP256K:
		return nil, errors.New("curves of type P-256K are not supported by this method")
	}

	// X coordinate
	if len(key.X) == 0 {
		return nil, errors.New("property X is empty")
	}
	res.X = &big.Int{}
	res.X.SetBytes(key.X)

	// Y coordinate
	if len(key.Y) == 0 {
		return nil, errors.New("property Y is empty")
	}
	res.Y = &big.Int{}
	res.Y.SetBytes(key.Y)

	return res, nil
}

// Private returns the private key included in the object, as a crypto.PrivateKey object.
// This method returns an error if it's invoked on a JSONWebKey object representing a symmetric key.
func (key JSONWebKey) Private() (crypto.PrivateKey, error) {
	if key.Kty == nil {
		return nil, errors.New("property Kty is nil")
	}

	switch {
	case key.Kty.IsRSAKey():
		return key.privateRSA()
	case key.Kty.IsECKey():
		return key.privateEC()
	}

	return nil, errors.New("unsupported key type")
}

func (key JSONWebKey) privateRSA() (*rsa.PrivateKey, error) {
	// Parse the public part first
	pk, err := key.publicRSA()
	if err != nil {
		return nil, err
	}

	res := &rsa.PrivateKey{
		PublicKey: *pk,
	}

	// D = private exponent
	if len(key.D) == 0 {
		return nil, errors.New("property D is empty")
	}
	res.D = &big.Int{}
	res.D.SetBytes(key.D)

	// Precomputed values
	// These are optional
	if len(key.DP) > 0 {
		res.Precomputed.Dp = &big.Int{}
		res.Precomputed.Dp.SetBytes(key.DP)
	}
	if len(key.DQ) > 0 {
		res.Precomputed.Dq = &big.Int{}
		res.Precomputed.Dq.SetBytes(key.DQ)
	}
	if len(key.QI) > 0 {
		res.Precomputed.Qinv = &big.Int{}
		res.Precomputed.Qinv.SetBytes(key.QI)
	}

	return res, nil
}

func (key JSONWebKey) privateEC() (*ecdsa.PrivateKey, error) {
	// Parse the public part first
	pk, err := key.publicEC()
	if err != nil {
		return nil, err
	}

	res := &ecdsa.PrivateKey{
		PublicKey: *pk,
	}

	// D = private key
	if len(key.D) == 0 {
		return nil, errors.New("property D is empty")
	}
	res.D = &big.Int{}
	res.D.SetBytes(key.D)

	return res, nil
}
