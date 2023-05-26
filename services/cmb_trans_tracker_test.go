package services

import (
	"testing"

	"github.com/web3-identity/conflux-pay/models"
)

func TestTracker(t *testing.T) {
	models.ConnectDB()
	StartCmbTransTracker()
}
