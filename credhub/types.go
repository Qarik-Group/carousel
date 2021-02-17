package credhub

type Value interface {
	Inspect() string
	Identidy() string
	ParentIdentiy() string
}

type Credential struct {
	ID                   string      `json:"id"`
	Metadata             interface{} `json:"metadata"`
	Name                 string      `json:"name"`
	Type                 string      `json:"type"`
	VersionCreatedAt     string      `json:"version_created_at"`
	CertificateAuthority bool        `json:"certificate_authority,omitempty"`
	ExpiryDate           string      `json:"expiry_date,omitempty"`
	Generated            bool        `json:"generated,omitempty"`
	SelfSigned           bool        `json:"self_signed,omitempty"`
	Transitional         bool        `json:"transitional,omitempty"`
	Value                Value       `json:"value"`
}

type PasswordValue string

type CertificateValue struct {
	Ca          string `json:"ca"`
	Certificate string `json:"certificate"`
	PrivateKey  string `json:"private_key"`
}

type SSHValue struct {
	PrivateKey           string `json:"private_key"`
	PublicKey            string `json:"public_key"`
	PublicKeyFingerprint string `json:"public_key_fingerprint"`
}

type RSAValue struct {
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
}

type UserValue struct {
	Password     string `json:"password"`
	PasswordHash string `json:"password_hash"`
	Username     string `json:"username"`
}
