package gocollectd

import (
	"bytes"
	"encoding/binary"
	"errors"
)

var ErrorUnsupported = errors.New("Unsupported collectd packet recieved")
var ErrorInvalid = errors.New("Invalid collectd packet recieved")

func Parse(b []byte) (*[]Value, error) {
	r := make([]Value, 0)

	buf := bytes.NewBuffer(b)
	var val Value
	var packetHeader struct {
		PartType   uint16
		PartLength uint16
	}
	var time uint64
	var valueCount uint16
	var err error

	for buf.Len() > 0 {
		err = binary.Read(buf, binary.BigEndian, &packetHeader)
		if err != nil {
			return nil, err
		}
		if packetHeader.PartLength < 5 {
			return nil, ErrorInvalid
		}

		partBytes := buf.Next(int(packetHeader.PartLength) - 4)
		if len(partBytes) < int(packetHeader.PartLength)-4 {
			return nil, ErrorInvalid
		}
		partBuffer := bytes.NewBuffer(partBytes)

		switch packetHeader.PartType {
		case 0:
			str := partBuffer.String()
			val.Hostname = str[0 : len(str)-1]
		case 1:
			err = binary.Read(partBuffer, binary.BigEndian, &time)
			if err != nil {
				return nil, err
			}
			val.CdTime = time << 30
		case 2:
			str := partBuffer.String()
			val.plugin = str[0 : len(str)-1]
		case 3:
			str := partBuffer.String()
			val.pluginInstance = str[0 : len(str)-1]
		case 4:
			str := partBuffer.String()
			val.pluginType = str[0 : len(str)-1]
		case 5:
			str := partBuffer.String()
			val.pluginTypeInstance = str[0 : len(str)-1]
		case 6:
			err = binary.Read(partBuffer, binary.BigEndian, &valueCount)
			if err != nil {
				return nil, err
			}

			for i := uint16(0); i < valueCount; i++ {

				valueBytes := make([]byte, 8, 8) // holds a copy so we lose reference to the underlying slice data

				// messy: collectd's protocol puts things in a weird order.
				err = binary.Read(partBuffer, binary.BigEndian, &val.DataType)
				if err != nil {
					return nil, err
				}
				copy(valueBytes, partBytes[2+valueCount+(i*8):2+valueCount+8+(i*8)])

				val.Number = i
				val.Bytes = valueBytes

				r = append(r, val)
			}
		case 7:
			// interval, ignore
		case 8:
			// high res time
			err = binary.Read(partBuffer, binary.BigEndian, &val.CdTime)
			if err != nil {
				return nil, err
			}
		case 9:
			// interval, ignore
		case 0x100:
			// message (notifications), ignore
		case 0x101:
			// severity, ignore
		case 0x200:
			// Signature (HMAC-SHA-256), todo
			return nil, ErrorUnsupported
		case 0x210:
			// Encryption (AES-256/OFB/SHA-1), todo
			return nil, ErrorUnsupported
		default:
			return nil, ErrorUnsupported
		}
	}
	return &r, nil
}
