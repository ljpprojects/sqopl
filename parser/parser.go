package parser

import (
	"fmt"
	"log"
	"maps"
	"os"
	"slices"
	"strconv"

	"ljpprojects.org/sqopl/lexer"
	"ljpprojects.org/sqopl/utils"
)

type ParserContext struct {
	Parsing     utils.Optional[ASTNodeKind]
	ErrorPrefix string
}

func DefaultParserContext() ParserContext {
	return ParserContext{
		Parsing:     utils.NoneOptional[ASTNodeKind](),
		ErrorPrefix: "[DEBUG]",
	}
}

func InitParserContext(error_prefix string) ParserContext {
	return ParserContext{
		Parsing:     utils.NoneOptional[ASTNodeKind](),
		ErrorPrefix: error_prefix,
	}
}

type Parser struct {
	lexer   *lexer.Lexer
	context ParserContext
}

func NewParser(lexer *lexer.Lexer, context ParserContext) *Parser {
	p := new(Parser)

	p.lexer = lexer
	p.context = context

	return p
}

func (p *Parser) SetPrefix(pref string) {
	p.context.ErrorPrefix = pref
}

func (p *Parser) WithPrefixAs(pref string, do func(p *Parser) any) any {
	old := p.context.ErrorPrefix
	p.context.ErrorPrefix = pref

	r := do(p)

	p.context.ErrorPrefix = old

	return r
}

func (p *Parser) DeclareParsing(kind ASTNodeKind) {
	p.context.Parsing = utils.SomeOptional(kind)
}

func (p *Parser) IsParsingNodeOfKind(kind ASTNodeKind) bool {
	k, err := p.context.Parsing.Value()
	if err != nil {
		return false
	}

	return k == kind
}

func (p *Parser) IsParsingNodeOfGroup(group *ASTNodeGroup) bool {
	k, err := p.context.Parsing.Value()
	if err != nil {
		return false
	}

	return slices.Contains(*group, k)
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
		return lexer.Token{}, p.GenerateExpectedCharacterError(
			lexer.InitToken(ofGroup, string(char), lexer.Location{}),
			tk,
		)
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
		return lexer.Token{}, p.GenerateExpectedTokenError(expect, tk)
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
		return lexer.Token{}, p.GenerateExpectedTokenError(lexer.InitToken(expectGroup, "ANYTHING", lexer.Location{}), tk)
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
					return lexer.Token{}, p.GenerateExpectedOneOfCharactersError(
						expectOneOf,
						lexer.InitToken(err.GotGroup, string(err.Got), lexer.Location{}),
					)
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
	p.DeclareParsing(NamedTypeASTNodeKind)

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
	p.DeclareParsing(ImmutableReferenceTypeASTNodeKind)

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
	p.DeclareParsing(SliceTypeASTNodeKind)

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
	p.DeclareParsing(NeverTypeASTNodeKind)

	if _, err := p.ExpectCharacter('!', &lexer.TokenGroupingGroup); err != nil {
		return NeverTypeASTNode{}, err
	}

	return NeverTypeASTNode{
		Loc: lexer.InitLocation(startpos, p.lexer.CurrentPos()),
	}, nil
}

