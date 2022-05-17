package randUtils

import "testing"

func TestRandInt(t *testing.T) {
	for i := 0; i < 100; i++ {
		t.Logf("%d", RandInt(0, 10))
	}
}
