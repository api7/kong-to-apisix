package kong

import (
	"testing"

	"github.com/api7/kong-to-apisix/pkg/apisix"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

type TestUpstreamExpect struct {
	Id        string
	Name      string
	Algorithm string
	Host      string
	Port      int
	Weight    int
	Retries   uint
	Protocol  string
}

type TestUpstream struct {
	Upstream Upstream
	Expect   TestUpstreamExpect
}

type TestFormatTarget struct {
	Hostname string
	Expect   A6UpstreamNodeHostPort
}

func TestMigrateUpstream(t *testing.T) {
	var kongConfig Config
	var apisixConfig apisix.Config

	serviceId := uuid.NewV4().String()
	upstreamId1 := uuid.NewV4().String()
	upstreamName1 := "up1"
	upstreamId2 := uuid.NewV4().String()
	upstreamName2 := "up2"

	testUpstreams := []TestUpstream{
		{
			Upstream: Upstream{
				ID:        upstreamId1,
				Name:      upstreamName1,
				Algorithm: KTAKongUpstreamAlgorithmRoundRobin,
				Targets: Targets{
					{
						Target: "127.0.0.1:1980",
						Weight: 100,
					},
				},
			},
			Expect: TestUpstreamExpect{
				Id:        upstreamId1,
				Name:      upstreamName1,
				Algorithm: KTASixUpstreamAlgorithmRoundRobin,
				Host:      "127.0.0.1",
				Port:      1980,
				Weight:    100,
			},
		},
		{
			Upstream: Upstream{
				ID:        upstreamId2,
				Name:      upstreamName2,
				Algorithm: KTAKongUpstreamAlgorithmConsistentHashing,
				Targets: Targets{
					{
						Target: "localhost",
						Weight: 100,
					},
				},
			},
			Expect: TestUpstreamExpect{
				Id:        upstreamId2,
				Name:      upstreamName2,
				Algorithm: KTASixUpstreamAlgorithmConsistentHashing,
				Host:      "localhost",
				Port:      80,
				Weight:    100,
			},
		},
	}

	for _, testUpstream := range testUpstreams {
		kongConfig.Upstreams = Upstreams{}
		kongConfig.Services = Services{}
		apisixConfig.Upstreams = apisix.Upstreams{}

		var kongConfigService Service
		kongConfigService.ID = serviceId
		kongConfigService.Name = "svc"
		kongConfigService.Retries = 10
		kongConfigService.Protocol = "http"
		kongConfigService.Host = testUpstream.Upstream.Name
		kongConfig.Services = append(kongConfig.Services, kongConfigService)
		err := MigrateService(&kongConfig, &apisixConfig)
		assert.NoError(t, err)
		kongConfig.Upstreams = append(kongConfig.Upstreams, testUpstream.Upstream)
		err = MigrateUpstream(&kongConfig, &apisixConfig)
		assert.NoError(t, err)
		assert.Equal(t, apisixConfig.Upstreams[1].ID, testUpstream.Expect.Id)
		assert.Equal(t, apisixConfig.Upstreams[1].Name, testUpstream.Expect.Name)
		assert.Equal(t, apisixConfig.Upstreams[1].Type, testUpstream.Expect.Algorithm)
		assert.Equal(t, apisixConfig.Upstreams[1].Retries, kongConfigService.Retries)
		assert.Equal(t, apisixConfig.Upstreams[1].Scheme, kongConfigService.Protocol)
		assert.Equal(t, apisixConfig.Upstreams[1].Nodes[0].Host, testUpstream.Expect.Host)
		assert.Equal(t, apisixConfig.Upstreams[1].Nodes[0].Port, testUpstream.Expect.Port)
		assert.Equal(t, apisixConfig.Upstreams[1].Nodes[0].Weight, testUpstream.Expect.Weight)
	}
}

func TestKTAFormatKongTarget(t *testing.T) {
	testTargets := []TestFormatTarget{
		{
			Hostname: "valid.name",
			Expect: A6UpstreamNodeHostPort{
				Host: "valid.name",
				Port: 80,
			},
		},
		{
			Hostname: "valid.name:8080",
			Expect: A6UpstreamNodeHostPort{
				Host: "valid.name",
				Port: 8080,
			},
		},
		{
			Hostname: "12.34.56.78",
			Expect: A6UpstreamNodeHostPort{
				Host: "12.34.56.78",
				Port: 80,
			},
		},
		{
			Hostname: "1.2.3.4:123",
			Expect: A6UpstreamNodeHostPort{
				Host: "1.2.3.4",
				Port: 123,
			},
		},
		{
			Hostname: "[A01F::0]:8000",
			Expect: A6UpstreamNodeHostPort{
				Host: "A01F::0",
				Port: 8000,
			},
		},
		{
			Hostname: "[A01F::0]",
			Expect: A6UpstreamNodeHostPort{
				Host: "A01F::0",
				Port: 80,
			},
		},
	}

	for _, testTarget := range testTargets {
		a6NodeHostPort, err := KTAFormatKongTarget(testTarget.Hostname)
		assert.NoError(t, err)
		assert.Equal(t, testTarget.Expect.Host, a6NodeHostPort.Host)
		assert.Equal(t, testTarget.Expect.Port, a6NodeHostPort.Port)
	}
}
