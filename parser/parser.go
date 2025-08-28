package parser

import (
	"fmt"
	"maps"
	"os"
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

	if slices.Compare(*ofGroup, *tk.Group) != 0 || char != []rune(tk.Characters)[0] {
		return lexer.Token{}, ParseErrorExpectedCharacter{
			Expected:      char,
			ExpectedGroup: ofGroup,
			Got:           []rune(tk.Characters)[0],
			GotGroup:      tk.Group,
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

	if expect.Group != tk.Group || expect.Characters != tk.Characters {
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
		Name:     tk.Characters,
		Generics: map[string]TypeGenericASTNode{},
	}, nil
}

func (p *Parser) ParseReferenceType(isEscaping bool, isDyn bool, startpos lexer.Position) (ImmutableReference, error) {
	if _, err := p.ExpectCharacter('&', &lexer.TokenOperatorGroup); err != nil {
		return ImmutableReference{}, err
	}

	typ, err := p.ParseType()

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

func (p *Parser) ParseSliceType(isEscaping bool, startpos lexer.Position) (SliceTypeASTNode, error) {
	if _, err := p.ExpectCharacter('[', &lexer.TokenGroupingGroup); err != nil {
		return SliceTypeASTNode{}, err
	}

	if _, err := p.ExpectCharacter(']', &lexer.TokenGroupingGroup); err != nil {
		return SliceTypeASTNode{}, err
	}

	typ, err := p.ParseType()
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

func (p *Parser) ParseNeverType(startpos lexer.Position) (NeverTypeASTNode, error) {
	if _, err := p.ExpectCharacter('!', &lexer.TokenGroupingGroup); err != nil {
		return NeverTypeASTNode{}, err
	}

	return NeverTypeASTNode{
		Loc: lexer.InitLocation(startpos, p.lexer.CurrentPos()),
	}, nil
}

func (p *Parser) ParseAndDetermineRefSubType(isEscaping bool, isDyn bool, startpos lexer.Position) (RefType, error) {
	if _, err := p.NextToken(); err != nil {
		return nil, err
	}

	mtk, err := p.PeekToken()
	if err != nil {
		return nil, err
	}

	tk, err := mtk.Value()
	if err != nil {
		return nil, ParseErrorUnexpectedEOF{
			WhileParsing: ImmutableReferenceTypeASTNodeKind,
		}
	}

	switch tk.Characters {
	case "[":
		return p.ParseSliceType(isEscaping, startpos)
	case "&":
		return p.ParseReferenceType(isEscaping, isDyn, startpos)
	}

	if _, err := p.ExpectOneOfCharacters(map[rune]*lexer.TokenGroup{
		'[': &lexer.TokenGroupingGroup,
		'&': &lexer.TokenOperatorGroup,
	}); err == nil {
		fmt.Println("How the fuck did this happen? It matched something that it said wasn't there!")
		os.Exit(1)
	}

	return nil, err
}

func (p *Parser) ParseAndDertermineRefMutability(isEscaping bool, isDyn bool, startpos utils.Optional[lexer.Position]) (RefType, error) {
	mtk, err := p.PeekToken()
	if err != nil {
		return nil, err
	}

	tk, err := mtk.Value()
	if err != nil {
		return nil, ParseErrorUnexpectedEOF{
			WhileParsing: ImmutableReferenceTypeASTNodeKind,
		}
	}

	switch tk.Characters {
	case "const", "mut":
		sp := tk.Startpos()
		if psp, err := startpos.Value(); err == nil {
			sp = psp
		}

		return p.ParseAndDetermineRefSubType(isEscaping, isDyn, sp)
	case "escaping":
		if _, err := p.NextToken(); err != nil {
			return nil, err
		}

		sp := tk.Startpos()

		return p.ParseReferenceType(true, false, sp)
	case "dyn":
		sp := tk.Startpos()
		if psp, err := startpos.Value(); err == nil {
			sp = psp
		}

		return p.ParseReferenceType(isEscaping, true, sp)
	}

	return nil, fmt.Errorf("Why did you do this? Why did you call this function without a reference type?")
}

func (p *Parser) ParseType() (Type, error) {
	mtk, err := p.PeekToken()
	if err != nil {
		return nil, err
	}

	tk, err := mtk.Value()
	if err != nil {
		return nil, ParseErrorUnexpectedEOF{
			WhileParsing: ImmutableReferenceTypeASTNodeKind,
		}
	}

	switch tk.Characters {
	case "escaping", "dyn", "const", "mut":
		return p.ParseAndDertermineRefMutability(false, false, utils.NoneOptional[lexer.Position]())
	case "!":
		return p.ParseNeverType(p.lexer.CurrentPos())
	default:
		return p.ParseNamedType()
	}
}

func (p *Parser) ParseFnDefinition() (FunctionDefinitionASTNode, error) {
	tk, err := p.ExpectToken(lexer.InitToken(&lexer.TokenIdentifierGroup, "fn", lexer.Location{}))
	if err != nil {
		return FunctionDefinitionASTNode{}, err
	}

	startpos := tk.Startpos()

	nametk, err := p.ExpectTokenOfGroup(&lexer.TokenIdentifierGroup)
	if err != nil {
		return FunctionDefinitionASTNode{}, nil
	}

	name := nametk.Characters

	arglist := map[string]Type{}

	// Parse arglist

	if _, err := p.ExpectCharacter('(', &lexer.TokenGroupingGroup); err != nil {
		return FunctionDefinitionASTNode{}, err
	}

	for {
		mnexttk, err := p.PeekToken()
		if err != nil {
			break
		}

		nexttk, err := mnexttk.Value()
		if err != nil {
			break
		}

		if nexttk.Characters == ")" {
			break
		}

		nametk, err := p.ExpectTokenOfGroup(&lexer.TokenIdentifierGroup)
		if err != nil {
			return FunctionDefinitionASTNode{}, nil
		}

		argname := nametk.Characters
		argtype, err := p.ParseType()
		if err != nil {
			return FunctionDefinitionASTNode{}, nil
		}

		arglist[argname] = argtype

		mnexttk, err = p.PeekToken()
		if err != nil {
			break
		}

		nexttk, err = mnexttk.Value()
		if err != nil {
			break
		}

		if nexttk.Characters != "," {
			break
		}

		p.NextToken()
	}

	rettype, err := p.ParseType()
	if err != nil {
		return FunctionDefinitionASTNode{}, nil
	}

	if _, err := p.ExpectCharacter(')', &lexer.TokenGroupingGroup); err != nil {
		return FunctionDefinitionASTNode{}, err
	}

	body, err := p.ParseBlock()
	if err != nil {
		return FunctionDefinitionASTNode{}, err
	}

	// TODO: generics

	return FunctionDefinitionASTNode{
		Loc:        lexer.InitLocation(startpos, p.lexer.CurrentPos()),
		Name:       name,
		ReturnType: rettype,
		Parameters: arglist,
		Generics:   map[string]TypeGenericASTNode{},
		Body:       body,
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

	path := []string{tk.Characters}

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

		path = append(path, tk.Characters)
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

	switch tk.Characters {
	case "import":
		n, err := p.ParseImportStatement()
		if err != nil {
			return utils.NoneOptional[Statement](), err
		}

		return utils.SomeOptional(Statement(n)), nil
	case "fn":
		n, err := p.ParseFnDefinition()
		if err != nil {
			return utils.NoneOptional[Statement](), err
		}

		return utils.SomeOptional(Statement(n)), nil
	}

	return utils.NoneOptional[Statement](), fmt.Errorf("Canot parse this node (it is either not supported or invalid) beginning with token %s", tk.ToDisplayString())
}

func (p *Parser) ParseBlock() (BlockASTNode, error) {
	startpos := p.lexer.CurrentPos()

	if _, err := p.ExpectCharacter('{', &lexer.TokenGroupingGroup); err != nil {
		return BlockASTNode{}, err
	}

	var nodes []ASTNode

	for {
		mnexttk, err := p.PeekToken()
		if err != nil {
			break
		}

		nexttk, err := mnexttk.Value()
		if err != nil {
			break
		}

		if nexttk.Characters == "}" {
			break
		}

		mstmt, err := p.ParseStatement()
		if err != nil {
			return BlockASTNode{}, err
		}

		stmt, err := mstmt.Value()
		if err != nil {
			return BlockASTNode{}, ParseErrorUnexpectedEOF{
				WhileParsing: BlockASTNodeKind,
			}
		}

		nodes = append(nodes, stmt)
	}

	if _, err := p.ExpectCharacter('}', &lexer.TokenGroupingGroup); err != nil {
		return BlockASTNode{}, err
	}

	return BlockASTNode{
		Loc:  lexer.InitLocation(startpos, p.lexer.CurrentPos()),
		Code: nodes,
	}, nil
}
