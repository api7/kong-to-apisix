package kong

import (
	"fmt"
	"github.com/api7/kong-to-apisix/pkg/apisix"
	"net/url"
	"strconv"
	"strings"

	"github.com/api7/kong-to-apisix/pkg/utils"
	"github.com/pkg/errors"
)

func MigrateUpstream(kongConfig *Config, configYamlAll *[]utils.YamlItem) (apisix.Upstreams, error) {
	kongUpstreams := kongConfig.Upstreams
	kongServices := kongConfig.Services
	upstreamsMap := make(map[string]Upstream)
	for _, u := range kongUpstreams {
		upstreamsMap[u.Name] = u
	}

	var apisixUpstreams apisix.Upstreams
	for i, s := range kongServices {
		kongConfig.Services[i].ID = strconv.Itoa(i)
		// TODO: gokong not support lbAlgorithm yet
		apisixUpstream := &apisix.Upstream{
			ID:      strconv.Itoa(i),
			Type:    "roundrobin",
			Scheme:  s.Protocol,
			Retries: uint(s.Retries),
			Timeout: apisix.UpstreamTimeout{
				Connect: uint(s.ConnectTimeout) / 1000,
				Send:    uint(s.WriteTimeout) / 1000,
				Read:    uint(s.ReadTimeout) / 1000,
			},
		}

		if s.Name != "" {
			apisixUpstream.Name = s.Name
		}

		var upstreamNodes []apisix.UpstreamNode

		// if service is bind to upstream
		if upstream, ok := upstreamsMap[s.Host]; ok {
			apisixUpstream.HashOn = upstream.HashOn
			switch upstream.HashOn {
			case "none":
				apisixUpstream.HashOn = "vars"
			case "ip":
				fmt.Println("upstream hashon parameter `ip` not supported in apisix")
				apisixUpstream.HashOn = "vars"
			default:
				apisixUpstream.HashOn = upstream.HashOn
			}

			targets := upstream.Targets

			for _, t := range targets {
				u, err := url.Parse(t.Target)
				if err != nil && !strings.Contains(err.Error(), "first path segment in URL cannot contain colon") {
					return nil, errors.Wrap(err, "url parse")
				}

				if u == nil || u.Host == "" {
					u, err = url.ParseRequestURI("http://" + t.Target)
					if err != nil {
						return nil, err
					}
				}

				port, err := strconv.Atoi(u.Port())
				if err != nil {
					return nil, err
				}
				upstreamNodes = append(upstreamNodes, apisix.UpstreamNode{
					Host:   u.Hostname(),
					Port:   port,
					Weight: t.Weight,
				})
			}
		} else {
			upstreamNodes = []apisix.UpstreamNode{{
				Host:   s.Host,
				Port:   s.Port,
				Weight: 1,
			}}
		}

		apisixUpstream.Nodes = upstreamNodes
		apisixUpstreams = append(apisixUpstreams, *apisixUpstream)
	}

	return apisixUpstreams, nil
}
