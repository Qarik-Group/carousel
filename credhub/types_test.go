package credhub_test

import (
	"encoding/json"

	. "github.com/starkandwayne/carousel/credhub"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Types", func() {
	Describe("Credential", func() {
		Specify("when decoding and encoding Certificate", func() {
			var cred Credential

			credJSON := `{
	"id": "some-id",
	"name": "/example-certificate",
	"type": "certificate",
	"value": {
		"ca": "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
		"certificate": "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
		"private_key": "-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----"
	},
	"metadata": {"some":"metadata"},
	"version_created_at": "2017-01-01T04:07:18Z",
	"certificate_authority": true,
	"expiry_date": "2018-01-01T04:07:18Z",
	"generated": true,
	"self_signed": true,
	"transitional": true
}`

			err := json.Unmarshal([]byte(credJSON), &cred)

			Expect(err).To(BeNil())

			Expect(cred.ID).To(Equal("some-id"))
			Expect(cred.Name).To(Equal("/example-certificate"))
			Expect(cred.Type).To(Equal(Certificate))
			Expect(cred.Metadata).To(Equal(Metadata{"some": "metadata"}))
			Expect(cred.CertificateAuthority).To(Equal(true))
			Expect(cred.ExpiryDate).To(Equal("2018-01-01T04:07:18Z"))
			Expect(cred.Generated).To(Equal(true))
			Expect(cred.SelfSigned).To(Equal(true))
			Expect(cred.Transitional).To(Equal(true))
			Expect(cred.Ca).To(Equal("-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----"))
			Expect(cred.Certificate).To(Equal("-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----"))
			Expect(cred.PrivateKey).To(Equal("-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----"))
			Expect(cred.VersionCreatedAt).To(Equal("2017-01-01T04:07:18Z"))

			jsonOutput, err := json.Marshal(cred)
			Expect(err).NotTo(HaveOccurred())
			Expect(jsonOutput).To(MatchJSON(credJSON))
		})

		Specify("when decoding and encoding SSH", func() {
			var cred Credential

			credJSON := `{
	"id": "some-id",
	"name": "/example-ssh",
	"type": "ssh",
	"value": {
		"public_key": "-----BEGIN RSA PUBLIC KEY-----\n...\n-----END RSA PUBLIC KEY-----",
		"private_key": "-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----",
		"public_key_fingerprint": "some fingerprint"
	},
	"metadata": {"some":"metadata"},
	"version_created_at": "2017-01-01T04:07:18Z"
}`

			err := json.Unmarshal([]byte(credJSON), &cred)

			Expect(err).To(BeNil())

			Expect(cred.ID).To(Equal("some-id"))
			Expect(cred.Name).To(Equal("/example-ssh"))
			Expect(cred.Type).To(Equal(SSH))
			Expect(cred.Metadata).To(Equal(Metadata{"some": "metadata"}))
			Expect(cred.PrivateKey).To(Equal("-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----"))
			Expect(cred.PublicKey).To(Equal("-----BEGIN RSA PUBLIC KEY-----\n...\n-----END RSA PUBLIC KEY-----"))
			Expect(cred.PublicKeyFingerprint).To(Equal("some fingerprint"))
			Expect(cred.VersionCreatedAt).To(Equal("2017-01-01T04:07:18Z"))

			jsonOutput, err := json.Marshal(cred)
			Expect(err).NotTo(HaveOccurred())
			Expect(jsonOutput).To(MatchJSON(credJSON))
		})

		Specify("when decoding and encoding RSA", func() {
			var cred Credential

			credJSON := `{
	"id": "some-id",
	"name": "/example-rsa",
	"type": "rsa",
	"value": {
		"public_key": "-----BEGIN RSA PUBLIC KEY-----\n...\n-----END RSA PUBLIC KEY-----",
		"private_key": "-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----"
	},
	"metadata": {"some":"metadata"},
	"version_created_at": "2017-01-01T04:07:18Z"
}`

			err := json.Unmarshal([]byte(credJSON), &cred)

			Expect(err).To(BeNil())

			Expect(cred.ID).To(Equal("some-id"))
			Expect(cred.Name).To(Equal("/example-rsa"))
			Expect(cred.Type).To(Equal(RSA))
			Expect(cred.Metadata).To(Equal(Metadata{"some": "metadata"}))
			Expect(cred.PrivateKey).To(Equal("-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----"))
			Expect(cred.PublicKey).To(Equal("-----BEGIN RSA PUBLIC KEY-----\n...\n-----END RSA PUBLIC KEY-----"))
			Expect(cred.VersionCreatedAt).To(Equal("2017-01-01T04:07:18Z"))

			jsonOutput, err := json.Marshal(cred)
			Expect(err).NotTo(HaveOccurred())
			Expect(jsonOutput).To(MatchJSON(credJSON))
		})

		Specify("when decoding and encoding Password", func() {
			var cred Credential

			credJSON := `{
	"id": "some-id",
	"name": "/example-password",
	"type": "password",
	"value": "some-password",
	"metadata": {"some":"metadata"},
	"version_created_at": "2017-01-01T04:07:18Z"
}`

			err := json.Unmarshal([]byte(credJSON), &cred)

			Expect(err).To(BeNil())

			Expect(cred.ID).To(Equal("some-id"))
			Expect(cred.Name).To(Equal("/example-password"))
			Expect(cred.Type).To(Equal(Password))
			Expect(cred.Metadata).To(Equal(Metadata{"some": "metadata"}))
			Expect(cred.Password).To(Equal("some-password"))
			Expect(cred.VersionCreatedAt).To(Equal("2017-01-01T04:07:18Z"))

			jsonOutput, err := json.Marshal(cred)
			Expect(err).NotTo(HaveOccurred())
			Expect(jsonOutput).To(MatchJSON(credJSON))
		})

		Specify("when decoding and encoding User", func() {
			var cred Credential

			credJSON := `{
	"id": "some-id",
	"name": "/example-user",
	"type": "user",
	"value": {
	  "username": "some-username",
	  "password": "some-password",
	  "password_hash": "foQzXY.HaydB."
	},
	"metadata": {"some":"metadata"},
	"version_created_at": "2017-01-01T04:07:18Z"
}`

			err := json.Unmarshal([]byte(credJSON), &cred)

			Expect(err).To(BeNil())

			Expect(cred.ID).To(Equal("some-id"))
			Expect(cred.Name).To(Equal("/example-user"))
			Expect(cred.Type).To(Equal(User))
			Expect(cred.Metadata).To(Equal(Metadata{"some": "metadata"}))
			Expect(cred.Password).To(Equal("some-password"))
			Expect(cred.Username).To(Equal("some-username"))
			Expect(cred.PasswordHash).To(Equal("foQzXY.HaydB."))
			Expect(cred.VersionCreatedAt).To(Equal("2017-01-01T04:07:18Z"))

			jsonOutput, err := json.Marshal(cred)
			Expect(err).NotTo(HaveOccurred())
			Expect(jsonOutput).To(MatchJSON(credJSON))
		})

		Specify("when decoding and encoding Value", func() {
			var cred Credential

			credJSON := `{
	"id": "some-id",
	"name": "/example-value",
	"type": "value",
	"value": "some-value",
	"metadata": {"some":"metadata"},
	"version_created_at": "2017-01-01T04:07:18Z"
}`

			err := json.Unmarshal([]byte(credJSON), &cred)

			Expect(err).To(BeNil())

			Expect(cred.ID).To(Equal("some-id"))
			Expect(cred.Name).To(Equal("/example-value"))
			Expect(cred.Type).To(Equal(Value))
			Expect(cred.Metadata).To(Equal(Metadata{"some": "metadata"}))
			Expect(cred.Value).To(Equal("some-value"))
			Expect(cred.VersionCreatedAt).To(Equal("2017-01-01T04:07:18Z"))

			jsonOutput, err := json.Marshal(cred)
			Expect(err).NotTo(HaveOccurred())
			Expect(jsonOutput).To(MatchJSON(credJSON))
		})

		Specify("when decoding and encoding JSON", func() {
			var cred Credential

			credJSON := `{
	"id": "some-id",
	"name": "/example-json",
	"type": "json",
	"value": {
	  "foo": "bar"
	},
	"metadata": {"some":"metadata"},
	"version_created_at": "2017-01-01T04:07:18Z"
}`

			err := json.Unmarshal([]byte(credJSON), &cred)

			Expect(err).To(BeNil())

			Expect(cred.ID).To(Equal("some-id"))
			Expect(cred.Name).To(Equal("/example-json"))
			Expect(cred.Type).To(Equal(JSON))
			Expect(cred.Metadata).To(Equal(Metadata{"some": "metadata"}))
			Expect(cred.JSON).To(Equal(map[string]interface{}{"foo": "bar"}))
			Expect(cred.VersionCreatedAt).To(Equal("2017-01-01T04:07:18Z"))

			jsonOutput, err := json.Marshal(cred)
			Expect(err).NotTo(HaveOccurred())
			Expect(jsonOutput).To(MatchJSON(credJSON))
		})

	})
})
