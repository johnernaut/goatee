package goatee

import (
	"testing"
)

func TestPubsubHub(t *testing.T) {
	client := setup(t)
	defer client.Close()
}
