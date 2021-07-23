package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/globocom/gokong"
	"github.com/onsi/ginkgo"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/api7/kong-to-apisix/pkg/apisix"
	"github.com/api7/kong-to-apisix/pkg/kong"
	"github.com/api7/kong-to-apisix/pkg/utils"
)

var (
	UpstreamAddr         = "http://172.17.0.1:7024"
	UpstreamAddr2        = "http://172.17.0.1:7025"
	ApisixConfigYamlPath = "../../repos/apisix-docker/example/apisix_conf/config.yaml"

	apisixAddr         = "http://127.0.0.1:9080"
	kongAddr           = "http://127.0.0.1:8000"
	kongAdminAddr      = "http://127.0.0.1:8001"
	apisixDeclYamlPath = "../../repos/apisix-docker/example/apisix_conf/apisix.yaml"
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

func GetKongConfig() ([]byte, error) {
	tmpStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	if err := kong.DumpKong(kongAdminAddr, ""); err != nil {
		return nil, err
	}

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = tmpStdout

	return out, nil
}

func TestMigrate() error {
	kongConfigBytes, err := GetKongConfig()
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

	apisixDecl, apisixConfig, err := kong.Migrate(kongConfig)
	if err != nil {
		return err
	}

	prettier, err = json.MarshalIndent(apisixDecl, "", "\t")
	if err != nil {
		return err
	}
	fmt.Fprintf(ginkgo.GinkgoWriter, "apisix yaml: %s\n", string(prettier))

	apisixYaml, err := apisix.MarshalYaml(apisixDecl)
	if err != nil {
		return err
	}

	if err := apisix.WriteToFile(apisixDeclYamlPath, apisixYaml); err != nil {
		return err
	}

	if err := utils.AppendToConfigYaml(apisixConfig, ApisixConfigYamlPath); err != nil {
		return err
	}

	// wait one second to make new config works
	time.Sleep(1500 * time.Millisecond)
	return nil
}

func GetResp(c *CompareCase) (string, string, error) {
	c.Url = kongAddr + c.Path
	kongResp, err := getBody(c)
	if err != nil {
		return "", "", errors.Wrap(err, "kong")
	}
	c.Url = apisixAddr + c.Path
	apisixResp, err := getBody(c)
	if err != nil {
		return "", "", errors.Wrap(err, "apisix")
	}

	ginkgo.GinkgoT().Logf("Kong: %s, APISIX: %s", kongResp, apisixResp)
	return apisixResp, kongResp, nil
}

func getBody(c *CompareCase) (string, error) {
	req, err := http.NewRequest("GET", c.Url, nil)
	if err != nil {
		return "", errors.Wrap(err, "http new request error")
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
		return "", errors.Wrap(err, "http get error")
	}
	defer resp.Body.Close()

	if resp.StatusCode == c.ExpectStatusCode {
		return "", nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "read body error")
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("%d ", resp.StatusCode), errors.Errorf("read body error: %s", string(body))
	}

	return string(body), nil
}
