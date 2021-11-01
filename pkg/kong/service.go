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

	for index, kongService := range kongServices {
		kongServiceId := kongService.ID
		if len(kongServiceId) <= 0 {
			kongServiceId = uuid.NewV4().String()
			kongConfig.Services[index].ID = kongServiceId
		}
		var apisixService apisix.Service
		apisixService.ID = kongServiceId
		apisixService.Name = kongService.Name
		if len(kongService.Host) > 0 {
			upstreamID := GenerateApisixServiceUpstream(kongService, apisixConfig)
			apisixService.UpstreamID = upstreamID
		}
		apisixServices = append(apisixServices, apisixService)
		fmt.Printf("Kong service [ %s ] to APISIX conversion completed\n", kongServiceId)
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
	apisixUpstream.Type = KTASixUpstreamAlgorithmRoundRobin
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
	apisixUpstreamTimeout.Send = KTAConversionKongUpstreamTimeout(kongService.WriteTimeout)
	apisixUpstreamTimeout.Read = KTAConversionKongUpstreamTimeout(kongService.ReadTimeout)
	apisixUpstreamTimeout.Connect = KTAConversionKongUpstreamTimeout(kongService.ConnectTimeout)
	apisixUpstream.Timeout = apisixUpstreamTimeout
	// apisix upstream scheme
	apisixUpstream.Scheme = kongService.Protocol
	// apisix upstream retries
	apisixUpstream.Retries = kongService.Retries

	apisixConfig.Upstreams = append(apisixConfig.Upstreams, apisixUpstream)

	return apisixUpstream.ID
}

func FindKongServiceById(kongServices *Services, serviceID string) (*Service, error) {
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

	if kongService == nil {
		return nil, fmt.Errorf("no service matching id /%s/ was found", serviceID)
	}

	return kongService, nil
}

func FindKongServiceByHost(kongServices *Services, serviceHost string) (*Service, error) {
	if kongServices == nil {
		return nil, errors.New("kong services is nil or invalid")
	}

	if len(serviceHost) <= 0 {
		return nil, errors.New("kong service id invalid")
	}

	var kongService *Service
	for _, service := range *kongServices {
		if service.Host == serviceHost {
			kongService = &service
			break
		}
	}

	if kongService == nil {
		return nil, fmt.Errorf("no service matching host /%s/ was found", serviceHost)
	}

	return kongService, nil
}
