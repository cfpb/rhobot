package healthcheck

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"

	"github.com/cfpb/rhobot/database"
)

var healthchecks = []byte(`---
healthchecks:
  -
    error: true
    expected: true
    name: "basic test 1 (should pass)"
    query: "select (select count(1) from information_schema.tables) > 0;"
  -
    error: false
    expected: true
    name: "basic test 2 (should pass)"
    query: "select (select count(1) from information_schema.tables) > 0;"
  -
    error: false
    expected: false
    name: "basic test 3 (should warn)"
    query: "select (select count(1) from information_schema.tables) > 0;"
  -
    error: true
    expected: false
    name: "basic test 4 (should error)"
    query: "select (select count(1) from information_schema.tables) > 0;"
  -
    error: true
    expected: 0
    name: "basic test 5 (should error)"
    query: "select count(1) from information_schema.tables;"
metadata:
  distribution:
    - stefan.fox@cfpb.gov
`)

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
	EvaluateHealthChecks(healthChecks) // this should fail
}
