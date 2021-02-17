package credhub

import (
	"encoding/json"
)

type CredentialType string

const (
	Certificate, SSH, RSA, Password, User, Value, JSON CredentialType = "certificate",
		"ssh", "rsa", "password", "user", "value", "json"
)

type Credential struct {
	ID                   string          `json:"id"`
	Metadata             Metadata        `json:"metadata"`
	Name                 string          `json:"name"`
	Type                 CredentialType  `json:"type"`
	VersionCreatedAt     string          `json:"version_created_at"`
	CertificateAuthority bool            `json:"certificate_authority,omitempty"`
	ExpiryDate           string          `json:"expiry_date,omitempty"`
	Generated            bool            `json:"generated,omitempty"`
	SelfSigned           bool            `json:"self_signed,omitempty"`
	Transitional         bool            `json:"transitional,omitempty"`
	RawValue             json.RawMessage `json:"value"`

	Ca                   string                 `json:"_"`
	Certificate          string                 `json:"_"`
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
	case Certificate, SSH, RSA, User:
		v := rawValue{}
		if err := json.Unmarshal(c.RawValue, &v); err != nil {
			return err
		}
		c.Ca = v.Ca
		c.Certificate = v.Certificate
		c.PrivateKey = v.PrivateKey
		c.PublicKey = v.PublicKey
		c.PublicKeyFingerprint = v.PublicKeyFingerprint
		c.Password = v.Password
		c.PasswordHash = v.PasswordHash
		c.Username = v.Username
	}

	return nil
}
