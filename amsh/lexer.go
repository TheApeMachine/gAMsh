package amsh

/*
Context is a state machine used for parsing raw string input.
*/
type Context int

const (
	// START indicates no defined state (yet).
	START Context = 0 + iota
	// ID indicates we have found a keyword, either built-in or external.
	ID
)

/*
Lexeme is a structured token created from the raw string input.
*/
type Lexeme struct {
	id string
}

/*
NewLexeme returns a reference to a new Lexeme instance.
*/
func NewLexeme() *Lexeme {
	return &Lexeme{
		id: "",
	}
}

/*
Lexer is a wrapper around the data and functionality of converting raw input to Lexemes.
*/
type Lexer struct {
	p      *Parser
	LCh    chan string
	delims []string
	ctx    Context
	buf    string
}

/*
NewLexer returns a reference to a new Lexer instance.
*/
func NewLexer() *Lexer {
	return &Lexer{
		p:      NewParser(),
		LCh:    make(chan string),
		delims: []string{" "},
		ctx:    START,
		buf:    "",
	}
}

/*
Run as a goroutine and accept incoming raw input on a channel to be converted to Lexemes, and passed to the parser.
*/
func (l *Lexer) Run() {
	// Start the parser as a goroutine.
	go l.p.Run()

	for {
		select {
		case cmd := <-l.LCh:
			l.convert(cmd)
		}
	}
}

func (l *Lexer) convert(cmd string) {
	le := NewLexeme()

	for _, r := range cmd {
		// Convert the rune to a string.
		c := string(r)

		switch l.ctx {
		case START:
			if l.sis(c, l.delims) {
				// Assign the keyword to the Lexeme and clear the transient buffer.
				le.id = l.buf
				l.buf = ""

				// Set the new context.
				l.ctx = ID

				continue
			}

			// Keep adding the current character to the transient buffer.
			l.buf += c
		}
	}

	// Handle case if we don't have any parameters.
	if le.id == "" {
		// Assign the keyword to the Lexeme and clear the transient buffer.
		le.id = l.buf
		l.buf = ""
	}

	// We have reached the end of the input, pass it to the parser.
	l.p.LCh <- le
}

/*
sis String In Slice check.
*/
func (l *Lexer) sis(str string, sl []string) bool {
	for _, v := range sl {
		if v == str {
			return true
		}
	}

	return false
}
