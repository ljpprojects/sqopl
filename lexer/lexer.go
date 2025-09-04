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

// The numerical base of a number used when lexing numbers.
// This is an enum.
type LexerNumericalBase uint8

const (
	// The binary base.
	Base2LexerNumericalBase LexerNumericalBase = iota

	// The octal base.
	Base8LexerNumericalBase

	// The decimal base.
	Base10LexerNumericalBase

	// The hexadecimal base.
	Base16LexerNumericalBase
)

// A structure that stores the position the lexer is currently at
type Position struct {
	Line   uint32
	Column uint32
}

// The lexer.
// Stores important state.
type Lexer struct {
	// The current Position of the lexer.
	currentPosition Position

	// The reader.
	reader *bufio.Reader

	// The file used to create the reader.
	// Used only when calling Lexer.Clone (to obtain a new, identical bufio.Reader)
	file *os.File

	// Whether or not the lexer just automatically skipped a newline.
	justSkippedNewline bool

	// How many bytes have been read by the Lexer.
	// Used only when calling Lexer.Clone (to reach the same position in the bufio.Reader)
	bytesRead uint64
}

// Allocates a new Lexer and assigns its properties to the appropriate values/defaults.
func NewLexer(file *os.File) *Lexer {
	l := new(Lexer)

	// Assign the starting position.
	l.currentPosition = Position{
		Line:   1,
		Column: 1,
	}

	// Assign the file.
	l.file = file

	// Create a new bufio.reader and assign it.
	l.reader = bufio.NewReader(file)

	// The fields Lexer.justSkippedNewline and Lexer.bytesRead need not be reinitialised to 0.

	return l
}

// Creates an exact clone of the Lexer.
func (l *Lexer) Clone() (*Lexer, error) {
	// Reopen the file.
	file, err := os.Open(l.file.Name())

	if err != nil {
		return nil, err
	}

	// Create a new lexer from the reopened file.
	lexer := NewLexer(file)

	// Set the position and justSkippedNewline fields.
	lexer.currentPosition = l.currentPosition
	lexer.justSkippedNewline = l.justSkippedNewline

	// Discard the appropriate amount of bytes to make it an exact clone.
	lexer.reader.Discard(int(l.bytesRead))

	return lexer, nil
}

// Checks if a given rune is a valid identifier start character.
func IsValidIdentStart(r rune) bool {
	allowedChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_"

	return strings.ContainsRune(allowedChars, r)
}

// Checks if a given rune is a valid identifier part character
func IsValidIdentPart(r rune) bool {
	allowedChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

	return strings.ContainsRune(allowedChars, r)
}

// Checks if a given rune is a valid character inside of a string.
func IsValidStringPart(r rune) bool {
	forbiddenChars := "\"\t\n\r"

	return !strings.ContainsRune(forbiddenChars, r)
}

// Given a rune and a LexerNumericalBase, checks if the rune is a valid character to compose an integer in that base.
func IsValidNumberPart(r rune, base LexerNumericalBase) bool {
	allowedChars := map[LexerNumericalBase]string{
		Base2LexerNumericalBase:  "01",
		Base8LexerNumericalBase:  "01234567",
		Base10LexerNumericalBase: "0123456789",
		Base16LexerNumericalBase: "0123456789ABCDEFabcdef",
	}

	return strings.ContainsRune(allowedChars[base], r)
}

// Returns the current Position of the Lexer.
// Primarily used in the parser.
func (l Lexer) CurrentPos() Position {
	return l.currentPosition
}

// Reads a rune from the Lexer's bufio.Reader.
// Can be configured to skip whitespace, handle newlines, or both.
// Use Lexer.readRuneDefault if you want  a shorter way of calling Lexer.readRune(true, true).
//
// Returns (None, nil) when EOF has been reached.
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
		l.currentPosition.Column++

		return l.readRune(skipWhitespace, autoHandleNewlines)
	case r == '\n' && autoHandleNewlines:
		l.currentPosition.Line++
		l.currentPosition.Column = 1

		return l.readRune(skipWhitespace, autoHandleNewlines)
	default:
		l.currentPosition.Column++
	}

	return utils.SomeOptional(r), nil
}

