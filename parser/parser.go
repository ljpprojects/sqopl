package parser

import (
	"fmt"
	"maps"
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

func (p *Parser) ExpectCharacter(char rune, ofGroup *lexer.TokenGroup) (lexer.Token, error) {
	mtk, err := p.NextToken()

	if err != nil {
		return lexer.Token{}, err
	}

	tk, err := mtk.Value()

	if err != nil {
		return lexer.Token{}, nil
	}

	if slices.Compare(*ofGroup, *tk.Group()) != 0 || char != []rune(tk.Characters())[0] {
		return lexer.Token{}, ParseErrorExpectedCharacter{
			Expected:      char,
			ExpectedGroup: ofGroup,
			Got:           []rune(tk.Characters())[0],
			GotGroup:      tk.Group(),
		}
	}

	return tk, nil
}

func (p *Parser) ExpectToken(expect lexer.Token) (lexer.Token, error) {
	mtk, err := p.NextToken()

	if err != nil {
		return lexer.Token{}, err
	}

	tk, err := mtk.Value()

	if err != nil {
		return lexer.Token{}, nil
	}

	if expect.Group() != tk.Group() || expect.Characters() != tk.Characters() {
		return lexer.Token{}, ParseErrorExpectedToken{
			Expected: expect,
			Got:      tk,
		}
	}

	return tk, nil
}

func (p *Parser) ExpectTokenOfGroup(expectGroup *lexer.TokenGroup) (lexer.Token, error) {
	mtk, err := p.NextToken()

	if err != nil {
		return lexer.Token{}, err
	}

	tk, err := mtk.Value()

	if err != nil {
		return lexer.Token{}, nil
	}

	if tk.Group != expectGroup {
		return lexer.Token{}, ParseErrorExpectedToken{
			Expected: lexer.InitToken(expectGroup, "[ANYTHING]", lexer.Location{}),
			Got:      tk,
		}
	}

	return tk, nil
}

func (p *Parser) ExpectOneOfCharacters(expectOneOf map[rune]*lexer.TokenGroup) (lexer.Token, error) {
	for r, g := range expectOneOf {
		tk, err := p.ExpectCharacter(r, g)

		if err != nil {
			switch err := err.(type) {
			case ParseErrorExpectedCharacter:
				if !slices.Contains(slices.Collect(maps.Keys(expectOneOf)), err.Got) {
					return lexer.Token{}, ParseErrorExpectedOneOfCharacters{
						ExpectedOneOf: expectOneOf,
						Got:           err.Got,
						GotGroup:      err.GotGroup,
					}
				} else {
					continue
				}
			default:
				return lexer.Token{}, err
			}
		}

		return tk, err
	}

	// This is to satisfy the compiler, it should never be executed.
	return lexer.Token{}, nil
}

func (p *Parser) ParseNamedType() (NamedTypeASTNode, error) {
	tk, err := p.ExpectTokenOfGroup(&lexer.TokenIdentifierGroup)

	if err != nil {
		return NamedTypeASTNode{}, err
	}

	startpos := tk.Startpos()

	return NamedTypeASTNode{
		Loc:      lexer.InitLocation(startpos, p.lexer.CurrentPos()),
		Name:     tk.Characters(),
		Generics: map[string]TypeGenericASTNode{},
	}, nil
}

func (p *Parser) ParseMutableReferenceType(isEscaping bool, isDyn bool) (MutableReference, error) {
	tk, err := p.ExpectToken(lexer.InitToken(&lexer.TokenIdentifierGroup, "mut", lexer.Location{}))

	if err != nil {
		return MutableReference{}, err
	}

	startpos := tk.Startpos()

	if _, err := p.ExpectCharacter('&', &lexer.TokenOperatorGroup); err != nil {
		return MutableReference{}, err
	}

	typ, err := p.ParseNamedType()

	if err != nil {
		return MutableReference{}, err
	}

	return MutableReference{
		Loc:        lexer.InitLocation(startpos, p.lexer.CurrentPos()),
		IsEscaping: isEscaping,
		Inner:      typ,
		IsDynamic:  isDyn,
	}, nil
}

func (p *Parser) ParseImmutableReferenceType(isEscaping bool, isDyn bool) (ImmutableReference, error) {
	tk, err := p.ExpectToken(lexer.InitToken(&lexer.TokenIdentifierGroup, "const", lexer.Location{}))

	if err != nil {
		return ImmutableReference{}, err
	}

	startpos := tk.Startpos()

	if _, err := p.ExpectCharacter('&', &lexer.TokenOperatorGroup); err != nil {
		return ImmutableReference{}, err
	}

	typ, err := p.ParseNamedType()

	if err != nil {
		return ImmutableReference{}, err
	}

	return ImmutableReference{
		Loc:        lexer.InitLocation(startpos, p.lexer.CurrentPos()),
		IsEscaping: isEscaping,
		Inner:      typ,
		IsDynamic:  isDyn,
	}, nil
}

func (p *Parser) ParseMutableSliceType(isEscaping bool) (SliceTypeASTNode, error) {
	tk, err := p.ExpectToken(lexer.InitToken(&lexer.TokenIdentifierGroup, "mut", lexer.Location{}))

	if err != nil {
		return SliceTypeASTNode{}, err
	}

	startpos := tk.Startpos()

	if _, err := p.ExpectCharacter('&', &lexer.TokenOperatorGroup); err != nil {
		return SliceTypeASTNode{}, err
	}

	typ, err := p.ParseNamedType()

	if err != nil {
		return SliceTypeASTNode{}, err
	}

	return SliceTypeASTNode{
		Loc:        lexer.InitLocation(startpos, p.lexer.CurrentPos()),
		IsEscaping: isEscaping,
		IsMutable:  true,
		ValueType:  typ,
	}, nil
}

func (p *Parser) ParseImmutableSliceType(isEscaping bool) (SliceTypeASTNode, error) {
	tk, err := p.ExpectToken(lexer.InitToken(&lexer.TokenIdentifierGroup, "const", lexer.Location{}))

	if err != nil {
		return SliceTypeASTNode{}, err
	}

	startpos := tk.Startpos()

	if _, err := p.ExpectCharacter('&', &lexer.TokenOperatorGroup); err != nil {
		return SliceTypeASTNode{}, err
	}

	typ, err := p.ParseNamedType()

	if err != nil {
		return SliceTypeASTNode{}, err
	}

	return SliceTypeASTNode{
		Loc:        lexer.InitLocation(startpos, p.lexer.CurrentPos()),
		IsEscaping: isEscaping,
		IsMutable:  false,
		ValueType:  typ,
	}, nil
}

func (p *Parser) ParseImportStatement() (ImportStatementASTNode, error) {
	tk, err := p.ExpectToken(lexer.InitToken(&lexer.TokenIdentifierGroup, "import", lexer.Location{}))

	if err != nil {
		return ImportStatementASTNode{}, err
	}

	startpos := tk.Startpos()

	tk, err = p.ExpectTokenOfGroup(&lexer.TokenIdentifierGroup)

	if err != nil {
		return ImportStatementASTNode{}, err
	}

	path := []string{tk.Characters()}

	for {
		_, err := p.ExpectCharacter(':', &lexer.TokenSeparatorGroup)

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

		tk, err = p.ExpectTokenOfGroup(&lexer.TokenIdentifierGroup)

		if err != nil {
			return ImportStatementASTNode{}, err
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
