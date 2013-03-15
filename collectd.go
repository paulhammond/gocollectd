package gocollectd

const (
	Counter  = 0
	Guage    = 1
	Derive   = 2
	Absolute = 3
)

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
