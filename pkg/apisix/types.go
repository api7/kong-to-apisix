package apisix

type Config struct {
	Services    Services    `json:"services,omitempty" yaml:"services,omitempty"`
	Routes      Routes      `json:"routes,omitempty" yaml:"routes,omitempty"`
	Upstreams   Upstreams   `json:"upstreams,omitempty" yaml:"upstreams,omitempty"`
	GlobalRules GlobalRules `json:"global_rules,omitempty" yaml:"global_rules,omitempty"`
	Consumers   Consumers   `json:"consumers,omitempty" yaml:"consumers,omitempty"`
}

// Route Configuration
// Route is the apisix route definition.
type Route struct {
	ID              string            `json:"id,omitempty" yaml:"id,omitempty"`
	Name            string            `json:"name,omitempty" yaml:"name,omitempty"`
	URI             string            `json:"uri,omitempty" yaml:"uri,omitempty"`
	URIs            []string          `json:"uris,omitempty" yaml:"uris,omitempty"`
	Host            string            `json:"host,omitempty" yaml:"host,omitempty"`
	Hosts           []string          `json:"hosts,omitempty" yaml:"hosts,omitempty"`
	RemoteAddress   string            `json:"remote_addr,omitempty" yaml:"remote_addr,omitempty"`
	RemoteAddresses []string          `json:"remote_addrs,omitempty" yaml:"remote_addrs,omitempty"`
	Methods         []string          `json:"methods,omitempty" yaml:"methods,omitempty"`
	Priority        uint              `json:"priority,omitempty" yaml:"priority,omitempty"`
	Vars            [][]string        `json:"vars,omitempty" yaml:"vars,omitempty"`
	FilterFunc      string            `json:"filter_func,omitempty" yaml:"filter_func,omitempty"`
	Plugins         Plugins           `json:"plugins,omitempty" yaml:"plugins,omitempty"`
	Script          string            `json:"script,omitempty" yaml:"script,omitempty"`
	Upstream        *Upstream         `json:"upstream,omitempty" yaml:"upstream,omitempty"`
	UpstreamID      string            `json:"upstream_id,omitempty" yaml:"upstream_id,omitempty"`
	ServiceID       string            `json:"service_id,omitempty" yaml:"service_id,omitempty"`
	PluginConfigID  string            `json:"plugin_config_id,omitempty" yaml:"plugin_config_id,omitempty"`
	Labels          map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Timeout         UpstreamTimeout   `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	EnableWebsocket bool              `json:"enable_websocket,omitempty" yaml:"enable_websocket,omitempty"`
	Status          uint              `json:"status,omitempty" yaml:"status,omitempty"`
}

type Routes []Route

// Service Configuration
// Service apisix route definition
type Service struct {
	ID              string            `json:"id,omitempty" yaml:"id,omitempty"`
	Name            string            `json:"name,omitempty" yaml:"name,omitempty"`
	Desc            string            `json:"desc,omitempty" yaml:"desc,omitempty"`
	Labels          map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	EnableWebsocket bool              `json:"enable_websocket,omitempty" yaml:"enable_websocket,omitempty"`
	Hosts           []string          `json:"hosts,omitempty" yaml:"hosts,omitempty"`
	Upstream        *Upstream         `json:"upstreams,omitempty" yaml:"upstreams,omitempty"`
	UpstreamID      string            `json:"upstream_id,omitempty" yaml:"upstream_id,omitempty"`
	Plugins         Plugins           `json:"plugins,omitempty" yaml:"plugins,omitempty"`
}

type Services []Service

// Consumer Configuration
// Consumer is the apisix consumer definition.
type Consumer struct {
	Username string            `json:"username" yaml:"username"`
	Desc     string            `json:"desc,omitempty" yaml:"desc,omitempty"`
	Labels   map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Plugins  ConsumerPlugins   `json:"plugins,omitempty" yaml:"plugins,omitempty"`
}

type ConsumerPlugins struct {
	KeyAuth KeyAuthCredential `json:"key-auth,omitempty" yaml:"key-auth,omitempty"`
}

type KeyAuthCredential struct {
	Key string `json:"key,omitempty" yaml:"key,omitempty"`
}

type Consumers []Consumer

// Upstream Configuration
// Upstream is the apisix upstream definition.
type Upstream struct {
	ID            string                `json:"id,omitempty" yaml:"id,omitempty"`
	Name          string                `json:"name,omitempty" yaml:"name,omitempty"`
	Desc          string                `json:"desc,omitempty" yaml:"desc,omitempty"`
	Type          string                `json:"type,omitempty" yaml:"type,omitempty"`
	Nodes         []UpstreamNode        `json:"nodes,omitempty" yaml:"nodes,omitempty"`
	ServiceName   string                `json:"service_name,omitempty" yaml:"service_name,omitempty"`
	DiscoveryType string                `json:"discovery_type,omitempty" yaml:"discovery_type,omitempty"`
	HashOn        string                `json:"hash_on,omitempty" yaml:"hash_on,omitempty"`
	Key           string                `json:"key,omitempty" yaml:"key,omitempty"`
	Checks        *UpstreamHealthCheck  `json:"checks,omitempty" yaml:"checks,omitempty"`
	Retries       uint                  `json:"retries,omitempty" yaml:"retries,omitempty"`
	RetryTimeout  uint                  `json:"retry_timeout,omitempty" yaml:"retry_timeout,omitempty"`
	Timeout       UpstreamTimeout       `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	PassHost      string                `json:"pass_host,omitempty" yaml:"pass_host,omitempty"`
	UpstreamHost  string                `json:"upstream_host,omitempty" yaml:"upstream_host,omitempty"`
	Scheme        string                `json:"scheme,omitempty" yaml:"scheme,omitempty"`
	Labels        map[string]string     `json:"labels,omitempty" yaml:"labels,omitempty"`
	TLS           UpstreamTLS           `json:"tls,omitempty" yaml:"tls,omitempty"`
	KeepalivePool UpstreamKeepalivePool `json:"keepalive_pool,omitempty" yaml:"keepalive_pool,omitempty"`
}

