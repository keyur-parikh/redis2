package parser

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

func ValidCheckParsing(readBufferPtr *[]byte) ([]string, int, error) {
	for {
		readBuffer := *readBufferPtr
		// Do the initial check
		fmt.Println("Inside valid check parsing")
		if len(readBuffer) == 0 || readBuffer[0] != '*' {
			return make([]string, 0), 2, errors.New("expecting * at the beginning")
		}
		// Find the first set of /r/n
		newline := bytes.Index(readBuffer, []byte("\r\n"))
		if newline == -1 {
			return make([]string, 0), 1, errors.New("couldnt find \\r\\n")
		}
		// Parse array count
		elementStr := string(readBuffer[1:newline])
		//Convert into an integer
		numberOfElements, err := strconv.Atoi(elementStr)
		if err != nil {
			return make([]string, 0), 2, errors.New("couldnt find the number of array elements")
		}
		request := make([]string, numberOfElements)
		offset := newline + 2
		for i := 0; i < numberOfElements; i++ {
			if offset > len(readBuffer) {
				return nil, 1, errors.New("couldnt complete parsing")
			}
			if readBuffer[offset] != '$' {
				return nil, 2, errors.New("all array elements must start with $")
			}
			// Find the next set of \r\n
			lengthEndRelative := bytes.Index(readBuffer[offset:], []byte("\r\n"))
			if lengthEndRelative == -1 {
				return nil, 1, errors.New("couldnt find the size of index")
			}
			lengthEnd := offset + lengthEndRelative
			elementSizeStr := string(readBuffer[offset+1 : lengthEnd])
			size_of_element, err := strconv.Atoi(elementSizeStr)
			if err != nil {
				return nil, 2, errors.New("couldn't find the number of array elements")
			}
			offset = lengthEnd + 2 // offset still before value of string
			if offset+size_of_element+2 > len(readBuffer) {
				return nil, 1, errors.New("couldnt find the element")
			}

			if readBuffer[offset+size_of_element] != '\r' || readBuffer[offset+size_of_element+1] != '\n' {
				return nil, 2, errors.New("Incorrect Sized Message or missing CRLF")
			}
			request[i] = string(readBuffer[offset : offset+size_of_element])
			offset += size_of_element + 2
		}
		fmt.Println("All Passed Successfully")
		fmt.Println("Messages", request)
		// Ideally after testing
		if offset == len(readBuffer) {
			fmt.Println("error is here in the read buffer zone")
			*readBufferPtr = readBuffer[:0]
			return request, 0, nil // fully consumed
		} else {
			fmt.Println("Its in the else after the offset thing")
			*readBufferPtr = readBuffer[offset:]
			continue // keep remaining
		}
	}
	// readBuffer = readBuffer[offset + 1: ]
}
