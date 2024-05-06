package main

import (
	"fmt"
	"log"
	"os"
	"sorbet"
	"strings"
	"time"

	"github.com/periaate/common"
	"github.com/periaate/meminfo"
	"golang.org/x/term"
)

var width int
var sb = strings.Builder{}
var k int

func main() {
	ticker := time.NewTicker(64 * time.Millisecond)
	defer ticker.Stop()
	defer showCursor()
	hideCursor()
	var i int
	var g int

	for range ticker.C {
		if i == 0 {
			w, _, err := term.GetSize(int(os.Stdout.Fd()))
			if err != nil {
				fmt.Println("Error getting terminal size:", err)
				return
			}
			if w != width {
				sb.WriteString("\033[2J")
				width = w
			}
		}

		i++
		i &= 15
		g++
		g &= 1
		if g == 0 {
			k++
		}

		render()
	}
}

func render() {
	currentTime := time.Now().Format("15:04:05")
	mi, err := meminfo.Get()
	if err != nil {
		log.Fatalln(err)
	}

	u := common.HumanizeBytes(1, mi.PhysUsed, 1, true)
	t := common.HumanizeBytes(1, mi.PhysTotal, 1, true)
	mem := fmt.Sprintf("%s / %s", u, t)
	sb.WriteString("\033[H")
	sb.WriteString(fmt.Sprintf("%s\n", pad(width, currentTime)))
	sb.WriteString(pad(width, mem))

	title := sorbet.GetTitle()

	animatedTitle := animateTitle(title, width, k)
	sb.WriteString("\033[J")
	sb.WriteString(fmt.Sprintf("\n%s", pad(width, animatedTitle)))

	fmt.Print(sb.String())
	sb.Reset()
}

func pad(width int, text string) string {
	// Calculate the left padding needed to center the text
	padding := (width - len(text)) / 2
	if padding < 0 {
		padding = 0
	}
	return fmt.Sprintf("%s%s", strings.Repeat(" ", padding), text)
}

// hideCursor hides the terminal cursor
func hideCursor() {
	fmt.Fprint(os.Stdout, "\033[?25l")
}

// showCursor shows the terminal cursor
func showCursor() {
	fmt.Fprint(os.Stdout, "\033[?25h")
}

func animateTitle(title string, width, index int) string {
	if l(title) <= width {
		return title
	}
	// Rotate through title
	index = index % l(title)
	if index+width < l(title) {
		return title[index : index+width]
	}
	// If the index goes beyond the title length, wrap around to the beginning
	part1 := title[index:]
	part2 := title[0 : width-l(part1)]
	return part1 + part2
}

func l(s string) int {
	return len(s)
}
