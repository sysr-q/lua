package lexer

import (
	"strings"
	"unicode"

	"github.com/sysr-q/lua/token"
)

const (
	eof        = -1
	alphabet   = "abcdefghijklymnopqrstuvwxyzABCDEFGHIJKLYMNOPQRSTUVWXYZ"
	digits     = "01234567890"
	underscore = "_"
)

// lexText is the starting point for a new Lua lexer.
func lexText(l *lexer) stateFn {
	l.skipSpace()

	if r := l.peek(); unicode.IsDigit(r) {
		return lexNumberLiteral
	} else if r == '.' {
		l.next() // must advance because we only peeked before
		r2 := l.peek()
		l.backup()
		if unicode.IsDigit(r2) {
			return lexNumberLiteral
		}
	}

	if strings.HasPrefix(l.input[l.pos:], "--[") {
		return lexBlockComment
	}

	if strings.HasPrefix(l.input[l.pos:], "--") {
		return lexLineComment
	}

	if l.peek() == eof {
		l.emit(token.EOF)
		return nil
	}

	if l.peek() == '\'' {
		return lexSingleQuotedStringLiteral
	}

	if l.peek() == '"' {
		return lexDoubleQuotedStringLiteral
	}

	if l.peek() == '[' {
		return lexDoubleBracketStringLiteral
	}

	for _, tokenString := range token.TokenList {
		t := token.TokenMap[tokenString]
		potentialToken := l.input[l.pos:]
		if len(potentialToken) > len(tokenString) {
			potentialToken = potentialToken[:len(tokenString)]
		}
		if strings.HasPrefix(strings.ToLower(potentialToken), tokenString) {
			l.pos += token.Pos(len(tokenString))
			if IsKeyword(t, tokenString) && l.accept(alphabet+underscore+digits) {
				l.backup() // to account for the character consumed by accept
				l.pos -= token.Pos(len(tokenString))
				break
			}
			l.emit(t)
			return lexText
		}
	}

	l.acceptRun(alphabet + underscore + digits)
	l.emit(token.Identifier)
	return lexText
}

func lexNumberLiteral(l *lexer) stateFn {
	if l.accept("0") && l.accept("xX") {
		// hexadecimal?
		l.acceptRun(digits + "abcdefABCDEF")
		l.emit(token.NumberLiteral)
		return lexText
	}

	// is decimal?
	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}

	// TODO: 3.1e-1 etc
	if l.accept("eE") {
		l.acceptRun(digits)
	}

	l.emit(token.NumberLiteral)
	return lexText
}

func lexBlockComment(l *lexer) stateFn {
	// Eat the --[
	l.acceptRun("-")
	l.next()

	equalsRun := l.acceptRun("=")
	if !l.accept("[") {
		return l.errorf("didn't find another opening bracket")
	}

	for {
		switch l.next() {
		case '\\':
			l.next()
			continue
		case ']':
			if l.acceptRun("=") != equalsRun {
				// Not the same amount as when we opened this string.
				continue
			}
			if !l.accept("]") {
				return l.errorf("didn't find closing bracket")
			}
			l.ignore()
			return lexText
		}
	}
}

func lexLineComment(l *lexer) stateFn {
	lineLength := strings.Index(l.input[l.pos:], "\n")
	if lineLength == 0 {
		// this is the last line, so lex until the end
		lineLength = len(l.input[l.pos:])
	}

	l.pos += token.Pos(lineLength)
	l.ignore()
	return lexText
}

func lexSingleQuotedStringLiteral(l *lexer) stateFn {
	l.next()
	for {
		switch l.next() {
		case '\\':
			l.next()
			continue
		case '\'':
			l.emit(token.StringLiteral)
			return lexText
		}
	}
}

func lexDoubleQuotedStringLiteral(l *lexer) stateFn {
	l.next()
	for {
		switch l.next() {
		case '\\':
			l.next()
			continue
		case '"':
			l.emit(token.StringLiteral)
			return lexText
		}
	}
}

func lexDoubleBracketStringLiteral(l *lexer) stateFn {
	// Eat the opening '['
	l.next()
	equalsRun := l.acceptRun("=")
	if !l.accept("[") {
		return l.errorf("didn't find another opening bracket")
	}

	// Potentially eat up a newline.
	l.accept("\n")

	for {
		switch l.next() {
		case '\\':
			l.next()
			continue
		case ']':
			if l.acceptRun("=") != equalsRun {
				// Not the same amount as when we opened this string.
				continue
			}
			if !l.accept("]") {
				return l.errorf("didn't find closing bracket")
			}
			l.emit(token.StringLiteral)
			return lexText
		}
	}
}
