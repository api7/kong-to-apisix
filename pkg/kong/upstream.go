package kong

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/api7/kong-to-apisix/pkg/apisix"

	uuid "github.com/satori/go.uuid"
)

const KTASixUpstreamAlgorithmRoundRobin = "roundrobin"
const KTASixUpstreamAlgorithmConsistentHashing = "chash"
const KTASixUpstreamAlgorithmLeastConnections = "least_conn"
const KTAKongUpstreamAlgorithmRoundRobin = "round-robin"
const KTAKongUpstreamAlgorithmConsistentHashing = "consistent-hashing"
const KTAKongUpstreamAlgorithmLeastConnections = "least-connections"

type SixServiceUpstreamMapping struct {
	ServiceId  string
	UpstreamId string
}

func MigrateUpstream(kongConfig *Config, apisixConfig *apisix.Config) error {
	kongServices := kongConfig.Services
	kongUpstreams := kongConfig.Upstreams
	var sixServiceUpstreamMappings []SixServiceUpstreamMapping
	for index, kongUpstream := range kongUpstreams {
		kongUpstreamId := kongUpstream.ID
		if len(kongUpstreamId) <= 0 {
			kongUpstreamId = uuid.NewV4().String()
			kongConfig.Upstreams[index].ID = kongUpstreamId
		}

		var sixServiceUpstreamMapping SixServiceUpstreamMapping
		var apisixUpstream apisix.Upstream

		apisixUpstream.ID = kongUpstreamId
		apisixUpstream.Name = kongUpstream.Name
		switch kongUpstream.Algorithm {
		case KTAKongUpstreamAlgorithmRoundRobin:
			apisixUpstream.Type = KTASixUpstreamAlgorithmRoundRobin
		case KTAKongUpstreamAlgorithmConsistentHashing:
			apisixUpstream.Type = KTASixUpstreamAlgorithmConsistentHashing
		case KTAKongUpstreamAlgorithmLeastConnections:
			apisixUpstream.Type = KTASixUpstreamAlgorithmLeastConnections
		}
		if len(kongUpstream.Targets) > 0 {
			// kong deck export
			apisixUpstream.Nodes = KTAConversionKongUpstreamTargets(kongUpstream.Targets, "")
		} else {
			// kong config export
			apisixUpstream.Nodes = KTAConversionKongUpstreamTargets(kongConfig.Targets, kongUpstreamId)
		}

		service, err := FindKongServiceByHost(&kongServices, kongUpstream.Name)
		if err != nil || service == nil {
			apisixConfig.Upstreams = append(apisixConfig.Upstreams, apisixUpstream)
			fmt.Printf("Kong upstream [ %s ] mapping service not found\n", kongUpstreamId)
			continue
		}

		apisixUpstream.Retries = service.Retries
		apisixUpstream.Scheme = service.Protocol
		apisixUpstream.Timeout.Connect = KTAConversionKongUpstreamTimeout(service.ConnectTimeout)
		apisixUpstream.Timeout.Send = KTAConversionKongUpstreamTimeout(service.WriteTimeout)
		apisixUpstream.Timeout.Read = KTAConversionKongUpstreamTimeout(service.ReadTimeout)
		apisixConfig.Upstreams = append(apisixConfig.Upstreams, apisixUpstream)
		sixServiceUpstreamMapping.ServiceId = service.ID
		sixServiceUpstreamMapping.UpstreamId = apisixUpstream.ID
		sixServiceUpstreamMappings = append(sixServiceUpstreamMappings, sixServiceUpstreamMapping)
		fmt.Printf("Kong upstream [ %s ] to APISIX conversion completed\n", apisixUpstream.ID)
	}

	// remapping service and upstream
	for serviceIndex, service := range apisixConfig.Services {
		for _, mapping := range sixServiceUpstreamMappings {
			if service.ID == mapping.ServiceId {
				apisixConfig.Services[serviceIndex].UpstreamID = mapping.UpstreamId
				break
			}
		}
	}

	fmt.Println("Kong to APISIX upstreams configuration conversion completed")
	return nil
}

func KTAConversionKongUpstreamTimeout(kongTime uint) float32 {
	return float32(kongTime) / float32(1000)
}

func KTAConversionKongUpstreamTargets(kongTargets Targets, kongUpstreamId string) apisix.UpstreamNodes {
	var targets Targets
	if len(kongUpstreamId) <= 0 {
		// kong deck export
		targets = kongTargets
	} else {
		// kong config export
		for _, target := range kongTargets {
			if target.UpstreamID == kongUpstreamId {
				targets = append(targets, target)
			}
		}
	}

	var sixUpstreamNodes apisix.UpstreamNodes
	for _, target := range targets {
		var sixUpstreamNode apisix.UpstreamNode
		targetResponse := strings.Split(target.Target, ":")
		switch len(targetResponse) {
		case 1:
			sixUpstreamNode.Host = targetResponse[0]
			sixUpstreamNode.Port = 80
			sixUpstreamNode.Weight = target.Weight
		case 2:
			port, err := strconv.Atoi(targetResponse[1])
			if err != nil {
				continue
			}
			sixUpstreamNode.Host = targetResponse[0]
			sixUpstreamNode.Port = port
			sixUpstreamNode.Weight = target.Weight
		default:
			continue
		}
		sixUpstreamNodes = append(sixUpstreamNodes, sixUpstreamNode)
	}
	return sixUpstreamNodes
}
