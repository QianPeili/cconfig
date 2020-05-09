package cconfig

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
)

type ConfigHandler func([]byte) error

type HandlerPair struct {
	key     string
	handler ConfigHandler
	index   uint64
}

var hGuard sync.RWMutex
var handlers = make(map[string]*HandlerPair)

func AddHandler(name string, h ConfigHandler) error {
	hGuard.Lock()
	defer hGuard.Unlock()

	if _, ok := handlers[name]; ok {
		return errors.New(fmt.Sprintf("handler %s duplicated.", name))
	}
	handlers[name] = &HandlerPair{
		key:     getConfigKey(name),
		handler: h,
	}
	return nil
}

func getConfigNameFromKey(key string) string {
	fn := filepath.Base(key)
	ext := filepath.Ext(fn)

	return strings.TrimSuffix(fn, ext)
}

func TriggerAll() error {
	pairs, _, err := kv.List(getConfigPrefix(), nil)
	if err != nil {
		return err
	}
	for _, pair := range pairs {
		name := getConfigNameFromKey(pair.Key)
		hGuard.RLock()
		nh, ok := handlers[name]
		hGuard.RUnlock()
		if !ok {
			continue
		}
		if nh.index == pair.ModifyIndex {
			continue
		}
		if nh.handler(pair.Value) == nil {
			nh.index = pair.ModifyIndex
		}
	}

	return nil
}
