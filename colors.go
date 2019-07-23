package main

import (
	"fmt"
	"github.com/fatih/color"
	"strings"
	"time"
	"unicode"
)

func main() {
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	var exampleTest = `
Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.
`

	f := func(c rune) bool {
		return unicode.IsSpace(c)
	}
	words := strings.FieldsFunc(exampleTest, f)

	mystring := fmt.Sprintf("%s", yellow(strings.Join(words, " ")))
	fmt.Printf("%s", mystring)
	time.Sleep(3 * time.Second)
	first := words[0]
	rest := words[1:]
	fmt.Printf("\r%s%s%s", red(first), " ", yellow(strings.Join(rest, " ")))
}
