package kong

import (
	"fmt"
	"net"
	"regexp"
	"strconv"

	"github.com/api7/kong-to-apisix/pkg/apisix"

	uuid "github.com/satori/go.uuid"
)

const KTASixUpstreamAlgorithmRoundRobin = "roundrobin"
const KTASixUpstreamAlgorithmConsistentHashing = "chash"
const KTASixUpstreamAlgorithmLeastConnections = "least_conn"
const KTAKongUpstreamAlgorithmRoundRobin = "round-robin"
const KTAKongUpstreamAlgorithmConsistentHashing = "consistent-hashing"
const KTAKongUpstreamAlgorithmLeastConnections = "least-connections"

type A6UpstreamNodeHostPort struct {
	Host string
	Port int
}

type A6ServiceUpstreamMapping struct {
	ServiceId  string
	UpstreamId string
}

func MigrateUpstream(kongConfig *Config, apisixConfig *apisix.Config) error {
	kongServices := kongConfig.Services
	kongUpstreams := kongConfig.Upstreams
	var a6ServiceUpstreamMappings []A6ServiceUpstreamMapping
	for index, kongUpstream := range kongUpstreams {
		kongUpstreamId := kongUpstream.ID
		if len(kongUpstreamId) <= 0 {
			kongUpstreamId = uuid.NewV4().String()
			kongConfig.Upstreams[index].ID = kongUpstreamId
		}

		var a6ServiceUpstreamMapping A6ServiceUpstreamMapping
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
		a6ServiceUpstreamMapping.ServiceId = service.ID
		a6ServiceUpstreamMapping.UpstreamId = apisixUpstream.ID
		a6ServiceUpstreamMappings = append(a6ServiceUpstreamMappings, a6ServiceUpstreamMapping)
		fmt.Printf("Kong upstream [ %s ] to APISIX conversion completed\n", apisixUpstream.ID)
	}

	// remapping service and upstream
	for serviceIndex, service := range apisixConfig.Services {
		for _, mapping := range a6ServiceUpstreamMappings {
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

	var a6UpstreamNodes apisix.UpstreamNodes
	for _, target := range targets {
		var a6UpstreamNode apisix.UpstreamNode
		a6NodeHostPort, err := KTAFormatKongTarget(target.Target)
		if a6NodeHostPort == nil {
			fmt.Printf("Kong upstream [ %s ] target [ %s ] conversion failure, %s\n",
				kongUpstreamId, target.Target, err)
			continue
		}
		a6UpstreamNode.Host = a6NodeHostPort.Host
		a6UpstreamNode.Port = a6NodeHostPort.Port
		a6UpstreamNode.Weight = target.Weight
		a6UpstreamNodes = append(a6UpstreamNodes, a6UpstreamNode)
	}
	return a6UpstreamNodes
}

func KTAFormatKongTarget(kongTarget string) (*A6UpstreamNodeHostPort, error) {
	reg := regexp.MustCompile(`\:([1-9]|[1-5]?[0-9]{2,4}|6[1-4][0-9]{3}|65[1-4][0-9]{2}|655[1-2][0-9]|6553[1-5])$`)
	if len(reg.FindAllString(kongTarget, -1)) < 1 {
		kongTarget = kongTarget + ":" + strconv.Itoa(80)
	}
	host, port, err := net.SplitHostPort(kongTarget)
	if err != nil {
		return nil, err
	}
	portNumber, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}
	var a6UpstreamNodeHostPort A6UpstreamNodeHostPort
	a6UpstreamNodeHostPort.Host = host
	a6UpstreamNodeHostPort.Port = portNumber
	return &a6UpstreamNodeHostPort, nil
}
