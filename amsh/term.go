package amsh

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	clear          = "\033[2J"
	moveTopLeft    = "\033[0;0H"
	blackText      = "\033[38;5;0m"
	purpleLabel    = "\033[48;5;57m"
	purpleOnBlack  = "\033[48;5;0m\033[38;5;57m"
	purpleOnYellow = "\033[48;5;220m\033[38;5;57m"
	yellowLabel    = "\033[48;5;220m"
	yellowOnBlack  = "\033[48;5;0m\033[38;5;220m"
	resetColor     = "\033[0m"
)

/*
Term is a wrapper around the data and functionality of the terminal.
*/
type Term struct {
	l   *Lexer
	pmt string
	wd  string
}

/*
NewTerm returns a reference to a new Term instance.
*/
func NewTerm() *Term {
	return &Term{
		l:   NewLexer(),
		pmt: "",
	}
}

/*
Run the instantiated terminal.
*/
func (t *Term) Run() {
	fmt.Printf(clear + moveTopLeft)

	r := bufio.NewReader(os.Stdin)

	go t.l.Run()

	t.wd, _ = os.Getwd()

	for {
		t.pmt += fmt.Sprintf("%s %s %s\ue0b0", purpleLabel, t.wd, purpleOnYellow)
		t.pmt += fmt.Sprintf("%s %s %s\ue0b0", yellowLabel, "master", yellowOnBlack)
		t.pmt += fmt.Sprintf("\n%s -> ", resetColor)

		fmt.Printf(t.pmt)

		// Read the raw string input and remove the newline.
		c, _ := r.ReadString('\n')
		c = strings.Replace(c, "\n", "", -1)

		// Send the raw string input to the Lexer over a channel.
		t.l.LCh <- c

		// Flush the prompt.
		t.pmt = ""
	}
}
