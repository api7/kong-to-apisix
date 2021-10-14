package e2e

import (
	"time"

	"github.com/api7/kong-to-apisix/test/e2e/utils"

	"github.com/globocom/gokong"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("plugins", func() {
	var kongCli gokong.KongAdminClient

	ginkgo.JustBeforeEach(func() {
		kongCli = gokong.NewClient(gokong.NewDefaultConfig())
		err := utils.PurgeAll(kongCli)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("default rate limit", func() {
		createdService := utils.DefaultService()
		createdService.Url = gokong.String(utils.UpstreamAddr)
		kongService, err := kongCli.Services().Create(createdService)
		gomega.Expect(err).To(gomega.BeNil())

		createdRoute := utils.DefaultRoute()
		createdRoute.Paths = gokong.StringSlice([]string{"/get"})
		createdRoute.Service = gokong.ToId(*kongService.Id)
		_, err = kongCli.Routes().Create(createdRoute)
		gomega.Expect(err).To(gomega.BeNil())

		createdPlugin := &gokong.PluginRequest{
			Name: "rate-limiting",
			Config: map[string]interface{}{
				"second": 1,
				"policy": "local",
			},
		}
		_, err = kongCli.Plugins().Create(createdPlugin)
		gomega.Expect(err).To(gomega.BeNil())

		err = utils.TestMigrate()
		gomega.Expect(err).To(gomega.BeNil())

		// first time to trigger rate limit
		utils.Compare(&utils.CompareCase{
			Path:              "/get/get",
			CompareStatusCode: 200,
		})

		utils.Compare(&utils.CompareCase{
			Path:              "/get/get",
			CompareStatusCode: 429,
		})

		// wait till rate limit expire
		time.Sleep(time.Second)
		utils.Compare(&utils.CompareCase{
			Path:              "/get/get",
			CompareStatusCode: 200,
		})
	})

	ginkgo.It("default proxy cache", func() {
		createdService := utils.DefaultService()
		createdService.Url = gokong.String(utils.UpstreamAddr)
		kongService, err := kongCli.Services().Create(createdService)
		gomega.Expect(err).To(gomega.BeNil())

		createdRoute := utils.DefaultRoute()
		createdRoute.Paths = gokong.StringSlice([]string{"/get"})
		createdRoute.Service = gokong.ToId(*kongService.Id)
		_, err = kongCli.Routes().Create(createdRoute)
		gomega.Expect(err).To(gomega.BeNil())

		createdPlugin := &gokong.PluginRequest{
			Name: "proxy-cache",
			Config: map[string]interface{}{
				"strategy":  "memory",
				"cache_ttl": 1,
			},
		}
		_, err = kongCli.Plugins().Create(createdPlugin)
		gomega.Expect(err).To(gomega.BeNil())

		err = utils.TestMigrate()
		gomega.Expect(err).To(gomega.BeNil())

		// first time to trigger cache
		utils.Compare(&utils.CompareCase{
			Path: "/get/get",
		})

		apisixResp, kongResp := utils.GetResps(&utils.CompareCase{
			Path: "/get/get",
		})
		apisixCacheStatus := apisixResp.Header.Get("Apisix-Cache-Status")
		kongCacheStatus := kongResp.Header.Get("X-Cache-Status")
		gomega.Ω(apisixCacheStatus).Should(gomega.Equal("HIT"))
		gomega.Ω(kongCacheStatus).Should(gomega.Equal("Hit"))

		/*
			// currently apisix need to reload to make cache_ttl in config.yaml works, so skip this test for now
			time.Sleep(2 * time.Second)
			apisixResp, kongResp = utils.GetResps(&utils.CompareCase{
				Path: "/get/get",
			})
			apisixCacheStatus = apisixResp.Header.Get("Apisix-Cache-Status")
			kongCacheStatus = kongResp.Header.Get("X-Cache-Status")
			gomega.Ω(kongCacheStatus).Should(gomega.Equal("Miss"))
			gomega.Ω(apisixCacheStatus).Should(gomega.Equal("EXPIRED"))
		*/
	})

	ginkgo.It("default key auth", func() {
		createdService := utils.DefaultService()
		createdService.Url = gokong.String(utils.UpstreamAddr)
		kongService, err := kongCli.Services().Create(createdService)
		gomega.Expect(err).To(gomega.BeNil())

		createdRoute := utils.DefaultRoute()
		createdRoute.Paths = gokong.StringSlice([]string{"/get"})
		createdRoute.Service = gokong.ToId(*kongService.Id)
		kongRoute, err := kongCli.Routes().Create(createdRoute)
		gomega.Expect(err).To(gomega.BeNil())

		createdPlugin := &gokong.PluginRequest{
			Name:    "key-auth",
			RouteId: (*gokong.Id)(kongRoute.Id),
		}
		_, err = kongCli.Plugins().Create(createdPlugin)
		gomega.Expect(err).To(gomega.BeNil())

		createdConsumer := utils.DefaultConsumer()
		createdConsumer.Username = "consumer"
		kongConsumer, err := kongCli.Consumers().Create(createdConsumer)
		gomega.Expect(err).To(gomega.BeNil())

		_, err = kongCli.Consumers().CreatePluginConfig(kongConsumer.Id, "key-auth", "{\"key\": \"apikey\"}")
		gomega.Expect(err).To(gomega.BeNil())

		err = utils.TestMigrate()
		gomega.Expect(err).To(gomega.BeNil())

		// without key
		utils.Compare(&utils.CompareCase{
			Path:              "/get/get",
			CompareStatusCode: 401,
		})

		// with key
		utils.Compare(&utils.CompareCase{
			Path:              "/get/get",
			Headers:           map[string]string{"apikey": "apikey"},
			CompareStatusCode: 200,
		})
	})
})
