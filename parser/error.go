package parser

import (
	"fmt"

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
