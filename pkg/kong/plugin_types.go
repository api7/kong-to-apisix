package kong

type ProxyCache struct {
	CacheControl    bool             `yaml:"cache_control"`
	CacheTTL        int              `yaml:"cache_ttl"`
	ContentType     []string         `yaml:"content_type"`
	Memory          ProxyCacheMemory `yaml:"memory"`
	RequestMethod   []string         `yaml:"request_method"`
	ResponseCode    []int            `yaml:"response_code"`
	StorageTTL      int              `yaml:"storage_ttl"`
	Strategy        string           `yaml:"strategy"`
	VaryHeaders     []string         `yaml:"vary_headers"`
	VaryQueryParams []string         `yaml:"vary_query_params"`
}

type ProxyCacheMemory struct {
	DictionaryName string `yaml:"dictionary_name"`
}

type RateLimiting struct {
	Day               int    `yaml:"day"`
	FaultTolerant     bool   `yaml:"fault_tolerant"`
	HeaderName        string `yaml:"header_name"`
	HideClientHeaders bool   `yaml:"hide_client_headers"`
	Hour              int    `yaml:"hour"`
	LimitBy           string `yaml:"limit_by"`
	Minute            int    `yaml:"minute"`
	Month             int    `yaml:"month"`
	Path              string `yaml:"path"`
	Policy            string `yaml:"policy"`
	RedisDatabase     int    `yaml:"redis_database"`
	RedisHost         string `yaml:"redis_host"`
	RedisPassword     string `yaml:"redis_password"`
	RedisPort         int    `yaml:"redis_port"`
	RedisTimeout      int    `yaml:"redis_timeout"`
	Second            int    `yaml:"second"`
	Year              int    `yaml:"year"`
}

type KeyAuth struct {
	Anonymous       string   `yaml:"anonymous"`
	HideCredentials bool     `yaml:"hide_credentials"`
	KeyInBody       bool     `yaml:"key_in_body"`
	KeyInHeader     bool     `yaml:"key_in_header"`
	KeyInQuery      bool     `yaml:"key_in_query"`
	KeyNames        []string `yaml:"key_names"`
	RunOnPreflight  bool     `yaml:"run_on_preflight"`
}
