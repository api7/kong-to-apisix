package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"

	"github.com/api7/kong-to-apisix/pkg/apisix"
	"github.com/api7/kong-to-apisix/pkg/kong"
	"github.com/api7/kong-to-apisix/pkg/utils"

	"github.com/globocom/gokong"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
)

var (
	UpstreamAddr         = "http://172.17.0.1:7024"
	UpstreamAddr2        = "http://172.17.0.1:7025"
	ApisixConfigYamlPath = "../../repos/apisix-docker/example/apisix_conf/config.yaml"

	apisixAddr            = "http://127.0.0.1:9080"
	kongAddr              = "http://127.0.0.1:8000"
	kongAdminAddr         = "http://127.0.0.1:8001"
	apisixDeclYamlPath    = "../../repos/apisix-docker/example/apisix_conf/apisix.yaml"
	kongDeclYamlPath      = "../../repos/kong-docker/compose/kong_conf/kong.yml"
	kongContainerName     = "/kong"
	kongDumpConfigCommand = []string{"/bin/sh", "-c", "rm -f kong.yml && kong config db_export"}
	TestKongDeckMode      = "kong/deck"
	TestKongConfigMode    = "kong/config"
)

type TestCase struct {
	RouteRequest   *gokong.RouteRequest
	ServiceRequest *gokong.ServiceRequest
}

type CompareCase struct {
	Path              string
	Url               string
	Headers           map[string]string
	CompareBody       bool
	CompareStatusCode int
}

func GetKongDeckExportData() ([]byte, error) {
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

func GetKongConfigExportData() ([]byte, error) {
	err := generateKongConfig()
	if err != nil {
		return nil, err
	}
	out, err := ioutil.ReadFile(kongDeclYamlPath)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func generateKongConfig() error {
	var kongContainerID string
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		return err
	}

	cli.NegotiateAPIVersion(ctx)

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return err
	}

	for _, container := range containers {
		if container.Names[0] == kongContainerName {
			kongContainerID = container.ID
		}
	}

	config := types.ExecConfig{
		AttachStdin:  true,
		AttachStderr: true,
		AttachStdout: true,
		Tty:          true,
		Cmd:          kongDumpConfigCommand,
	}

	exec, err := cli.ContainerExecCreate(ctx, kongContainerID, config)
	if err != nil {
		return err
	}

	resp, err := cli.ContainerExecAttach(ctx, exec.ID, types.ExecStartCheck{})
	if err != nil {
		return err
	}

	defer func() {
		err = resp.Conn.Close()
		if err != nil {
			log.Panic(err)
		}
		err = cli.Close()
		if err != nil {
			log.Panic(err)
		}
	}()

	buf := make([]byte, 1024)
	var n int
	n, err = resp.Reader.Read(buf)
	if err != nil && err != io.EOF {
		return err
	}
	if len(buf[:n]) > 0 {
		var sb strings.Builder
		sb.Write(buf[:n])
		return errors.New(sb.String())
	}
	return nil
}

func TestMigrate(mode string) error {
	var err error
	var kongConfigBytes []byte
	if mode == TestKongConfigMode {
		kongConfigBytes, err = GetKongConfigExportData()
	} else {
		kongConfigBytes, err = GetKongDeckExportData()
	}
	if err != nil {
		return err
	}
	var kongConfig *kong.Config
	err = yaml.Unmarshal(kongConfigBytes, &kongConfig)
	if err != nil {
		return err
	}

	prettier, err := json.MarshalIndent(*kongConfig, "", "\t")
	if err == nil {
		fmt.Fprintf(ginkgo.GinkgoWriter, "kong yaml: %s\n", string(prettier))
	}

	apisixDecl, apisixConfig, err := kong.Migrate(kongConfig)
	if err != nil {
		return err
	}

	prettier, err = json.MarshalIndent(*apisixDecl, "", "\t")
	if err == nil {
		fmt.Fprintf(ginkgo.GinkgoWriter, "apisix yaml: %s\n", string(prettier))
	}

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

func GetResps(c *CompareCase) (*http.Response, *http.Response) {
	c.Url = kongAddr + c.Path
	kongResp, err := getResp(c)
	gomega.Expect(err).To(gomega.BeNil())
	kongResp.Body.Close()

	c.Url = apisixAddr + c.Path
	apisixResp, err := getResp(c)
	gomega.Expect(err).To(gomega.BeNil())
	apisixResp.Body.Close()

	return apisixResp, kongResp
}

func GetBodys(c *CompareCase) (string, string) {
	c.Url = kongAddr + c.Path
	kongResp, err := getResp(c)
	gomega.Expect(err).To(gomega.BeNil())
	defer kongResp.Body.Close()

	c.Url = apisixAddr + c.Path
	apisixResp, err := getResp(c)
	gomega.Expect(err).To(gomega.BeNil())
	defer apisixResp.Body.Close()

	kongBody, err := getBody(kongResp)
	gomega.Expect(err).To(gomega.BeNil())
	apisixBody, err := getBody(apisixResp)
	gomega.Expect(err).To(gomega.BeNil())

	return kongBody, apisixBody
}

// do compare here
func Compare(c *CompareCase) {
	c.Url = kongAddr + c.Path
	kongResp, err := getResp(c)
	gomega.Expect(err).To(gomega.BeNil())
	defer kongResp.Body.Close()

	c.Url = apisixAddr + c.Path
	apisixResp, err := getResp(c)
	gomega.Expect(err).To(gomega.BeNil())
	defer apisixResp.Body.Close()

	if c.CompareStatusCode != 0 {
		gomega.Ω(kongResp.StatusCode).Should(gomega.Equal(c.CompareStatusCode))
		gomega.Ω(apisixResp.StatusCode).Should(gomega.Equal(c.CompareStatusCode))
	}

	if c.CompareBody {
		kongBody, err := getBody(kongResp)
		gomega.Expect(err).To(gomega.BeNil())
		apisixBody, err := getBody(apisixResp)
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Ω(apisixBody).Should(gomega.Equal(kongBody))
	}
}

func getResp(c *CompareCase) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.Url, nil)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return resp, nil
}

func getBody(resp *http.Response) (string, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
