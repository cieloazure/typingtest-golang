package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func main() {
	var exampleTest = `
Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.
`
	f := func(c rune) bool {
		return unicode.IsSpace(c)
	}
	words := strings.FieldsFunc(exampleTest, f)
	reader := bufio.NewReader(os.Stdin)
	for _, word := range words {
		fmt.Printf("Enter: %s\n", word)
		for {
			message, _ := reader.ReadString('\n')
			fmt.Printf("Entered: %s\n", message)
			message = strings.TrimRight(message, "\n")
			if message == word {
				break
			}
		}
	}
	//message, _ := reader.ReadString('\n')
}
