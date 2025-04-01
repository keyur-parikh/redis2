package parser

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/keyur-parikh/redis2/internal/definitions"
	"github.com/keyur-parikh/redis2/internal/response"
	"strconv"
)

func ValidCheckParsing(readBufferPtr *[]byte, requestContext *definitions.RequestContext) (int, error) {
	for {
		readBuffer := *readBufferPtr
		// Do the initial check
		if len(readBuffer) == 0 || readBuffer[0] != '*' {
			return 2, errors.New("expecting * at the beginning")
		}
		// Find the first set of /r/n
		newline := bytes.Index(readBuffer, []byte("\r\n"))
		if newline == -1 {
			return 1, errors.New("couldnt find \\r\\n")
		}
		// Parse array count
		elementStr := string(readBuffer[1:newline])
		//Convert into an integer
		numberOfElements, err := strconv.Atoi(elementStr)
		if err != nil {
			return 2, errors.New("couldnt find the number of array elements")
		}
		request := make([]string, numberOfElements)
		offset := newline + 2
		for i := 0; i < numberOfElements; i++ {
			if offset > len(readBuffer) {
				return 1, errors.New("couldnt complete parsing")
			}
			if readBuffer[offset] != '$' {
				return 2, errors.New("all array elements must start with $")
			}
			// Find the next set of \r\n
			lengthEndRelative := bytes.Index(readBuffer[offset:], []byte("\r\n"))
			if lengthEndRelative == -1 {
				return 1, errors.New("couldnt find the size of index")
			}
			lengthEnd := offset + lengthEndRelative
			elementSizeStr := string(readBuffer[offset+1 : lengthEnd])
			size_of_element, err := strconv.Atoi(elementSizeStr)
			if err != nil {
				return 2, errors.New("couldn't find the number of array elements")
			}
			offset = lengthEnd + 2 // offset still before value of string
			if offset+size_of_element+2 > len(readBuffer) {
				return 1, errors.New("couldnt find the element")
			}

			if readBuffer[offset+size_of_element] != '\r' || readBuffer[offset+size_of_element+1] != '\n' {
				return 2, errors.New("Incorrect Sized Message or missing CRLF")
			}
			request[i] = string(readBuffer[offset : offset+size_of_element])
			offset += size_of_element + 2
		}
		fmt.Println("All Passed Successfully")
		fmt.Println("Messages", request)
		if requestContext.Connection != nil {
			err = response.RespondCommand(request, requestContext)
		}
		if err != nil {
			return 2, err
		}
		// Ideally after testing
		if offset == len(readBuffer) {
			*readBufferPtr = readBuffer[:0]
			return 0, nil // fully consumed
		} else {
			*readBufferPtr = readBuffer[offset:]
			continue // keep remaining
		}
	}
	// readBuffer = readBuffer[offset + 1: ]
}
