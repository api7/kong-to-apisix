package kong

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/api7/kong-to-apisix/pkg/apisix"

	uuid "github.com/satori/go.uuid"
)

func MigrateRoute(kongConfig *Config, apisixConfig *apisix.Config) error {
	kongServices := kongConfig.Services
	kongRoutes := kongConfig.Routes

	// Kong deck export
	for ksIndex, kongService := range kongServices {
		if len(kongService.Routes) <= 0 {
			continue
		}

		for krIndex, kongRoute := range kongService.Routes {
			kongRouteId := kongRoute.ID
			if len(kongRouteId) <= 0 {
				kongRouteId = uuid.NewV4().String()
				kongConfig.Services[ksIndex].Routes[krIndex].ID = kongRouteId
				kongRoute.ID = kongRouteId
			}
			a6Routes := ConvertKongRouteToA6(kongService, kongRoute)
			if len(a6Routes) <= 0 {
				fmt.Printf("Kong route [ %s ] to APISIX conversion failure\n", kongRouteId)
				continue
			}
			apisixConfig.Routes = append(apisixConfig.Routes, a6Routes...)
			fmt.Printf("Kong service [ %s ] to APISIX conversion completed\n", kongRouteId)
		}
	}

	// Kong config export
	for krIndex, kongRoute := range kongRoutes {
		kongRouteId := kongRoute.ID
		// reset route id
		if len(kongRouteId) <= 0 {
			kongRouteId = uuid.NewV4().String()
			kongConfig.Routes[krIndex].ID = kongRouteId
		}

		if len(kongRoute.ServiceID) <= 0 {
			fmt.Printf("Kong route [ %s ] not setting service\n", kongRouteId)
			continue
		}

		kongService, err := FindKongServiceById(&kongConfig.Services, kongRoute.ServiceID)
		if err != nil {
			fmt.Printf("Kong route [ %s ] mapping service not found\n", kongRouteId)
			continue
		}
		a6Routes := ConvertKongRouteToA6(*kongService, kongRoute)
		if len(a6Routes) <= 0 {
			fmt.Printf("Kong route [ %s ] to APISIX conversion failure\n", kongRouteId)
			continue
		}
		apisixConfig.Routes = append(apisixConfig.Routes, a6Routes...)
		fmt.Printf("Kong service [ %s ] to APISIX conversion completed\n", kongRouteId)
	}

	fmt.Println("Kong to APISIX routes configuration conversion completed")
	return nil
}

func ConvertKongRouteToA6(kongService Service, kongRoute Route) apisix.Routes {
	var a6Routes apisix.Routes
	isPathGroup := len(kongRoute.Paths) > 1
	for krpIndex, kongRoutePath := range kongRoute.Paths {
		var a6Route apisix.Route
		if isPathGroup {
			a6Route.ID = kongRoute.ID + "-" + strconv.Itoa(krpIndex+1)
			a6Route.Name = kongRoute.Name + "-" + strconv.Itoa(krpIndex+1)
		} else {
			a6Route.ID = kongRoute.ID
			a6Route.Name = kongRoute.Name
		}
		a6Route.URI = kongRoutePath + "*"
		a6Route.Hosts = kongRoute.Hosts
		a6Route.Methods = kongRoute.Methods
		a6Route.Priority = kongRoute.RegexPriority
		a6Route.ServiceID = kongService.ID
		proxyRewrite := GenerateProxyRewritePluginConfig(kongService.Path, kongRoutePath,
			kongRoute.StripPath, kongRoute.PathHandling)
		// mapping kong to apisix upstream request URI
		a6Route.Plugins.ProxyRewrite = proxyRewrite

		for _, kongRoutePlugin := range kongRoute.Plugins {
			if !kongRoutePlugin.Enabled {
				continue
			}
			switch kongRoutePlugin.Name {
			case PluginKeyAuth:
				a6Route.Plugins.KeyAuth = KTAConversionKongPluginKeyAuth(kongRoutePlugin)
			case PluginProxyCache:
				a6Route.Plugins.ProxyCache = KTAConversionKongPluginProxyCache(kongRoutePlugin)
			case PluginRateLimiting:
				a6Route.Plugins.LimitCount = KTAConversionKongPluginRateLimiting(kongRoutePlugin)
			default:
				fmt.Printf("Kong route [%s] plugin %s not supported by apisix yet\n", a6Route.ID,
					kongRoutePlugin.Name)
			}
		}

		a6Routes = append(a6Routes, a6Route)
	}
	return a6Routes
}

// GenerateProxyRewritePluginConfig Generate routing and forwarding rules
// https://docs.konghq.com/gateway-oss/2.4.x/admin-api/#path-handling-algorithms
func GenerateProxyRewritePluginConfig(servicePath string, routerPath string, stripPath bool,
	pathHandling string) *apisix.ProxyRewrite {
	if len(servicePath) == 0 {
		servicePath = "/"
	}

	var pathRegex string
	if stripPath {
		pathRegex = fmt.Sprintf(`^%s/?(.*)`, routerPath)
	} else {
		pathRegex = `^/?(.*)`
	}
	pathRegex = strings.Replace(pathRegex, "//", "/", -1)

	var pathPattern string
	if pathHandling == "v1" {
		pathPattern = fmt.Sprintf(`%s$1`, servicePath)
	} else { // pathHandling == "v0"
		pathPattern = fmt.Sprintf(`%s/$1`, servicePath)
	}
	pathPattern = strings.Replace(pathPattern, "//", "/", -1)

	var proxyRewrite apisix.ProxyRewrite
	proxyRewrite.RegexURI = []string{pathRegex, pathPattern}
	return &proxyRewrite
}
