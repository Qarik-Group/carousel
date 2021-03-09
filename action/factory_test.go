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
		})

		Context("given a up-to-date credential", func() {
			It("finds the next action", func() {
				credential.PathVersion()
				actions := factory.NextAction(credential)
				Expect(actions).To(ContainElements(HaveName("up-to-date")))
			})

		})

		Context("given a latest credential which has not been deployed yet", func() {
			JustBeforeEach(func() {
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
				JustBeforeEach(func() {
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
	})
})

func HaveName(n string) types.GomegaMatcher {
	return WithTransform(func(a Action) string {
		return a.Name()
	}, ContainSubstring(n))
}
