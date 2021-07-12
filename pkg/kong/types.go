package kong

type KongConfig struct {
	Consumers Consumers `yaml:"consumers"`
	Plugins   Plugins   `yaml:"plugins"`
	Services  Services  `yaml:"services"`
	Upstreams Upstreams `yaml:"upstreams"`
}

type Routes []struct {
	HTTPSRedirectStatusCode int      `yaml:"https_redirect_status_code"`
	Name                    string   `yaml:"name"`
	PathHandling            string   `yaml:"path_handling"`
	Paths                   []string `yaml:"paths"`
	Plugins                 Plugins  `yaml:"plugins"`
	PreserveHost            bool     `yaml:"preserve_host"`
	Protocols               []string `yaml:"protocols"`
	RegexPriority           int      `yaml:"regex_priority"`
	RequestBuffering        bool     `yaml:"request_buffering"`
	ResponseBuffering       bool     `yaml:"response_buffering"`
	StripPath               bool     `yaml:"strip_path"`
}

type Services []struct {
	ConnectTimeout int    `yaml:"connect_timeout"`
	Host           string `yaml:"host"`
	Name           string `yaml:"name"`
	Port           int    `yaml:"port"`
	Protocol       string `yaml:"protocol"`
	ReadTimeout    int    `yaml:"read_timeout"`
	Retries        int    `yaml:"retries"`
	Routes         Routes `yaml:"routes"`
	WriteTimeout   int    `yaml:"write_timeout"`
}

type Upstreams []struct {
	Algorithm        string       `yaml:"algorithm"`
	HashFallback     string       `yaml:"hash_fallback"`
	HashOn           string       `yaml:"hash_on"`
	HashOnCookiePath string       `yaml:"hash_on_cookie_path"`
	Healthchecks     Healthchecks `yaml:"healthchecks"`
	Name             string       `yaml:"name"`
	Slots            int          `yaml:"slots"`
	Targets          Targets      `yaml:"targets"`
}

type Targets []struct {
	Target string `yaml:"target"`
	Weight int    `yaml:"weight"`
}

type Healthchecks struct {
	Active struct {
		Concurrency int `yaml:"concurrency"`
		Healthy     struct {
			HTTPStatuses []int `yaml:"http_statuses"`
			Interval     int   `yaml:"interval"`
			Successes    int   `yaml:"successes"`
		} `yaml:"healthy"`
		HTTPPath               string `yaml:"http_path"`
		HTTPSVerifyCertificate bool   `yaml:"https_verify_certificate"`
		Timeout                int    `yaml:"timeout"`
		Type                   string `yaml:"type"`
		Unhealthy              struct {
			HTTPFailures int   `yaml:"http_failures"`
			HTTPStatuses []int `yaml:"http_statuses"`
			Interval     int   `yaml:"interval"`
			TCPFailures  int   `yaml:"tcp_failures"`
			Timeouts     int   `yaml:"timeouts"`
		} `yaml:"unhealthy"`
	} `yaml:"active"`
	Passive struct {
		Healthy struct {
			HTTPStatuses []int `yaml:"http_statuses"`
			Successes    int   `yaml:"successes"`
		} `yaml:"healthy"`
		Type      string `yaml:"type"`
		Unhealthy struct {
			HTTPFailures int   `yaml:"http_failures"`
			HTTPStatuses []int `yaml:"http_statuses"`
			TCPFailures  int   `yaml:"tcp_failures"`
			Timeouts     int   `yaml:"timeouts"`
		} `yaml:"unhealthy"`
	} `yaml:"passive"`
	Threshold int `yaml:"threshold"`
}

type Consumers []struct {
	CustomId string `json:"custom_id,omitempty" yaml:"custom_id,omitempty"`
	Username string `json:"username,omitempty" yaml:"username,omitempty"`
}

type Plugins []struct {
	Name     string                 `json:"name" yaml:"name"`
	Config   map[string]interface{} `json:"config,omitempty" yaml:"config,omitempty"`
	Enabled  bool                   `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Protocol []string               `json:"protocols,omitempty" yaml:"protocol,omitempty"`
}
