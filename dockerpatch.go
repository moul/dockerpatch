package dockerpatch

import (
	"bytes"
	"io"
	"strings"

	"github.com/docker/docker/builder/parser"
)

type Dockerfile struct {
	root *parser.Node
}

// DockerfileNew returns an empty Dockerfile
func DockerfileNew() *Dockerfile {
	dockerfile, _ := DockerfileFromString("")
	return dockerfile
}

// DockerfileFromString reads a Dockerfiler as string
func DockerfileFromString(input string) (*Dockerfile, error) {
	payload := new(bytes.Buffer)
	payload.Write([]byte(input))
	return DockerfileRead(payload)
}

// DockerfileRead reads a Dockerfile as io.Reader
func DockerfileRead(input io.Reader) (*Dockerfile, error) {
	dockerfile := Dockerfile{}

	root, err := parser.Parse(input)
	if err != nil {
		return nil, err
	}
	dockerfile.root = root

	return &dockerfile, nil
}

// String returns a docker-readable Dockerfile
func (d *Dockerfile) String() string {
	lines := []string{}
	for _, child := range d.root.Children {
		lines = append(lines, child.Original)
	}
	return strings.Join(lines, "\n")
}
