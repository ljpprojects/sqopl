package parser

import (
	"fmt"

	"ljpprojects.org/sqopl/lexer"
)

func ParseErrorExpectedCharacter(
	expected rune,
	expectedGroup *lexer.TokenGroup,
	got rune,
	gotGroup *lexer.TokenGroup,
) error {
	return fmt.Errorf(
		"Expected character '%#U' of group %s, but got character '%#U' of group %s",
		expected,
		expectedGroup.ToDisplayString(),
		got,
		gotGroup.ToDisplayString(),
	)
}
