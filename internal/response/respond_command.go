package response

import (
	"errors"
	"fmt"
	"github.com/keyur-parikh/redis2/internal/definitions"
	"github.com/keyur-parikh/redis2/internal/query_handling"
	"strconv"
	"strings"
)

func RespondCommand(request []string, KVStore *map[string]string) ([]string, error) {
	// Plan Change for this function: Return whatever it is supposed to be writing
	// So all this function does it look at the first value
	// Call the function that is equipped to handle the thing I want it to handle
	// And return whatever response I want the connection to give back to the client
	command := strings.ToLower(request[0])
	length := len(request)
	switch command {
	case "echo":
		if length != 2 {
			return nil, errors.New("Invalid Length")
		}
		responseString := request[1]
		return []string{responseString}, nil
	case "get":
		if length > 2 {
			return nil, errors.New("get must only have one value")
		}
		// GET will return what it wants to write
		return query_handling.Get(request[1], KVStore) // Change here
	case "set":
		px := false
		timer := 0
		if length == 5 {
			if strings.ToLower(request[3]) != "px" {
				return nil, errors.New("invalid request")
			} else {
				tempTimer, conversionError := strconv.Atoi(request[4])
				if conversionError != nil {
					return nil, errors.New("invalid request")
				} else {
					timer = tempTimer
					px = true
				}
			}
		} else if length != 2 {
			return nil, errors.New("invalid number of arguments")
		}
		query_handling.Set(request[1], request[2], px, timer) // Change Here

	default:
		return errors.New("can't understand the command") // Change Here
	}
	return nil // Change Here
}
