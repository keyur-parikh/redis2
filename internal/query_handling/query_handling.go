package query_handling

import (
	"errors"
	"fmt"
	"github.com/keyur/redis2/internal/definitions"
	"github.com/keyur/redis2/internal/writer"
	"time"
)

func Get(key string, requestContext *definitions.RequestContext) error {
	dictionary := requestContext.KVStore
	value, ok := dictionary[key]
	if !ok {
		if requestContext.Connection != nil {
			err := writer.InvalidResponseWriter(requestContext)
			if err != nil {
				fmt.Println("error writing failure")
				return err
			} else {
				return errors.New("not a valid value")
			}
		}
	} else {
		responseContent := []string{value}
		err := writer.ArrayResponseWriter(responseContent, requestContext)
		if err != nil {
			fmt.Println("Error writing: ", err)
			return err
		}
	}
	return nil
}

func Set(key string, value string, px bool, timer int, requestContext *definitions.RequestContext) error {
	dictionary := requestContext.KVStore
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
