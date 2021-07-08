package e2e

import (
	"fmt"
	"strings"

	"github.com/apache/apisix-ingress-controller/pkg/apisix"
	"github.com/globocom/gokong"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/yiyiyimu/kongtoapisix/pkg/apisixcli"
	"github.com/yiyiyimu/kongtoapisix/pkg/route"
	"github.com/yiyiyimu/kongtoapisix/pkg/upstream"
)

var _ = Describe("route", func() {
	var (
		apisixCli apisix.Cluster
		kongCli   gokong.KongAdminClient
	)

	BeforeEach(func() {
		var err error
		apisixCli, err = apisixcli.CreateAPISIXCli()
		Expect(err).To(BeNil())
		kongCli = gokong.NewClient(gokong.NewDefaultConfig())
	})

	JustBeforeEach(func() {
		err := purgeAll(apisixCli, kongCli)
		Expect(err).To(BeNil())
	})

	It("default service with upstream", func() {
		createdUpstream := defaultUpstream()
		createdUpstream.Name = "upstream"
		kongUpstream, err := kongCli.Upstreams().Create(createdUpstream)
		Expect(err).To(BeNil())

		createdTarget := defaultTarget()
		createdTarget.Target = strings.TrimPrefix(upstreamAddr, "http://")
		_, err = kongCli.Targets().CreateFromUpstreamId(kongUpstream.Id, createdTarget)
		Expect(err).To(BeNil())

		createdTarget.Target = strings.TrimPrefix(upstreamAddr2, "http://")
		_, err = kongCli.Targets().CreateFromUpstreamId(kongUpstream.Id, createdTarget)
		Expect(err).To(BeNil())

		createdService := defaultService()
		createdService.Host = gokong.String(kongUpstream.Name)
		kongService, err := kongCli.Services().Create(createdService)
		Expect(err).To(BeNil())

		createdRoute := defaultRoute()
		createdRoute.Paths = gokong.StringSlice([]string{"/get"})
		createdRoute.Service = gokong.ToId(*kongService.Id)
		_, err = kongCli.Routes().Create(createdRoute)
		Expect(err).To(BeNil())

		err = upstream.MigrateUpstream(apisixCli, kongCli)
		Expect(err).To(BeNil())

		err = route.MigrateRoute(apisixCli, kongCli)
		Expect(err).To(BeNil())

		apisixRespMap := make(map[string]int)
		kongRespMap := make(map[string]int)
		cc := &CompareCase{Path: "/get/get"}
		for range [10]int{} {
			apisixResp, kongResp, err := getResp(cc)
			Expect(err).To(BeNil())
			apisixRespMap[string(apisixResp)]++
			kongRespMap[string(kongResp)]++
		}
		for k, count := range apisixRespMap {
			fmt.Printf("apisix: %s - %d\n", k, count)
			Expect(count > 0).To(BeTrue())
		}
		for k, count := range kongRespMap {
			fmt.Printf("kong: %s - %d\n", k, count)
			Expect(count > 0).To(BeTrue())
		}
	})
})
