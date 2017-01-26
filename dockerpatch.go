package dockerpatch

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/builder/dockerfile/command"
	"github.com/docker/docker/builder/dockerfile/parser"
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

	d := parser.Directive{LookingForDirectives: true}
	parser.SetEscapeToken(parser.DefaultEscapeToken, &d)
	root, err := parser.Parse(input, &d)
	if err != nil {
		return nil, err
	}
	dockerfile.root = root

	return &dockerfile, nil
}

func (d *Dockerfile) RemoveAt(i int) error {
	if i >= 0 && i < d.Length() {
		d.root.Children = append(d.root.Children[:i], d.root.Children[i+1:]...)
		return nil
	}
	return fmt.Errorf("Cannot remove %d: index error", i)
}

// String returns a docker-readable Dockerfile
func (d *Dockerfile) String() string {
	lines := []string{}
	for _, child := range d.root.Children {
		lines = append(lines, child.Original)
	}
	return strings.Join(lines, "\n")
}

// PrependNode attach a new node on first position of the AST
func (d *Dockerfile) PrependNode(node *parser.Node) error {
	d.root.Children = append([]*parser.Node{node}, d.root.Children...)
	return nil
}

// AppendNode attach a new node on last position of the AST
func (d *Dockerfile) AppendNode(node *parser.Node) error {
	d.root.Children = append(d.root.Children, node)
	return nil
}

// RemoveNodesByType removes all nodes of a specific type from the AST
func (d *Dockerfile) RemoveNodesByType(nodeType string) error {
	newChildren := []*parser.Node{}
	for _, node := range d.root.Children {
		if node.Value != nodeType {
			newChildren = append(newChildren, node)
		}
	}
	d.root.Children = newChildren
	return nil
}

// SetFrom sets the current FROM
func (d *Dockerfile) SetFrom(from string) error {
	if err := d.RemoveNodesByType(command.From); err != nil {
		return err
	}

	return d.PrependLine(fmt.Sprintf("FROM %s", from))
}

// ParseLine returns an AST node based on a line
func ParseLine(line string) (*parser.Node, error) {
	tmp, err := DockerfileFromString(line)
	if err != nil {
		return nil, err
	}
	return tmp.root.Children[0], nil
}

// GetFrom returns the current FROM
func (d *Dockerfile) From() string {
	for _, node := range d.root.Children {
		if node.Value == command.From {
			return strings.Split(node.Original, " ")[1]
		}
	}
	return ""
}

// Length returns length of the AST
func (d *Dockerfile) Length() int {
	return len(d.root.Children)
}

// AppendLine parses and appends a new line to the AST
func (d *Dockerfile) AppendLine(line string) error {
	node, err := ParseLine(line)
	if err != nil {
		return err
	}

	return d.AppendNode(node)
}

// PrependLine parses and prepends a new line to the AST
func (d *Dockerfile) PrependLine(line string) error {
	node, err := ParseLine(line)
	if err != nil {
		return err
	}

	return d.PrependNode(node)
}

// AddLineAfterFrom parses and add a line after from in the AST
func (d *Dockerfile) AddLineAfterFrom(line string) error {
	node, err := ParseLine(line)
	if err != nil {
		return err
	}
	return d.AddNodeAfterFrom(node)
}

// AddNodeAfterFrom adds a node after from in the AST
func (d *Dockerfile) AddNodeAfterFrom(node *parser.Node) error {
	if d.Length() == 0 {
		return d.AppendNode(node)
	}

	firstNode := d.root.Children[0]
	if firstNode.Value != command.From {
		return d.PrependNode(node)
	}

	newChildren := []*parser.Node{firstNode, node}
	newChildren = append(newChildren, d.root.Children[1:]...)
	d.root.Children = newChildren
	return nil
}

// Replace tries to replace a string with another in each lines
func (d *Dockerfile) Replace(from, to string) error {
	for i, node := range d.root.Children {
		if strings.Contains(node.Original, from) {
			newNode, err := ParseLine(strings.Replace(node.Original, from, to, -1))
			if err != nil {
				return err
			}
			d.root.Children[i] = newNode
		}
	}
	return nil
}

// GetNodesByType returns nodes matching a type
func (d *Dockerfile) GetNodesByType(nodeType string) []*parser.Node {
	output := []*parser.Node{}
	for _, node := range d.root.Children {
		if node.Value == nodeType {
			output = append(output, node)
		}
	}
	return output
}

// NodeGetArgs returns the arguments of a node
func NodeGetArgs(node *parser.Node) []string {
	return strings.Split(node.Original, " ")[1:]
}

// NodeGetLine returns the arguments of a node as a string
func NodeGetLine(node *parser.Node) string {
	return strings.Join(NodeGetArgs(node), " ")
}
