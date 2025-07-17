package lexer

import "slices"

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
}

func (t Token) Characters() string {
	return t.characters
}

func (t Token) Group() *TokenGroup {
	return t.group
}

func (g *TokenGroup) ToDisplayString() string {
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
