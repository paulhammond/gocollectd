package gocollectd

type Value struct {
	Hostname       string
	Plugin         string
	PluginInstance string
	Type           string
	TypeInstance   string
	Number         uint16
	CdTime         uint64
	DataType       uint8
	Bytes          []byte
}
