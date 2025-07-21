package lexer

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"ljpprojects.org/sqopl/utils"
)

type LexerNumericalBase uint8

const (
	Base2LexerNumericalBase LexerNumericalBase = iota
	Base8LexerNumericalBase
	Base10LexerNumericalBase
	Base16LexerNumericalBase
)

type Position struct {
	line   uint32
	column uint32
}

type Lexer struct {
	currentPosition    Position
	reader             *bufio.Reader
	file               *os.File
	justSkippedNewline bool
	bytesRead          uint64
}

func NewLexer(file *os.File) *Lexer {
	l := new(Lexer)

	l.currentPosition = Position{
		line:   1,
		column: 1,
	}

	l.file = file
	l.reader = bufio.NewReader(file)
	l.justSkippedNewline = false
	l.bytesRead = 0

	return l
}

func (l *Lexer) Clone() (*Lexer, error) {
	file, err := os.Open(l.file.Name())

	if err != nil {
		return nil, err
	}

	lexer := NewLexer(file)

	lexer.currentPosition = l.currentPosition
	lexer.justSkippedNewline = l.justSkippedNewline
	lexer.reader.Discard(int(l.bytesRead))

	return lexer, nil
}

func IsValidIdentStart(r rune) bool {
	allowedChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_"

	return strings.ContainsRune(allowedChars, r)
}

func IsValidIdentPart(r rune) bool {
	allowedChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

	return strings.ContainsRune(allowedChars, r)
}

func IsValidStringPart(r rune) bool {
	forbiddenChars := "\"\t\n\r"

	return !strings.ContainsRune(forbiddenChars, r)
}

func IsValidNumberPart(r rune, base LexerNumericalBase) bool {
	allowedChars := map[LexerNumericalBase]string{
		Base2LexerNumericalBase:  "01",
		Base8LexerNumericalBase:  "01234567",
		Base10LexerNumericalBase: "0123456789",
		Base16LexerNumericalBase: "0123456789ABCDEFabcdef",
	}

	return strings.ContainsRune(allowedChars[base], r)
}

func (l Lexer) CurrentPos() Position {
	return l.currentPosition
}

// Function to read a character from the lexer's reader.
// By default, skips whitespace and automatically handles newlines.
// Use Lexer.readRuneDefault if you want  ashorthand way of calling Lexer.readRune(true, true)
// Returns (None, nil) when EOF
func (l *Lexer) readRune(skipWhitespace bool, autoHandleNewlines bool) (utils.Optional[rune], error) {
	r, s, err := l.reader.ReadRune()
	l.bytesRead += uint64(s)

	if err == io.EOF {
		return utils.NoneOptional[rune](), nil
	} else if err != nil {
		return utils.NoneOptional[rune](), err
	}

	l.justSkippedNewline = false

	switch {
	case r == ' ' && skipWhitespace:
		l.currentPosition.column++

		return l.readRune(skipWhitespace, autoHandleNewlines)
	case r == '\n' && autoHandleNewlines:
		l.currentPosition.line++
		l.currentPosition.column = 1

		log.Println("NEWLINE", l.currentPosition)

		return l.readRune(skipWhitespace, autoHandleNewlines)
	default:
		l.currentPosition.column++
	}

	return utils.SomeOptional(r), nil
}

func (l *Lexer) peekBytes(n int) (utils.Optional[[]byte], error) {
	a, err := l.reader.Peek(n)

	// Check for EOF
	if err == io.EOF {
		return utils.NoneOptional[[]byte](), nil
	} else if err != nil {
		return utils.NoneOptional[[]byte](), err
	}

	return utils.SomeOptional(a), nil
}

func (l *Lexer) readRuneDefault() (utils.Optional[rune], error) {
	return l.readRune(true, true)
}

func (l *Lexer) PeekToken() (utils.Optional[Token], error) {
	lexer, err := l.Clone()

	if err != nil {
		return utils.NoneOptional[Token](), err
	}

	return lexer.NextToken()
}

