package lexer

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/sysr-q/lua/token"
)

type stateFn func(*lexer) stateFn

// lexer holds the state of the scanner.
// Based off of github.com/stephens2424/php and Rob Pike's talk on lexers:
// https://www.youtube.com/watch?v=HxaD_trXwRE
type lexer struct {
	name    string          // the name of the input; used only for error reports
	input   string          // the string being scanned
	state   stateFn         // the next lexing function to enter
	pos     token.Pos       // current position in the input
	start   token.Pos       // start position of this item
	width   token.Pos       // width of last rune read from input
	lastPos token.Pos       // position of most recent item returned by Next
	items   chan token.Item // channel of scanned items
}

// lex creates a new scanner for the input string.
func NewLexer(name, input string) token.Stream {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan token.Item),
	}
	go l.run()
	return l
}

// run runs the state machine for the lexer.
func (l *lexer) run() {
	for l.state = lexText; l.state != nil; {
		l.state = l.state(l)
	}
}

// next returns the next rune in the input.
func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = token.Pos(w)
	l.pos += l.width
	return r
}

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
}

// emit passes an item back to the client.
func (l *lexer) emit(t token.Token) {
	l.items <- token.Item{
		Typ:   t,
		Begin: l.start,
		Val:   l.input[l.start:l.pos],
	}
	l.start = l.pos
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.start = l.pos
}

// skipSpace skips over a run of whitespace: [ \t]*
func (l *lexer) skipSpace() {
	r := l.next()
	for isSpace(r) {
		r = l.next()
	}
	l.backup()
	l.ignore()
}

// accept consumes the next rune if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) int {
	runLength := 0
	for strings.IndexRune(valid, l.next()) >= 0 {
		runLength += 1
	}
	l.backup()
	return runLength
}

// lineNumber reports which line we're on, based on the position of
// the previous item returned by Next. Doing it this way
// means we don't have to worry about peek double counting.
func (l *lexer) lineNumber() int {
	return 1 + strings.Count(l.input[:l.lastPos], "\n")
}

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.Next.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- token.Item{
		Typ:   token.Error,
		Begin: l.start,
		Val:   fmt.Sprintf(format, args...),
	}
	return nil
}

// Next returns the next item from the input.
func (l *lexer) Next() token.Item {
	item := <-l.items
	l.lastPos = item.Pos()
	return item
}

/// Stealing code from stephens2424/php: awesome.

// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	return unicode.IsSpace(r)
}

func IsKeyword(i token.Token, tokenString string) bool {
	_, ok := keywordMap[i]
	return ok && !isNonAlphaOperator(tokenString)
}

var nonalpha *regexp.Regexp

func init() {
	nonalpha = regexp.MustCompile(`^[^a-zA-Z0-9]*$`)
}

func isNonAlphaOperator(s string) bool {
	return nonalpha.MatchString(s)
}

// keywordMap lists all keywords that should be ignored as a prefix to a longer
// identifier.
var keywordMap = map[token.Token]bool{}

func init() {
	re := regexp.MustCompile("^[a-zA-Z]+")
	for keyword, t := range token.TokenMap {
		if re.MatchString(keyword) {
			keywordMap[t] = true
		}
	}
}
