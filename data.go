package cconfig

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

func GetConfigDataByName(name string) (*api.KVPair, error) {
	key := getConfigKey(name)
	pair, _, err := kv.Get(key, nil)
	return pair, err
}

func getConfigKey(name string) string {
	return fmt.Sprintf("%s/%s/%s", config.KeyPath, config.Env, name)
}

func getConfigPrefix() string {
	return fmt.Sprintf("%s/%s", config.KeyPath, config.Env)
}
