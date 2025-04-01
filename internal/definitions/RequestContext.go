package definitions

import "net"

type RequestContext struct {
	Connection net.Conn
	KVStore    map[string]string
}
