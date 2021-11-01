package kong

type Config struct {
	Services             Services             `json:"services,omitempty" yaml:"services,omitempty"`
	Consumers            Consumers            `json:"consumers,omitempty" yaml:"consumers,omitempty"`
	Plugins              Plugins              `json:"plugins,omitempty" yaml:"plugins,omitempty"`
	Upstreams            Upstreams            `json:"upstreams,omitempty" yaml:"upstreams,omitempty"`
	Routes               Routes               `json:"routes,omitempty" yaml:"routes,omitempty"`
	Targets              Targets              `json:"targets,omitempty" yaml:"targets,omitempty"`
	KeyAuthCredentials   KeyAuthCredentials   `json:"keyauth_credentials,omitempty" yaml:"keyauth_credentials,omitempty"`
	BasicAuthCredentials BasicAuthCredentials `json:"basicauth_credentials,omitempty" yaml:"basicauth_credentials,omitempty"`
	HmacAuthCredentials  HmacAuthCredentials  `json:"hmacauth_credentials,omitempty" yaml:"hmacauth_credentials,omitempty"`
	JwtSecrets           JwtSecrets           `json:"jwt_secrets,omitempty" yaml:"jwt_secrets,omitempty"`
}

// Common Configuration

// CIDRPort is the kong cidr port definition.
type CIDRPort struct {
	IP   string `json:"ip,omitempty" yaml:"ip,omitempty"`
	Port int    `json:"port,omitempty" yaml:"port,omitempty"`
}

// Configuration is the kong plugin config definition.
type Configuration map[string]interface{}

// Service Configuration

// Service is the kong service definition.
type Service struct {
	ID                string      `json:"id,omitempty" yaml:"id,omitempty"`
	Name              string      `json:"name,omitempty" yaml:"name,omitempty"`
	Retries           uint        `json:"retries,omitempty" yaml:"retries,omitempty"`
	Protocol          string      `json:"protocol,omitempty" yaml:"protocol,omitempty"`
	Host              string      `json:"host,omitempty" yaml:"host,omitempty"`
	Path              string      `json:"path,omitempty" yaml:"path,omitempty"`
	Port              int         `json:"port,omitempty" yaml:"port,omitempty"`
	ConnectTimeout    uint        `json:"connect_timeout,omitempty" yaml:"connect_timeout,omitempty"`
	ReadTimeout       uint        `json:"read_timeout,omitempty" yaml:"read_timeout,omitempty"`
	WriteTimeout      uint        `json:"write_timeout,omitempty" yaml:"write_timeout,omitempty"`
	Tags              []string    `json:"tags,omitempty" yaml:"tags,omitempty"`
	ClientCertificate Certificate `json:"client_certificate,omitempty" yaml:"client_certificate,omitempty"`
	TLSVerify         bool        `json:"tls_verify,omitempty" yaml:"tls_verify,omitempty"`
	TLSVerifyDepth    int         `json:"tls_verify_depth,omitempty" yaml:"tls_verify_depth,omitempty"`
	CACertificates    []string    `json:"ca_certificates,omitempty" yaml:"ca_certificates,omitempty"`
	URL               string      `json:"url,omitempty" yaml:"url,omitempty"`
	Routes            Routes      `json:"routes,omitempty" yaml:"routes,omitempty"`
	Plugins           Plugins     `json:"plugins,omitempty" yaml:"plugins,omitempty"`
}

type Services []Service

// Route Configuration

// Route is the kong route definition.
type Route struct {
	ID                      string              `json:"id,omitempty" yaml:"id,omitempty"`
	Name                    string              `json:"name,omitempty" yaml:"name,omitempty"`
	Protocols               []string            `json:"protocols,omitempty" yaml:"protocols,omitempty"`
	Methods                 []string            `json:"methods,omitempty" yaml:"methods,omitempty"`
	Hosts                   []string            `json:"hosts,omitempty" yaml:"hosts,omitempty"`
	Paths                   []string            `json:"paths,omitempty" yaml:"paths,omitempty"`
	Headers                 map[string][]string `json:"headers,omitempty" yaml:"headers,omitempty"`
	HTTPSRedirectStatusCode int                 `json:"https_redirect_status_code,omitempty" yaml:"https_redirect_status_code,omitempty"`
	RegexPriority           uint                `json:"regex_priority,omitempty" yaml:"regex_priority,omitempty"`
	StripPath               bool                `json:"strip_path,omitempty" yaml:"strip_path,omitempty"`
	PathHandling            string              `json:"path_handling,omitempty" yaml:"path_handling,omitempty"`
	PreserveHost            bool                `json:"preserve_host,omitempty" yaml:"preserve_host,omitempty"`
	RequestBuffering        bool                `json:"request_buffering,omitempty" yaml:"request_buffering,omitempty"`
	ResponseBuffering       bool                `json:"response_buffering,omitempty" yaml:"response_buffering,omitempty"`
	SNIs                    []string            `json:"snis,omitempty" yaml:"snis,omitempty"`
	Sources                 []CIDRPort          `json:"sources,omitempty" yaml:"sources,omitempty"`
	Destinations            []CIDRPort          `json:"destinations,omitempty" yaml:"destinations,omitempty"`
	Tags                    []string            `json:"tags,omitempty" yaml:"tags,omitempty"`
	ServiceID               string              `json:"service,omitempty" yaml:"service,omitempty"`
	Plugins                 Plugins             `json:"plugins,omitempty" yaml:"plugins,omitempty"`
}

