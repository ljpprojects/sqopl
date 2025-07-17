package parser

import (
	"slices"

	"ljpprojects.org/sqopl/lexer"
	"ljpprojects.org/sqopl/utils"
)

type Parser struct {
	lexer *lexer.Lexer
}

func NewParser(lexer *lexer.Lexer) *Parser {
	p := new(Parser)

	p.lexer = lexer

	return p
}

func (p *Parser) NextToken() (utils.Optional[lexer.Token], error) {
	return p.lexer.NextToken()
}

func (p *Parser) ExpectCharacter(char rune, ofGroup *lexer.TokenGroup) (utils.Optional[lexer.Token], error) {
	mtk, err := p.NextToken()

	if err != nil {
		return utils.NoneOptional[lexer.Token](), err
	}

	tk, err := mtk.Value()

	if err != nil {
		return utils.NoneOptional[lexer.Token](), nil
	}

	if slices.Compare(*ofGroup, *tk.Group()) != 0 || char != []rune(tk.Characters())[0] {
		return utils.NoneOptional[lexer.Token](), ParseErrorExpectedCharacter(char, ofGroup, []rune(tk.Characters())[0], tk.Group())
	}

	return utils.SomeOptional(tk), nil
}

func (p *Parser) NextNode() (utils.Optional[ASTNode[any]], error) {

}
