package e2e

import (
	"strings"

	"github.com/globocom/gokong"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("upstream", func() {
	var kongCli gokong.KongAdminClient

	ginkgo.JustBeforeEach(func() {
		kongCli = gokong.NewClient(gokong.NewDefaultConfig())
		err := purgeAll(kongCli)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("default route, service and upstream", func() {
		createdUpstream := defaultUpstream()
		createdUpstream.Name = "upstream"
		kongUpstream, err := kongCli.Upstreams().Create(createdUpstream)
		gomega.Expect(err).To(gomega.BeNil())

		createdTarget := defaultTarget()
		createdTarget.Target = strings.TrimPrefix(upstreamAddr, "http://")
		_, err = kongCli.Targets().CreateFromUpstreamId(kongUpstream.Id, createdTarget)
		gomega.Expect(err).To(gomega.BeNil())

		createdTarget.Target = strings.TrimPrefix(upstreamAddr2, "http://")
		_, err = kongCli.Targets().CreateFromUpstreamId(kongUpstream.Id, createdTarget)
		gomega.Expect(err).To(gomega.BeNil())

		createdService := defaultService()
		createdService.Host = gokong.String(kongUpstream.Name)
		kongService, err := kongCli.Services().Create(createdService)
		gomega.Expect(err).To(gomega.BeNil())

		createdRoute := defaultRoute()
		createdRoute.Paths = gokong.StringSlice([]string{"/get"})
		createdRoute.Service = gokong.ToId(*kongService.Id)
		_, err = kongCli.Routes().Create(createdRoute)
		gomega.Expect(err).To(gomega.BeNil())

		err = testMigrate()
		gomega.Expect(err).To(gomega.BeNil())

		apisixRespMap := make(map[string]int)
		kongRespMap := make(map[string]int)
		cc := &CompareCase{Path: "/get/get"}
		for range [10]int{} {
			apisixResp, kongResp, err := getResp(cc)
			gomega.Expect(err).To(gomega.BeNil())
			apisixRespMap[string(apisixResp)]++
			kongRespMap[string(kongResp)]++
		}
		for _, count := range apisixRespMap {
			gomega.Expect(count > 0).To(gomega.BeTrue())
		}
		for _, count := range kongRespMap {
			gomega.Expect(count > 0).To(gomega.BeTrue())
		}
	})

})
