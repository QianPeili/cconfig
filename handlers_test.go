package cconfig

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GetNameFromKey(t *testing.T) {
	key := "games/scmj/beta/a.json"
	expt := "a"
	actual := getConfigNameFromKey(key)
	require.Equal(t, expt, actual)
}