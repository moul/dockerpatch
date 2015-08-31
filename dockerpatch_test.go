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

func TestDockerfile_AppendNode(t *testing.T) {
	Convey("Testing Dockerfile.AppendNode()", t, func() {
		dockerfile := DockerfileNew()
		So(dockerfile, ShouldNotBeNil)
		So(dockerfile.Length(), ShouldEqual, 0)

		for i := 0; i < 10; i++ {
			node, err := ParseLine(fmt.Sprintf("RUN echo %d", i))
			So(err, ShouldBeNil)

			err = dockerfile.AppendNode(node)
			So(err, ShouldBeNil)
			So(dockerfile.Length(), ShouldEqual, i+1)
		}

		So(dockerfile.Length(), ShouldEqual, 10)
		expected := `RUN echo 0
RUN echo 1
RUN echo 2
RUN echo 3
RUN echo 4
RUN echo 5
RUN echo 6
RUN echo 7
RUN echo 8
RUN echo 9`
		So(dockerfile.String(), ShouldEqual, expected)
	})
}

func TestDockerfile_PrependNode(t *testing.T) {
	Convey("Testing Dockerfile.PrependNode()", t, func() {
		dockerfile := DockerfileNew()
		So(dockerfile, ShouldNotBeNil)
		So(dockerfile.Length(), ShouldEqual, 0)

		for i := 0; i < 10; i++ {
			node, err := ParseLine(fmt.Sprintf("RUN echo %d", i))
			So(err, ShouldBeNil)

			err = dockerfile.PrependNode(node)
			So(err, ShouldBeNil)
			So(dockerfile.Length(), ShouldEqual, i+1)
		}

		So(dockerfile.Length(), ShouldEqual, 10)
		expected := `RUN echo 9
RUN echo 8
RUN echo 7
RUN echo 6
RUN echo 5
RUN echo 4
RUN echo 3
RUN echo 2
RUN echo 1
RUN echo 0`
		So(dockerfile.String(), ShouldEqual, expected)
	})
}

func TestDockerfile_AppendLine(t *testing.T) {
	Convey("Testing Dockerfile.AppendLine()", t, func() {
		dockerfile := DockerfileNew()
		So(dockerfile, ShouldNotBeNil)
		So(dockerfile.Length(), ShouldEqual, 0)

		for i := 0; i < 10; i++ {
			err := dockerfile.AppendLine(fmt.Sprintf("RUN echo %d", i))
			So(err, ShouldBeNil)
			So(dockerfile.Length(), ShouldEqual, i+1)
		}

		So(dockerfile.Length(), ShouldEqual, 10)
		expected := `RUN echo 0
RUN echo 1
RUN echo 2
RUN echo 3
RUN echo 4
RUN echo 5
RUN echo 6
RUN echo 7
RUN echo 8
RUN echo 9`
		So(dockerfile.String(), ShouldEqual, expected)
	})
}

func TestDockerfile_PrependLine(t *testing.T) {
	Convey("Testing Dockerfile.PrependLine()", t, func() {
		dockerfile := DockerfileNew()
		So(dockerfile, ShouldNotBeNil)
		So(dockerfile.Length(), ShouldEqual, 0)

		for i := 0; i < 10; i++ {
			err := dockerfile.PrependLine(fmt.Sprintf("RUN echo %d", i))
			So(err, ShouldBeNil)
			So(dockerfile.Length(), ShouldEqual, i+1)
		}

		So(dockerfile.Length(), ShouldEqual, 10)
		expected := `RUN echo 9
RUN echo 8
RUN echo 7
RUN echo 6
RUN echo 5
RUN echo 4
RUN echo 3
RUN echo 2
RUN echo 1
RUN echo 0`
		So(dockerfile.String(), ShouldEqual, expected)
	})
}

func TestDockerfile_RemoveNodesByType(t *testing.T) {
	Convey("Testing Dockerfile.RemoveNodesByType", t, func() {
		dockerfile, err := DockerfileFromString(ExampleDockerfile)
		So(err, ShouldBeNil)
		So(dockerfile.root, ShouldNotBeNil)
		So(dockerfile.root.Dump(), ShouldEqual, ExampleDockerfileDump)
	})
}

func TestDockerfile(t *testing.T) {
	Convey("Testing Dockerfile", t, func() {
		dockerfile := DockerfileNew()
		So(dockerfile, ShouldNotBeNil)
		So(dockerfile.String(), ShouldEqual, "")
		So(dockerfile.Length(), ShouldEqual, 0)
		So(dockerfile.From(), ShouldEqual, "")

		err := dockerfile.SetFrom("ubuntu:latest")
		So(err, ShouldBeNil)
		So(dockerfile.String(), ShouldEqual, "FROM ubuntu:latest")
		So(dockerfile.Length(), ShouldEqual, 1)
		So(dockerfile.From(), ShouldEqual, "ubuntu:latest")

		err = dockerfile.SetFrom("debian")
		So(err, ShouldBeNil)
		So(dockerfile.String(), ShouldEqual, "FROM debian")
		So(dockerfile.Length(), ShouldEqual, 1)
		So(dockerfile.From(), ShouldEqual, "debian")

		err = dockerfile.AppendLine("RUN echo hello world")
		So(err, ShouldBeNil)
		So(dockerfile.String(), ShouldEqual, "FROM debian\nRUN echo hello world")
		So(dockerfile.Length(), ShouldEqual, 2)

		err = dockerfile.AppendLine("RUN echo goodbye world")
		So(err, ShouldBeNil)
		So(dockerfile.String(), ShouldEqual, "FROM debian\nRUN echo hello world\nRUN echo goodbye world")
		So(dockerfile.Length(), ShouldEqual, 3)

		err = dockerfile.AddLineAfterFrom("RUN echo after from")
		So(err, ShouldBeNil)
		So(dockerfile.String(), ShouldEqual, "FROM debian\nRUN echo after from\nRUN echo hello world\nRUN echo goodbye world")
		So(dockerfile.Length(), ShouldEqual, 4)
	})
}
