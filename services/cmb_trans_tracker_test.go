package services

import (
	"testing"

	"github.com/web3-identity/conflux-pay/models"
)

// This is a function only for test environment and won't stop
func TestTracker(t *testing.T) {
	models.ConnectDB()
	StartCmbTransTracker()
}
