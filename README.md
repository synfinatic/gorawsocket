# gorawsocket

A cross-platform library for writing packets at the IP layer/layer 3.

## Why?

I wanted to use Go to send packets with forged source IP addresses and 
historically to do this I used the [google/gopacket](
https://pkg.go.dev/github.com/google/gopacket) library which necessiated 
using CGO due to the dependency on the [libpcap](https://www.tcpdump.org)
library.

Also, I wanted the ability to send packets to the local system without
relying on the loopback interface- something that is not possible with
injecting frames at Layer 2 via libpcap.

## Overview

gorawsocket is designed to work well with the [google/gopacket/layers](
https://pkg.go.dev/github.com/google/gopacket/layers) library.



## Resources

* [rawip.txt](https://www.digiater.nl/openvms/decus/vmslt01b/sec/rawip.txt)
* [SOCK_RAW Demystified](https://sock-raw.org/papers/sock_raw)
* [Introduction to RAW-sockets](https://tuprints.ulb.tu-darmstadt.de/6243/1/TR-18.pdf)
* [traceroute source code](ftp://ftp.ee.lbl.gov/traceroute-1.4a12.tar.gz)
* [raw sockets in go ip layer](https://darkcoding.net/uncategorized/raw-sockets-in-go-ip-layer/)
