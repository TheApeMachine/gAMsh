package amsh

import (
	"os"
	"os/exec"
)

/*
Parser is a wrapper around the data and functionality of converting Lexemes to an Abstract Syntax Tree.
*/
type Parser struct {
	LCh chan *Lexeme
}

/*
NewParser returns a reference to a new Parser instance.
*/
func NewParser() *Parser {
	return &Parser{
		LCh: make(chan *Lexeme),
	}
}

/*
Run as a goroutine and accept incoming Lexemes on a channel to be converted to an AST.
*/
func (p *Parser) Run() {
	for {
		select {
		case le := <-p.LCh:
			p.convert(le)
		}
	}
}

func (p *Parser) convert(le *Lexeme) {
	c := exec.Command(le.id)
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	c.Run()
}
