package e2e

import (
	"net/http"

	"github.com/globocom/gokong"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("route", func() {
	var kongCli gokong.KongAdminClient

	ginkgo.JustBeforeEach(func() {
		kongCli = gokong.NewClient(gokong.NewDefaultConfig())
		err := purgeAll(kongCli)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("default route and service", func() {
		createdService := defaultService()
		createdService.Url = gokong.String(upstreamAddr)
		kongService, err := kongCli.Services().Create(createdService)
		gomega.Expect(err).To(gomega.BeNil())

		createdRoute := defaultRoute()
		createdRoute.Paths = gokong.StringSlice([]string{"/get"})
		createdRoute.Service = gokong.ToId(*kongService.Id)
		_, err = kongCli.Routes().Create(createdRoute)
		gomega.Expect(err).To(gomega.BeNil())

		err = testMigrate()
		gomega.Expect(err).To(gomega.BeNil())

		apisixResp, kongResp, err := getResp(&CompareCase{
			Path: "/get/get",
		})
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(apisixResp == kongResp).To(gomega.BeTrue())
	})

	ginkgo.It("kong route disable strip_path", func() {
		createdService := defaultService()
		createdService.Url = gokong.String(upstreamAddr)
		kongService, err := kongCli.Services().Create(createdService)
		gomega.Expect(err).To(gomega.BeNil())

		createdRoute := defaultRoute()
		createdRoute.Paths = gokong.StringSlice([]string{"/get"})
		createdRoute.Service = gokong.ToId(*kongService.Id)
		createdRoute.StripPath = gokong.Bool(false)
		_, err = kongCli.Routes().Create(createdRoute)
		gomega.Expect(err).To(gomega.BeNil())

		err = testMigrate()
		gomega.Expect(err).To(gomega.BeNil())

		apisixResp, kongResp, err := getResp(&CompareCase{
			Path: "/get",
		})
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(apisixResp == kongResp).To(gomega.BeTrue())
	})

	ginkgo.It("kong route with host", func() {
		createdService := defaultService()
		createdService.Url = gokong.String(upstreamAddr)
		kongService, err := kongCli.Services().Create(createdService)
		gomega.Expect(err).To(gomega.BeNil())

		createdRoute := defaultRoute()
		createdRoute.Paths = gokong.StringSlice([]string{"/get"})
		createdRoute.Service = gokong.ToId(*kongService.Id)
		createdRoute.Hosts = gokong.StringSlice([]string{"foo.com"})
		_, err = kongCli.Routes().Create(createdRoute)
		gomega.Expect(err).To(gomega.BeNil())

		err = testMigrate()
		gomega.Expect(err).To(gomega.BeNil())

		apisixResp, kongResp, err := getResp(&CompareCase{
			Path:             "/get/get",
			ExpectStatusCode: http.StatusNotFound,
		})
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(apisixResp == kongResp).To(gomega.BeTrue())

		apisixResp, kongResp, err = getResp(&CompareCase{
			Path:    "/get/get",
			Headers: map[string]string{"Host": "foo.com"},
		})
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(apisixResp == kongResp).To(gomega.BeTrue())
	})

})