func (l *Lexer) NextToken() (utils.Optional[Token], error) {
	maybeRune, err := l.readRuneDefault()

	// Check for EOF
	if err == io.EOF {
		return utils.NoneOptional[Token](), nil
	} else if err != nil {
		return utils.NoneOptional[Token](), err
	}

	r, err := maybeRune.Value()

	startpos := Position{
		line:   l.currentPosition.line,
		column: l.currentPosition.column - 1,
	}

	if err != nil {
		return utils.NoneOptional[Token](), nil
	}

	switch {
	case strings.ContainsRune(string(TokenOperatorGroup), r):
		return utils.SomeOptional(InitToken(&TokenOperatorGroup, string(r), InitLocation(startpos, l.currentPosition))), nil
	case strings.ContainsRune(string(TokenGroupingGroup), r):
		return utils.SomeOptional(InitToken(&TokenGroupingGroup, string(r), InitLocation(startpos, l.currentPosition))), nil
	case strings.ContainsRune(string(TokenSeparatorGroup), r):
		return utils.SomeOptional(InitToken(&TokenSeparatorGroup, string(r), InitLocation(startpos, l.currentPosition))), nil
	case r == '#':
		line, _, err := l.reader.ReadLine()
		l.bytesRead += uint64(len(line))

		if err != nil {
			return utils.NoneOptional[Token](), err
		}

		l.currentPosition.line++
		l.currentPosition.column = 1

		return l.NextToken()
	case r == '"':
		str := ""
		mp, err := l.readRune(false, true)

		if err != nil {
			return utils.NoneOptional[Token](), err
		}

		p, err := mp.Value()

		if err != nil {
			return utils.NoneOptional[Token](), err
		}

		for {
			str += string(p)

			mb, err := l.peekBytes(1)

			if err != nil {
				return utils.NoneOptional[Token](), err
			}

			b, err := mb.Value()

			if err != nil {
				return utils.NoneOptional[Token](), err
			}

			if !IsValidStringPart(rune(b[0])) {
				break
			}

			mp, err := l.readRune(false, false)

			if err != nil {
				return utils.NoneOptional[Token](), err
			}

			p, err = mp.Value()

			if err != nil {
				return utils.NoneOptional[Token](), err
			}
		}

		if _, err := l.readRuneDefault(); err != nil {
			return utils.NoneOptional[Token](), err
		}

		return utils.SomeOptional(InitToken(&TokenStringGroup, str, InitLocation(startpos, l.currentPosition))), nil
	case r == '1' || r == '2' || r == '3' || r == '4' || r == '5' || r == '6' || r == '7' || r == '8' || r == '9':
		var num int64 = 0

		str := string(r)
		mp, err := l.readRune(false, false)

		if err != nil {
			return utils.NoneOptional[Token](), err
		}

		p, err := mp.Value()

		if err != nil {
			return utils.NoneOptional[Token](), err
		}

		for {
			if !IsValidNumberPart(p, Base10LexerNumericalBase) {
				break
			}

			str += string(p)

			mb, err := l.peekBytes(1)

			if err != nil {
				return utils.NoneOptional[Token](), err
			}

			b, err := mb.Value()

			if err != nil {
				return utils.NoneOptional[Token](), err
			}

			if !IsValidNumberPart(rune(b[0]), Base10LexerNumericalBase) {
				break
			}

			mp, err := l.readRune(false, false)

			if err != nil {
				return utils.NoneOptional[Token](), err
			}

			p, err = mp.Value()

			if err != nil {
				return utils.NoneOptional[Token](), err
			}
		}

		n, err := strconv.ParseInt(str, 10, 64)

		if err != nil {
			return utils.NoneOptional[Token](), err
		}

		num = n

		return utils.SomeOptional(InitToken(&TokenIntegerGroup, strconv.FormatInt(num, 10), InitLocation(startpos, l.currentPosition))), nil
	case r == '0':
		maybeRune, err := l.readRuneDefault()

		if err != nil {
			return utils.NoneOptional[Token](), err
		}

		r, err := maybeRune.Value()

		var num int64 = 0

		if err != nil {
			return utils.NoneOptional[Token](), err
		}

		switch r {
		case 'b':
			str := ""
			mp, err := l.readRune(false, false)

			if err != nil {
				return utils.NoneOptional[Token](), err
			}

			p, err := mp.Value()

			if err != nil {
				return utils.NoneOptional[Token](), err
			}

			for {
				str += string(p)

				mb, err := l.peekBytes(1)

				if err != nil {
					return utils.NoneOptional[Token](), err
				}

				b, err := mb.Value()

				if err != nil {
					return utils.NoneOptional[Token](), err
				}

				if !IsValidNumberPart(rune(b[0]), Base2LexerNumericalBase) {
					break
				}

				mp, err := l.readRune(false, false)

				if err != nil {
					return utils.NoneOptional[Token](), err
				}

				p, err = mp.Value()

				if err != nil {
					return utils.NoneOptional[Token](), err
				}
			}

			n, err := strconv.ParseInt(str, 2, 64)

			if err != nil {
				return utils.NoneOptional[Token](), err
			}

			num = n
		case 'o':
			str := ""
			mp, err := l.readRune(false, false)

			if err != nil {
				return utils.NoneOptional[Token](), err
			}

			p, err := mp.Value()

			if err != nil {
				return utils.NoneOptional[Token](), err
			}

			for {
				str += string(p)

				mb, err := l.peekBytes(1)

				if err != nil {
					return utils.NoneOptional[Token](), err
				}

				b, err := mb.Value()

				if err != nil {
					return utils.NoneOptional[Token](), err
				}

				if !IsValidNumberPart(rune(b[0]), Base8LexerNumericalBase) {
					break
				}

				mp, err := l.readRune(false, false)

				if err != nil {
					return utils.NoneOptional[Token](), err
				}

				p, err = mp.Value()

				if err != nil {
					return utils.NoneOptional[Token](), err
				}
			}

			n, err := strconv.ParseInt(str, 8, 64)

			if err != nil {
				return utils.NoneOptional[Token](), err
			}

			num = n
		case 'x':
			str := ""
			mp, err := l.readRune(false, false)

			if err != nil {
				return utils.NoneOptional[Token](), err
			}

			p, err := mp.Value()

			if err != nil {
				return utils.NoneOptional[Token](), err
			}

			for {
				str += string(p)

				mb, err := l.peekBytes(1)

				if err != nil {
					return utils.NoneOptional[Token](), err
				}

				b, err := mb.Value()

				if err != nil {
					return utils.NoneOptional[Token](), err
				}

				log.Println(string(b))

				if !IsValidNumberPart(rune(b[0]), Base16LexerNumericalBase) {
					break
				}

				mp, err := l.readRune(false, false)

				if err != nil {
					return utils.NoneOptional[Token](), err
				}

				p, err = mp.Value()

				if err != nil {
					return utils.NoneOptional[Token](), err
				}
			}

			n, err := strconv.ParseInt(str, 16, 64)

			if err != nil {
				return utils.NoneOptional[Token](), err
			}

			num = n
		}

		return utils.SomeOptional(InitToken(&TokenIntegerGroup, strconv.FormatInt(num, 10), InitLocation(startpos, l.currentPosition))), nil
	}

	if IsValidIdentStart(r) {
		ident := string(r)
		mp, err := l.readRune(false, true)

		if err != nil {
			return utils.NoneOptional[Token](), err
		}

		p, err := mp.Value()

		if err != nil {
			return utils.NoneOptional[Token](), err
		}

		for {
			if !IsValidIdentPart(p) {
				break
			}

			ident += string(p)

			mb, err := l.peekBytes(1)

			if err != nil {
				return utils.NoneOptional[Token](), err
			}

			b, err := mb.Value()

			if err != nil {
				return utils.NoneOptional[Token](), err
			}

			if !IsValidIdentPart(rune(b[0])) {
				break
			}

			mp, err := l.readRune(false, false)

			if err != nil {
				return utils.NoneOptional[Token](), err
			}

			p, err = mp.Value()

			if err != nil {
				return utils.NoneOptional[Token](), err
			}
		}

		return utils.SomeOptional(InitToken(&TokenIdentifierGroup, ident, InitLocation(startpos, l.currentPosition))), nil
	}

	return utils.NoneOptional[Token](), nil
}
