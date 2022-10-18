package config

import (
	"fmt"
	"testing"
)

func TestReadCompany(t *testing.T) {
	c := getCompany()
	fmt.Println(c.MchPrivateKey)
}
