package parser

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"ljpprojects.org/sqopl/lexer"
)

type ParseErrorExpectedCharacter struct {
	Expected      rune
	ExpectedGroup *lexer.TokenGroup

	Got      rune
	GotGroup *lexer.TokenGroup
}

func (e ParseErrorExpectedCharacter) Error() string {
	return fmt.Sprintf(
		"Expected character '%#U' of group %s, but got character '%#U' of group %s",
		e.Expected,
		e.ExpectedGroup.ToDisplayString(),
		e.Got,
		e.GotGroup.ToDisplayString(),
	)
}

type ParseErrorExpectedOneOfCharacters struct {
	ExpectedOneOf map[rune]*lexer.TokenGroup
	Got           rune
	GotGroup      *lexer.TokenGroup
}

func (e ParseErrorExpectedOneOfCharacters) Error() string {
	var builder strings.Builder

	builder.WriteString("Expected one of the following characters: ")

	keys := slices.Collect(maps.Keys(e.ExpectedOneOf))

	i := 0

	for c, g := range e.ExpectedOneOf {
		if i == len(keys)-1 {
			builder.WriteString(", or")
			builder.WriteRune(c)
			builder.WriteString(" (of group ")
			builder.WriteString(g.ToDisplayString())
			builder.WriteString(")")
		} else if i == 0 {
			builder.WriteRune(c)
			builder.WriteString(" (of group ")
			builder.WriteString(g.ToDisplayString())
			builder.WriteString(")")
		} else {
			builder.WriteString(", ")
			builder.WriteRune(c)
			builder.WriteString(" (of group ")
			builder.WriteString(g.ToDisplayString())
			builder.WriteString(")")
		}

		i++
	}

	builder.WriteString(". Got ")
	builder.WriteRune(e.Got)
	builder.WriteString(" (of group ")
	builder.WriteString(e.GotGroup.ToDisplayString())
	builder.WriteString(").")

	return builder.String()
}

type ParseErrorExpectedToken struct {
	Expected lexer.Token
	Got      lexer.Token
}

func (e ParseErrorExpectedToken) Error() string {
	return fmt.Sprintf(
		"Expected token %s, but got token %s",
		e.Expected.ToDisplayString(),
		e.Got.ToDisplayString(),
	)
}

type ParseErrorUnexpectedEOF struct {
	WhileParsing ASTNodeKind
}

func (e ParseErrorUnexpectedEOF) Error() string {
	return fmt.Sprintf(
		"Unexpected EOF while parsing node %s",
		e.WhileParsing.ToDisplayString(),
	)
}
