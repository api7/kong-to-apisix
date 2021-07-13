package kong

import (
	"fmt"
	"net/url"
	"strconv"

	v1 "github.com/apache/apisix-ingress-controller/pkg/types/apisix/v1"
)

func MigrateUpstream(kongConfig *KongConfig) (*[]v1.Upstream, error) {
	kongUpstreams := kongConfig.Upstreams
	kongServices := kongConfig.Services
	upstreamsMap := make(map[string]Upstream)
	for _, u := range kongUpstreams {
		upstreamsMap[u.Name] = u
	}

	var apisixUpstreams []v1.Upstream
	for i, s := range kongServices {
		kongConfig.Services[i].ID = strconv.Itoa(i)
		// TODO: gokong not support lbAlgorithm yet
		apisixUpstream := &v1.Upstream{
			Metadata: v1.Metadata{
				ID: strconv.Itoa(i),
			},
			Type:    "roundrobin",
			Scheme:  s.Protocol,
			Retries: s.Retries,
			Timeout: &v1.UpstreamTimeout{
				Connect: s.ConnectTimeout / 1000,
				Send:    s.WriteTimeout / 1000,
				Read:    s.ReadTimeout / 1000,
			},
		}

		if s.Name != "" {
			apisixUpstream.Metadata.Name = s.Name
		}

		var upstreamNodes []v1.UpstreamNode

		// if service is bind to upstream
		if upstream, ok := upstreamsMap[s.Host]; ok {
			apisixUpstream.HashOn = upstream.HashOn
			switch upstream.HashOn {
			case "none":
				apisixUpstream.HashOn = "vars"
			case "ip":
				fmt.Println("upstream hashon parameter `hashon` not supported in apisix")
				apisixUpstream.HashOn = "vars"
			default:
				apisixUpstream.HashOn = upstream.HashOn
			}

			targets := upstream.Targets

			for _, t := range targets {
				u, err := url.Parse(t.Target)
				if err != nil {
					return nil, err
				}

				if u.Host == "" {
					u, err = url.ParseRequestURI("http://" + t.Target)
					if err != nil {
						return nil, err
					}
				}

				port, err := strconv.Atoi(u.Port())
				if err != nil {
					return nil, err
				}
				upstreamNodes = append(upstreamNodes, v1.UpstreamNode{
					Host:   u.Hostname(),
					Port:   port,
					Weight: t.Weight,
				})
			}
		} else {
			upstreamNodes = []v1.UpstreamNode{{
				Host:   s.Host,
				Port:   s.Port,
				Weight: 1,
			}}
		}

		apisixUpstream.Nodes = upstreamNodes
		apisixUpstreams = append(apisixUpstreams, *apisixUpstream)
	}

	return &apisixUpstreams, nil
}
