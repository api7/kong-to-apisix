package e2e

import (
	"github.com/api7/kong-to-apisix/test/e2e/utils"
	"github.com/globocom/gokong"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("consumer", func() {
	var kongCli gokong.KongAdminClient

	ginkgo.JustBeforeEach(func() {
		kongCli = gokong.NewClient(gokong.NewDefaultConfig())
		err := utils.PurgeAll(kongCli)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("test multiple consumers", func() {
		createdService := utils.DefaultService()
		createdService.Url = gokong.String(utils.UpstreamAddr)
		kongService, err := kongCli.Services().Create(createdService)
		gomega.Expect(err).To(gomega.BeNil())

		createdRoute := utils.DefaultRoute()
		createdRoute.Paths = gokong.StringSlice([]string{"/test/consumer"})
		createdRoute.Service = gokong.ToId(*kongService.Id)
		kongRoute, err := kongCli.Routes().Create(createdRoute)
		gomega.Expect(err).To(gomega.BeNil())

		createdPlugin := &gokong.PluginRequest{
			Name:    "key-auth",
			RouteId: (*gokong.Id)(kongRoute.Id),
		}
		_, err = kongCli.Plugins().Create(createdPlugin)
		gomega.Expect(err).To(gomega.BeNil())

		createdConsumer1 := utils.DefaultConsumer()
		createdConsumer1.Username = "consumer1"
		kongConsumer1, err := kongCli.Consumers().Create(createdConsumer1)
		gomega.Expect(err).To(gomega.BeNil())

		_, err = kongCli.Consumers().CreatePluginConfig(kongConsumer1.Id, "key-auth",
			"{\"key\": \"apikey1\"}")
		gomega.Expect(err).To(gomega.BeNil())

		createdConsumer2 := utils.DefaultConsumer()
		createdConsumer2.Username = "consumer2"
		kongConsumer2, err := kongCli.Consumers().Create(createdConsumer2)
		gomega.Expect(err).To(gomega.BeNil())

		_, err = kongCli.Consumers().CreatePluginConfig(kongConsumer2.Id, "key-auth",
			"{\"key\": \"apikey2\"}")
		gomega.Expect(err).To(gomega.BeNil())

		// kong deck export mode data test
		err = utils.TestMigrate(utils.TestKongDeckMode)
		gomega.Expect(err).To(gomega.BeNil())

		// without key
		utils.Compare(&utils.CompareCase{
			Path:              "/test/consumer/apikey",
			CompareStatusCode: 401,
		})

		// with key consumer1
		utils.Compare(&utils.CompareCase{
			Path:              "/test/consumer/apikey1",
			Headers:           map[string]string{"apikey": "apikey1"},
			CompareStatusCode: 200,
		})

		// with key consumer2
		utils.Compare(&utils.CompareCase{
			Path:              "/test/consumer/apikey2",
			Headers:           map[string]string{"apikey": "apikey2"},
			CompareStatusCode: 200,
		})

		// kong config export mode data test
		err = utils.TestMigrate(utils.TestKongConfigMode)
		gomega.Expect(err).To(gomega.BeNil())

		// without key
		utils.Compare(&utils.CompareCase{
			Path:              "/test/consumer/apikey",
			CompareStatusCode: 401,
		})

		// with key consumer1
		utils.Compare(&utils.CompareCase{
			Path:              "/test/consumer/apikey1",
			Headers:           map[string]string{"apikey": "apikey1"},
			CompareStatusCode: 200,
		})

		// with key consumer2
		utils.Compare(&utils.CompareCase{
			Path:              "/test/consumer/apikey2",
			Headers:           map[string]string{"apikey": "apikey2"},
			CompareStatusCode: 200,
		})
	})
})