type Upstreams []Upstream

// UpstreamTimeout is the apisix upstream.timeout definition.
type UpstreamTimeout struct {
	Connect float32 `json:"connect,omitempty" yaml:"connect,omitempty"`
	Send    float32 `json:"send,omitempty" yaml:"send,omitempty"`
	Read    float32 `json:"read,omitempty" yaml:"read,omitempty"`
}

// UpstreamNode is the apisix upstream[index].node definition.
type UpstreamNode struct {
	Host   string `json:"host,omitempty" yaml:"host,omitempty"`
	Port   int    `json:"port,omitempty" yaml:"port,omitempty"`
	Weight int    `json:"weight,omitempty" yaml:"weight,omitempty"`
}

// UpstreamTLS is the apisix upstream.tls definition.
type UpstreamTLS struct {
	ClientCert string `json:"client_cert,omitempty" yaml:"client_cert,omitempty"`
	ClientKey  string `json:"client_key,omitempty" yaml:"client_key,omitempty"`
}

// UpstreamKeepalivePool is the apisix upstream.keepalive_pool definition.
type UpstreamKeepalivePool struct {
	Size        uint `json:"size,omitempty" yaml:"size,omitempty"`
	IdleTimeout uint `json:"idle_timeout,omitempty" yaml:"idle_timeout,omitempty"`
	Requests    uint `json:"requests,omitempty" yaml:"requests,omitempty"`
}

// UpstreamHealthCheck is the apisix upstream.checks definition.
type UpstreamHealthCheck struct {
	Active  *UpstreamActiveHealthCheck  `json:"active,omitempty" yaml:"active,omitempty"`
	Passive *UpstreamPassiveHealthCheck `json:"passive,omitempty" yaml:"passive,omitempty"`
}

