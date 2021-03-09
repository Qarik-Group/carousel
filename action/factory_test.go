package action_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"

	. "github.com/starkandwayne/carousel/action"
	"github.com/starkandwayne/carousel/credhub"
	"github.com/starkandwayne/carousel/state"
)

// var _ = Describe("ActionFactory", func() {

// })

var _ = Describe("ConcreteActionFactory", func() {
	var (
		factory       ActionFactory
		olderThan     time.Time
		expiresBefore time.Time
		credential    *state.Credential
	)

	BeforeEach(func() {
		olderThan = time.Now()
		expiresBefore = time.Now()
		factory = &ConcreteActionFactory{
			OlderThan:     olderThan,
			ExpiresBefore: expiresBefore,
		}
	})

	Describe("NextAction", func() {
		BeforeEach(func() {
			vca := olderThan.Add(time.Hour)
			fooDeployment := state.Deployment{Name: "foo-deployment"}
			credential = &state.Credential{
				Latest:      true,
				Deployments: state.Deployments{&fooDeployment},
				Credential: &credhub.Credential{
					VersionCreatedAt: &vca,
					ID:               "foo-id",
					Name:             "/foo-name",
				},
				Path: &state.Path{
					Deployments: state.Deployments{&fooDeployment},
				},
			}
			credential.Path.Versions = state.Credentials{credential}
		})

		Context("given a up-to-date credential", func() {
			It("finds the next action", func() {
				credential.PathVersion()
				actions := factory.NextAction(credential)
				Expect(actions).To(ContainElements(HaveName("up-to-date")))
			})

		})

		Context("given a latest credential which has not been deployed yet", func() {
			BeforeEach(func() {
				credential.Latest = true
				credential.Path.Deployments = append(
					credential.Path.Deployments, &state.Deployment{Name: "bar-deployment"})
			})

			It("finds the next action", func() {
				credential.PathVersion()
				actions := factory.NextAction(credential)
				Expect(actions).To(ContainElements(
					HaveName("deploy(bar-deployment)"),
				))
				Expect(len(actions)).To(Equal(1))
			})

			Context("which is to old", func() {
				BeforeEach(func() {
					vca := olderThan.Add(-10 * time.Minute)
					credential.VersionCreatedAt = &vca
				})

				It("finds the next action", func() {
					credential.PathVersion()
					actions := factory.NextAction(credential)
					Expect(actions).To(ContainElements(
						HaveName("regenerate"),
					))
					Expect(len(actions)).To(Equal(1))
				})
			})
		})

		Context("given a credential which has not been deployed", func() {
			BeforeEach(func() {
				credential.Deployments = make(state.Deployments, 0)
				credential.Path.Deployments = make(state.Deployments, 0)
			})

			It("finds the next action", func() {
				credential.PathVersion()
				actions := factory.NextAction(credential)
				Expect(actions).To(ContainElements(
					HaveName("clean-up"),
				))
				Expect(len(actions)).To(Equal(1))
			})

			Context("which is still being referenced by a deployed credential", func() {
				BeforeEach(func() {
					credential.ReferencedBy = state.Credentials{&state.Credential{
						Deployments: make(state.Deployments, 1),
					}}
				})

				It("finds the next action", func() {
					credential.PathVersion()
					actions := factory.NextAction(credential)
					Expect(actions).To(ContainElements(
						HaveName("up-to-date"),
					))
					Expect(len(actions)).To(Equal(1))
				})
			})
		})

		Context("given a credential which is expiring", func() {
			BeforeEach(func() {
				ed := expiresBefore.Add(-time.Hour)
				credential.ExpiryDate = &ed
				edca := expiresBefore.Add(+time.Hour)
				credential.SignedBy = &state.Credential{
					Credential: &credhub.Credential{
						ExpiryDate: &edca,
					},
				}
			})

			It("finds the next action", func() {
				credential.PathVersion()
				actions := factory.NextAction(credential)
				Expect(actions).To(ContainElements(
					HaveName("regenerate"),
				))
				Expect(len(actions)).To(Equal(1))
			})

			Context("which is signed by an expiring ca", func() {
				BeforeEach(func() {
					ed := expiresBefore.Add(-time.Hour)
					credential.SignedBy = &state.Credential{
						Credential: &credhub.Credential{
							ExpiryDate: &ed,
						},
					}
				})

				It("finds the next action", func() {
					credential.PathVersion()
					actions := factory.NextAction(credential)
					Expect(actions).To(ContainElements(
						HaveName("up-to-date"),
					))
					Expect(len(actions)).To(Equal(1))
				})
			})
		})

		Context("given a signing credential with an latest active transitional sibling", func() {
			BeforeEach(func() {
				signing := true
				vca := olderThan.Add(time.Hour)
				credential.Signing = &signing
				credential.Latest = false
				credential.Path.Versions = append(credential.Path.Versions,
					&state.Credential{
						Deployments: credential.Path.Deployments,
						Latest:      true,
						Credential: &credhub.Credential{
							VersionCreatedAt: &vca,
							Transitional:     true,
						},
						Path: credential.Path,
					})
				credential.Path.Versions.SortByCreatedAt()
			})

			It("finds the next action", func() {
				credential.PathVersion()
				actions := factory.NextAction(credential)
				Expect(actions).To(ContainElements(
					HaveName("mark-transitional"),
				))
				Expect(len(actions)).To(Equal(1))
			})
		})
	})
})

func HaveName(n string) types.GomegaMatcher {
	return WithTransform(func(a Action) string {
		return a.Name()
	}, ContainSubstring(n))
}
