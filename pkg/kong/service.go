package kong

import (
	"errors"
	"fmt"

	"github.com/api7/kong-to-apisix/pkg/apisix"

	uuid "github.com/satori/go.uuid"
)

func MigrateService(kongConfig *Config, apisixConfig *apisix.Config) error {
	kongServices := kongConfig.Services
	apisixServices := apisixConfig.Services

	for _, kongService := range kongServices {
		if len(kongService.ID) <= 0 {
			kongService.ID = uuid.NewV4().String()
		}
		var apisixService apisix.Service
		apisixService.ID = kongService.ID
		apisixService.Name = kongService.Name
		upstreamID := GenerateApisixServiceUpstream(kongService, apisixConfig)
		apisixService.UpstreamID = upstreamID
		apisixServices = append(apisixServices, apisixService)
		fmt.Printf("Kong service [ %s ] to APISIX conversion completed\n", kongService.ID)
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
	apisixUpstream.Type = "roundrobin"
	// apisix upstream nodes
	var apisixUpstreamNode apisix.UpstreamNode
	apisixUpstreamNode.Weight = 100
	apisixUpstreamNode.Host = kongService.Host
	if kongService.Port > 0 {
		apisixUpstreamNode.Port = kongService.Port
	} else {
		apisixUpstreamNode.Port = 80
	}
	apisixUpstream.Nodes = append(apisixUpstream.Nodes, apisixUpstreamNode)
	// apisix upstream timeout
	var apisixUpstreamTimeout apisix.UpstreamTimeout
	apisixUpstreamTimeout.Send = float32(kongService.WriteTimeout) / float32(1000)
	apisixUpstreamTimeout.Read = float32(kongService.ReadTimeout) / float32(1000)
	apisixUpstreamTimeout.Connect = float32(kongService.ConnectTimeout) / float32(1000)
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