// UpstreamActiveHealthCheck is the apisix upstream.checks.active definition.
type UpstreamActiveHealthCheck struct {
	Type            string                             `json:"type,omitempty" yaml:"type,omitempty"`
	Timeout         int                                `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Concurrency     int                                `json:"concurrency,omitempty" yaml:"concurrency,omitempty"`
	Host            string                             `json:"host,omitempty" yaml:"host,omitempty"`
	Port            int32                              `json:"port,omitempty" yaml:"port,omitempty"`
	HTTPPath        string                             `json:"http_path,omitempty" yaml:"http_path,omitempty"`
	HTTPSVerifyCert bool                               `json:"https_verify_certificate,omitempty" yaml:"https_verify_certificate,omitempty"`
	RequestHeaders  []string                           `json:"req_headers,omitempty" yaml:"req_headers,omitempty"`
	Healthy         UpstreamActiveHealthCheckHealthy   `json:"healthy,omitempty" yaml:"healthy,omitempty"`
	Unhealthy       UpstreamActiveHealthCheckUnhealthy `json:"unhealthy,omitempty" yaml:"unhealthy,omitempty"`
}

// UpstreamActiveHealthCheckHealthy is the apisix upstream.checks.active.healthy definition.
type UpstreamActiveHealthCheckHealthy struct {
	UpstreamPassiveHealthCheckHealthy `json:",inline" yaml:",inline"`
	Interval                          int `json:"interval,omitempty" yaml:"interval,omitempty"`
}

// UpstreamActiveHealthCheckUnhealthy is the apisix upstream.checks.active.unhealthy definition.
type UpstreamActiveHealthCheckUnhealthy struct {
	UpstreamPassiveHealthCheckUnhealthy `json:",inline" yaml:",inline"`
	Interval                            int `json:"interval,omitempty" yaml:"interval,omitempty"`
}

// UpstreamPassiveHealthCheck is the apisix upstream.checks.passive definition.
type UpstreamPassiveHealthCheck struct {
	Type      string                              `json:"type,omitempty" yaml:"type,omitempty"`
	Healthy   UpstreamPassiveHealthCheckHealthy   `json:"healthy,omitempty" yaml:"healthy,omitempty"`
	Unhealthy UpstreamPassiveHealthCheckUnhealthy `json:"unhealthy,omitempty" yaml:"unhealthy,omitempty"`
}

// UpstreamPassiveHealthCheckHealthy is the apisix upstream.checks.passive.healthy definition.
type UpstreamPassiveHealthCheckHealthy struct {
	HTTPStatuses []int `json:"http_statuses,omitempty" yaml:"http_statuses,omitempty"`
	Successes    int   `json:"successes,omitempty" yaml:"successes,omitempty"`
}

// UpstreamPassiveHealthCheckUnhealthy is the apisix upstream.checks.passive.unhealthy definition.
type UpstreamPassiveHealthCheckUnhealthy struct {
	HTTPStatuses []int   `json:"http_statuses,omitempty" yaml:"http_statuses,omitempty"`
	HTTPFailures int     `json:"http_failures,omitempty" yaml:"http_failures,omitempty"`
	TCPFailures  int     `json:"tcp_failures,omitempty" yaml:"tcp_failures,omitempty"`
	Timeouts     float64 `json:"timeouts,omitempty" yaml:"timeouts,omitempty"`
}

// SSL Configuration
// SSL is the apisix ssl definition.
type SSL struct {
	ID     string            `json:"id,omitempty" yaml:"id,omitempty"`
	Cert   string            `json:"cert,omitempty" yaml:"cert,omitempty"`
	Key    string            `json:"key,omitempty" yaml:"key,omitempty"`
	Certs  []string          `json:"certs,omitempty" yaml:"certs,omitempty"`
	Keys   []string          `json:"keys,omitempty" yaml:"keys,omitempty"`
	Client SSLClient         `json:"client,omitempty" yaml:"client,omitempty"`
	SNIs   []string          `json:"snis,omitempty" yaml:"snis,omitempty"`
	Labels map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Status uint              `json:"status,omitempty" yaml:"status,omitempty"`
}

type SSLs []SSL

// SSLClient is the apisix ssl.client definition.
type SSLClient struct {
	CA    string `json:"ca,omitempty" yaml:"ca,omitempty"`
	Depth uint   `json:"depth,omitempty" yaml:"depth,omitempty"`
}

// Plugins Configuration
// Plugins is the apisix plugins definition
type Plugins struct {
	KeyAuth      KeyAuth      `json:"key-auth,omitempty" yaml:"key-auth,omitempty"`
	ProxyRewrite ProxyRewrite `json:"proxy-rewrite,omitempty" yaml:"proxy-rewrite,omitempty"`
}

type ProxyRewrite struct {
	Scheme   string            `json:"scheme,omitempty" yaml:"scheme,omitempty"`
	URI      string            `json:"uri,omitempty" yaml:"uri,omitempty"`
	RegexURI []string          `json:"regex_uri,omitempty" yaml:"regex_uri,omitempty"`
	Host     string            `json:"host,omitempty" yaml:"host,omitempty"`
	Headers  map[string]string `json:"headers,omitempty" yaml:"headers,omitempty"`
}

type KeyAuth struct {
	Header string `json:"header,omitempty" yaml:"header,omitempty"`
	Query  string `json:"query,omitempty" yaml:"query,omitempty"`
}

// GlobalRule Configuration
// GlobalRule is the apisix global_rules definition.
type GlobalRule struct {
	ID      string  `json:"id,omitempty" yaml:"id,omitempty"`
	Plugins Plugins `json:"plugins,omitempty" yaml:"plugins,omitempty"`
}

type GlobalRules []GlobalRule
