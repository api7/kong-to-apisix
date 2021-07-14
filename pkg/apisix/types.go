package apisix

import v1 "github.com/apache/apisix-ingress-controller/pkg/types/apisix/v1"

type Config struct {
	Routes      *[]v1.Route      `yaml:"routes"`
	Upstreams   *[]v1.Upstream   `yaml:"upstreams"`
	GlobalRules *[]v1.GlobalRule `yaml:"global_rules"`
	Consumers   *[]v1.Consumer   `yaml:"consumers"`
}
