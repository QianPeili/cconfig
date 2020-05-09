package cconfig

import (
	consulApi "github.com/hashicorp/consul/api"
)

var client *consulApi.Client
var kv *consulApi.KV
var reference string

type Config struct {
	Addr    string
	Token   string
	KeyPath string
	Env     string
}

var config Config

func InitKV(c Config) error {
	var err error
	config = c
	client, err = consulApi.NewClient(&consulApi.Config{
		Address: config.Addr,
		Token:   config.Token,
	})
	if err != nil {
		return err
	}
	kv = client.KV()
	return nil
}

func Start() {
	go startWatchRef()
}
