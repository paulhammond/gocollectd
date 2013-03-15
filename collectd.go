package gocollectd

type Value struct {
	Hostname           string
	plugin             string
	pluginInstance     string
	pluginType         string
	pluginTypeInstance string
	Number             uint16
	CdTime             uint64
	DataType           uint8
	Bytes              []byte
}
