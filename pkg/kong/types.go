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
	Hosts                   []string `yaml:"hosts"`
	Methods                 []string `yaml:"methods"`
	PathHandling            string   `yaml:"path_handling"`
	Paths                   []string `yaml:"paths"`
	Plugins                 Plugins  `yaml:"plugins,omitempty"`
	PreserveHost            bool     `yaml:"preserve_host"`
	Protocols               []string `yaml:"protocols"`
	RegexPriority           int      `yaml:"regex_priority"`
	RequestBuffering        bool     `yaml:"request_buffering"`
	ResponseBuffering       bool     `yaml:"response_buffering"`
	StripPath               bool     `yaml:"strip_path"`
}

type Services []struct {
	ID             string  `yaml:"id,omitempty"`
	ConnectTimeout int     `yaml:"connect_timeout"`
	Host           string  `yaml:"host"`
	Name           string  `yaml:"name"`
	Port           int     `yaml:"port"`
	Protocol       string  `yaml:"protocol"`
	ReadTimeout    int     `yaml:"read_timeout"`
	Retries        int     `yaml:"retries"`
	Routes         Routes  `yaml:"routes"`
	WriteTimeout   int     `yaml:"write_timeout"`
	Plugins        Plugins `yaml:"Plugins,omitempty"`
}

type Upstreams []Upstream

type Upstream struct {
	Algorithm        string       `yaml:"algorithm"`
	HashFallback     string       `yaml:"hash_fallback"`
	HashOn           string       `yaml:"hash_on"`
	HashOnCookiePath string       `yaml:"hash_on_cookie_path"`
	Healthchecks     Healthchecks `yaml:"healthchecks"`
	Name             string       `yaml:"name"`
	Slots            int          `yaml:"slots"`
	Targets          Targets      `yaml:"targets"`
	Plugins          Plugins      `yaml:"Plugins,omitempty"`
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
	CustomId           string `yaml:"custom_id,omitempty"`
	Username           string `yaml:"username,omitempty"`
	KeyAuthCredentials []struct {
		Key string `yaml:"key"`
	} `yaml:"keyauth_credentials,omitempty"`
}

type Plugins []Plugin

type Plugin struct {
	Name     string                 `yaml:"name"`
	Config   map[string]interface{} `yaml:"config,omitempty"`
	Enabled  bool                   `yaml:"enabled,omitempty"`
	Protocol []string               `yaml:"protocol,omitempty"`
}
