package e2e

import (
	"net/http"

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

	It("default route and service", func() {
		createdService := defaultService()
		createdService.Url = gokong.String(upstreamAddr)
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

		same, err := compareURL(&CompareCase{
			Path: "/get/get",
		})
		Expect(err).To(BeNil())
		Expect(same).To(BeTrue())
	})

	It("kong route disable strip_path", func() {
		createdService := defaultService()
		createdService.Url = gokong.String(upstreamAddr)
		kongService, err := kongCli.Services().Create(createdService)
		Expect(err).To(BeNil())

		createdRoute := defaultRoute()
		createdRoute.Paths = gokong.StringSlice([]string{"/get"})
		createdRoute.Service = gokong.ToId(*kongService.Id)
		createdRoute.StripPath = gokong.Bool(false)
		_, err = kongCli.Routes().Create(createdRoute)
		Expect(err).To(BeNil())

		err = upstream.MigrateUpstream(apisixCli, kongCli)
		Expect(err).To(BeNil())

		err = route.MigrateRoute(apisixCli, kongCli)
		Expect(err).To(BeNil())

		same, err := compareURL(&CompareCase{
			Path: "/get",
		})
		Expect(err).To(BeNil())
		Expect(same).To(BeTrue())
	})

	It("kong route with host", func() {
		createdService := defaultService()
		createdService.Url = gokong.String(upstreamAddr)
		kongService, err := kongCli.Services().Create(createdService)
		Expect(err).To(BeNil())

		createdRoute := defaultRoute()
		createdRoute.Paths = gokong.StringSlice([]string{"/get"})
		createdRoute.Service = gokong.ToId(*kongService.Id)
		createdRoute.Hosts = gokong.StringSlice([]string{"foo.com"})
		_, err = kongCli.Routes().Create(createdRoute)
		Expect(err).To(BeNil())

		err = upstream.MigrateUpstream(apisixCli, kongCli)
		Expect(err).To(BeNil())

		err = route.MigrateRoute(apisixCli, kongCli)
		Expect(err).To(BeNil())

		same, err := compareURL(&CompareCase{
			Path:             "/get/get",
			ExpectStatusCode: http.StatusNotFound,
		})
		Expect(err).To(BeNil())
		Expect(same).To(BeTrue())

		same, err = compareURL(&CompareCase{
			Path:    "/get/get",
			Headers: map[string]string{"Host": "foo.com"},
		})
		Expect(err).To(BeNil())
		Expect(same).To(BeTrue())
	})

})
