package utils

import "github.com/globocom/gokong"

func PurgeAll(kongCli gokong.KongAdminClient) error {
	if err := DeleteRoute(kongCli); err != nil {
		return err
	}
	if err := DeleteService(kongCli); err != nil {
		return err
	}
	if err := DeleteUpstream(kongCli); err != nil {
		return err
	}
	if err := DeleteConsumer(kongCli); err != nil {
		return err
	}
	if err := DeletePlugin(kongCli); err != nil {
		return err
	}
	return nil
}

func DeleteRoute(kongCli gokong.KongAdminClient) error {
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

func DeleteService(kongCli gokong.KongAdminClient) error {
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

func DeleteUpstream(kongCli gokong.KongAdminClient) error {
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

func DeleteConsumer(kongCli gokong.KongAdminClient) error {
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

func DeletePlugin(kongCli gokong.KongAdminClient) error {
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
