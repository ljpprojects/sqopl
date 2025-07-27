package lexer

import (
	"fmt"
	"slices"

	"ljpprojects.org/sqopl/utils"
)

type Location utils.Range[Position]

func InitLocation(start Position, end Position) Location {
	return Location(utils.InitRange(start, end))
}

type TokenGroup []rune

var (
	TokenOperatorGroup  TokenGroup = TokenGroup("~!^&*-+=|/.?<>%")
	TokenSeparatorGroup TokenGroup = TokenGroup(";:,")
	TokenGroupingGroup  TokenGroup = TokenGroup("([{}])")

	TokenIdentifierGroup TokenGroup = TokenGroup{}
	TokenStringGroup     TokenGroup = TokenGroup{}
	TokenIntegerGroup    TokenGroup = TokenGroup{}
	TokenDecimalGroup    TokenGroup = TokenGroup{}
)

type Token struct {
	group      *TokenGroup
	characters string
	loc        Location
}

func (t Token) Characters() string {
	return t.characters
}

func (t Token) Group() *TokenGroup {
	return t.group
}

func (t Token) Startpos() Position {
	return t.loc.Start
}

func (t Token) Endpos() Position {
	return t.loc.End
}

func (t Token) ToDisplayString() string {
	return fmt.Sprintf(
		"%s(%s) @ (%d:%d)-(%d:%d)",
		t.group.ToDisplayString(),
		t.characters,
		t.loc.Start.line,
		t.loc.Start.column,
		t.loc.End.line,
		t.loc.End.column,
	)
}

func (g *TokenGroup) ToDisplayString() string {
	if g == nil {
		return "Invalid"
	}

	if slices.Compare(*g, TokenOperatorGroup) == 0 {
		return "Operators"
	} else if slices.Compare(*g, TokenSeparatorGroup) == 0 {
		return "Separators"
	} else if slices.Compare(*g, TokenGroupingGroup) == 0 {
		return "Grouping"
	} else if slices.Compare(*g, TokenIdentifierGroup) == 0 {
		return "Identifiers"
	} else if slices.Compare(*g, TokenStringGroup) == 0 {
		return "Strings"
	} else if slices.Compare(*g, TokenIntegerGroup) == 0 {
		return "Integers"
	} else if slices.Compare(*g, TokenDecimalGroup) == 0 {
		return "Decimals"
	} else {
		return "Unknown"
	}
}

func InitToken(group *TokenGroup, characters string, loc Location) Token {
	return Token{
		group,
		characters,
		loc,
	}
}
