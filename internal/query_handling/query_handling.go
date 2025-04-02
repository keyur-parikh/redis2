package query_handling

import (
	"errors"
	"fmt"
	"github.com/keyur-parikh/redis2/internal/definitions"
	"github.com/keyur-parikh/redis2/internal/writer"
	"time"
)

func HandleGet([]string, KVStore *map[string]string) ([]string, error) {
	dictionary := *KVStore
	value, ok := dictionary[key]
	if !ok {
		return nil, errors.New("Invalid Response")
	} else {
		responseContent := []string{value}
		return responseContent, nil
	}
}

func HandleSet(key string, value string, px bool, timer int, KVStore *map[string]string) error {
	dictionary := *KVStore
	dictionary[key] = value
	if requestContext.Connection != nil {
		err := writer.SuccessResponseWriter(requestContext)
		if err != nil {
			fmt.Println("error writing success: ", err)
			return err
		}
		if px {
			go func() {
				time.Sleep(time.Duration(timer) * time.Millisecond)
				err := Delete(key, requestContext)
				if err != nil {
					fmt.Printf("Couldnt delete in px")
				}
			}()
		}
	}

	return nil
}

func Delete(key string, requestContext *definitions.RequestContext) error {
	dictionary := requestContext.KVStore
	_, ok := dictionary[key]
	if !ok {
		return errors.New("key was never in the dictionary")
	} else {
		delete(dictionary, key)
		return nil
	}
}
