package response

import (
	"errors"
	"fmt"
	"github.com/keyur/redis2/internal/definitions"
	"github.com/keyur/redis2/internal/query_handling"
	"strconv"
	"strings"
)

func RespondCommand(request []string, requestContext *definitions.RequestContext) error {
	command := strings.ToLower(request[0])
	length := len(request)
	switch command {
	case "echo":
		responseString := request[1] + "\n"
		response := []byte(responseString)
		if requestContext.Connection != nil {
			_, err := requestContext.Connection.Write(response)
			if err != nil {
				return err
			}
		} else {
			return nil
		}
	case "get":
		if length > 2 {
			return errors.New("get must only have one value")
		}
		err := query_handling.Get(request[1], requestContext)
		if err != nil {
			fmt.Println("Error Handling GET ", err)
		}
		return err

	case "set":
		px := false
		timer := 0
		if length == 5 {
			if strings.ToLower(request[3]) != "px" {
				return errors.New(fmt.Sprintf("invalid flag: %v", request[3]))
			} else {
				tempTimer, conversionError := strconv.Atoi(request[4])
				if conversionError != nil {
					return errors.New(fmt.Sprintf("Not a valid integer %v", request[4]))
				} else {
					timer = tempTimer
					px = true
				}
			}
		} else if length != 2 {
			return errors.New("invalid number of arguments")
		}
		query_handling.Set(request[1], request[2], px, timer, requestContext)

	default:
		return errors.New("can't understand the command")
	}
	return nil
}