type Routes []Route

// Consumer Configuration

type KeyAuthCredential struct {
	ID         string `json:"id,omitempty" yaml:"id,omitempty"`
	ConsumerID string `json:"consumer,omitempty" yaml:"consumer,omitempty"`
	Key        string `json:"key,omitempty" yaml:"key,omitempty"`
}

type KeyAuthCredentials []KeyAuthCredential

type BasicAuthCredential struct {
	ID         string `json:"id,omitempty" yaml:"id,omitempty"`
	ConsumerID string `json:"consumer,omitempty" yaml:"consumer,omitempty"`
	Username   string `json:"username,omitempty" yaml:"username,omitempty"`
	Password   string `json:"password,omitempty" yaml:"password,omitempty"`
}

type BasicAuthCredentials []BasicAuthCredential

type HmacAuthCredential struct {
	ID         string `json:"id,omitempty" yaml:"id,omitempty"`
	ConsumerID string `json:"consumer,omitempty" yaml:"consumer,omitempty"`
	Username   string `json:"username,omitempty" yaml:"username,omitempty"`
	Secret     string `json:"secret,omitempty" yaml:"secret,omitempty"`
}

type HmacAuthCredentials []HmacAuthCredential

type JwtSecret struct {
	ID         string `json:"id,omitempty" yaml:"id,omitempty"`
	ConsumerID string `json:"consumer,omitempty" yaml:"consumer,omitempty"`
	Algorithm  string `json:"algorithm,omitempty" yaml:"algorithm,omitempty"`
	Key        string `json:"key,omitempty" yaml:"key,omitempty"`
	Secret     string `json:"secret,omitempty" yaml:"secret,omitempty"`
}

type JwtSecrets []JwtSecret

// Consumer is the kong consumer definition.
type Consumer struct {
	ID                   string               `json:"id,omitempty" yaml:"id,omitempty"`
	CustomID             string               `json:"custom_id,omitempty" yaml:"custom_id,omitempty"`
	Username             string               `json:"username,omitempty" yaml:"username,omitempty"`
	Tags                 []string             `json:"tags,omitempty" yaml:"tags,omitempty"`
	KeyAuthCredentials   KeyAuthCredentials   `json:"keyauth_credentials,omitempty" yaml:"keyauth_credentials,omitempty"`
	BasicAuthCredentials BasicAuthCredentials `json:"basicauth_credentials,omitempty" yaml:"basicauth_credentials,omitempty"`
	HmacAuthCredentials  HmacAuthCredentials  `json:"hmacauth_credentials,omitempty" yaml:"hmacauth_credentials,omitempty"`
	JwtSecrets           JwtSecrets           `json:"jwt_secrets,omitempty" yaml:"jwt_secrets,omitempty"`
}

type Consumers []Consumer

// Plugin Configuration

