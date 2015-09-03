# dockerpatch
:whale: Read, write, manipulate, convert &amp; apply filters to Dockerfiles

[![GoDoc](https://godoc.org/github.com/moul/dockerpatch?status.svg)](https://godoc.org/github.com/moul/dockerpatch)

## Web demo

A demo is available on [dockerpatch.appspot.com](https://dockerpatch.appspot.com)

[![Screenshot](https://raw.githubusercontent.com/moul/dockerpatch/master/assets/screen.png)](https://dockerpatch.appspot.com)

## Install

```bash
go get github.com/moul/dockerpatch/...
```


## Examples

```console
$ cat input
```

```Dockerfile
FROM ubuntu:14.04

RUN apt-get update && apt-get install wget -y
RUN wget http://s3.amazonaws.com/influxdb/influxdb_latest_amd64.deb
RUN dpkg -i influxdb_latest_amd64.deb
RUN rm -r /opt/influxdb/shared

VOLUME /opt/influxdb/shared

CMD /usr/bin/influxdb --pidfile /var/run/influxdb.pid -config /opt/influxdb/shared/config.toml

EXPOSE 8083
EXPOSE 8086
EXPOSE 8090
EXPOSE 8099
```

---

```console
$ dockerpatch --to-arm --optimize ./input
```

```Dockerfile
FROM armbuild/ubuntu:14.04
RUN apt-get update && apt-get install wget -y && wget http://s3.amazonaws.com/influxdb/influxdb_latest_armhf.deb && dpkg -i influxdb_latest_armhf.deb && rm -r /opt/influxdb/shared
VOLUME /opt/influxdb/shared
CMD /usr/bin/influxdb --pidfile /var/run/influxdb.pid -config /opt/influxdb/shared/config.toml
EXPOSE 8083 8086 8090 8099
```

---

```console
cat input | dockerpatch --to-arm -
```

```Dockerfile
FROM armbuild/ubuntu:14.04
RUN apt-get update && apt-get install wget -y
RUN wget http://s3.amazonaws.com/influxdb/influxdb_latest_armhf.deb
RUN dpkg -i influxdb_latest_armhf.deb
RUN rm -r /opt/influxdb/shared
VOLUME /opt/influxdb/shared
CMD /usr/bin/influxdb --pidfile /var/run/influxdb.pid -config /opt/influxdb/shared/config.toml
EXPOSE 8083
EXPOSE 8086
EXPOSE 8090
EXPOSE 8099
```

---

```console
cat input | dockerpatch --disable-network -
```

```Dockerfile
FROM ubuntu:14.04
RUN apt-get update && apt-get install wget -y
RUN wget http://s3.amazonaws.com/influxdb/influxdb_latest_amd64.deb
RUN dpkg -i influxdb_latest_amd64.deb
RUN rm -r /opt/influxdb/shared
VOLUME /opt/influxdb/shared
CMD /usr/bin/influxdb --pidfile /var/run/influxdb.pid -config /opt/influxdb/shared/config.toml
```

## Usage

```console
$ dockerpatch -h
NAME:
   dockerpatch - Read, write, manipulate, convert & apply filters to Dockerfiles

USAGE:
   dockerpatch [global options] command [command options] [arguments...]

VERSION:
   0.1.0

AUTHOR(S):
   Manfred Touron <https://github.com/moul/dockerpatch>

COMMANDS:
   patch	Patch a Dockerfile with filters
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h		show help
   --version, -v	print the version
```

```console
$ dockerpatch patch -h
NAME:
   patch - Patch a Dockerfile with filters

USAGE:
   command patch [command options] [arguments...]

DESCRIPTION:
   patch [--to-arm]

OPTIONS:
   --to-arm		Convert to armhf architecture
   --disable-network	Remove network rules
   --optimize		Optimize Dockerfile
```

## License

MIT
