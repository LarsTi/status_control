package main
import(
	"testing"
)
func TestReadEnv(t *testing.T) {
	got := readEnv()
	if len(got) != 0 {
		t.Errorf("len(got) = %d; want 0", len(got))
	}
}
