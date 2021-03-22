package main
import(
	"testing"
	"os"
)
func checkEnv(expected int, t *testing.T){
	got := readEnv()
	os.Unsetenv("DEFAULT_1")
	os.Unsetenv("SECRET_1")
	os.Unsetenv("STATUS_1")
	os.Unsetenv("DEFAULT_2")
	os.Unsetenv("SECRET_2")
	os.Unsetenv("STATUS_2")
	if len(got) != expected {
		t.Errorf("len(got) = %d; want %d", len(got), expected)
	}
}
func TestReadEnvEmpty(t *testing.T) {
	checkEnv(0, t)
}
func TestReadEnv1(t *testing.T) {
	os.Setenv("DEFAULT_1", "free")
	os.Setenv("SECRET_1", "randomsecret")
	os.Setenv("STATUS_1", "test")
	checkEnv(1, t)
	
}
func TestReadEnvMissingSecret(t *testing.T) {
	os.Setenv("DEFAULT_1", "free")
	//os.Setenv("SECRET_1", "randomsecret")
	os.Setenv("STATUS_1", "test")
	checkEnv(0, t)
}
func TestReadEnvMissingStatus(t *testing.T) {
	os.Setenv("DEFAULT_1", "free")
	os.Setenv("SECRET_1", "randomsecret")
	//os.Setenv("STATUS_1", "test")
	checkEnv(0, t)
}
func TestReadEnvMissingDefault(t *testing.T) {
	//os.Setenv("DEFAULT_1", "free")
	os.Setenv("SECRET_1", "randomsecret")
	os.Setenv("STATUS_1", "test")
	checkEnv(1, t)
}


