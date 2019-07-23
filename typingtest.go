package main

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	term "github.com/nsf/termbox-go"
	"time"
)

var exampleTest = `Lorem Ipsum is simply dummy text.`

func reset(textUpdateChannel chan string, currSpeedChannel chan float64, prevSpeed float64) float64 {
	newText := <-textUpdateChannel
	term.Sync()
	fmt.Println("Typing...:", newText)
	select {
	case currSpeed := <-currSpeedChannel:
		fmt.Println("Current Speed:", currSpeed)
		return currSpeed
	default:
		fmt.Println("Current Speed:", prevSpeed)
		return prevSpeed

	}
}

func sendErrorStroke(keyStrokeChannel chan int) {
	keyStrokeChannel <- -1
}

func sendValidStroke(keyStrokeChannel chan int, stroke int) {
	keyStrokeChannel <- stroke
}

func getchar(keyStrokeChannel chan int, textUpdateChannel chan string, currSpeedChannel chan float64, wordCountTerminator chan bool) {
	err := term.Init()
	defer term.Close()
	if err != nil {
		panic(err)
	}
	fmt.Println("Type:", exampleTest)
	var prevSpeed float64
	prevSpeed = 0.0
keyPressListenerLoop:
	for {
		select {
		case <-wordCountTerminator:
			break keyPressListenerLoop
		default:
			ev := term.PollEvent()
			switch ev.Type {
			case term.EventKey:
				switch ev.Key {
				case term.KeyEsc:
					break keyPressListenerLoop
				case term.KeyF1:
					fallthrough
				case term.KeyF2:
					fallthrough
				case term.KeyF3:
					fallthrough
				case term.KeyF4:
					fallthrough
				case term.KeyF5:
					fallthrough
				case term.KeyF6:
					fallthrough
				case term.KeyF7:
					fallthrough
				case term.KeyF8:
					fallthrough
				case term.KeyF9:
					fallthrough
				case term.KeyF10:
					fallthrough
				case term.KeyF11:
					fallthrough
				case term.KeyF12:
					fallthrough
				case term.KeyInsert:
					fallthrough
				case term.KeyHome:
					fallthrough
				case term.KeyEnd:
					fallthrough
				case term.KeyPgup:
					fallthrough
				case term.KeyPgdn:
					fallthrough
				case term.KeyArrowUp:
					fallthrough
				case term.KeyArrowDown:
					fallthrough
				case term.KeyArrowLeft:
					fallthrough
				case term.KeyArrowRight:
					fallthrough
				case term.KeyTab:
					sendErrorStroke(keyStrokeChannel)
					prevSpeed = reset(textUpdateChannel, currSpeedChannel, prevSpeed)
				case term.KeyBackspace:
					fallthrough
				case term.KeyBackspace2:
					fallthrough
				case term.KeySpace:
					fallthrough
				case term.KeyEnter:
					sendValidStroke(keyStrokeChannel, int(ev.Key))
					prevSpeed = reset(textUpdateChannel, currSpeedChannel, prevSpeed)
				default:
					sendValidStroke(keyStrokeChannel, int(ev.Ch))
					prevSpeed = reset(textUpdateChannel, currSpeedChannel, prevSpeed)

				}
			case term.EventError:
				panic(ev.Err)
			}
		}
	}
	fmt.Println("getChar: Ending...")
}

func textChecker(keyStrokeChannel chan int, textUpdateChannel chan string, wordCountChannel chan bool, wordCountTerminator chan bool, getCharTerminator chan bool) {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	chars := []rune(exampleTest)
	var curr bytes.Buffer
	incorrectChars := 0
	incorrectString := ""

	for idx, nextChar := range chars {

		for {
			keyStroke := <-keyStrokeChannel
			fmt.Println("nextChar:", nextChar)
			fmt.Printf("nextChar:%c\n", rune(nextChar))
			fmt.Println("idx:", idx)
			fmt.Println("keyStroke:", keyStroke)
			fmt.Println("IncorrectChars:", incorrectChars)
			fmt.Println("curr:", curr.String())
			var rest []rune

			var nextCorrectChar rune
			if incorrectChars > 0 {
				nextCorrectChar = 127
			} else {
				nextCorrectChar = nextChar
			}

			fmt.Println("nextCorrectChar:", nextCorrectChar)
			if int(nextCorrectChar) == int(keyStroke) {
				if incorrectChars > 0 {
					incorrectChars--
					rest = chars[idx:]
				} else {
					// do nothing
					// maybe do something
					curr.WriteRune(nextCorrectChar)
					rest = chars[idx+1:]
				}
			} else {
				incorrectChars++
				rest = chars[idx:]
				nextCorrectChar = 127
			}

			if incorrectChars > 0 {
				incorrectString = fmt.Sprintf("You have entered %d incorrect characters. Backspace to get it all and then enter the correct character '%c'\n", incorrectChars, rune(nextChar))
			} else {
				incorrectString = ""
			}

			update := fmt.Sprintf("%s%s%s", incorrectString, green(curr.String()), red(string(rest)))
			textUpdateChannel <- update

			if int(nextChar) == int(keyStroke) && incorrectChars == 0 {
				if nextChar == 32 {
					wordCountChannel <- true
				}
				break
			}
		}
	}
	fmt.Println("textChecker: Ending...")
	wordCountTerminator <- true
	getCharTerminator <- true
}

func wordCounter(wordCountChannel chan bool, currSpeedChannel chan float64, wordCountTerminator chan bool, result chan float64) {
	totalWords := 0
	<-wordCountChannel
	totalWords++
	var speed float64
	start := time.Now()
outer:
	for {
		select {
		case <-wordCountChannel:
			totalWords++
			curr := time.Now()
			elapsed := curr.Sub(start).Minutes()
			speed = float64(totalWords) / elapsed
			currSpeedChannel <- speed
		case <-wordCountTerminator:
			result <- speed
			break outer

		}
	}
	fmt.Println("wordCounter: Ending.....")
}

func main() {
	keyStrokeChannel := make(chan int)
	textUpdateChannel := make(chan string)
	wordCountChannel := make(chan bool)
	currSpeedChannel := make(chan float64)
	wordCountTerminator := make(chan bool)
	getCharTerminator := make(chan bool)
	resultChannel := make(chan float64)
	go textChecker(keyStrokeChannel, textUpdateChannel, wordCountChannel, wordCountTerminator, getCharTerminator)
	go wordCounter(wordCountChannel, currSpeedChannel, wordCountTerminator, resultChannel)
	getchar(keyStrokeChannel, textUpdateChannel, currSpeedChannel, getCharTerminator)
	speed := <-resultChannel
	fmt.Println("Result:", speed)
}
