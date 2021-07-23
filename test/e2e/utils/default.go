package utils

import "github.com/globocom/gokong"

func DefaultService() *gokong.ServiceRequest {
	return &gokong.ServiceRequest{
		Retries:        gokong.Int(5),
		Protocol:       gokong.String("http"),
		Port:           gokong.Int(80),
		ConnectTimeout: gokong.Int(60000),
		WriteTimeout:   gokong.Int(60000),
		ReadTimeout:    gokong.Int(60000),
	}
}

func DefaultRoute() *gokong.RouteRequest {
	return &gokong.RouteRequest{
		Protocols:     gokong.StringSlice([]string{"http", "https"}),
		RegexPriority: gokong.Int(0),
		StripPath:     gokong.Bool(true),
		PreserveHost:  gokong.Bool(true),
	}
}

func DefaultUpstream() *gokong.UpstreamRequest {
	return &gokong.UpstreamRequest{
		HashOn:           "none",
		HashFallback:     "none",
		HashOnCookiePath: "/",
		Slots:            10000,
		HealthChecks: &gokong.UpstreamHealthCheck{
			Active: &gokong.UpstreamHealthCheckActive{
				Type:        "http",
				Concurrency: 10,
				Unhealthy: &gokong.ActiveUnhealthy{
					Interval:     0,
					TcpFailures:  0,
					HttpStatuses: []int{429, 404, 500, 501, 502, 503, 504, 505},
					HttpFailures: 0,
					Timeouts:     0,
				},
				Healthy: &gokong.ActiveHealthy{
					HttpStatuses: []int{200, 302},
					Interval:     0,
					Successes:    0,
				},
			},
		},
	}
}

func DefaultTarget() *gokong.TargetRequest {
	return &gokong.TargetRequest{
		Weight: 100,
	}
}
