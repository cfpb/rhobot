package gocd

import (
	"bytes"
	"encoding/json"
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"io"
	"os"
	"reflect"
	"testing"
)

var gocd_pipeline_config []byte

func init() {
	buf := bytes.NewBuffer(nil)
	f, _ := os.Open("./test.json")
	io.Copy(buf, f)
	f.Close()
	gocd_pipeline_config = buf.Bytes()
}

func TestMarshalJSONHAL(t *testing.T) {
	var pipeline Pipeline = ReadPipelineJSONFromFile("./test.json")
	//spew.Dump(pipeline)

	fmt.Printf("Pipeline Name: %+v\n", pipeline.Name)
	fmt.Printf("Pipeline Git URL: %v:%v\n", pipeline.Materials[0].Attributes.URL, pipeline.Materials[0].Attributes.Branch)
}

// TestUnmarshalFidelityLoss checks that data can be reserielized without fidelity loss
func TestUnmarshalFidelityLoss(t *testing.T) {
	data := UnmarshalPipeline(gocd_pipeline_config)
	gocd_pipeline_config2, _ := json.Marshal(data)
	data2 := UnmarshalPipeline(gocd_pipeline_config2)
	if !reflect.DeepEqual(data, data2) {
		t.Error("not the same")
	}
	fmt.Println("TestUnmarshalFidelityLoss")
}

func TestGocdPOST(t *testing.T) {
	pipeline := ReadPipelineJSONFromFile("./test.json")
	pipelineConfig := PipelineConfig{"Dev", pipeline}
	pipelineConfigPOST("http://localhost:8153", pipelineConfig)
	fmt.Println("TestGocdPOST")
}

func TestGocdGET(t *testing.T) {
	pipelineGET("http://localhost:8153", "test")
	fmt.Println("TestGocdGET")
}
