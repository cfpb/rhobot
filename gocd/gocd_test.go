package gocd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

var gocd_pipeline_config []byte

func init() {
	buf := bytes.NewBuffer(nil)
	f, _ := os.Open("gocd_pipeline_config.json")
	io.Copy(buf, f)
	f.Close()
	gocd_pipeline_config = buf.Bytes()
}

func TestMarshalJSONHAL(t *testing.T) {
	file, e := ioutil.ReadFile("./gocd_pipeline_config.json")
	if e != nil {
		t.Error("file can't be loaded")
	}

	var prettyJSON bytes.Buffer
	e = json.Indent(&prettyJSON, file, "", "\t")
	if e != nil {
		t.Error("JSON parse error")
	}
	fmt.Println("gocd_pipeline_config:", string(prettyJSON.Bytes()))

	var pipeline Pipeline
	e = json.Unmarshal(file, &pipeline)
	if e != nil {
		t.Error("Unmarshaling error")
	}
	spew.Dump(pipeline)

	fmt.Printf("Pipeline Name: %+v\n", pipeline.Name)
	fmt.Printf("Pipeline Git URL: %v:%v\n", pipeline.Materials[0].Attributes.URL, pipeline.Materials[0].Attributes.Branch)
}
