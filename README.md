# Connbeat

[![Join the chat at https://gitter.im/connbeat/Lobby](https://badges.gitter.im/connbeat/Lobby.svg)](https://gitter.im/connbeat/Lobby?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

[![Build Status](https://travis-ci.org/raboof/connbeat.svg?branch=master)](https://travis-ci.org/raboof/connbeat)

Connbeat, short for 'Connectionbeat', is an open source agent that monitors TCP connection metadata and
ships the data to Kafka or Elasticsearch, or an HTTP endpoint.

The main distinction from [Packetbeat](https://www.elastic.co/products/beats/packetbeat)
is that Connbeat is intended to be able to monitor all connections on a
machine (rather than just selected protocols), and does not inspect the
package/connection contents, only metadata.

## Credits

Development of connbeat was funded by [StackState](http://www.stackstate.com).
Collecting connection data is only part of the puzzle: [StackState](http://www.stackstate.com)
combines it with information from many other sources, presenting it in a way that
provides actionable insights.

![StackState logo](http://www.stackstate.com/wp-content/uploads/2016/12/Sts_LOGO_RGB_Full_Horizontal.png)

## Building

### On linux and mac osx

Connbeat is built with 'make'. You need at least golang 1.7.3.

    # Make sure $GOPATH is set
    go get github.com/raboof/connbeat
    cd $GOPATH/src/github.com/raboof/connbeat
    make

It is possible to build connbeat on OSX. However, no integrations are implemented at this
point. It is possible to run the unit tests.

### Building for linux on OSX

To build a linux binary on OSX you can use docker:

    docker run --rm -v "$PWD":/go/src/github.com/raboof/connbeat -w /go/src/github.com/raboof/connbeat golang:1.7.4 make

This will produce a dynamically linked connbeat executable in the current
directory.

To create linux packages, use `make package`

## Running

Edit the configuration (connbeat.yml) to specify where you want your events to go (e.g. Kafka, Elasticsearch, the console).

You need to be root if you want to see the process for processes other than your own:

    sudo ./connbeat

You can view the events on kafka with something like kafkacat:

    kafkacat -C -b localhost -t connbeat

## Docker

You can use connbeat to monitor TCP connections from docker instances - see
[here](docker#readme) for details.

## Performance overhead

We tested the overhead of running the connbeat agent using the
[TechEmpower web framework benchmarks](https://www.techempower.com/benchmarks/).

After deploying to AWS, we ran the [query](https://www.techempower.com/benchmarks/#test=query)
benchmark workload against the Spring Boot framework.

The result was encouraging: the total requests throughput took a hit of only
0.47% (58 fewer requests on a total of 12312). The average latency was in fact
a little better in the test runs with connbeat - which must of course be caused
by noise, but inspires confidence that connbeat introduce no noticable degredation.

The complete test results can be found in the /tests/performance folder of this repo.

Of course performance impact may vary due to all kinds of circumstances and
differences in workload. We're aware of several potential further
optimizations, which can be applied when a situation comes up where connbeat
does have a noticable impact. If you encounter such a situation, be sure to
file an issue.

## Events

For connections where the agent is the server:

    {
      "@timestamp": "2016-05-20T14:54:29.442Z",
      "beat": {
        "hostname": "yinka",
        "name": "yinka",
        "local_ips": [
          "192.168.2.243"
        ]
      },
      "local_port": 80,
      "local_process": {
        "binary": "dnsmasq",
        "cmdline": ""/usr/sbin/dnsmasq -x /var/run/dnsmasq/dnsmasq.pid -u dnsmasq -7 /etc/dnsmasq.d,.dpkg-dist,.dpkg-old,.dpkg-new --local-service",
        "environ": [
        "LANGUAGE=en_US:en",
        "PATH=/usr/local/sbin:/usr/local/bin:/sbin:/bin:/usr/sbin:/usr/bin",
        "LANG=en_US.UTF-8",
        "_SYSTEMCTL_SKIP_REDIRECT=true",
        "PWD=/",

        ]
      },
      "type": "connbeat"
    }

For connections where the agent appears to be the client:

    {
      "@timestamp": "2016-05-20T14:54:29.506Z",
      "beat": {
        "hostname": "yinka",
        "name": "yinka",
        "local_ips": [
          "192.168.2.243"
        ]
      },
      "local_ip": "192.168.2.243",
      "local_port": 40074,
      "local_process": {
        "binary": "chromium",
        "cmdline": "/usr/lib/chromium/chromium --show-component-extension-options --ignore-gpu-blacklist --ppapi-flash-path=/usr/lib/pepperflashplugin-nonfree/libpepflashplayer.so --ppapi-flash-version=20.0.0.228",
        "environ": [
          ""
        ]
      },
      "remote_ip": "52.91.150.74",
      "remote_port": 443,
      "type": "connbeat"
    }

## Testing

To run the regular go unit test, run 'make unit'.

To also run docker-based system tests, run 'make testsuite'

## Packaging

Preliminary packaging is available, but the resulting packages are not yet
intended for general consumption.

'make package' should be sufficient to produce a deb, rpm and a binary .tar.gz

## Elastic Beat Upgrade 

Currently `elastic\beats` package is set to 5.6.9 and there is a manual change in the `vendor/github/com/elastic/beats/libbeat/script/Makefile` for parameter `TESTIFY_TOOL_REPO`. The value is set from `github.com/elastic/beats/vendor/github.com/stretchr/testify` to `github.com/elastic/beats/vendor/github.com/stretchr/testify/assert` because it tries to download master and the repo doesn't contain the actual code. and also this paramter can't be overriden in this version 5.6.9. So please check this parameter if it can be overriden or not and then change that in your `Makefile` when you update the `elastic/beats` library. Mostly in the latest library you don't need to change this parameter becuase it is fixed. 

## Contributing

Contributions are welcome! Feel free to [submit issues](https://github.com/raboof/connbeat/issues) to discuss problems and propose solutions, or send a [pull request](https://github.com/raboof/connbeat/pulls).

Pull requests are expected to include tests (which are run on Travis). We strive to merge any reasonable features, though features that might increase the load on the machine will likely have to be behind a feature switch that is off by default.

## Security

We take great care to ensure connbeat is secure. If despite our efforts you
have found what looks like a vulnerability, please contact us privately at
aengelen@xebia.com. For extra safety the email may be encrypted with the
public key which can be found at https://keybase.io/raboof
