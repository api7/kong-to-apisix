package e2e

import (
	"net/http"

	"github.com/globocom/gokong"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"

	"github.com/api7/kong-to-apisix/test/e2e/utils"
)

var _ = ginkgo.Describe("route", func() {
	var kongCli gokong.KongAdminClient

	ginkgo.JustBeforeEach(func() {
		kongCli = gokong.NewClient(gokong.NewDefaultConfig())
		err := utils.PurgeAll(kongCli)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("default route and service", func() {
		createdService := utils.DefaultService()
		createdService.Url = gokong.String(utils.UpstreamAddr)
		kongService, err := kongCli.Services().Create(createdService)
		gomega.Expect(err).To(gomega.BeNil())

		createdRoute := utils.DefaultRoute()
		createdRoute.Paths = gokong.StringSlice([]string{"/get"})
		createdRoute.Service = gokong.ToId(*kongService.Id)
		_, err = kongCli.Routes().Create(createdRoute)
		gomega.Expect(err).To(gomega.BeNil())

		err = utils.TestMigrate()
		gomega.Expect(err).To(gomega.BeNil())

		apisixResp, kongResp, err := utils.GetResp(&utils.CompareCase{
			Path: "/get/get",
		})
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(apisixResp == kongResp).To(gomega.BeTrue())
	})

	ginkgo.It("kong route disable strip_path", func() {
		createdService := utils.DefaultService()
		createdService.Url = gokong.String(utils.UpstreamAddr)
		kongService, err := kongCli.Services().Create(createdService)
		gomega.Expect(err).To(gomega.BeNil())

		createdRoute := utils.DefaultRoute()
		createdRoute.Paths = gokong.StringSlice([]string{"/get"})
		createdRoute.Service = gokong.ToId(*kongService.Id)
		createdRoute.StripPath = gokong.Bool(false)
		_, err = kongCli.Routes().Create(createdRoute)
		gomega.Expect(err).To(gomega.BeNil())

		err = utils.TestMigrate()
		gomega.Expect(err).To(gomega.BeNil())

		apisixResp, kongResp, err := utils.GetResp(&utils.CompareCase{
			Path: "/get",
		})
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(apisixResp == kongResp).To(gomega.BeTrue())
	})

	ginkgo.It("kong route with host", func() {
		createdService := utils.DefaultService()
		createdService.Url = gokong.String(utils.UpstreamAddr)
		kongService, err := kongCli.Services().Create(createdService)
		gomega.Expect(err).To(gomega.BeNil())

		createdRoute := utils.DefaultRoute()
		createdRoute.Paths = gokong.StringSlice([]string{"/get"})
		createdRoute.Service = gokong.ToId(*kongService.Id)
		createdRoute.Hosts = gokong.StringSlice([]string{"foo.com"})
		_, err = kongCli.Routes().Create(createdRoute)
		gomega.Expect(err).To(gomega.BeNil())

		err = utils.TestMigrate()
		gomega.Expect(err).To(gomega.BeNil())

		apisixResp, kongResp, err := utils.GetResp(&utils.CompareCase{
			Path:             "/get/get",
			ExpectStatusCode: http.StatusNotFound,
		})
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(apisixResp == kongResp).To(gomega.BeTrue())

		apisixResp, kongResp, err = utils.GetResp(&utils.CompareCase{
			Path:    "/get/get",
			Headers: map[string]string{"Host": "foo.com"},
		})
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(apisixResp == kongResp).To(gomega.BeTrue())
	})

})
