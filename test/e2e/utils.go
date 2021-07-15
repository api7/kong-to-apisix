package e2e

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/api7/kongtoapisix/pkg/apisix"
	"github.com/api7/kongtoapisix/pkg/kong"
	"github.com/globocom/gokong"
	"github.com/onsi/ginkgo"
	"gopkg.in/yaml.v2"
)

var (
	upstreamAddr = "http://172.17.0.1:8088"
	apisixAddr   = "http://127.0.0.1:9080"
	kongAddr     = "http://127.0.0.1:8000"
)

type TestCase struct {
	RouteRequest   *gokong.RouteRequest
	ServiceRequest *gokong.ServiceRequest
}

type CompareCase struct {
	Path             string
	Url              string
	Headers          map[string]string
	ExpectStatusCode int
}

func purgeAll(kongCli gokong.KongAdminClient) error {
	if err := deleteRoute(kongCli); err != nil {
		return err
	}
	if err := deleteService(kongCli); err != nil {
		return err
	}
	if err := deleteUpstream(kongCli); err != nil {
		return err
	}
	if err := deleteConsumer(kongCli); err != nil {
		return err
	}
	if err := deletePlugin(kongCli); err != nil {
		return err
	}
	return nil
}

func deleteRoute(kongCli gokong.KongAdminClient) error {
	kongRoutes, err := kongCli.Routes().List(&gokong.RouteQueryString{})
	if err != nil {
		return err
	}
	for _, r := range kongRoutes {
		if err := kongCli.Routes().DeleteById(*r.Id); err != nil {
			return err
		}
	}

	return nil
}

func deleteService(kongCli gokong.KongAdminClient) error {
	kongServices, err := kongCli.Services().GetServices(&gokong.ServiceQueryString{})
	if err != nil {
		return err
	}
	for _, s := range kongServices {
		if err := kongCli.Services().DeleteServiceById(*s.Id); err != nil {
			return err
		}
	}
	return nil
}

func deleteUpstream(kongCli gokong.KongAdminClient) error {
	kongUpstreams, err := kongCli.Upstreams().List()
	if err != nil {
		return err
	}
	for _, u := range kongUpstreams.Results {
		err := kongCli.Upstreams().DeleteById(u.Id)
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteConsumer(kongCli gokong.KongAdminClient) error {
	kongConsumers, err := kongCli.Consumers().List(&gokong.ConsumerQueryString{})
	if err != nil {
		return err
	}
	for _, c := range kongConsumers {
		err := kongCli.Consumers().DeleteById(c.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func deletePlugin(kongCli gokong.KongAdminClient) error {
	kongPlugins, err := kongCli.Plugins().List(&gokong.PluginQueryString{})
	if err != nil {
		return err
	}
	for _, p := range kongPlugins {
		err := kongCli.Plugins().DeleteById(p.Id)
		if err != nil {
			return err
		}
	}

	return nil
}

func defaultService() *gokong.ServiceRequest {
	return &gokong.ServiceRequest{
		Retries:        gokong.Int(5),
		Protocol:       gokong.String("http"),
		Port:           gokong.Int(80),
		ConnectTimeout: gokong.Int(60000),
		WriteTimeout:   gokong.Int(60000),
		ReadTimeout:    gokong.Int(60000),
	}
}

func defaultRoute() *gokong.RouteRequest {
	return &gokong.RouteRequest{
		Protocols:     gokong.StringSlice([]string{"http", "https"}),
		RegexPriority: gokong.Int(0),
		StripPath:     gokong.Bool(true),
		PreserveHost:  gokong.Bool(true),
	}
}

func getKongConfig() ([]byte, error) {
	tmpStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	if err := kong.DumpKong(""); err != nil {
		return nil, err
	}

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = tmpStdout

	return out, nil
}

func testMigrate() error {
	kongConfigBytes, err := getKongConfig()
	if err != nil {
		return err
	}
	var kongConfig *kong.KongConfig
	err = yaml.Unmarshal(kongConfigBytes, &kongConfig)
	if err != nil {
		return err
	}

	prettier, err := json.MarshalIndent(kongConfig, "", "\t")
	if err != nil {
		return err
	}
	fmt.Fprintf(ginkgo.GinkgoWriter, "kong yaml: %s\n", string(prettier))

	apisixConfig, err := kong.Migrate(kongConfig)
	if err != nil {
		return err
	}

	prettier, err = json.MarshalIndent(apisixConfig, "", "\t")
	if err != nil {
		return err
	}
	fmt.Fprintf(ginkgo.GinkgoWriter, "apisix yaml: %s\n", string(prettier))

	if err := apisix.WriteToFile(apisixConfig); err != nil {
		return err
	}
	// wait one second to make new config works
	time.Sleep(1500 * time.Millisecond)
	return nil
}
