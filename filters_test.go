package dockerpatch

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	ExampleDockerfileArmString = `FROM armbuild/ubuntu:14.04
RUN apt-get update && apt-get install wget -y
RUN wget http://s3.amazonaws.com/influxdb/influxdb_latest_armhf.deb
RUN dpkg -i influxdb_latest_armhf.deb
RUN rm -r /opt/influxdb/shared
VOLUME /opt/influxdb/shared
CMD /usr/bin/influxdb --pidfile /var/run/influxdb.pid -config /opt/influxdb/shared/config.toml
EXPOSE 8083
EXPOSE 8086
EXPOSE 8090
EXPOSE 8099`

	ExampleDockerfileDisableNetworkString = `FROM ubuntu:14.04
RUN apt-get update && apt-get install wget -y
RUN wget http://s3.amazonaws.com/influxdb/influxdb_latest_amd64.deb
RUN dpkg -i influxdb_latest_amd64.deb
RUN rm -r /opt/influxdb/shared
VOLUME /opt/influxdb/shared
CMD /usr/bin/influxdb --pidfile /var/run/influxdb.pid -config /opt/influxdb/shared/config.toml`

	ExampleDockerfileOptimizeString = `FROM ubuntu:14.04
RUN apt-get update && apt-get install wget -y
RUN wget http://s3.amazonaws.com/influxdb/influxdb_latest_amd64.deb
RUN dpkg -i influxdb_latest_amd64.deb
RUN rm -r /opt/influxdb/shared
VOLUME /opt/influxdb/shared
CMD /usr/bin/influxdb --pidfile /var/run/influxdb.pid -config /opt/influxdb/shared/config.toml
EXPOSE 8083 8086 8090 8099`
)

func TestDockerfile_FilterToArm(t *testing.T) {
	Convey("Testing Dockerfile.FilterToArm", t, func() {
		dockerfile, err := DockerfileFromString(ExampleDockerfile)
		So(err, ShouldBeNil)

		err = dockerfile.FilterToArm("armhf")
		So(err, ShouldBeNil)
		So(dockerfile.String(), ShouldEqual, ExampleDockerfileArmString)
	})
}

func TestDockerfile_FilterDisableNetwork(t *testing.T) {
	Convey("Testing Dockerfile.FilterDisableNetwork", t, func() {
		dockerfile, err := DockerfileFromString(ExampleDockerfile)
		So(err, ShouldBeNil)

		err = dockerfile.FilterDisableNetwork()
		So(err, ShouldBeNil)
		So(dockerfile.String(), ShouldEqual, ExampleDockerfileDisableNetworkString)
	})
}

func TestDockerfile_FilterOptimize(t *testing.T) {
	Convey("Testing Dockerfile.FilterOptimize", t, func() {
		dockerfile, err := DockerfileFromString(ExampleDockerfile)
		So(err, ShouldBeNil)

		err = dockerfile.FilterOptimize()
		So(err, ShouldBeNil)
		So(dockerfile.String(), ShouldEqual, ExampleDockerfileOptimizeString)
	})
}
