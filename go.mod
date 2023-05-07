module github.com/synfinatic/gorawsocket

go 1.20

require (
	github.com/alecthomas/kong v0.7.1
	github.com/google/gopacket v1.1.19
	github.com/sirupsen/logrus v1.9.0
)

require golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect

replace github.com/google/gopacket v1.1.19 => github.com/synfinatic/gopacket v1.1.19-0.20230429033913-a1848bb584fe
