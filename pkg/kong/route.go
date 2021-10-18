package kong

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	uuid "github.com/satori/go.uuid"

	apisix "github.com/api7/kong-to-apisix/pkg/apisix"
)

func MigrateRoute(kongConfig *Config, apisixConfig *apisix.Config) error {
	kongServices := kongConfig.Services
	kongRoutes := kongConfig.Routes

	var err error
	// compatible kong deck and kong config mode
	for _, service := range kongServices {
		var kongServiceRoutes *Routes
		if len(service.Routes) > 0 {
			kongServiceRoutes = &service.Routes
		} else {
			kongServiceRoutes, err = GetKongRoutesByServiceID(&kongRoutes, service.ID)
			if kongServiceRoutes == nil || err != nil {
				continue
			}
		}

		KTAConversionKongRouter(apisixConfig, service, *kongServiceRoutes)
	}

	fmt.Println("Kong to APISIX routes configuration conversion completed")
	return nil
}

func KTAConversionKongRouter(apisixConfig *apisix.Config, kongService Service, kongRoutes Routes) {
	var apisixRoute apisix.Route
	var kongRoute Route
	// Kong and apisix plugin structure and routing rules are different, so split routing
	for routeIndex := range kongRoutes {
		kongRoute = kongRoutes[routeIndex]
		isPathGroup := len(kongRoute.Paths) > 1
		for pathIndex, path := range kongRoute.Paths {
			if len(kongRoute.ID) <= 0 {
				apisixRoute.ID = uuid.NewV4().String()
				apisixRoute.Name = kongRoute.Name
			} else {
				if isPathGroup {
					apisixRoute.ID = kongRoute.ID + "-" + strconv.Itoa(pathIndex+1)
					apisixRoute.Name = kongRoute.Name + "-" + strconv.Itoa(pathIndex+1)
				} else {
					apisixRoute.ID = kongRoute.ID
					apisixRoute.Name = kongRoute.Name
				}
			}
			apisixRoute.URI = path + "*"
			apisixRoute.Hosts = kongRoute.Hosts
			apisixRoute.Methods = kongRoute.Methods
			apisixRoute.Priority = kongRoute.RegexPriority
			apisixRoute.ServiceID = kongService.ID
			proxyRewrite := GenerateProxyRewritePluginConfig(kongService.Path, path,
				kongRoute.StripPath, kongRoute.PathHandling)
			// mapping kong to apisix upstream request URI
			apisixRoute.Plugins.ProxyRewrite = proxyRewrite
			apisixConfig.Routes = append(apisixConfig.Routes, apisixRoute)
		}
		fmt.Fprintf(os.Stdout, "Kong route [ %s ] to APISIX conversion completed\n", kongRoute.ID)
	}
}

func GetKongRoutesByServiceID(kongRouters *Routes, kongServiceID string) (*Routes, error) {
	if kongRouters == nil {
		return nil, errors.New("kong routers is nil or invalid")
	}

	if len(kongServiceID) <= 0 {
		return nil, errors.New("kong service id invalid")
	}

	var kongServiceRoutes Routes
	for _, router := range *kongRouters {
		if router.Service == kongServiceID {
			kongServiceRoutes = append(kongServiceRoutes, router)
		}
	}

	return &kongServiceRoutes, nil
}

// GenerateProxyRewritePluginConfig Generate routing and forwarding rules
// https://docs.konghq.com/gateway-oss/2.4.x/admin-api/#path-handling-algorithms
func GenerateProxyRewritePluginConfig(servicePath string, routerPath string, stripPath bool,
	pathHandling string) apisix.ProxyRewrite {
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
	return proxyRewrite
}
