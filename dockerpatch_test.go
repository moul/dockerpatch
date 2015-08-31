package dockerpatch

import (
	"bytes"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	ExampleDockerfile = `FROM ubuntu:14.04
# COMMENT

RUN apt-get update && apt-get install wget -y
RUN wget http://s3.amazonaws.com/influxdb/influxdb_latest_amd64.deb
RUN dpkg -i influxdb_latest_amd64.deb
RUN rm -r /opt/influxdb/shared

VOLUME /opt/influxdb/shared

CMD /usr/bin/influxdb --pidfile /var/run/influxdb.pid -config /opt/influxdb/shared/config.toml

EXPOSE 8083
EXPOSE 8086
EXPOSE 8090
EXPOSE 8099`
	ExampleDockerfileDump = `(from "ubuntu:14.04")
(run "apt-get update && apt-get install wget -y")
(run "wget http://s3.amazonaws.com/influxdb/influxdb_latest_amd64.deb")
(run "dpkg -i influxdb_latest_amd64.deb")
(run "rm -r /opt/influxdb/shared")
(volume "/opt/influxdb/shared")
(cmd "/usr/bin/influxdb --pidfile /var/run/influxdb.pid -config /opt/influxdb/shared/config.toml")
(expose "8083")
(expose "8086")
(expose "8090")
(expose "8099")`
	ExampleDockerfileString = `FROM ubuntu:14.04
RUN apt-get update && apt-get install wget -y
RUN wget http://s3.amazonaws.com/influxdb/influxdb_latest_amd64.deb
RUN dpkg -i influxdb_latest_amd64.deb
RUN rm -r /opt/influxdb/shared
VOLUME /opt/influxdb/shared
CMD /usr/bin/influxdb --pidfile /var/run/influxdb.pid -config /opt/influxdb/shared/config.toml
EXPOSE 8083
EXPOSE 8086
EXPOSE 8090
EXPOSE 8099`
)

func TestDockerfileFromString(t *testing.T) {
	Convey("Testing DockerfileFromString", t, func() {
		dockerfile, err := DockerfileFromString(ExampleDockerfile)
		So(err, ShouldBeNil)
		So(dockerfile.root, ShouldNotBeNil)
		So(dockerfile.root.Dump(), ShouldEqual, ExampleDockerfileDump)
	})
}

func TestDockerfileRead(t *testing.T) {
	Convey("Testing DockerfileFromString", t, func() {
		payload := new(bytes.Buffer)
		payload.Write([]byte(ExampleDockerfile))
		dockerfile, err := DockerfileRead(payload)
		So(err, ShouldBeNil)
		So(dockerfile.root, ShouldNotBeNil)
		So(dockerfile.root.Dump(), ShouldEqual, ExampleDockerfileDump)
	})
}

func TestDockerfile_String(t *testing.T) {
	Convey("Testing Dockerfile.String", t, func() {
		dockerfile, err := DockerfileFromString(ExampleDockerfile)
		So(err, ShouldBeNil)
		So(dockerfile.root, ShouldNotBeNil)
		So(dockerfile.String(), ShouldEqual, ExampleDockerfileString)
		So(fmt.Sprintf("%s", dockerfile), ShouldEqual, ExampleDockerfileString)
	})
}

func TestDockerfile_New(t *testing.T) {
	Convey("Testing Dockerfile.New", t, func() {
		dockerfile := DockerfileNew()
		So(dockerfile, ShouldNotBeNil)
	})
}
