package upstream

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/apache/apisix-ingress-controller/pkg/apisix"
	v1 "github.com/apache/apisix-ingress-controller/pkg/types/apisix/v1"
	"github.com/globocom/gokong"
)

func MigrateUpstream(apisixCli apisix.Cluster, kongCli gokong.KongAdminClient) error {
	upstreams, err := kongCli.Upstreams().List()
	if err != nil {
		return err
	}
	upstreamsMap := make(map[string]*gokong.Upstream)
	for _, u := range upstreams.Results {
		upstreamsMap[u.Name] = u
		//fmt.Printf("got upstream: %#v\n", u)
	}

	services, err := kongCli.Services().GetServices(&gokong.ServiceQueryString{})
	if err != nil {
		return err
	}

	for _, s := range services {
		//fmt.Printf("got service: %#v\n", s)

		// TODO: gokong not support lbAlgorithm yet
		apisixUpstream := &v1.Upstream{
			Metadata: v1.Metadata{
				ID: *s.Id,
			},
			Type:    "roundrobin",
			Scheme:  *s.Protocol,
			Retries: *s.Retries,
			Timeout: &v1.UpstreamTimeout{
				Connect: *s.ConnectTimeout / 1000,
				Send:    *s.WriteTimeout / 1000,
				Read:    *s.ReadTimeout / 1000,
			},
		}

		if s.Name != nil {
			apisixUpstream.Metadata.Name = *s.Name
		}

		var upstreamNodes []v1.UpstreamNode

		// if service is bind to upstream
		if upstream, ok := upstreamsMap[*s.Host]; ok {
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

			targets, err := kongCli.Targets().GetTargetsFromUpstreamId(upstream.Id)
			if err != nil {
				return err
			}

			for _, t := range targets {
				u, err := url.Parse(*t.Target)
				if err != nil {
					return err
				}

				if u.Host == "" {
					u, err = url.ParseRequestURI("http://" + *t.Target)
					if err != nil {
						return err
					}
				}

				port, err := strconv.Atoi(u.Port())
				if err != nil {
					return err
				}
				upstreamNodes = append(upstreamNodes, v1.UpstreamNode{
					Host:   u.Hostname(),
					Port:   port,
					Weight: *t.Weight,
				})
			}
		} else {
			upstreamNodes = []v1.UpstreamNode{{
				Host:   *s.Host,
				Port:   *s.Port,
				Weight: 1,
			}}
		}

		apisixUpstream.Nodes = upstreamNodes
		_, err := apisixCli.Upstream().Create(context.Background(), apisixUpstream)
		if err != nil {
			return err
		}

		var printName string
		if s.Name != nil {
			printName = *s.Name
		} else {
			printName = *s.Id
		}
		fmt.Printf("migrate service %s to upstream succeeds\n", printName)
	}

	return nil
}
