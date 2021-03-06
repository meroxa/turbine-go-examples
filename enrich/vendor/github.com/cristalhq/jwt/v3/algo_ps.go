package jwt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
)

// NewSignerPS returns a new RSA-PSS-based signer.
func NewSignerPS(alg Algorithm, key *rsa.PrivateKey) (Signer, error) {
	if key == nil {
		return nil, ErrNilKey
	}
	hash, opts, err := getParamsPS(alg, key.Size())
	if err != nil {
		return nil, err
	}
	return &psAlg{
		alg:        alg,
		hash:       hash,
		privateKey: key,
		opts:       opts,
	}, nil
}

// NewVerifierPS returns a new RSA-PSS-based signer.
func NewVerifierPS(alg Algorithm, key *rsa.PublicKey) (Verifier, error) {
	if key == nil {
		return nil, ErrNilKey
	}
	hash, opts, err := getParamsPS(alg, key.Size())
	if err != nil {
		return nil, err
	}
	return &psAlg{
		alg:       alg,
		hash:      hash,
		publicKey: key,
		opts:      opts,
	}, nil
}

func getParamsPS(alg Algorithm, size int) (crypto.Hash, *rsa.PSSOptions, error) {
	var hash crypto.Hash
	var opts *rsa.PSSOptions
	switch alg {
	case PS256:
		hash, opts = crypto.SHA256, optsPS256
	case PS384:
		hash, opts = crypto.SHA384, optsPS384
	case PS512:
		hash, opts = crypto.SHA512, optsPS512
	default:
		return 0, nil, ErrUnsupportedAlg
	}

	if alg.keySize() != size {
		return 0, nil, ErrInvalidKey
	}
	return hash, opts, nil
}

var (
	optsPS256 = &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthAuto,
		Hash:       crypto.SHA256,
	}

	optsPS384 = &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthAuto,
		Hash:       crypto.SHA384,
	}

	optsPS512 = &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthAuto,
		Hash:       crypto.SHA512,
	}
)

type psAlg struct {
	alg        Algorithm
	hash       crypto.Hash
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
	opts       *rsa.PSSOptions
}

func (ps *psAlg) SignSize() int {
	return ps.privateKey.Size()
}

func (ps *psAlg) Algorithm() Algorithm {
	return ps.alg
}

func (ps *psAlg) Sign(payload []byte) ([]byte, error) {
	digest, err := hashPayload(ps.hash, payload)
	if err != nil {
		return nil, err
	}

	signature, errSign := rsa.SignPSS(rand.Reader, ps.privateKey, ps.hash, digest, ps.opts)
	if errSign != nil {
		return nil, errSign
	}
	return signature, nil
}

func (ps *psAlg) VerifyToken(token *Token) error {
	if constTimeAlgEqual(token.Header().Algorithm, ps.alg) {
		return ps.Verify(token.Payload(), token.Signature())
	}
	return ErrAlgorithmMismatch
}

func (ps *psAlg) Verify(payload, signature []byte) error {
	digest, err := hashPayload(ps.hash, payload)
	if err != nil {
		return err
	}

	errVerify := rsa.VerifyPSS(ps.publicKey, ps.hash, digest, signature, ps.opts)
	if errVerify != nil {
		return ErrInvalidSignature
	}
	return nil
}
