package healthcheck

import (
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"
)

var healthchecks = `---
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
metadata:
  distribution:
    - stefan.fox@cfpb.gov
`

func TestUnmarshal(t *testing.T) {

	unmarshalHealthChecks(healthchecks)

}

// TestUnmarshalFidelityLoss checks that data can be reserielized without fidelity loss
func TestUnmarshalFidelityLoss(t *testing.T) {

	data := unmarshalHealthChecks(healthchecks)
	healthchecks2, _ := yaml.Marshal(data)
	data2 := unmarshalHealthChecks(string(healthchecks2))
	if !reflect.DeepEqual(data, data2) {
		t.Error("not the same")
	}
}
