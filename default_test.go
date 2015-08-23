package statsd

import (
	"testing"
)

func TestDefaultClient(t *testing.T) {
	if DefaultClient == nil {
		t.Errorf("default client not set")
	}
}
