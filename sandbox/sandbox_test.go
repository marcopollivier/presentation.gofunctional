package sabdbox

import "testing"

func TestSample(t *testing.T) {
	if got := Sample(); got != 1 {
		t.Errorf("Sample() = %d, esperado 1", got)
	}
}