func (p *Parser) ParseAndDetermineRefSubType(isEscaping bool, isDyn bool, startpos lexer.Position) (RefType, error) {
	p.DeclareParsing(MutableReferenceTypeASTNodeKind)

	if _, err := p.NextToken(); err != nil {
		return nil, err
	}

	mtk, err := p.PeekToken()
	if err != nil {
		return nil, err
	}

	tk, err := mtk.Value()
	if err != nil {
		return nil, p.GenerateUnexpectedEOFError()
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
	p.DeclareParsing(RawPointerTypeASTNodeKind)

	mtk, err := p.PeekToken()
	if err != nil {
		return nil, err
	}

	tk, err := mtk.Value()
	if err != nil {
		return nil, p.GenerateUnexpectedEOFError()
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
	p.DeclareParsing(NamedTypeASTNodeKind)

	mtk, err := p.PeekToken()
	if err != nil {
		return nil, err
	}

	tk, err := mtk.Value()
	if err != nil {
		return nil, p.GenerateUnexpectedEOFError()
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
	p.DeclareParsing(FunctionDefinitionASTNodeKind)

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

	log.Printf("fn %s(...) ...\n", name)

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

		log.Printf("fn %s(%s %s) ...\n", name, argname, argtype.Kind().ToDisplayString())

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

	if _, err := p.ExpectCharacter(')', &lexer.TokenGroupingGroup); err != nil {
		return FunctionDefinitionASTNode{}, err
	}

	rettype, err := p.ParseType()
	if err != nil {
		return FunctionDefinitionASTNode{}, nil
	}

	body, err := p.ParseBlock()
	if err != nil {
		return FunctionDefinitionASTNode{}, err
	}

	// TODO: generics

	ret := FunctionDefinitionASTNode{
		Loc:        lexer.InitLocation(startpos, p.lexer.CurrentPos()),
		Name:       name,
		ReturnType: rettype,
		Parameters: arglist,
		Generics:   map[string]TypeGenericASTNode{},
		Body:       body,
	}

	log.Printf("fn %s(...) ...\n", ret.Name)

	return ret, nil
}

func (p *Parser) ParsePaths(startpos lexer.Position, initial []string, separator rune) ([][]string, lexer.Location, error) {
	p.DeclareParsing(ModulePathASTNodeKind)

	rootpath := initial
	fullpaths := [][]string{}

	for {
		_, err := p.ExpectCharacter(':', &lexer.TokenSeparatorGroup)
		if err != nil {
			switch err := err.(type) {
			case ParseErrorExpectedCharacter:
				if err.Got == ';' {
					endpos := p.lexer.CurrentPos()

					return [][]string{rootpath}, lexer.InitLocation(startpos, endpos), nil
				}
			default:
				return [][]string{}, lexer.Location{}, err
			}
		}

		tk, err := p.ExpectTokenOfGroup(&lexer.TokenIdentifierGroup)
		if err != nil {
			fmt.Printf("[WARN] %s\n", err.Error())

			switch err := err.(type) {
			case ParseErrorExpectedToken:
				if err.Got.Characters == "{" { // Parse shorthand imports

					// Just parse mutliple paths
					for {
						if mtk, err := p.lexer.PeekToken(); err == nil {
							if tk, err := mtk.Value(); err == nil {
								if tk.Characters == "}" {
									break
								} else if tk.Characters == "," {
									p.lexer.NextToken()
								}
							} else {
								return [][]string{}, lexer.Location{}, nil
							}
						} else {
							log.Panicf("%s", err.Error())
						}

						subpaths, _, err := p.ParsePaths(startpos, rootpath, separator)
						if err != nil {
							return [][]string{}, lexer.Location{}, err
						}

						fullpaths = append(fullpaths, subpaths...)
					}

					if _, err := p.ExpectCharacter('}', &lexer.TokenGroupingGroup); err != nil {
						return [][]string{}, lexer.Location{}, err
					}

					continue
				}
			default:
				return [][]string{}, lexer.Location{}, err
			}
		}

		rootpath = append(rootpath, tk.Characters)
	}
}

func (p *Parser) ParseImportStatement() (ImportStatementASTNode, error) {
	p.DeclareParsing(ImportStatementASTNodeKind)

	tk, err := p.ExpectToken(lexer.InitToken(&lexer.TokenIdentifierGroup, "import", lexer.Location{}))
	if err != nil {
		return ImportStatementASTNode{}, err
	}

	startpos := tk.Startpos()

	tk, err = p.ExpectTokenOfGroup(&lexer.TokenIdentifierGroup)
	if err != nil {
		return ImportStatementASTNode{}, err
	}

	paths, loc, err := p.ParsePaths(startpos, []string{}, ':')
	if err != nil {
		return ImportStatementASTNode{}, err
	}

	return ImportStatementASTNode{
		Loc:   loc,
		Paths: paths,
	}, nil
}

func (p *Parser) ParseExplciitReturnStatement() (ExplicitReturnASTNode, error) {
	p.DeclareParsing(ExplicitReturnASTNodeKind)

	tk, err := p.ExpectToken(lexer.InitToken(&lexer.TokenIdentifierGroup, "return", lexer.Location{}))
	if err != nil {
		return ExplicitReturnASTNode{}, err
	}

	startpos := tk.Startpos()

	// Handle void return case
	if mtk, err := p.PeekToken(); err == nil {
		if tk, err := mtk.Value(); err == nil {
			if tk.Characters == ";" {
				if _, err := p.ExpectToken(lexer.InitToken(&lexer.TokenSeparatorGroup, ";", lexer.Location{})); err != nil {
					return ExplicitReturnASTNode{}, err
				}

				return ExplicitReturnASTNode{
					Loc:   lexer.InitLocation(startpos, p.lexer.CurrentPos()),
					Value: utils.NoneOptional[Expression](),
				}, nil
			}
		}
	}

	val, err := p.ParseExpression()
	if err != nil {
		return ExplicitReturnASTNode{}, err
	}

	if _, err := p.ExpectCharacter(';', &lexer.TokenSeparatorGroup); err != nil {
		return ExplicitReturnASTNode{}, err
	}

	return ExplicitReturnASTNode{
		Loc:   lexer.InitLocation(startpos, p.lexer.CurrentPos()),
		Value: utils.SomeOptional(val),
	}, nil
}

func (p *Parser) ParseVarDefinition() (VarDefinitionASTNode, error) {
	p.DeclareParsing(VarDefinitionASTNodeKind)

	tk, err := p.ExpectToken(lexer.InitToken(&lexer.TokenIdentifierGroup, "var", lexer.Location{}))
	if err != nil {
		return VarDefinitionASTNode{}, err
	}

	startpos := tk.Startpos()

	nametk, err := p.ExpectTokenOfGroup(&lexer.TokenIdentifierGroup)
	if err != nil {
		return VarDefinitionASTNode{}, err
	}

	name := nametk.Characters

	typ, err := p.ParseType()
	if err != nil {
		return VarDefinitionASTNode{}, err
	}

	if _, err := p.ExpectCharacter(';', &lexer.TokenSeparatorGroup); err == nil {
		return VarDefinitionASTNode{
			Loc:   lexer.InitLocation(startpos, p.lexer.CurrentPos()),
			Name:  name,
			Value: utils.NoneOptional[Expression](),
			Type:  typ,
		}, nil
	} else {
		switch err := err.(type) {
		case ParseErrorExpectedCharacter:
			if err.Got != '=' {
				return VarDefinitionASTNode{}, err
			}
		default:
			return VarDefinitionASTNode{}, err
		}
	}

	val, err := p.ParseExpression()
	if err != nil {
		return VarDefinitionASTNode{}, err
	}

	if _, err := p.ExpectCharacter(';', &lexer.TokenSeparatorGroup); err != nil {
		return VarDefinitionASTNode{}, err
	}

	return VarDefinitionASTNode{
		Loc:   lexer.InitLocation(startpos, p.lexer.CurrentPos()),
		Name:  name,
		Value: utils.SomeOptional(val),
		Type:  typ,
	}, nil
}

func (p *Parser) ParseLetDefinition() (LetDefinitionASTNode, error) {
	p.DeclareParsing(LetDefinitionASTNodeKind)

	tk, err := p.ExpectToken(lexer.InitToken(&lexer.TokenIdentifierGroup, "let", lexer.Location{}))
	if err != nil {
		return LetDefinitionASTNode{}, err
	}

	startpos := tk.Startpos()

	nametk, err := p.ExpectTokenOfGroup(&lexer.TokenIdentifierGroup)
	if err != nil {
		return LetDefinitionASTNode{}, err
	}

	name := nametk.Characters

	typ, err := p.ParseType()
	if err != nil {
		return LetDefinitionASTNode{}, err
	}

	if _, err := p.ExpectCharacter('=', &lexer.TokenOperatorGroup); err != nil {
		return LetDefinitionASTNode{}, err
	}

	val, err := p.ParseExpression()
	if err != nil {
		return LetDefinitionASTNode{}, err
	}

	if _, err := p.ExpectCharacter(';', &lexer.TokenSeparatorGroup); err != nil {
		return LetDefinitionASTNode{}, err
	}

	return LetDefinitionASTNode{
		Loc:   lexer.InitLocation(startpos, p.lexer.CurrentPos()),
		Name:  name,
		Value: val,
		Type:  typ,
	}, nil
}

func (p *Parser) ParseStatement() (Statement, error) {
	p.DeclareParsing(BlockASTNodeKind)

	mtk, err := p.PeekToken()
	if err != nil {
		return nil, err
	}

	tk, err := mtk.Value()
	if err != nil {
		return nil, nil
	}

	switch tk.Characters {
	case "import":
		n, err := p.ParseImportStatement()
		if err != nil {
			return nil, err
		}

		return n, nil
	case "fn":
		n, err := p.ParseFnDefinition()
		if err != nil {
			return nil, err
		}

		return n, nil
	case "return":
		n, err := p.ParseExplciitReturnStatement()
		if err != nil {
			return nil, err
		}

		return n, nil
	case "var":
		n, err := p.ParseVarDefinition()
		if err != nil {
			return nil, err
		}

		return n, nil
	case "let":
		n, err := p.ParseLetDefinition()
		if err != nil {
			return nil, err
		}

		return n, nil
	}

	return p.ParseExpression()
}

func (p *Parser) ParseExpression() (Expression, error) {
	return p.ParsePrimaryExpression()
}

func (p *Parser) ParsePrimaryExpression() (Expression, error) {
	return p.ParseLiteral()
}

func (p *Parser) ParseLiteral() (Literal, error) {
	p.DeclareParsing(IdentifierLiteralASTNodeKind)

	mtk, err := p.NextToken()
	if err != nil {
		return nil, err
	}

	tk, err := mtk.Value()
	if err != nil {
		return nil, p.GenerateUnexpectedEOFError()
	}

	startpos := tk.Startpos()

	switch tk.Group {
	case &lexer.TokenIdentifierGroup:
		return IdentifierLiteralASTNode{
			Loc:  lexer.InitLocation(startpos, p.lexer.CurrentPos()),
			Name: tk.Characters,
		}, nil
	case &lexer.TokenStringGroup:
		return StringLiteralASTNode{
			Loc:    lexer.InitLocation(startpos, p.lexer.CurrentPos()),
			String: tk.Characters,
		}, nil
	case &lexer.TokenIntegerGroup:
		val, _ := strconv.ParseInt(tk.Characters, 10, 64)

		return IntegerLiteralASTNode{
			Loc:   lexer.InitLocation(startpos, p.lexer.CurrentPos()),
			Value: val,
		}, nil
	case &lexer.TokenDecimalGroup:
		val, _ := strconv.ParseFloat(tk.Characters, 64)

		return DecimalLiteralASTNode{
			Loc:   lexer.InitLocation(startpos, p.lexer.CurrentPos()),
			Value: val,
		}, nil
	}

	return nil, fmt.Errorf("Canot parse this node (it is either not supported or invalid) beginning with token %s", tk.ToDisplayString())
}

func (p *Parser) ParseBlock() (BlockASTNode, error) {
	p.DeclareParsing(BlockASTNodeKind)

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

		stmt, err := p.ParseStatement()
		if err != nil {
			return BlockASTNode{}, err
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
