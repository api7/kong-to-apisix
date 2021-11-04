package e2e

import (
	"net/http"
	"strings"

	"github.com/api7/kong-to-apisix/pkg/kong"

	"github.com/api7/kong-to-apisix/test/e2e/utils"

	"github.com/globocom/gokong"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
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

		// kong deck export mode data test
		err = utils.TestMigrate(utils.TestKongDeckMode)
		gomega.Expect(err).To(gomega.BeNil())

		utils.Compare(&utils.CompareCase{
			Path:        "/get",
			CompareBody: true,
		})

		// kong config export mode data test
		err = utils.TestMigrate(utils.TestKongConfigMode)
		gomega.Expect(err).To(gomega.BeNil())

		utils.Compare(&utils.CompareCase{
			Path:        "/get",
			CompareBody: true,
		})
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

		// kong deck export mode data test
		err = utils.TestMigrate(utils.TestKongDeckMode)
		gomega.Expect(err).To(gomega.BeNil())

		utils.Compare(&utils.CompareCase{
			Path:              "/get/get",
			CompareStatusCode: http.StatusNotFound,
		})

		utils.Compare(&utils.CompareCase{
			Path:        "/get/get",
			Headers:     map[string]string{"Host": "foo.com"},
			CompareBody: true,
		})

		// kong config export mode data test
		err = utils.TestMigrate(utils.TestKongConfigMode)
		gomega.Expect(err).To(gomega.BeNil())

		utils.Compare(&utils.CompareCase{
			Path:              "/get/get",
			CompareStatusCode: http.StatusNotFound,
		})

		utils.Compare(&utils.CompareCase{
			Path:        "/get/get",
			Headers:     map[string]string{"Host": "foo.com"},
			CompareBody: true,
		})
	})

	ginkgo.It("path handling algorithms", func() {
		createdUpstream := utils.DefaultUpstream()
		createdUpstream.Name = "upstream"
		kongUpstream, err := kongCli.Upstreams().Create(createdUpstream)
		gomega.Expect(err).To(gomega.BeNil())

		createdTarget := utils.DefaultTarget()
		createdTarget.Target = strings.TrimPrefix(utils.UpstreamAddr7026, "http://")
		_, err = kongCli.Targets().CreateFromUpstreamId(kongUpstream.Id, createdTarget)
		gomega.Expect(err).To(gomega.BeNil())

		createdService := utils.DefaultService()
		createdService.Host = gokong.String(kongUpstream.Name)
		createdService.Path = gokong.String("/s")
		kongService, err := kongCli.Services().Create(createdService)
		gomega.Expect(err).To(gomega.BeNil())

		type TestExpect struct {
			RequestPath string
			ExpectPath  string
		}

		type testRouteExpect struct {
			Route  kong.Route
			Expect TestExpect
		}

		// https://docs.konghq.com/gateway-oss/2.4.x/admin-api/#path-handling-algorithms
		testRoutes := []testRouteExpect{
			{
				Route: kong.Route{
					Name:         "s-fv0req",
					Paths:        []string{"/fv0"},
					StripPath:    false,
					PathHandling: "v0",
					ServiceID:    *kongService.Id,
				},
				Expect: TestExpect{
					RequestPath: "/fv0req",
					ExpectPath:  "/s/fv0req",
				},
			},
			{
				Route: kong.Route{
					Name:         "sfv1req",
					Paths:        []string{"/fv1"},
					StripPath:    false,
					PathHandling: "v1",
					ServiceID:    *kongService.Id,
				},
				Expect: TestExpect{
					RequestPath: "/fv1req",
					ExpectPath:  "/sfv1req",
				},
			},
			{
				Route: kong.Route{
					Name:         "s-req",
					Paths:        []string{"/tv0"},
					StripPath:    true,
					PathHandling: "v0",
					ServiceID:    *kongService.Id,
				},
				Expect: TestExpect{
					RequestPath: "/tv0req",
					ExpectPath:  "/s/req",
				},
			},
			{
				Route: kong.Route{
					Name:         "sreq",
					Paths:        []string{"/tv1"},
					StripPath:    true,
					PathHandling: "v1",
					ServiceID:    *kongService.Id,
				},
				Expect: TestExpect{
					RequestPath: "/tv1req",
					ExpectPath:  "/sreq",
				},
			},
			{
				Route: kong.Route{
					Name:         "s-fv0-req",
					Paths:        []string{"/fv0/"},
					StripPath:    false,
					PathHandling: "v0",
					ServiceID:    *kongService.Id,
				},
				Expect: TestExpect{
					RequestPath: "/fv0/req",
					ExpectPath:  "/s/fv0/req",
				},
			},
			{
				Route: kong.Route{
					Name:         "sfv1-req",
					Paths:        []string{"/fv1/"},
					StripPath:    false,
					PathHandling: "v1",
					ServiceID:    *kongService.Id,
				},
				Expect: TestExpect{
					RequestPath: "/fv1/req",
					ExpectPath:  "/sfv1/req",
				},
			},
			{
				Route: kong.Route{
					Name:         "s-req",
					Paths:        []string{"/tv0/"},
					StripPath:    true,
					PathHandling: "v0",
					ServiceID:    *kongService.Id,
				},
				Expect: TestExpect{
					RequestPath: "/tv0/req",
					ExpectPath:  "/s/req",
				},
			},
			{
				Route: kong.Route{
					Name:         "sreq",
					Paths:        []string{"/tv1/"},
					StripPath:    true,
					PathHandling: "v1",
					ServiceID:    *kongService.Id,
				},
				Expect: TestExpect{
					RequestPath: "/tv1/req",
					ExpectPath:  "/sreq",
				},
			},
		}

		for _, testRoute := range testRoutes {
			// The gokong project cannot complete the construction of strip_path and path_handling parameters,
			// and realizes the creation of routing functions separately
			ok, err := utils.CreateRoute(testRoute.Route)
			gomega.Expect(ok).To(gomega.BeTrue())
			gomega.Expect(err).To(gomega.BeNil())
			// kong deck export mode data test
			err = utils.TestMigrate(utils.TestKongDeckMode)
			gomega.Expect(err).To(gomega.BeNil())
			req := &utils.CompareCase{Path: testRoute.Expect.RequestPath}
			apisixResp, kongResp := utils.GetBodys(req)
			gomega.Ω(apisixResp).Should(gomega.Equal(kongResp))
			gomega.Ω(apisixResp).Should(gomega.Equal(testRoute.Expect.ExpectPath))
			// kong deck config mode data test
			err = utils.TestMigrate(utils.TestKongConfigMode)
			gomega.Expect(err).To(gomega.BeNil())
			apisixResp, kongResp = utils.GetBodys(req)
			gomega.Ω(apisixResp).Should(gomega.Equal(kongResp))
			gomega.Ω(apisixResp).Should(gomega.Equal(testRoute.Expect.ExpectPath))
			// remove test route
			err = kongCli.Routes().DeleteByName(testRoute.Route.Name)
			gomega.Expect(err).To(gomega.BeNil())
		}
	})
})
