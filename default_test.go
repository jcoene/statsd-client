package statsd

import (
	"testing"
)

func TestDefaultClient(t *testing.T) {
	if defaultClient == nil {
		t.Errorf("default client not set")
	}
}
