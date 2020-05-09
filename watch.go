package cconfig

import (
	"fmt"

	"github.com/hashicorp/consul/api"
	consulApi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
)

func startWatchRef() error {
	c := map[string]interface{}{
		"type": "key",
		"key":  getRefKey(),
	}
	plan, err := watch.Parse(c)
	if err != nil {
		return err
	}
	plan.Handler = func(idx uint64, raw interface{}) {
		var v *consulApi.KVPair
		if raw == nil {
			return
		}
		var ok bool
		if v, ok = raw.(*api.KVPair); !ok {
			// todo error handle
			return
		}
		newRef := string(v.Value)
		if reference == newRef {
			return
		}
		reference = newRef
		// check all config
		TriggerAll()
	}
	return plan.RunWithClientAndHclog(client, nil)
}

func getRefKey() string {
	return fmt.Sprintf("%s/%s.ref", config.KeyPath, config.Env)
}
