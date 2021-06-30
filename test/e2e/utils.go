package e2e

import (
	"context"

	"github.com/apache/apisix-ingress-controller/pkg/apisix"
	"github.com/globocom/gokong"
)

var (
	upstreamAddr    = "http://host.docker.internal:8088"
	apisixAddr      = "http://127.0.0.1:9080"
	kongAddr        = "http://127.0.0.1:8000"
	upstreamPathSet = "/get"
	upstreamPathGet = "/get/get"
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
