package lexer

import (
	"bufio"
	"io"
	"log"
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
	justSkippedNewline bool
}

func NewLexer(reader *bufio.Reader) *Lexer {
	l := new(Lexer)

	l.currentPosition = Position{
		line:   1,
		column: 1,
	}

	l.reader = reader

	return l
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

// Function to read a character from the lexer's reader.
// By default, skips whitespace and automatically handles newlines.
// Use Lexer.readRuneDefault if you want  ashorthand way of calling Lexer.readRune(true, true)
// Returns (None, nil) when EOF
func (l *Lexer) readRune(skipWhitespace bool, autoHandleNewlines bool) (utils.Optional[rune], error) {
	r, _, err := l.reader.ReadRune()

	// Check for EOF
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
	case r == '\n' || r == '\r' && autoHandleNewlines:
		l.currentPosition.line++
		l.currentPosition.column = 1

		return l.readRune(skipWhitespace, autoHandleNewlines)
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

func (l *Lexer) NextToken() (utils.Optional[Token], error) {
	maybeRune, err := l.readRuneDefault()

	// Check for EOF
	if err == io.EOF {
		return utils.NoneOptional[Token](), nil
	} else if err != nil {
		return utils.NoneOptional[Token](), err
	}

	r, err := maybeRune.Value()

	if err != nil {
		return utils.NoneOptional[Token](), err
	}

	switch {
	case strings.ContainsRune(string(TokenOperatorGroup), r):
		return utils.SomeOptional(Token{
			characters: string(r),
			group:      &TokenOperatorGroup,
		}), nil
	case strings.ContainsRune(string(TokenGroupingGroup), r):
		return utils.SomeOptional(Token{
			characters: string(r),
			group:      &TokenGroupingGroup,
		}), nil
	case strings.ContainsRune(string(TokenSeparatorGroup), r):
		return utils.SomeOptional(Token{
			characters: string(r),
			group:      &TokenSeparatorGroup,
		}), nil
	case r == '#':
		_, _, err := l.reader.ReadLine()

		if err != nil {
			return utils.NoneOptional[Token](), err
		}

		return l.NextToken()
	case r == '"':
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

		return utils.SomeOptional(Token{
			characters: str,
			group:      &TokenStringGroup,
		}), nil
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

		return utils.SomeOptional[Token](Token{
			characters: strconv.FormatInt(num, 10),
			group:      &TokenIntegerGroup,
		}), nil
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

		return utils.SomeOptional[Token](Token{
			characters: strconv.FormatInt(num, 10),
			group:      &TokenIntegerGroup,
		}), nil
	}

	if IsValidIdentStart(r) {
		ident := string(r)
		mp, err := l.readRune(false, false)

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

		return utils.SomeOptional(Token{
			characters: ident,
			group:      &TokenIdentifierGroup,
		}), nil
	}

	return utils.NoneOptional[Token](), nil
}
