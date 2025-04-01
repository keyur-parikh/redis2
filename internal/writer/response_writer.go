package writer

import (
	"errors"
	"fmt"
	"github.com/keyur/redis2/internal/definitions"
)

func ArrayResponseWriter(response []string, requestContext *definitions.RequestContext) error {
	lengthOfElements := len(response)
	writeString := fmt.Sprintf("*%d\\r\\n", lengthOfElements)
	for _, word := range response {
		size := len(word)
		writeString += fmt.Sprintf("$\\r\\n%d\\r\\n%s\\r\\n", size, word)
	}
	if requestContext != nil {
		_, err := requestContext.Connection.Write([]byte(writeString))
		if err != nil {
			return err
		}
	} else {
		fmt.Println(writeString)
		return errors.New(writeString)
	}
	return nil
}

func InvalidResponseWriter(requestContext *definitions.RequestContext) error {
	_, err := requestContext.Connection.Write([]byte("$-1\\r\\n"))
	if err != nil {
		return err
	}
	return nil
}

func SuccessResponseWriter(requestContext *definitions.RequestContext) error {
	_, err := requestContext.Connection.Write([]byte("+OK\\r\\n"))
	if err != nil {
		return err
	}
	return nil
}
