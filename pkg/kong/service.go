package kong

import (
	"errors"
	"fmt"
	"os"

	uuid "github.com/satori/go.uuid"

	"github.com/api7/kong-to-apisix/pkg/apisix"
)

func MigrateService(kongConfig *Config, apisixConfig *apisix.Config) error {
	kongServices := kongConfig.Services
	apisixServices := apisixConfig.Services

	for _, kongService := range kongServices {
		var apisixService apisix.Service
		apisixService.ID = kongService.ID
		apisixService.Name = kongService.Name
		upstreamID := GenerateApisixServiceUpstream(kongService, apisixConfig)
		apisixService.UpstreamID = upstreamID
		apisixServices = append(apisixServices, apisixService)
		fmt.Fprintf(os.Stdout, "Kong service [ %s ] to APISIX conversion completed\n", kongService.ID)
	}

	apisixConfig.Services = apisixServices
	fmt.Println("Kong to APISIX services configuration conversion completed")
	return nil
}

func GenerateApisixServiceUpstream(kongService Service, apisixConfig *apisix.Config) string {
	var apisixUpstream apisix.Upstream
	// apisix upstream id
	if len(kongService.ID) > 0 {
		apisixUpstream.ID = "svc-" + kongService.ID + "-ups"
	} else {
		apisixUpstream.ID = uuid.NewV4().String()
	}
	// apisix upstream nodes
	var apisixUpstreamNode apisix.UpstreamNode
	apisixUpstreamNode.Weight = 100
	apisixUpstreamNode.Host = kongService.Host
	if kongService.Port > 0 {
		apisixUpstreamNode.Port = kongService.Port
	} else {
		switch kongService.Protocol {
		case "https":
			apisixUpstreamNode.Port = 443
		default:
			apisixUpstreamNode.Port = 80
		}
	}
	apisixUpstream.Nodes = append(apisixUpstream.Nodes, apisixUpstreamNode)
	// apisix upstream timeout
	var apisixUpstreamTimeout apisix.UpstreamTimeout
	apisixUpstreamTimeout.Send = kongService.WriteTimeout / 1000
	apisixUpstreamTimeout.Read = kongService.ReadTimeout / 1000
	apisixUpstreamTimeout.Connect = kongService.ConnectTimeout / 1000
	apisixUpstream.Timeout = apisixUpstreamTimeout
	// apisix upstream scheme
	apisixUpstream.Scheme = kongService.Protocol
	// apisix upstream retries
	apisixUpstream.Retries = kongService.Retries

	apisixConfig.Upstreams = append(apisixConfig.Upstreams, apisixUpstream)

	return apisixUpstream.ID
}

func GetKongServiceByID(kongServices *Services, serviceID string) (*Service, error) {
	if kongServices == nil {
		return nil, errors.New("kong services is nil or invalid")
	}

	if len(serviceID) <= 0 {
		return nil, errors.New("kong service id invalid")
	}

	var kongService *Service
	for _, service := range *kongServices {
		if service.ID == serviceID {
			kongService = &service
			break
		}
	}

	if kongService != nil {
		return kongService, nil
	}

	return nil, fmt.Errorf("kong service id /%s/ not found", serviceID)
}