// Peeks n bytes from the *Lexer's *bufio.Reader.
// Returns (None, nil) when EOF is reached.
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

// Equivalent to *Lexer.readRune(true, true).
func (l *Lexer) readRuneDefault() (utils.Optional[rune], error) {
	return l.readRune(true, true)
}

// Peeks the next Token from the Lexer.
// Used primarily in the parser.
// Functionally equivalent to calling Lexer.Clone().NextToken() and handling any errors.
func (l *Lexer) PeekToken() (utils.Optional[Token], error) {
	// Create a clone
	cloned, err := l.Clone()

	if err != nil {
		return utils.NoneOptional[Token](), err
	}

	// Get the next token from the clone
	return cloned.NextToken()
}

// Consumes any amount of runes or lines from the *Lexer's *bufio.Reader to produce a Token.
// Returns (None, nil) when EOF is reached.
func (l *Lexer) NextToken() (utils.Optional[Token], error) {
	// Read the next character
	maybeRune, err := l.readRuneDefault()

	// Check for EOF
	if err == io.EOF {
		return utils.NoneOptional[Token](), nil
	} else if err != nil {
		return utils.NoneOptional[Token](), err
	}

	// Get the character read (maybeRune is None when EOF is reached).
	r, err := maybeRune.Value()

	// Return (None, nil) because we have reached the EOF
	if err != nil {
		return utils.NoneOptional[Token](), nil
	}

	startpos := l.currentPosition

	// A bifg switch for handling characters, comments, strings, and integers of any base
	switch {
	// See if the character we just read is in the Operator token group.
	case strings.ContainsRune(string(TokenOperatorGroup), r):
		return utils.SomeOptional(InitToken(&TokenOperatorGroup, string(r), InitLocation(startpos, l.currentPosition))), nil

	// See if the character we just read is in the Grouping token group.
	case strings.ContainsRune(string(TokenGroupingGroup), r):
		return utils.SomeOptional(InitToken(&TokenGroupingGroup, string(r), InitLocation(startpos, l.currentPosition))), nil

	// See if the character we just read is in the Separator token group.
	case strings.ContainsRune(string(TokenSeparatorGroup), r):
		return utils.SomeOptional(InitToken(&TokenSeparatorGroup, string(r), InitLocation(startpos, l.currentPosition))), nil

	// Skip past comments
	case r == '#':
		line, _, err := l.reader.ReadLine()
		l.bytesRead += uint64(len(line))

		if err != nil {
			return utils.NoneOptional[Token](), err
		}

		l.currentPosition.Line++
		l.currentPosition.Column = 1

		return l.NextToken()

	// Lex a string
	case r == '"':
		str := ""

		// Read the next character without skipping whitespace or handling newlines.
		mp, err := l.readRune(false, false)
		if err != nil {
			return utils.NoneOptional[Token](), err
		}

		// Check for EOF (not allowed here).
		p, err := mp.Value()
		if err != nil {
			return utils.NoneOptional[Token](), err
		}

		for {
			// Add the previosuly read character to the string
			str += string(p)

			// Peek the next byte
			mb, err := l.peekBytes(1)
			if err != nil {
				return utils.NoneOptional[Token](), err
			}

			// Check for EOF (not allowed here).
			b, err := mb.Value()
			if err != nil {
				return utils.NoneOptional[Token](), err
			}

			// If the next character (that we peeked) is not a valid part of a string, exit the loop now.
			if !IsValidStringPart(rune(b[0])) {
				break
			}

			// Since the character we peeked was a valid part of a string,
			// consume it and add it to the string.
			mp, err := l.readRune(false, false)
			if err != nil {
				return utils.NoneOptional[Token](), err
			}

			// Check for EOF (not allowed here).
			// Also assign `p` to the character so that the next loop iteration will add it to the string.
			p, err = mp.Value()
			if err != nil {
				return utils.NoneOptional[Token](), err
			}
		}

		// Blindly assume the next character is a quote and consume it.
		if _, err := l.readRuneDefault(); err != nil {
			return utils.NoneOptional[Token](), err
		}

		// Return a string token.
		return utils.SomeOptional(InitToken(&TokenStringGroup, str, InitLocation(startpos, l.currentPosition))), nil

	// Lex a base-10 number.
	case r == '1' || r == '2' || r == '3' || r == '4' || r == '5' || r == '6' || r == '7' || r == '8' || r == '9':
		// The actual number lexed.
		num := int64(0)

		// A string storing the digits of the lexed number.
		str := string(r)

		// Read the next character
		mp, err := l.readRune(false, false)
		if err != nil {
			return utils.NoneOptional[Token](), err
		}

		// Check for EOF (not allowed here).
		p, err := mp.Value()
		if err != nil {
			return utils.NoneOptional[Token](), err
		}

		for {
			// If the previously read character is not a valid base-10 character, exit the loop.
			// I am unsure if this check here is required, but let's assume it is.
			if !IsValidNumberPart(p, Base10LexerNumericalBase) {
				break
			}

			// Add the digit to the string
			str += string(p)

			// Peek the next byte
			mb, err := l.peekBytes(1)
			if err != nil {
				return utils.NoneOptional[Token](), err
			}

			// Check for EOF (not allowed here).
			b, err := mb.Value()
			if err != nil {
				return utils.NoneOptional[Token](), err
			}

			// If the next character is not a valid digit, exit the loop
			if !IsValidNumberPart(rune(b[0]), Base10LexerNumericalBase) {
				break
			}

			// Consume the character
			mp, err := l.readRune(false, false)
			if err != nil {
				return utils.NoneOptional[Token](), err
			}

			// Check for EOF (not allowed here) and set `p` to the character
			p, err = mp.Value()
			if err != nil {
				return utils.NoneOptional[Token](), err
			}
		}

		// Convert the string of digits to a number
		num, err = strconv.ParseInt(str, 10, 64)
		if err != nil {
			return utils.NoneOptional[Token](), err
		}

		// Return a number token
		return utils.SomeOptional(InitToken(&TokenIntegerGroup, strconv.FormatInt(num, 10), InitLocation(startpos, l.currentPosition))), nil

	// Lex a number of base 2, 8, or 16
	case r == '0':
		// Read the next character so we know which base to use
		maybeRune, err := l.readRuneDefault()
		if err != nil {
			return utils.NoneOptional[Token](), err
		}

		// Check for EOF (not allowed here).
		r, err := maybeRune.Value()
		if err != nil {
			return utils.NoneOptional[Token](), err
		}

		// The number that will be lexed
		num := int64(0)

		// Parse numbers based on the prefix
		switch r {

		// the 0b prefix indicates a base-2 number
		// The code in here is almost identical to the base-10 integer code
		// so I will only explain the important differences
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

			// Parse the digits as base-2 into a 64 bit integer
			num, err = strconv.ParseInt(str, 2, 64)
			if err != nil {
				return utils.NoneOptional[Token](), err
			}

		// The 0o prefix indicates an octal number
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

			// Parse the digits as base-8 into a 64 bit integer
			n, err := strconv.ParseInt(str, 8, 64)

			if err != nil {
				return utils.NoneOptional[Token](), err
			}

			num = n

		// The 0x prefix indicates a hexadecimal number
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

			// Parse the digits as base-16 into a 64 bit integer
			n, err := strconv.ParseInt(str, 16, 64)

			if err != nil {
				return utils.NoneOptional[Token](), err
			}

			num = n
		}

		return utils.SomeOptional(InitToken(&TokenIntegerGroup, strconv.FormatInt(num, 10), InitLocation(startpos, l.currentPosition))), nil
	}

	// Check if `r` is a valid identifier start
	// Again, this is almost identical to the number parsing code.
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

	// TODO: have an error for the invalid token case
	return utils.NoneOptional[Token](), nil
}
