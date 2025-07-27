package lexer

import (
	"fmt"
	"slices"

	"ljpprojects.org/sqopl/utils"
)

// An alias for utils.Range[Position] primarily for specifying the locations of tokens and AST nodes.
type Location utils.Range[Position]

// Intialises a Location by casting the result of utils.InitRange(start, end)
func InitLocation(start Position, end Position) Location {
	return Location(utils.InitRange(start, end))
}

// An alias for []rune, but serves to separate characters into groups.
type TokenGroup []rune

var (
	// A TokenGroup with the characters making up operators.
	TokenOperatorGroup TokenGroup = TokenGroup("~!^&*-+=|/.?<>%")

	// A TokenGroup with the characters used to separate.
	TokenSeparatorGroup TokenGroup = TokenGroup(";:,")

	// A TokenGroup with the characters that group expressions or statements.
	TokenGroupingGroup TokenGroup = TokenGroup("([{}])")

	// A TokenGroup representing identifiers (but not containing the characters of them).
	TokenIdentifierGroup TokenGroup = TokenGroup{}

	// A TokenGroup representing strings (but not containing the characters of them).
	TokenStringGroup TokenGroup = TokenGroup{}

	// A TokenGroup representing integers (but not containing the characters of them).
	TokenIntegerGroup TokenGroup = TokenGroup{}

	// A TokenGroup representing decimals (but not containing the characters of them).
	//
	// UNUSED
	TokenDecimalGroup TokenGroup = TokenGroup{}
)

// A Token, which stores important information used to identify and locate it.
type Token struct {
	// A pointer to the group of the token.
	Group *TokenGroup

	// A string storing the characters of the token.
	Characters string

	// The location of the token in the source.
	Loc Location
}

// A method that returns a Token's .Loc.Start
func (t Token) Startpos() Position {
	return t.Loc.Start
}

// A method that returns a Token's .Loc.End
func (t Token) Endpos() Position {
	return t.Loc.End
}

// Formats the Token into a string for printing out.
func (t Token) ToDisplayString() string {
	return fmt.Sprintf(
		"%s(%s) @ (%d:%d)-(%d:%d)",
		t.Group.ToDisplayString(),
		t.Characters,
		t.Loc.Start.Line,
		t.Loc.Start.Column,
		t.Loc.End.Line,
		t.Loc.End.Column,
	)
}

// Formats a TokenGroup into a string for printing out.
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

// Creates a new token.
// Equivalent to just using the structure initialisation syntax.
func InitToken(group *TokenGroup, characters string, loc Location) Token {
	return Token{
		Group:      group,
		Characters: characters,
		Loc:        loc,
	}
}
