package parser

import (
	"fmt"
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

func (p *Parser) PeekToken() (utils.Optional[lexer.Token], error) {
	return p.lexer.PeekToken()
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
		return utils.NoneOptional[lexer.Token](), ParseErrorExpectedCharacter{
			Expected:      char,
			ExpectedGroup: ofGroup,
			Got:           []rune(tk.Characters())[0],
			GotGroup:      tk.Group(),
		}
	}

	return utils.SomeOptional(tk), nil
}

func (p *Parser) ExpectToken(expect lexer.Token) (utils.Optional[lexer.Token], error) {
	mtk, err := p.NextToken()

	if err != nil {
		return utils.NoneOptional[lexer.Token](), err
	}

	tk, err := mtk.Value()

	if err != nil {
		return utils.NoneOptional[lexer.Token](), nil
	}

	if expect.Group() != tk.Group() || expect.Characters() != tk.Characters() {
		return utils.NoneOptional[lexer.Token](), ParseErrorExpectedToken{
			Expected: expect,
			Got:      tk,
		}
	}

	return utils.SomeOptional(tk), nil
}

func (p *Parser) ExpectTokenOfGroup(expectGroup *lexer.TokenGroup) (utils.Optional[lexer.Token], error) {
	mtk, err := p.NextToken()

	if err != nil {
		return utils.NoneOptional[lexer.Token](), err
	}

	tk, err := mtk.Value()

	if err != nil {
		return utils.NoneOptional[lexer.Token](), nil
	}

	if tk.Group() != expectGroup {
		return utils.NoneOptional[lexer.Token](), ParseErrorExpectedToken{
			Expected: lexer.InitToken(expectGroup, "[ANYTHING]", lexer.Location{}),
			Got:      tk,
		}
	}

	return utils.SomeOptional(tk), nil
}

func (p *Parser) ParseNamedType() (NamedTypeASTNode, error) {
	mtk, err := p.ExpectToken(lexer.InitToken(&lexer.TokenIdentifierGroup, "import", lexer.Location{}))

	if err != nil {
		return NamedTypeASTNode{}, err
	}

	tk, err := mtk.Value()

	if err != nil {
		return NamedTypeASTNode{}, ParseErrorUnexpectedEOF{
			WhileParsing: ImportStatementASTNodeKind,
		}
	}

	startpos := tk.Startpos()

	mtk, err = p.ExpectTokenOfGroup(&lexer.TokenIdentifierGroup)

	if err != nil {
		return NamedTypeASTNode{}, err
	}

	tk, err = mtk.Value()

	if err != nil {
		return NamedTypeASTNode{}, ParseErrorUnexpectedEOF{
			WhileParsing: ImportStatementASTNodeKind,
		}
	}

	return NamedTypeASTNode{
		Loc:      lexer.InitLocation(startpos, p.lexer.CurrentPos()),
		Name:     tk.Characters(),
		Generics: map[string]TypeGenericASTNode{},
	}, nil
}

func (p *Parser) ParseImportStatement() (ImportStatementASTNode, error) {
	mtk, err := p.ExpectToken(lexer.InitToken(&lexer.TokenIdentifierGroup, "import", lexer.Location{}))

	if err != nil {
		return ImportStatementASTNode{}, err
	}

	tk, err := mtk.Value()

	if err != nil {
		return ImportStatementASTNode{}, ParseErrorUnexpectedEOF{
			WhileParsing: ImportStatementASTNodeKind,
		}
	}

	startpos := tk.Startpos()

	mtk, err = p.ExpectTokenOfGroup(&lexer.TokenIdentifierGroup)

	if err != nil {
		return ImportStatementASTNode{}, err
	}

	tk, err = mtk.Value()

	if err != nil {
		return ImportStatementASTNode{}, ParseErrorUnexpectedEOF{
			WhileParsing: ImportStatementASTNodeKind,
		}
	}

	path := []string{tk.Characters()}

	for {
		mtk, err := p.ExpectCharacter(':', &lexer.TokenSeparatorGroup)

		if err != nil {
			switch err := err.(type) {
			case ParseErrorExpectedCharacter:
				if err.Got == ';' {
					endpos := p.lexer.CurrentPos()

					return ImportStatementASTNode{
						Loc:  lexer.InitLocation(startpos, endpos),
						Path: path,
					}, nil
				}
			default:
				return ImportStatementASTNode{}, err
			}
		}

		_, err = mtk.Value()

		if err != nil {
			return ImportStatementASTNode{}, ParseErrorUnexpectedEOF{
				WhileParsing: ImportStatementASTNodeKind,
			}
		}

		mtk, err = p.ExpectTokenOfGroup(&lexer.TokenIdentifierGroup)

		if err != nil {
			return ImportStatementASTNode{}, err
		}

		tk, err = mtk.Value()

		if err != nil {
			return ImportStatementASTNode{}, ParseErrorUnexpectedEOF{
				WhileParsing: ImportStatementASTNodeKind,
			}
		}

		path = append(path, tk.Characters())
	}
}

func (p *Parser) ParseStatement() (utils.Optional[Statement], error) {
	mtk, err := p.PeekToken()

	if err != nil {
		return utils.NoneOptional[Statement](), err
	}

	tk, err := mtk.Value()

	if err != nil {
		return utils.NoneOptional[Statement](), nil
	}

	switch tk.Characters() {
	case "import":
		n, err := p.ParseImportStatement()

		if err != nil {
			return utils.NoneOptional[Statement](), nil
		}

		return utils.SomeOptional(Statement(n)), nil
	case "fn":
		n, err := p.ParseImportStatement()

		if err != nil {
			return utils.NoneOptional[Statement](), nil
		}

		return utils.SomeOptional(Statement(n)), nil
	}

	return utils.NoneOptional[Statement](), fmt.Errorf("Canot parse this node (it is either not supported or invalid) beginning with token %s", tk.ToDisplayString())
}
