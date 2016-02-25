package healthcheck

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"

	"github.com/cfpb/rhobot/database"
)

var healthchecks []byte

func init() {
	buf := bytes.NewBuffer(nil)
	f, _ := os.Open("test.yml")
	io.Copy(buf, f)
	f.Close()
	healthchecks = buf.Bytes()
}

func TestUnmarshal(t *testing.T) {

	unmarshalHealthChecks(healthchecks)

}

// TestUnmarshalFidelityLoss checks that data can be reserielized without fidelity loss
func TestUnmarshalFidelityLoss(t *testing.T) {

	data := unmarshalHealthChecks(healthchecks)
	healthchecks2, _ := yaml.Marshal(data)
	data2 := unmarshalHealthChecks(healthchecks2)
	if !reflect.DeepEqual(data, data2) {
		t.Error("not the same")
	}
}

func TestRunningBasicChecks(t *testing.T) {

	host := os.Getenv("PGHOST")
	db := os.Getenv("PGDATABASE")
	user := os.Getenv("PGUSER")
	pass := os.Getenv("PGPASSWORD")
	uri := fmt.Sprintf("postgres://%s:%s@%s/%s", user, pass, host, db)

	cxn := database.GetPGConnection(uri)
	healthChecks := unmarshalHealthChecks(healthchecks)
	RunHealthChecks(healthChecks, cxn)

}

func TestEvaluatingBasicChecks(t *testing.T) {

	host := os.Getenv("PGHOST")
	db := os.Getenv("PGDATABASE")
	user := os.Getenv("PGUSER")
	pass := os.Getenv("PGPASSWORD")
	uri := fmt.Sprintf("postgres://%s:%s@%s/%s", user, pass, host, db)

	cxn := database.GetPGConnection(uri)
	healthChecks := unmarshalHealthChecks(healthchecks)
	healthChecks = RunHealthChecks(healthChecks, cxn)
	// EvaluateHealthChecks(healthChecks) // this should fail
}
