package e2e

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/globocom/gokong"
	"github.com/icza/dyno"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

var _ = Describe("route", func() {
	var kongCli gokong.KongAdminClient

	JustBeforeEach(func() {
		kongCli = gokong.NewClient(gokong.NewDefaultConfig())
		err := purgeAll(kongCli)
		Expect(err).To(BeNil())
	})

	It("default route and service", func() {
		createdService := defaultService()
		createdService.Url = gokong.String(upstreamAddr)
		kongService, err := kongCli.Services().Create(createdService)
		Expect(err).To(BeNil())

		createdRoute := defaultRoute()
		createdRoute.Paths = gokong.StringSlice([]string{"/get"})
		createdRoute.Service = gokong.ToId(*kongService.Id)
		_, err = kongCli.Routes().Create(createdRoute)
		Expect(err).To(BeNil())

		err = testMigrate()
		Expect(err).To(BeNil())

		same, err := compareResp(&CompareCase{
			Path: "/get/get",
		})
		Expect(err).To(BeNil())
		Expect(same).To(BeTrue())
	})

	It("kong route disable strip_path", func() {
		createdService := defaultService()
		createdService.Url = gokong.String(upstreamAddr)
		kongService, err := kongCli.Services().Create(createdService)
		Expect(err).To(BeNil())

		createdRoute := defaultRoute()
		createdRoute.Paths = gokong.StringSlice([]string{"/get"})
		createdRoute.Service = gokong.ToId(*kongService.Id)
		createdRoute.StripPath = gokong.Bool(false)
		_, err = kongCli.Routes().Create(createdRoute)
		Expect(err).To(BeNil())

		err = testMigrate()
		Expect(err).To(BeNil())

		same, err := compareResp(&CompareCase{
			Path: "/get",
		})
		Expect(err).To(BeNil())
		Expect(same).To(BeTrue())
	})

	It("kong route with host", func() {
		createdService := defaultService()
		createdService.Url = gokong.String(upstreamAddr)
		kongService, err := kongCli.Services().Create(createdService)
		Expect(err).To(BeNil())

		createdRoute := defaultRoute()
		createdRoute.Paths = gokong.StringSlice([]string{"/get"})
		createdRoute.Service = gokong.ToId(*kongService.Id)
		createdRoute.Hosts = gokong.StringSlice([]string{"foo.com"})
		_, err = kongCli.Routes().Create(createdRoute)
		Expect(err).To(BeNil())

		err = testMigrate()
		Expect(err).To(BeNil())

		same, err := compareResp(&CompareCase{
			Path:             "/get/get",
			ExpectStatusCode: http.StatusNotFound,
		})
		Expect(err).To(BeNil())
		Expect(same).To(BeTrue())

		same, err = compareResp(&CompareCase{
			Path:    "/get/get",
			Headers: map[string]string{"Host": "foo.com"},
		})
		Expect(err).To(BeNil())
		Expect(same).To(BeTrue())
	})

})

func compareResp(c *CompareCase) (bool, error) {
	c.Url = apisixAddr + c.Path
	apisixResp, err := getBody(c)
	if err != nil {
		return false, errors.Wrap(err, "apisix")
	}

	c.Url = kongAddr + c.Path
	kongResp, err := getBody(c)
	if err != nil {
		return false, errors.Wrap(err, "kong")
	}

	GinkgoT().Logf("Kong: %s, APISIX: %s", kongResp, apisixResp)
	return kongResp == apisixResp, nil
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

	v := make(map[string]interface{})
	err = json.Unmarshal(body, &v)
	if err != nil {
		return "", errors.Wrapf(err, "unmarshal error: %s", string(body))
	}

	value, err := dyno.Get(v, "url")
	if err != nil {
		return "", errors.Wrap(err, "get url from interface error")
	}
	return value.(string), nil
}
