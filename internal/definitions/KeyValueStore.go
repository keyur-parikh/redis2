package definitions

import (
	"net"
	"time"
)

type ValueType int

const (
	StringType ValueType = iota
	ListType
	MapType
)

type Key struct {
	Type     ValueType
	Value    interface{}
	TimeSet  time.Time
	Duration time.Duration
}

type CommandInfo struct {
	Command    string
	Connection net.Conn
	KeyData    Key
}
