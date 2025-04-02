package writer

import (
	"fmt"
)

func ArrayResponseWriter(response []string) []byte {
	lengthOfElements := len(response)
	writeString := fmt.Sprintf("*%d\r\n", lengthOfElements)
	for _, word := range response {
		size := len(word)
		writeString += fmt.Sprintf("$\r\n%d\r\n%s\r\n", size, word)
	}

	return []byte(writeString)
}
