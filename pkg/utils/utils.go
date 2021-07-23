package utils

import (
	"bytes"
	"io/ioutil"

	"github.com/icza/dyno"
	"gopkg.in/yaml.v2"
)

type YamlItem struct {
	Value interface{}
	Path  []interface{}
}

var (
	// from kong to apisix
	WordMap = map[string]string{
		"round-robin":        "roundrobin",
		"consistent-hashing": "chash",
	}
)

func AppendToConfigYaml(items *[]YamlItem, filePath string) error {
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	out, err := addValueToYaml(yamlFile, items)
	if err != nil {
		return err
	}

	ioutil.WriteFile(filePath, out, 0644)
	return nil
}

func ShowConfigYaml(items *[]YamlItem) ([]byte, error) {
	out, err := addValueToYaml(nil, items)
	if err != nil {
		return nil, err
	}

	return out, err
}

func addValueToYaml(src []byte, items *[]YamlItem) ([]byte, error) {
	v := make(map[interface{}]interface{})

	if src != nil {
		if err := yaml.Unmarshal(src, &v); err != nil {
			return nil, err
		}
	}

	for _, item := range *items {
		value, path := item.Value, item.Path
		for i := 1; i <= len(path); i++ {
			if i == len(path) {
				if err := dyno.Set(v, value, path...); err != nil {
					return nil, err
				}
			} else {
				if value, err := dyno.Get(v, path[:i]...); err == nil && value != nil {
					// layer already have content
					value.(map[interface{}]interface{})[path[i]] = nil
					if err := dyno.Set(v, value, path[:i]...); err != nil {
						return nil, err
					}
				} else {
					if err := dyno.Set(v, map[interface{}]interface{}{path[i]: nil}, path[:i]...); err != nil {
						return nil, err
					}
				}
			}
		}
	}

	out, err := yaml.Marshal(v)
	if err != nil {
		return nil, err
	}
	// wait till new lua-tinyyaml in newer
	oldEtcdAddr := "http://etcd:2379"
	newEtcdAddr := "\"http://etcd:2379\""
	out = bytes.Replace(out, []byte(oldEtcdAddr), []byte(newEtcdAddr), 1)

	return out, nil
}
