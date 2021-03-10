package credhub_test

import (
	"fmt"
	"log"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	. "github.com/onsi/gomega/gstruct"

	. "github.com/starkandwayne/carousel/credhub"

	chcli "code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/auth"
)

var _ = Describe("Credhub", func() {
	It("Implements the CredHub interface", func() {
		logger := log.New(GinkgoWriter, "", 0)
		var ch CredHub
		ch = NewCredHub(nil)
		logger.Println(ch) // use client so it compiles
	})

	var (
		server     *ghttp.Server
		credhub    CredHub
		apiAddress string
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		apiAddress = server.URL()
		header := http.Header{}
		header.Add("Content-Type", "application/json")
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/info"),
				ghttp.RespondWith(http.StatusOK, fmt.Sprintf(
					`{"auth-server":{"url":"%s/uaa"}}`, apiAddress), header),
			),
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/uaa/oauth/token"),
				ghttp.RespondWith(http.StatusOK, `
							{"access_token":"token","token_type":"bearer","expires_in":"3600"}
				`, header),
			),
		)

		ch, err := chcli.New(
			apiAddress,
			chcli.SkipTLSValidation(true),
			chcli.Auth(auth.UaaClientCredentials("foo-client", "bar-secert")),
		)
		Expect(err).ToNot(HaveOccurred())
		credhub = NewCredHub(ch)
	})

	Describe("FindAll", func() {
		JustBeforeEach(func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v1/data", "path="),
					ghttp.RespondWith(http.StatusOK, `{
	"credentials" : [ {
		"version_created_at" : "2019-02-01T20:37:52Z",
		"name" : "/some-unique-name"
	}, {
		"version_created_at" : "2019-02-01T20:37:52Z",
		"name" : "/some-name"
	} ]
}`,
					),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v1/certificates/"),
					ghttp.RespondWith(http.StatusOK, `{
	"certificates" : [ {
		"name" : "/some-name",
		"versions" : [ {
			"id" : "b386c4cc-abfb-4150-95e5-449d7655e62d",
			"expiry_date" : "2020-02-01T20:37:52Z",
			"transitional" : true,
			"certificate_authority" : false,
			"self_signed" : false,
			"generated" : false
		}, {
			"id" : "86bfcd3a-aceb-4ec6-bf67-efd5932a9bf2",
			"expiry_date" : "2019-02-01T20:37:52Z",
			"transitional" : false,
			"certificate_authority" : false,
			"self_signed" : false,
			"generated" : false
		} ],
		"signed_by" : "/testCa",
		"signs" : [ "/cert1", "/cert2" ],
		"id" : "f6f1da12-03a3-4db9-93c5-26aa1346785b"
	} ]
}`,
					),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", ContainSubstring("/api/v1/data")),
					ghttp.RespondWith(http.StatusOK, `{
	"data" : [ {
		"type" : "value",
		"version_created_at" : "2019-02-01T20:37:52Z",
		"id" : "eaebb03f-21a9-41f7-beb0-af4c60aa38d6",
		"name" : "/some-name",
		"metadata" : {
			"description" : "example metadata"
		},
		"value" : "some-value"
	} ]
}`,
					),
				),
			)
		})

		Context("when all requests return valid data", func() {
			JustBeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", ContainSubstring("/api/v1/data")),
						ghttp.RespondWith(http.StatusOK, `{
	"data" : [ {
		"type" : "ssh",
		"version_created_at" : "2019-02-01T20:37:52Z",
		"id" : "355d920a-5f2b-4e99-81f2-47562d7db5d4",
		"name" : "/some-ssh",
		"value" : {
			 "public_key": "-----BEGIN RSA PUBLIC KEY-----\n...\n-----END RSA PUBLIC KEY-----",
			 "private_key": "-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----",
			 "public_key_fingerprint": "some fingerprint"
		}},{
		"type" : "ssh",
		"version_created_at" : "2020-02-01T20:37:52Z",
		"id" : "6f7b19fc-3098-485a-b724-0c7d8788f9a5",
		"name" : "/some-ssh",
		"value" : {
			 "public_key": "-----BEGIN RSA PUBLIC KEY-----\n...\n-----END RSA PUBLIC KEY-----",
			 "private_key": "-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----",
			 "public_key_fingerprint": "some fingerprint"
	 }}]
}`,
						),
					),
				)
			})

			It("finds all credentials", func() {
				creds, err := credhub.FindAll()
				Expect(err).ToNot(HaveOccurred())
				Expect(len(creds)).To(Equal(3))

				id := func(element interface{}) string {
					return element.(*Credential).ID
				}
				Expect(creds).To(MatchElements(id, IgnoreExtras, Elements{
					"355d920a-5f2b-4e99-81f2-47562d7db5d4": Not(BeZero()),
					"eaebb03f-21a9-41f7-beb0-af4c60aa38d6": Not(BeZero()),
					"6f7b19fc-3098-485a-b724-0c7d8788f9a5": Not(BeZero()),
				}))

			})
		})

		Context("when a request returns invalid data", func() {
			JustBeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", ContainSubstring("/api/v1/data")),
						ghttp.RespondWith(http.StatusOK, `{`),
					),
				)
			})

			It("returns an error", func() {
				_, err := credhub.FindAll()
				Expect(err).To(MatchError("unexpected EOF"))
			})
		})
	})
})