// Plugin is the kong plugin definition.
type Plugin struct {
	ID         string        `json:"id,omitempty" yaml:"id,omitempty"`
	Name       string        `json:"name,omitempty" yaml:"name,omitempty"`
	RouteID    string        `json:"route,omitempty" yaml:"route,omitempty"`
	ServiceID  string        `json:"service,omitempty" yaml:"service,omitempty"`
	ConsumerID string        `json:"consumer,omitempty" yaml:"consumer,omitempty"`
	Config     Configuration `json:"config,omitempty" yaml:"config,omitempty"`
	Protocols  []string      `json:"protocols,omitempty" yaml:"protocols,omitempty"`
	Enabled    bool          `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Tags       []string      `json:"tags,omitempty" yaml:"tags,omitempty"`
}

type Plugins []Plugin

// Certificate Configuration

// Certificate is the kong certificate definition.
type Certificate struct {
	ID      string   `json:"id,omitempty" yaml:"id,omitempty"`
	Cert    string   `json:"cert,omitempty" yaml:"cert,omitempty"`
	Key     string   `json:"key,omitempty" yaml:"key,omitempty"`
	CertALT string   `json:"cert_alt,omitempty" yaml:"cert_alt,omitempty"`
	KeyALT  string   `json:"key_alt,omitempty" yaml:"key_alt,omitempty"`
	SNIs    []string `json:"snis,omitempty" yaml:"snis,omitempty"`
	Tags    []string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

type Certificates []Certificate

// SNI Configuration

// SNI is the kong sni definition.
type SNI struct {
	ID          string      `json:"id,omitempty" yaml:"id,omitempty"`
	Name        string      `json:"name,omitempty" yaml:"name,omitempty"`
	Certificate Certificate `json:"certificate,omitempty" yaml:"certificate,omitempty"`
	Tags        []string    `json:"tags,omitempty" yaml:"tags,omitempty"`
}

type SNIs []SNI

// Upstream Configuration

// Upstream is the kong upstream definition.
type Upstream struct {
	ID                 string      `json:"id,omitempty" yaml:"id,omitempty"`
	Name               string      `json:"name,omitempty" yaml:"name,omitempty"`
	Algorithm          string      `json:"algorithm,omitempty" yaml:"algorithm,omitempty"`
	HashOn             string      `json:"hash_on,omitempty" yaml:"hash_on,omitempty"`
	HashFallback       string      `json:"hash_fallback,omitempty" yaml:"hash_fallback,omitempty"`
	HashOnHeader       string      `json:"hash_on_header,omitempty" yaml:"hash_on_header,omitempty"`
	HashFallbackHeader string      `json:"hash_fallback_header,omitempty" yaml:"hash_fallback_header,omitempty"`
	HashOnCookie       string      `json:"hash_on_cookie,omitempty" yaml:"hash_on_cookie,omitempty"`
	HashOnCookiePath   string      `json:"hash_on_cookie_path,omitempty" yaml:"hash_on_cookie_path,omitempty"`
	Slots              int         `json:"slots,omitempty" yaml:"slots,omitempty"`
	HealthChecks       HealthCheck `json:"healthchecks,omitempty" yaml:"healthchecks,omitempty"`
	Tags               []string    `json:"tags,omitempty" yaml:"tags,omitempty"`
	HostHeader         string      `json:"host_header,omitempty" yaml:"host_header,omitempty"`
	ClientCertificate  Certificate `json:"client_certificate,omitempty" yaml:"client_certificate,omitempty"`
	Targets            Targets     `json:"targets,omitempty" yaml:"targets,omitempty"`
}

type HealthCheck struct {
	Active    ActiveHealthCheck  `json:"active,omitempty" yaml:"active,omitempty"`
	Passive   PassiveHealthCheck `json:"passive,omitempty" yaml:"passive,omitempty"`
	Threshold float64            `json:"threshold,omitempty" yaml:"threshold,omitempty"`
}

type ActiveHealthCheck struct {
	Concurrency            int       `json:"concurrency,omitempty" yaml:"concurrency,omitempty"`
	Healthy                Healthy   `json:"healthy,omitempty" yaml:"healthy,omitempty"`
	HTTPPath               string    `json:"http_path,omitempty" yaml:"http_path,omitempty"`
	HTTPSSni               string    `json:"https_sni,omitempty" yaml:"https_sni,omitempty"`
	HTTPSVerifyCertificate bool      `json:"https_verify_certificate,omitempty" yaml:"https_verify_certificate,omitempty"`
	Type                   string    `json:"type,omitempty" yaml:"type,omitempty"`
	Timeout                int       `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Unhealthy              Unhealthy `json:"unhealthy,omitempty" yaml:"unhealthy,omitempty"`
}

type PassiveHealthCheck struct {
	Healthy   Healthy   `json:"healthy,omitempty" yaml:"healthy,omitempty"`
	Type      string    `json:"type,omitempty" yaml:"type,omitempty"`
	Unhealthy Unhealthy `json:"unhealthy,omitempty" yaml:"unhealthy,omitempty"`
}

type Healthy struct {
	HTTPStatuses []int `json:"http_statuses,omitempty" yaml:"http_statuses,omitempty"`
	Interval     int   `json:"interval,omitempty" yaml:"interval,omitempty"`
	Successes    int   `json:"successes,omitempty" yaml:"successes,omitempty"`
}

type Unhealthy struct {
	HTTPFailures int   `json:"http_failures,omitempty" yaml:"http_failures,omitempty"`
	HTTPStatuses []int `json:"http_statuses,omitempty" yaml:"http_statuses,omitempty"`
	TCPFailures  int   `json:"tcp_failures,omitempty" yaml:"tcp_failures,omitempty"`
	Timeouts     int   `json:"timeouts,omitempty" yaml:"timeouts,omitempty"`
	Interval     int   `json:"interval,omitempty" yaml:"interval,omitempty"`
}

type Upstreams []Upstream

type Target struct {
	ID         string   `json:"id,omitempty" yaml:"id,omitempty"`
	Target     string   `json:"target,omitempty" yaml:"target,omitempty"`
	UpstreamID string   `json:"upstream,omitempty" yaml:"upstream,omitempty"`
	Weight     int      `json:"weight,omitempty" yaml:"weight,omitempty"`
	Tags       []string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

type Targets []Target
