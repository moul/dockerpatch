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

func (d *Dockerfile) FilterOptimize() error {
	globalExposedPorts := []string{}
	for _, node := range d.GetNodesByType(command.Expose) {
		nodeExposedPorts := strings.Split(node.Original, " ")[1:]
		globalExposedPorts = append(globalExposedPorts, nodeExposedPorts...)
	}
	d.RemoveNodesByType(command.Expose)
	d.AppendLine(fmt.Sprintf("EXPOSE %s", strings.Join(globalExposedPorts, " ")))
	return nil
}
