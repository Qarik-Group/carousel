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
	  "ca": "-----BEGIN CERTIFICATE-----\nMIIDCjCCAfKgAwIBAgIUYjrEAQDlq1eWqQbCKndS9kI1l/AwDQYJKoZIhvcNAQEL\nBQAwFjEUMBIGA1UEAxMLZXhhbXBsZS5jb20wHhcNMjEwMjE4MDkwNjEyWhcNMjIw\nMjE4MDkwNjEyWjAWMRQwEgYDVQQDEwtleGFtcGxlLmNvbTCCASIwDQYJKoZIhvcN\nAQEBBQADggEPADCCAQoCggEBAIx75lo/fol53qVExXbQaAJ1qdv/NIcPezlEdmWW\n44tb+pE12j+gd2PP7+iPif7eMT6U3n5DjN4q/VyPI8ebwb4LU4Blz2MGLbI/hiA2\n0stFteLR66tP35gODo75s0WYVjTYqE39rTXrErEMUvWl8q0MbqRKGWj4+cEywgVy\nW+5jcDDI5t9CKYYt/IHqMX3r0b7Pwcjp60ozTFxKWSoXQlDz1szw0g+jpJNKmlM7\nEJ0Mm8XElncph8beCTk2exRTxb3fvy9oIWA8Kud4HxM9ZxTKHoV23dROL2uPQoUx\ngSHt2FWM69NA81zkF575YuJV5+mmlpHpVXIAXGhKWLxSQBMCAwEAAaNQME4wHQYD\nVR0OBBYEFIkF/e3zb/wfLRU3X3Va4dFSX1r7MB8GA1UdIwQYMBaAFIkF/e3zb/wf\nLRU3X3Va4dFSX1r7MAwGA1UdEwEB/wQCMAAwDQYJKoZIhvcNAQELBQADggEBAGVt\nT7kQpflQJIwb8QydU04Q0CQJ+O2sTMf2Wmbe/+73mRbkzAhD0oKCkvK3TJ4Xl89O\n5tCCmCIS+rsF3iepS+EIrjA/cZ39Zgo3/B39IMvEyL96GSXCeuHgWys7yNHuDvmh\n9qK0eZ4YEfl9mU57lG8EeP2BVLE2RoAKWbzNanDPJkXLvoUxdphpj6Ne9GxKkl9g\nXqmggEFqlw7G8nScJT/RYK0h0QmaGZ7TLCZ6yNyUki+Ps3S15h4xaxnxQcAp0Udj\npxeZG0vcqJ/5gLszr/llaBw4Rv/ysDev43IRmyY2erpal6MUbk++1Hmo7uifMgK6\nFsDWROtc+z5HPoZZgm8=\n-----END CERTIFICATE-----\n",
	  "certificate": "-----BEGIN CERTIFICATE-----\nMIIDCjCCAfKgAwIBAgIUYjrEAQDlq1eWqQbCKndS9kI1l/AwDQYJKoZIhvcNAQEL\nBQAwFjEUMBIGA1UEAxMLZXhhbXBsZS5jb20wHhcNMjEwMjE4MDkwNjEyWhcNMjIw\nMjE4MDkwNjEyWjAWMRQwEgYDVQQDEwtleGFtcGxlLmNvbTCCASIwDQYJKoZIhvcN\nAQEBBQADggEPADCCAQoCggEBAIx75lo/fol53qVExXbQaAJ1qdv/NIcPezlEdmWW\n44tb+pE12j+gd2PP7+iPif7eMT6U3n5DjN4q/VyPI8ebwb4LU4Blz2MGLbI/hiA2\n0stFteLR66tP35gODo75s0WYVjTYqE39rTXrErEMUvWl8q0MbqRKGWj4+cEywgVy\nW+5jcDDI5t9CKYYt/IHqMX3r0b7Pwcjp60ozTFxKWSoXQlDz1szw0g+jpJNKmlM7\nEJ0Mm8XElncph8beCTk2exRTxb3fvy9oIWA8Kud4HxM9ZxTKHoV23dROL2uPQoUx\ngSHt2FWM69NA81zkF575YuJV5+mmlpHpVXIAXGhKWLxSQBMCAwEAAaNQME4wHQYD\nVR0OBBYEFIkF/e3zb/wfLRU3X3Va4dFSX1r7MB8GA1UdIwQYMBaAFIkF/e3zb/wf\nLRU3X3Va4dFSX1r7MAwGA1UdEwEB/wQCMAAwDQYJKoZIhvcNAQELBQADggEBAGVt\nT7kQpflQJIwb8QydU04Q0CQJ+O2sTMf2Wmbe/+73mRbkzAhD0oKCkvK3TJ4Xl89O\n5tCCmCIS+rsF3iepS+EIrjA/cZ39Zgo3/B39IMvEyL96GSXCeuHgWys7yNHuDvmh\n9qK0eZ4YEfl9mU57lG8EeP2BVLE2RoAKWbzNanDPJkXLvoUxdphpj6Ne9GxKkl9g\nXqmggEFqlw7G8nScJT/RYK0h0QmaGZ7TLCZ6yNyUki+Ps3S15h4xaxnxQcAp0Udj\npxeZG0vcqJ/5gLszr/llaBw4Rv/ysDev43IRmyY2erpal6MUbk++1Hmo7uifMgK6\nFsDWROtc+z5HPoZZgm8=\n-----END CERTIFICATE-----\n",
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
			Expect(cred.ExpiryDate.String()).To(Equal("2018-01-01 04:07:18 +0000 UTC"))
			Expect(cred.Generated).To(Equal(true))
			Expect(cred.SelfSigned).To(Equal(true))
			Expect(cred.Transitional).To(Equal(true))
			Expect(cred.Ca.Issuer.CommonName).To(Equal("example.com"))
			Expect(cred.Certificate.Issuer.CommonName).To(Equal("example.com"))
			Expect(cred.PrivateKey).To(Equal("-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----"))
			Expect(cred.VersionCreatedAt.String()).To(Equal("2017-01-01 04:07:18 +0000 UTC"))

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
			Expect(cred.VersionCreatedAt.String()).To(Equal("2017-01-01 04:07:18 +0000 UTC"))

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
			Expect(cred.VersionCreatedAt.String()).To(Equal("2017-01-01 04:07:18 +0000 UTC"))

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
			Expect(cred.VersionCreatedAt.String()).To(Equal("2017-01-01 04:07:18 +0000 UTC"))

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
			Expect(cred.VersionCreatedAt.String()).To(Equal("2017-01-01 04:07:18 +0000 UTC"))

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
			Expect(cred.VersionCreatedAt.String()).To(Equal("2017-01-01 04:07:18 +0000 UTC"))

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
			Expect(cred.VersionCreatedAt.String()).To(Equal("2017-01-01 04:07:18 +0000 UTC"))

			jsonOutput, err := json.Marshal(cred)
			Expect(err).NotTo(HaveOccurred())
			Expect(jsonOutput).To(MatchJSON(credJSON))
		})

	})
})
