package consumer

import (
	"context"
	"fmt"

	"github.com/apache/apisix-ingress-controller/pkg/apisix"
	v1 "github.com/apache/apisix-ingress-controller/pkg/types/apisix/v1"
	"github.com/globocom/gokong"
)

func MigrateConsumer(apisixCli apisix.Cluster, kongCli gokong.KongAdminClient) error {
	consumers, err := kongCli.Consumers().List(&gokong.ConsumerQueryString{})
	if err != nil {
		return err
	}

	for _, c := range consumers {
		//fmt.Printf("got consumer: %#v\n", c)
		username := c.Username
		if username == "" {
			username = c.CustomId
		}
		apisixConsumer := &v1.Consumer{
			Username: username,
		}

		_, err := apisixCli.Consumer().Create(context.Background(), apisixConsumer)
		if err != nil {
			return err
		}
		fmt.Printf("migrate consumer %s succeeds\n", c.Username)
	}

	return nil
}
