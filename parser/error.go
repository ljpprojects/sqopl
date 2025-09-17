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

	prefix        string
	while_parsing ASTNodeKind
}

func (e ParseErrorExpectedCharacter) Error() string {
	return fmt.Sprintf(
		"(during parsing %s) [%s] Expected character '%#U' of group %s, but got character '%#U' of group %s",
		e.while_parsing.ToDisplayString(),
		e.prefix,
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
	prefix        string
	while_parsing ASTNodeKind
}

func (e ParseErrorExpectedOneOfCharacters) Error() string {
	var builder strings.Builder

	builder.WriteString("(during parsing " + e.while_parsing.ToDisplayString() + ")")
	builder.WriteString(" [" + e.prefix + "]")
	builder.WriteString(" Expected one of the following characters: ")

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
	Expected      lexer.Token
	Got           lexer.Token
	prefix        string
	while_parsing ASTNodeKind
}

func (e ParseErrorExpectedToken) Error() string {
	return fmt.Sprintf(
		"(during parsing %s) [%s] Expected token %s, but got token %s",
		e.while_parsing.ToDisplayString(),
		e.prefix,
		e.Expected.ToDisplayString(),
		e.Got.ToDisplayString(),
	)
}

type ParseErrorUnexpectedEOF struct {
	WhileParsing ASTNodeKind
	prefix       string
}

func (e ParseErrorUnexpectedEOF) Error() string {
	return fmt.Sprintf(
		"(during parsing %s) [%s] Unexpected EOF",
		e.WhileParsing.ToDisplayString(),
		e.prefix,
	)
}

func (p *Parser) GenerateExpectedCharacterError(
	expected lexer.Token,
	got lexer.Token,
) ParseErrorExpectedCharacter {
	return ParseErrorExpectedCharacter{
		Expected:      []rune(expected.Characters)[0],
		ExpectedGroup: expected.Group,
		Got:           []rune(got.Characters)[0],
		GotGroup:      got.Group,
		prefix:        p.context.ErrorPrefix,
		while_parsing: p.context.Parsing.Unwrap(),
	}
}

func (p *Parser) GenerateExpectedOneOfCharactersError(
	expected_one_of map[rune]*lexer.TokenGroup,
	got lexer.Token,
) ParseErrorExpectedOneOfCharacters {
	return ParseErrorExpectedOneOfCharacters{
		ExpectedOneOf: expected_one_of,
		Got:           []rune(got.Characters)[0],
		GotGroup:      got.Group,
		prefix:        p.context.ErrorPrefix,
		while_parsing: p.context.Parsing.Unwrap(),
	}
}

func (p *Parser) GenerateExpectedTokenError(expected lexer.Token, got lexer.Token) ParseErrorExpectedToken {
	return ParseErrorExpectedToken{
		Expected:      expected,
		Got:           got,
		prefix:        p.context.ErrorPrefix,
		while_parsing: p.context.Parsing.Unwrap(),
	}
}

func (p *Parser) GenerateUnexpectedEOFError() ParseErrorUnexpectedEOF {
	return ParseErrorUnexpectedEOF{
		WhileParsing: p.context.Parsing.Unwrap(),
		prefix:       p.context.ErrorPrefix,
	}
}
