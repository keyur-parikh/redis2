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

type RedisKey struct {
	Creation time.Time
	Duration time.Duration
}
type RedisValue struct {
	Type  ValueType
	Value interface{}
}

type CommandInfo struct {
	ParsedCommand []string
	Connection    net.Conn
}
