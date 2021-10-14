package kong

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/api7/kong-to-apisix/pkg/apisix"
)

func MigrateRoute(kongConfig *Config, apisixConfig *apisix.Config) error {
	kongServices := kongConfig.Services
	kongRoutes := kongConfig.Routes
	apisixRoutes := apisixConfig.Routes

	for _, kongRoute := range kongRoutes {
		var err error
		var kongService *Service
		if len(kongRoute.Service) > 0 {
			kongService, err = GetKongServiceByID(&kongServices, kongRoute.Service)
			if err != nil {
				fmt.Printf("Migrate Route ID: %s, err: %s", kongRoute.ID, err.Error())
				continue
			}
		}

		var apisixRoute apisix.Route
		// Kong and apisix plugin structure and routing rules are different, so split routing
		for index, path := range kongRoute.Paths {
			apisixRoute.ID = kongRoute.ID + "-" + strconv.Itoa(index)
			apisixRoute.Name = kongRoute.Name + "-" + strconv.Itoa(index)
			apisixRoute.URI = path + "*"
			apisixRoute.Hosts = kongRoute.Hosts
			apisixRoute.Methods = kongRoute.Methods
			apisixRoute.Priority = kongRoute.RegexPriority
			apisixRoute.ServiceID = kongRoute.Service
			proxyRewrite := GenerateProxyRewritePluginConfig(kongService.Path, path,
				kongRoute.StripPath, kongRoute.PathHandling)
			// mapping kong to apisix upstream request URI
			apisixRoute.Plugins = make(apisix.Plugins)
			apisixRoute.Plugins["proxy-rewrite"] = proxyRewrite
			apisixRoutes = append(apisixRoutes, apisixRoute)
			fmt.Fprintf(os.Stdout, "Kong route [ %s ] to APISIX conversion completed\n", kongRoute.ID)
		}
	}
	apisixConfig.Routes = apisixRoutes
	fmt.Println("Kong to APISIX routes configuration conversion completed")
	return nil
}

// GenerateProxyRewritePluginConfig Generate routing and forwarding rules
// https://docs.konghq.com/gateway-oss/2.4.x/admin-api/#path-handling-algorithms
func GenerateProxyRewritePluginConfig(servicePath string, routerPath string, stripPath bool,
	pathHandling string) map[string]interface{} {
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

	config := make(map[string]interface{})
	config["regex_uri"] = []string{pathRegex, pathPattern}

	return config
}
