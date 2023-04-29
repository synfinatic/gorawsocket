package gorawsocket

import (
	"encoding/binary"
	"errors"
	"net"

	"github.com/google/gopacket/layers"
)

type IPv4 struct {
	layers.IPv4
	Version    uint8
	IHL        uint8
	TOS        uint8
	Length     uint16
	Id         uint16
	Flags      layers.IPv4Flag
	FragOffset uint16
	TTL        uint8
	Protocol   layers.IPProtocol
	Checksum   uint16
	SrcIP      net.IP
	DstIP      net.IP
	Options    []layers.IPv4Option
	Padding    []byte
}

// SerializeTo writes the serialized form of this layer into the
// SerializationBuffer, implementing gopacket.SerializableLayer.
func (ip *IPv4) SerializeTo(b SerializeBuffer, opts SerializeOptions) error {
	optionLength := ip.getIPv4OptionSize()
	bytes, err := b.PrependBytes(20 + int(optionLength))
	if err != nil {
		return err
	}
	if opts.FixLengths {
		ip.IHL = 5 + (optionLength / 4)
		ip.Length = uint16(len(b.Bytes()))
	}
	bytes[0] = (ip.Version << 4) | ip.IHL
	bytes[1] = ip.TOS

	if opts.HostByteOrderLength {
		endian().PutUint16(bytes[2:], ip.Length)
	} else {
		binary.BigEndian.PutUint16(bytes[2:], ip.Length)
	}

	binary.BigEndian.PutUint16(bytes[4:], ip.Id)
	binary.BigEndian.PutUint16(bytes[6:], ip.flagsfrags())
	bytes[8] = ip.TTL
	bytes[9] = byte(ip.Protocol)
	if err := ip.AddressTo4(); err != nil {
		return err
	}
	copy(bytes[12:16], ip.SrcIP)
	copy(bytes[16:20], ip.DstIP)

	curLocation := 20
	// Now, we will encode the options
	for _, opt := range ip.Options {
		switch opt.OptionType {
		case 0:
			// this is the end of option lists
			bytes[curLocation] = 0
			curLocation++
		case 1:
			// this is the padding
			bytes[curLocation] = 1
			curLocation++
		default:
			bytes[curLocation] = opt.OptionType
			bytes[curLocation+1] = opt.OptionLength

			// sanity checking to protect us from buffer overrun
			if len(opt.OptionData) > int(opt.OptionLength-2) {
				return errors.New("option length is smaller than length of option data")
			}
			copy(bytes[curLocation+2:curLocation+int(opt.OptionLength)], opt.OptionData)
			curLocation += int(opt.OptionLength)
		}
	}

	if opts.ComputeChecksums {
		ip.Checksum = checksum(bytes)
	}
	binary.BigEndian.PutUint16(bytes[10:], ip.Checksum)
	return nil
}

// for the current ipv4 options, return the number of bytes (including
// padding that the options used)
func (ip *IPv4) getIPv4OptionSize() uint8 {
	optionSize := uint8(0)
	for _, opt := range ip.Options {
		switch opt.OptionType {
		case 0:
			// this is the end of option lists
			optionSize++
		case 1:
			// this is the padding
			optionSize++
		default:
			optionSize += opt.OptionLength

		}
	}
	// make sure the options are aligned to 32 bit boundary
	if (optionSize % 4) != 0 {
		optionSize += 4 - (optionSize % 4)
	}
	return optionSize
}

func checksum(bytes []byte) uint16 {
	// Clear checksum bytes
	bytes[10] = 0
	bytes[11] = 0

	// Compute checksum
	var csum uint32
	for i := 0; i < len(bytes); i += 2 {
		csum += uint32(bytes[i]) << 8
		csum += uint32(bytes[i+1])
	}
	for {
		// Break when sum is less or equals to 0xFFFF
		if csum <= 65535 {
			break
		}
		// Add carry to the sum
		csum = (csum >> 16) + uint32(uint16(csum))
	}
	// Flip all the bits
	return ^uint16(csum)
}

func (ip *IPv4) flagsfrags() (ff uint16) {
	ff |= uint16(ip.Flags) << 13
	ff |= ip.FragOffset
	return
}
