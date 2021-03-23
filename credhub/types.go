package credhub

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"time"
)

//
//go:generate go run github.com/alvaroloes/enumer -type=CredentialType -json -transform=snake

type CredentialType int

const (
	Certificate CredentialType = iota
	SSH
	RSA
	Password
	User
	Value
	JSON
)

func CredentialTypeStringValues() []string {
	out := make([]string, 0)
	for _, c := range CredentialTypeValues() {
		out = append(out, c.String())
	}
	return out
}

type Credential struct {
	ID                   string          `json:"id"`
	Metadata             Metadata        `json:"metadata,omitempty"`
	Name                 string          `json:"name"`
	Type                 CredentialType  `json:"type"`
	VersionCreatedAt     *time.Time      `json:"version_created_at"`
	CertificateAuthority bool            `json:"certificate_authority,omitempty"`
	ExpiryDate           *time.Time      `json:"expiry_date,omitempty"`
	Generated            bool            `json:"generated,omitempty"`
	SelfSigned           bool            `json:"self_signed,omitempty"`
	Transitional         bool            `json:"transitional,omitempty"`
	RawValue             json.RawMessage `json:"value,omitempty"`

	Ca                   []*x509.Certificate    `json:"_"`
	PEMCa                string                 `json:"_"`
	Certificate          *x509.Certificate      `json:"_"`
	PEMCertificate       string                 `json:"_"`
	PrivateKey           string                 `json:"_"`
	PublicKey            string                 `json:"_"`
	PublicKeyFingerprint string                 `json:"_"`
	Password             string                 `json:"_"`
	PasswordHash         string                 `json:"_"`
	Username             string                 `json:"_"`
	JSON                 map[string]interface{} `json:"_"`
	Value                string                 `json:"_"`
}

type Metadata map[string]string

type rawValue struct {
	Ca                   string `json:"ca,omitempty"`
	Certificate          string `json:"certificate,omitempty"`
	PrivateKey           string `json:"private_key,omitempty"`
	PublicKey            string `json:"public_key,omitempty"`
	PublicKeyFingerprint string `json:"public_key_fingerprint,omitempty"`
	Password             string `json:"password,omitempty"`
	PasswordHash         string `json:"password_hash,omitempty"`
	Username             string `json:"username,omitempty"`
}

func (c *Credential) UnmarshalJSON(b []byte) error {
	type credential Credential

	if err := json.Unmarshal(b, (*credential)(c)); err != nil {
		return err
	}

	switch c.Type {
	case Value:
		if err := json.Unmarshal(c.RawValue, &c.Value); err != nil {
			return err
		}
	case JSON:
		if err := json.Unmarshal(c.RawValue, &c.JSON); err != nil {
			return err
		}
	case Password:
		if err := json.Unmarshal(c.RawValue, &c.Password); err != nil {
			return err
		}
	case Certificate:
		v := rawValue{}
		if err := json.Unmarshal(c.RawValue, &v); err != nil {
			return err
		}
		ca, err := parseCAs(v.Ca)
		if err != nil {
			return err
		}
		c.Ca = ca
		c.PEMCa = v.Ca
		cert, err := parseCertificate(v.Certificate)
		if err != nil {
			return err
		}
		c.Certificate = cert
		c.PEMCertificate = v.Certificate
		c.PrivateKey = v.PrivateKey

	case SSH, RSA, User:
		v := rawValue{}
		if err := json.Unmarshal(c.RawValue, &v); err != nil {
			return err
		}

		c.PrivateKey = v.PrivateKey
		c.PublicKey = v.PublicKey
		c.PublicKeyFingerprint = v.PublicKeyFingerprint
		c.Password = v.Password
		c.PasswordHash = v.PasswordHash
		c.Username = v.Username
	}

	return nil
}

func parseCAs(raw string) ([]*x509.Certificate, error) {
	out := make([]*x509.Certificate, 0)
	for cb, r := pem.Decode([]byte(raw)); cb != nil; cb, r = pem.Decode([]byte(r)) {
		cert, err := x509.ParseCertificate(cb.Bytes)
		if err != nil {
			return nil, err
		}
		out = append(out, cert)
	}
	return out, nil
}

func parseCertificate(raw string) (*x509.Certificate, error) {
	certBlock, _ := pem.Decode([]byte(raw))
	return x509.ParseCertificate(certBlock.Bytes)
}

func (c *Credential) ToStaticVariable() interface{} {
	switch c.Type {
	case Password:
		return c.Password
	case Value, JSON:
		return c.Value
	default:
		return map[interface{}]interface{}{
			"ca":                     c.PEMCa,
			"certificate":            c.PEMCertificate,
			"private_key":            c.PrivateKey,
			"public_key":             c.PublicKey,
			"public_key_fingerprint": c.PublicKeyFingerprint,
			"password":               c.Password,
			"password_hash":          c.PasswordHash,
			"username":               c.Username,
		}
	}
}

func (c *Credential) summary() string {
	return fmt.Sprintf("name: %s\nversion: %s\ncreated_at: %s",
		c.Name, c.ID, c.VersionCreatedAt.Format(time.RFC3339))
}

func (c *Credential) ToStaticVariableMetaOnly() interface{} {
	switch c.Type {
	case Certificate:
		return map[interface{}]interface{}{
			"ca":          c.summary(),
			"certificate": c.summary(),
			"private_key": c.summary(),
		}
	case SSH:
		return map[interface{}]interface{}{
			"private_key":            c.summary(),
			"public_key":             c.summary(),
			"public_key_fingerprint": c.summary(),
		}
	case RSA:
		return map[interface{}]interface{}{
			"private_key":            c.summary(),
			"public_key":             c.summary(),
			"public_key_fingerprint": c.summary(),
		}
	case User:
		return map[interface{}]interface{}{
			"password":      c.summary(),
			"password_hash": c.summary(),
			"username":      c.summary(),
		}
	default:
		return c.summary()
	}
}
