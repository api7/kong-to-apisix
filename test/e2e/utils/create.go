package utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	uuid "github.com/satori/go.uuid"

	"github.com/api7/kong-to-apisix/pkg/kong"
)

func CreateRoute(route kong.Route) (bool, error) {

	if len(route.Name) == 0 {
		route.Name = uuid.NewV4().String()
	}

	if len(route.Paths) == 0 {
		route.Paths = []string{"/"}
	}

	if len(route.Methods) == 0 {
		route.Methods = []string{"GET", "POST"}
	}

	buf, err := json.Marshal(route)
	if err != nil {
		return false, err
	}

	if len(route.ServiceID) > 0 {
		var result map[string]interface{}
		err = json.Unmarshal(buf, &result)
		if err != nil {
			return false, err
		}
		result["service"] = map[string]string{"id": route.ServiceID}
		buf, err = json.Marshal(result)
		if err != nil {
			return false, err
		}
	}

	kongAdminRoute := KongAdminAddress + KongAdminRouteURI
	resp, err := http.Post(kongAdminRoute, "application/json", strings.NewReader(string(buf)))
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	if resp.StatusCode > 299 {
		return false, errors.New(string(body))
	}

	return true, nil
}
