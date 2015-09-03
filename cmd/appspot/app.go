package dockerpatchapp

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/moul/dockerpatch"
)

func init() {
	http.HandleFunc("/", handler)
}

type ConvertRequest struct {
	Dockerfile string          `json:"Dockerfile,omitempty"`
	Options    map[string]bool `json:"Options,omitempty"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var convertRequest ConvertRequest
	err := decoder.Decode(&convertRequest)

	if err != nil {
		fmt.Fprintf(w, "POST parsing error: %v\n", err)
		return
	}
	dockerfile, err := dockerpatch.DockerfileFromString(convertRequest.Dockerfile)
	if err != nil {
		fmt.Fprintf(w, "Invalid Dockerfile: %v", err)
		return
	}

	if convertRequest.Options["ToArm"] {
		dockerfile.FilterToArm("armhf")
	}
	if convertRequest.Options["DisableNetwork"] {
		dockerfile.FilterDisableNetwork()
	}
	if convertRequest.Options["Optimize"] {
		dockerfile.FilterOptimize()
	}

	fmt.Fprintf(w, "%v\n", dockerfile)
}
