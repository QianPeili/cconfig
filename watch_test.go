package cconfig

import (
	"testing"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/require"
)

func Test_StartWatch(t *testing.T) {
	err := InitKV(Config{
		Addr:    "192.168.88.236:8500",
		Env:     "beta",
		KeyPath: "games/scmj",
	})
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	_, err = kv.Put(&api.KVPair{
		Key:   getConfigKey("a"),
		Value: []byte("a1"),
	}, nil)
	require.Nil(t, err)
	_, err = kv.Put(&api.KVPair{
		Key:   getConfigKey("b"),
		Value: []byte("b1"),
	}, nil)
	require.Nil(t, err)
	_, err = kv.Put(&api.KVPair{
		Key:   getRefKey(),
		Value: []byte("ref1"),
	}, nil)
	require.Nil(t, err)

	aReceiver := make([]string, 0)
	bReceiver := make([]string, 0)
	AddHandler("a", func(data []byte) error {
		aReceiver = append(aReceiver, string(data))
		return nil
	})
	AddHandler("b", func(data []byte) error {
		bReceiver = append(bReceiver, string(data))
		return nil
	})

	go func() {
		err = startWatchRef()
		if err != nil {
			t.Error(err)
			t.Fail()
		}
	}()
	time.Sleep(1 * time.Second)
	kv.Put(&api.KVPair{
		Key:   getConfigKey("a"),
		Value: []byte("a2"),
	}, nil)
	kv.Put(&api.KVPair{
		Key:   getRefKey(),
		Value: []byte("ref2"),
	}, nil)

	time.Sleep(1 * time.Second)

	require.Len(t, aReceiver, 2)
	require.Equal(t, "a2", aReceiver[1])
	require.Len(t, bReceiver, 1)
}
