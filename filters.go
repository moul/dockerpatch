package dockerpatch

import (
	"fmt"
	"strings"

	"github.com/docker/docker/builder/command"
)

func stdFromToArm(input string) string {
	return fmt.Sprintf("armbuild/%s", strings.Replace(input, "/", "-", -1))
}

func (d *Dockerfile) FilterToArm(destArchitecture string) error {
	d.SetFrom(stdFromToArm(d.From()))
	d.Replace("amd64", destArchitecture)
	d.Replace("x86_64", destArchitecture)
	d.Replace("i386", destArchitecture)
	return nil
}

func (d *Dockerfile) FilterDisableNetwork() error {
	d.RemoveNodesByType(command.Expose)
	return nil
}

func (d *Dockerfile) combineNodesByType(nodeType string) error {
	combinedArgs := []string{}
	for _, node := range d.GetNodesByType(nodeType) {
		nodeArgs := strings.Split(node.Original, " ")[1:]
		combinedArgs = append(combinedArgs, nodeArgs...)
	}
	d.RemoveNodesByType(nodeType)
	if len(combinedArgs) > 0 {
		d.AppendLine(fmt.Sprintf("%s %s", strings.ToUpper(nodeType), strings.Join(combinedArgs, " ")))
	}
	return nil
}

func (d *Dockerfile) combineFollowingRunNodes() error {
	hasChanged := true
	for hasChanged {
		hasChanged = false
		for i, node := range d.root.Children {
			switch node.Value {
			case command.Run:
				if i >= d.Length()-1 {
					continue
				}
				next := d.root.Children[i+1]
				if next.Value == command.Run {
					nodeCommand := NodeGetLine(node)
					nextCommand := NodeGetLine(next)
					combined := fmt.Sprintf("RUN %s && %s", nodeCommand, nextCommand)
					newNode, err := ParseLine(combined)
					if err != nil {
						return err
					}
					d.root.Children[i] = newNode
					d.RemoveAt(i + 1)
					hasChanged = true
					break
				}
			}
		}
	}
	return nil
}

func (d *Dockerfile) FilterOptimize() error {
	if err := d.combineNodesByType(command.Expose); err != nil {
		return err
	}

	if err := d.combineFollowingRunNodes(); err != nil {
		return err
	}

	return nil
}
