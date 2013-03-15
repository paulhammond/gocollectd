package gocollectd

import (
	"bytes"
	"encoding/binary"
)

func Parse(b []byte) []Value {
	r := make([]Value, 0)

	buf := bytes.NewBuffer(b)
	var val Value
	var partType uint16
	var partLength uint16
	var time uint64
	var valueCount uint16

	for buf.Len() > 0 {
		binary.Read(buf, binary.BigEndian, &partType)
		binary.Read(buf, binary.BigEndian, &partLength)
		partBytes := buf.Next(int(partLength) - 4)
		partBuffer := bytes.NewBuffer(partBytes)
		switch {
		case partType == 0:
			str := partBuffer.String()
			val.Hostname = str[0 : len(str)-1]
		case partType == 1:
			binary.Read(partBuffer, binary.BigEndian, &time)
			val.CdTime = time << 30
		case partType == 2:
			str := partBuffer.String()
			val.plugin = str[0 : len(str)-1]
		case partType == 3:
			str := partBuffer.String()
			val.pluginInstance = str[0 : len(str)-1]
		case partType == 4:
			str := partBuffer.String()
			val.pluginType = str[0 : len(str)-1]
		case partType == 5:
			str := partBuffer.String()
			val.pluginTypeInstance = str[0 : len(str)-1]
		case partType == 6:
			binary.Read(partBuffer, binary.BigEndian, &valueCount)

			for i := uint16(0); i < valueCount; i++ {

				valueBytes := make([]byte, 8, 8) // holds a copy so we lose reference to the underlying slice data

				// messy: collectd's protocol puts things in a weird order.
				binary.Read(partBuffer, binary.BigEndian, &val.DataType)
				copy(valueBytes, partBytes[2+valueCount+(i*8):2+valueCount+8+(i*8)])

				val.Number = i
				val.Bytes = valueBytes

				r = append(r, val)
			}
		case partType == 7:
			// interval, ignore
		case partType == 8:
			// high res time
			binary.Read(partBuffer, binary.BigEndian, &val.CdTime)
		case partType == 9:
			// interval, ignore
		case partType == 0x100:
			// message (notifications), ignore
		case partType == 0x100:
			// severity, ignore
		case partType == 0x200:
			// Signature (HMAC-SHA-256), todo
		case partType == 0x210:
			// Encryption (AES-256/OFB/SHA-1), todo
		default:
			// todo: log unexpected type here
		}
	}
	return r
}