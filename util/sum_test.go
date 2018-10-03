package util

import "testing"

func TestSum(t *testing.T) {
    total := Sum(1, 2)
    if total != 3 {
       t.Errorf("Sum was incorrect, got: %d, want: %d.", total, 3)
    }
}