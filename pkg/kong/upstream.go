package kong

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/api7/kong-to-apisix/pkg/apisix"
	"github.com/pkg/errors"
)

func MigrateUpstream(kongConfig *Config, apisixConfig *apisix.Config) error {
	kongUpstreams := kongConfig.Upstreams
	kongServices := kongConfig.Services
	upstreamsMap := make(map[string]Upstream)
	for _, u := range kongUpstreams {
		upstreamsMap[u.Name] = u
	}

	// TODO: Temporarily compatible, this module will be refactored
	apisixUpstreams := apisixConfig.Upstreams
	for i, s := range kongServices {
		kongConfig.Services[i].ID = strconv.Itoa(i)
		// TODO: gokong not support lbAlgorithm yet
		apisixUpstream := &apisix.Upstream{
			ID:      strconv.Itoa(i),
			Type:    "roundrobin",
			Scheme:  s.Protocol,
			Retries: uint(s.Retries),
			Timeout: apisix.UpstreamTimeout{
				Connect: KTATimeoutConversion(s.ConnectTimeout),
				Send:    KTATimeoutConversion(s.WriteTimeout),
				Read:    KTATimeoutConversion(s.ReadTimeout),
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
					return errors.Wrap(err, "url parse")
				}

				if u == nil || u.Host == "" {
					u, err = url.ParseRequestURI("http://" + t.Target)
					if err != nil {
						return err
					}
				}

				port, err := strconv.Atoi(u.Port())
				if err != nil {
					return err
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

	apisixConfig.Upstreams = apisixUpstreams

	return nil
}
