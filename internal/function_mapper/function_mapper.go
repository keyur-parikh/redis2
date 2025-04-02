package function_mapper

import (
	"errors"
	"github.com/keyur-parikh/redis2/internal/command_functions"
	"github.com/keyur-parikh/redis2/internal/definitions"
	"strings"
)

type CommandFunc func(args []string, KVStore map[definitions.RedisKey]definitions.RedisValue, stringToKeys map[string]definitions.RedisKey) ([]string, error)

var commandTable = map[string]CommandFunc{
	"GET": command_functions.HandleGetString,
	"SET": command_functions.HandleSetString,
	"DEL": command_functions.HandleDelete,
	// etc.
}

func FunctionMapper(args []string) (CommandFunc, error) {
	if len(args) < 1 {
		return nil, errors.New("must be least length of 1")
	}
	command := args[0]
	function, ok := commandTable[strings.ToUpper(command)]
	if !ok {
		return nil, errors.New("not a valid command")
	}
	return function, nil
}
