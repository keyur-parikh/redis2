package command_functions

import (
	"errors"
	"github.com/keyur-parikh/redis2/internal/definitions"
	"strconv"
	"time"
)

// This package will have all the functions that have to do with strings

func HandleGetString(args []string, KVStore map[definitions.RedisKey]definitions.RedisValue, stringToKeys map[string]definitions.RedisKey) ([]string, error) {
	// This function is only called when the function mapper maps this to the get command
	if len(args) != 1 {
		return nil, errors.New("get must be followed by key")
	}
	stringKey := args[0]
	keyDictionary := stringToKeys
	redisKey, ok := keyDictionary[stringKey]
	if !ok {
		return nil, errors.New("no key found")
	}
	if redisKey.Duration != 0 && time.Now().After(redisKey.Creation.Add(redisKey.Duration)) {
		_, _ = HandleDelete(args, KVStore, stringToKeys)
		return nil, errors.New("key expired")
	}
	dictionary := KVStore
	value, ok := dictionary[redisKey]
	if !ok {
		return nil, errors.New("value not found")
	}
	if value.Type != definitions.StringType {
		return nil, errors.New("value not of string type")
	}

	strValue, ok := value.Value.(string)
	if !ok {
		return nil, errors.New("value for some reason was stored wrong")
	}
	return []string{strValue}, nil
}

func HandleSetString(args []string, KVStore map[definitions.RedisKey]definitions.RedisValue, stringToKeys map[string]definitions.RedisKey) ([]string, error) {
	stringKey := args[0]
	stringValue := args[1]
	var expiration int
	if len(args) == 4 {
		if args[2] != "px" {
			return nil, errors.New(`expected "px" as third argument`)
		}
		var err error
		expiration, err = strconv.Atoi(args[3])
		if err != nil {
			return nil, errors.New("expiration must be a valid integer")
		}
	} else if len(args) != 2 {
		return nil, errors.New("wrong number of arguments")
	}
	redisKey := definitions.RedisKey{Creation: time.Now(), Duration: time.Duration(expiration) * time.Millisecond}
	keyDictionary := stringToKeys
	keyDictionary[stringKey] = redisKey
	redisValue := definitions.RedisValue{Type: definitions.StringType, Value: stringValue}
	dictionary := KVStore
	dictionary[redisKey] = redisValue
	return nil, nil
}

func HandleDelete(args []string, KVStore map[definitions.RedisKey]definitions.RedisValue, stringToKeys map[string]definitions.RedisKey) ([]string, error) {
	if len(args) != 1 {
		return nil, errors.New("del must be followed by one key")
	}
	stringKey := args[0]
	keyDictionary := stringToKeys
	redisKey, _ := keyDictionary[stringKey]
	dictionary := KVStore
	delete(keyDictionary, stringKey)
	delete(dictionary, redisKey)
	return nil, nil
}
