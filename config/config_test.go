package config

import (
	"fmt"
	"log"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func init() {
	viper.SetConfigName("config")             // name of config file (without extension)
	viper.SetConfigType("yaml")               // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/conflux-pay/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.conflux-pay") // call multiple times to add many search paths
	viper.AddConfigPath(".")                  // optionally look for config in the working directory
	viper.AddConfigPath("..")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatalln(fmt.Errorf("fatal error config file: %w", err))
	}
}

func TestReadCompany(t *testing.T) {
	var v Company
	err := viper.UnmarshalKey("company", &v)
	assert.NoError(t, err)

	fmt.Printf("%+v", v)
}
