package e2e

import (
	"strings"

	"github.com/api7/kong-to-apisix/test/e2e/utils"
	"github.com/globocom/gokong"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("service", func() {
	var kongCli gokong.KongAdminClient

	ginkgo.JustBeforeEach(func() {
		kongCli = gokong.NewClient(gokong.NewDefaultConfig())
		err := utils.PurgeAll(kongCli)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("test only service", func() {
		createdService := utils.DefaultService()
		createdService.Url = gokong.String(utils.UpstreamAddr)
		kongService, err := kongCli.Services().Create(createdService)
		gomega.Expect(err).To(gomega.BeNil())

		createdRoute := utils.DefaultRoute()
		createdRoute.Paths = gokong.StringSlice([]string{"/get"})
		createdRoute.Service = gokong.ToId(*kongService.Id)
		_, err = kongCli.Routes().Create(createdRoute)
		gomega.Expect(err).To(gomega.BeNil())

		// kong deck export mode data test
		err = utils.TestMigrate(utils.TestKongDeckMode)
		gomega.Expect(err).To(gomega.BeNil())

		utils.Compare(&utils.CompareCase{
			Path:        "/get/get",
			CompareBody: true,
		})

		// kong config export mode data test
		err = utils.TestMigrate(utils.TestKongConfigMode)
		gomega.Expect(err).To(gomega.BeNil())

		utils.Compare(&utils.CompareCase{
			Path:        "/get/get",
			CompareBody: true,
		})
	})

	ginkgo.It("test service and upstream", func() {
		kongUpstreamName := "upstream"

		createdService := utils.DefaultService()
		createdService.Host = gokong.String(kongUpstreamName)
		kongService, err := kongCli.Services().Create(createdService)
		gomega.Expect(err).To(gomega.BeNil())

		createdUpstream := utils.DefaultUpstream()
		createdUpstream.Name = kongUpstreamName
		kongUpstream, err := kongCli.Upstreams().Create(createdUpstream)
		gomega.Expect(err).To(gomega.BeNil())

		createdTarget := utils.DefaultTarget()
		createdTarget.Target = strings.TrimPrefix(utils.UpstreamAddr, "http://")
		_, err = kongCli.Targets().CreateFromUpstreamId(kongUpstream.Id, createdTarget)
		gomega.Expect(err).To(gomega.BeNil())

		createdTarget.Target = strings.TrimPrefix(utils.UpstreamAddr2, "http://")
		_, err = kongCli.Targets().CreateFromUpstreamId(kongUpstream.Id, createdTarget)
		gomega.Expect(err).To(gomega.BeNil())

		createdRoute := utils.DefaultRoute()
		createdRoute.Paths = gokong.StringSlice([]string{"/hello"})
		createdRoute.Service = gokong.ToId(*kongService.Id)
		_, err = kongCli.Routes().Create(createdRoute)
		gomega.Expect(err).To(gomega.BeNil())

		// kong deck export mode data test
		err = utils.TestMigrate(utils.TestKongDeckMode)
		gomega.Expect(err).To(gomega.BeNil())

		req := &utils.CompareCase{Path: "/hello/world"}

		a6RespMaps := make(map[string]int)
		kongRespMaps := make(map[string]int)
		for i := 0; i < 10; i++ {
			apisixResp, kongResp := utils.GetBodys(req)
			a6RespMaps[apisixResp]++
			kongRespMaps[kongResp]++
		}

		for key, count := range a6RespMaps {
			gomega.Expect(kongRespMaps[key]).To(gomega.Equal(count))
		}

		// kong config export mode data test
		err = utils.TestMigrate(utils.TestKongConfigMode)
		gomega.Expect(err).To(gomega.BeNil())

		a6RespMaps = make(map[string]int)
		kongRespMaps = make(map[string]int)
		for i := 0; i < 10; i++ {
			apisixResp, kongResp := utils.GetBodys(req)
			a6RespMaps[apisixResp]++
			kongRespMaps[kongResp]++
		}

		for key, count := range a6RespMaps {
			gomega.Expect(kongRespMaps[key]).To(gomega.Equal(count))
		}
	})

	ginkgo.It("test service and plugin", func() {
		createdService := utils.DefaultService()
		createdService.Url = gokong.String(utils.UpstreamAddr)
		kongService, err := kongCli.Services().Create(createdService)
		gomega.Expect(err).To(gomega.BeNil())

		createdPlugin := &gokong.PluginRequest{
			Name:      "key-auth",
			ServiceId: (*gokong.Id)(kongService.Id),
		}
		_, err = kongCli.Plugins().Create(createdPlugin)
		gomega.Expect(err).To(gomega.BeNil())

		createdRoute := utils.DefaultRoute()
		createdRoute.Paths = gokong.StringSlice([]string{"/get"})
		createdRoute.Service = gokong.ToId(*kongService.Id)
		_, err = kongCli.Routes().Create(createdRoute)
		gomega.Expect(err).To(gomega.BeNil())

		createdConsumer := utils.DefaultConsumer()
		createdConsumer.Username = "consumer"
		kongConsumer, err := kongCli.Consumers().Create(createdConsumer)
		gomega.Expect(err).To(gomega.BeNil())

		_, err = kongCli.Consumers().CreatePluginConfig(kongConsumer.Id, "key-auth", "{\"key\": \"apikey\"}")
		gomega.Expect(err).To(gomega.BeNil())

		// kong deck export mode data test
		err = utils.TestMigrate(utils.TestKongDeckMode)
		gomega.Expect(err).To(gomega.BeNil())

		utils.Compare(&utils.CompareCase{
			Path:              "/get/get",
			CompareStatusCode: 401,
		})

		utils.Compare(&utils.CompareCase{
			Path:              "/get/get",
			Headers:           map[string]string{"apikey": "apikey"},
			CompareStatusCode: 200,
		})

		// kong config export mode data test
		err = utils.TestMigrate(utils.TestKongConfigMode)
		gomega.Expect(err).To(gomega.BeNil())

		utils.Compare(&utils.CompareCase{
			Path:              "/get/get",
			CompareStatusCode: 401,
		})

		utils.Compare(&utils.CompareCase{
			Path:              "/get/get",
			Headers:           map[string]string{"apikey": "apikey"},
			CompareStatusCode: 200,
		})
	})
})
