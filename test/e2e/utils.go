package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/apache/apisix-ingress-controller/pkg/apisix"
	"github.com/globocom/gokong"
	"github.com/icza/dyno"
	. "github.com/onsi/ginkgo"
	"github.com/pkg/errors"
)

var (
	upstreamAddr  = "http://172.17.0.1:8088"
	upstreamAddr2 = "http://172.17.0.1:8001"
	apisixAddr    = "http://127.0.0.1:9080"
	kongAddr      = "http://127.0.0.1:8000"
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

func purgeAll(apisixCli apisix.Cluster, kongCli gokong.KongAdminClient) error {
	if err := deleteRoute(apisixCli, kongCli); err != nil {
		return err
	}
	if err := deleteService(apisixCli, kongCli); err != nil {
		return err
	}
	if err := deleteUpstream(apisixCli, kongCli); err != nil {
		return err
	}
	if err := deleteConsumer(apisixCli, kongCli); err != nil {
		return err
	}
	if err := deletePlugin(apisixCli, kongCli); err != nil {
		return err
	}
	return nil
}

func deleteRoute(apisixCli apisix.Cluster, kongCli gokong.KongAdminClient) error {
	ctx := context.Background()
	kongRoutes, err := kongCli.Routes().List(&gokong.RouteQueryString{})
	if err != nil {
		return err
	}
	for _, r := range kongRoutes {
		if err := kongCli.Routes().DeleteById(*r.Id); err != nil {
			return err
		}
	}

	apisixRoutes, err := apisixCli.Route().List(ctx)
	if err != nil {
		return err
	}
	for _, r := range apisixRoutes {
		if err := apisixCli.Route().Delete(ctx, r); err != nil {
			return err
		}
	}
	return nil
}

func deleteService(apisixCli apisix.Cluster, kongCli gokong.KongAdminClient) error {
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

func deleteUpstream(apisixCli apisix.Cluster, kongCli gokong.KongAdminClient) error {
	ctx := context.Background()
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

	apisixUpstreams, err := apisixCli.Upstream().List(ctx)
	if err != nil {
		return err
	}
	for _, u := range apisixUpstreams {
		if err := apisixCli.Upstream().Delete(ctx, u); err != nil {
			return err
		}
	}
	return nil
}

func deleteConsumer(apisixCli apisix.Cluster, kongCli gokong.KongAdminClient) error {
	ctx := context.Background()
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

	apisixConsumers, err := apisixCli.Consumer().List(ctx)
	if err != nil {
		return err
	}
	for _, c := range apisixConsumers {
		if err := apisixCli.Consumer().Delete(ctx, c); err != nil {
			return err
		}
	}
	return nil
}

func deletePlugin(apisixCli apisix.Cluster, kongCli gokong.KongAdminClient) error {
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

func defaultUpstream() *gokong.UpstreamRequest {
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

func defaultTarget() *gokong.TargetRequest {
	return &gokong.TargetRequest{
		Weight: 100,
	}
}

func compareURL(c *CompareCase) (bool, error) {
	kongResp, apisixResp, err := getResp(c)
	if err != nil {
		return false, err
	}

	kongURL, err := getFromJson(kongResp, "url")
	if err != nil {
		return false, err
	}
	apisixURL, err := getFromJson(apisixResp, "url")
	if err != nil {
		return false, err
	}

	return kongURL == apisixURL, nil
}

func getResp(c *CompareCase) ([]byte, []byte, error) {
	c.Url = apisixAddr + c.Path
	apisixResp, err := getBody(c)
	if err != nil {
		return nil, nil, errors.Wrap(err, "apisix")
	}

	c.Url = kongAddr + c.Path
	kongResp, err := getBody(c)
	if err != nil {
		return nil, nil, errors.Wrap(err, "kong")
	}

	GinkgoT().Logf("Kong: %s, APISIX: %s", kongResp, apisixResp)
	return kongResp, apisixResp, nil
}

func getBody(c *CompareCase) ([]byte, error) {
	req, err := http.NewRequest("GET", c.Url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "http new request error")
	}
	for k, v := range c.Headers {
		if k == "Host" {
			req.Host = c.Headers["Host"]
		} else {
			req.Header.Set(k, v)
		}
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "http get error")
	}
	defer resp.Body.Close()

	if resp.StatusCode == c.ExpectStatusCode {
		return nil, nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read body error")
	}

	if resp.StatusCode != http.StatusOK {
		return []byte(fmt.Sprintf("%d ", resp.StatusCode)), errors.Errorf("read body error: %s", string(body))
	}

	return body, err
}

func getFromJson(body []byte, key string) (string, error) {
	v := make(map[string]interface{})
	err := json.Unmarshal(body, &v)
	if err != nil {
		return "", errors.Wrap(err, "unmarshal error")
	}

	value, err := dyno.Get(v, key)
	if err != nil {
		return "", errors.Wrap(err, "get url from interface error")
	}
	return value.(string), nil
}
