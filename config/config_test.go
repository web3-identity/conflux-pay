package config

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestReadCompany(t *testing.T) {
	c := getCompany()
	fmt.Println(c)
}

func TestGetCns(t *testing.T) {
	var apps map[string]App
	err := viper.UnmarshalKey("apps", &apps)
	assert.NoError(t, err)
	fmt.Println(apps)
}
