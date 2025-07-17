package parser

import (
	"ljpprojects.org/sqopl/lexer"
	"ljpprojects.org/sqopl/utils"
)

type ASTNodeKind uint8

const (
	ImportStatementASTNodeKind ASTNodeKind = iota
	ConstDefinitionASTNodeKind
	VarDefinitionASTNodeKind
	LetDefinitionASTNodeKind
	VariableUsageExpresionASTNodeKind
	BinaryExpressionASTNodeKind
	PostfixUnaryExpressionASTNodeKind
	PrefixUnaryExpressionASTNodeKind
	StructureDefinitionASTNodeKind
	ClassDefinitionASTNodeKind
	FunctionDefinitionASTNodeKind
	MethodDefinitionASTNodeKind
	OperatorOverloadASTNodeKind
	AssignmentStatementASTNodeKind
	StructureInitilisationExpressionASTNodeKind
	StructureRefInitilisationExpressionASTNodeKind
	ImplicitReturnASTNodeKind
	ExplicitReturnASTNodeKind
	FunctionCallExpressionASTNodeKind
	MethodCallExpressionASTNodeKind
	MemberExpressionASTNodeKind
	ModulePathASTNodeKind
	LambdaExpressionASTNodeKind
	IfExpressionASTNodeKind
	IfStatementASTNodeKind
	SwitchStatementASTNodeKind
	MatchExpressionASTNodeKind
	WhenExpressionASTNodeKind
	InterfaceDefinitionASTNodeKind
	StringLiteralASTNodeKind
	ArrayLiteralASTNodeKind
	IntegerLiteralASTNodeKind
	DecimalLiteralASTNodeKind
	CStyleEnumDefinitionASTNodeKind
	SumTypeEnumDefinitionASTNodeKind
	NamespaceDefinitionASTNodeKind
	ExternalFnDeclarationASTNodeKind
	CStyleForLoopStatementASTNodeKind
	ForInLoopStatementASTNodeKind
	WhileLoopStatementASTNodeKind
	ForeverLoopStatementASTNodeKind
	TernaryExpressionASTNodeKind
	OptionalChainingASTNodeKind
	OptionalChainingCallASTNodeKind
	TypeCastableQueryExpressionASTNodeKind
	TypeCastExpressionASTNodeKind
	RuntimeTypeCastExpressionASTNodeKind
	InternalMacroDeclarationASTNodeKind
	CustomMacroDefinitionASTNodeKind
	MacroUsageASTNodeKind
	MacroVariableUsageASTNodeKind
	BlockASTNodeKind
	TupleDestructuringASTNodeKind
	ArrayCompTimeDestructuringASTNodeKind
	ArrayRuntimeDestructuringASTNodeKind
	StructOrClassDestructuringASTNodeKind
	ReferenceClassDestructuringASTNodeKind
	ConstraintASTNodeKind
	IfLetStatementASTNodeKind
	IfVarStatementASTNodeKind
	IfLetExpressionASTNodeKind
	IfVarExpressionASTNodeKind
	NullCoalesceExpressionASTNodeKind
	BubbleValueToReturnASTNodeKind
)

type ASTNodeGroup []ASTNodeKind

var (
	MetaASTNodeGroup ASTNodeGroup = ASTNodeGroup{
		MacroUsageASTNodeKind,
		MacroVariableUsageASTNodeKind,
	}

	SegmentASTNodeGroup ASTNodeGroup = ASTNodeGroup{
		TupleDestructuringASTNodeKind,
		ArrayCompTimeDestructuringASTNodeKind,
		ArrayRuntimeDestructuringASTNodeKind,
		StructOrClassDestructuringASTNodeKind,
		ReferenceClassDestructuringASTNodeKind,
		ConstraintASTNodeKind,
	}

	StatementASTNodeGroup ASTNodeGroup = ASTNodeGroup{
		ImportStatementASTNodeKind,
		AssignmentStatementASTNodeKind,
		ImplicitReturnASTNodeKind,
		ExplicitReturnASTNodeKind,
		IfStatementASTNodeKind,
		SwitchStatementASTNodeKind,
		CStyleForLoopStatementASTNodeKind,
		ForInLoopStatementASTNodeKind,
		WhileLoopStatementASTNodeKind,
		ForInLoopStatementASTNodeKind,
		IfLetStatementASTNodeKind,
		IfVarStatementASTNodeKind,
	}

	DefinitionASTNodeGroup ASTNodeGroup = ASTNodeGroup{
		ConstDefinitionASTNodeKind,
		VarDefinitionASTNodeKind,
		LetDefinitionASTNodeKind,
		StructureDefinitionASTNodeKind,
		ClassDefinitionASTNodeKind,
		FunctionDefinitionASTNodeKind,
		MethodDefinitionASTNodeKind,
		OperatorOverloadASTNodeKind,
		InterfaceDefinitionASTNodeKind,
		CStyleEnumDefinitionASTNodeKind,
		SumTypeEnumDefinitionASTNodeKind,
		NamespaceDefinitionASTNodeKind,
		CustomMacroDefinitionASTNodeKind,
	}

	DeclarationASTNodeGroup ASTNodeGroup = ASTNodeGroup{
		ExternalFnDeclarationASTNodeKind,
		InternalMacroDeclarationASTNodeKind,
	}

	LiteralASTNodeGroup ASTNodeGroup = ASTNodeGroup{
		StringLiteralASTNodeKind,
		ArrayLiteralASTNodeKind,
		IntegerLiteralASTNodeKind,
		DecimalLiteralASTNodeKind,
	}

	ExpressionASTNodeGroup ASTNodeGroup = ASTNodeGroup{
		VariableUsageExpresionASTNodeKind,
		BinaryExpressionASTNodeKind,
		PostfixUnaryExpressionASTNodeKind,
		PrefixUnaryExpressionASTNodeKind,
		StructureInitilisationExpressionASTNodeKind,
		StructureRefInitilisationExpressionASTNodeKind,
		FunctionCallExpressionASTNodeKind,
		MethodCallExpressionASTNodeKind,
		MemberExpressionASTNodeKind,
		ModulePathASTNodeKind,
		LambdaExpressionASTNodeKind,
		IfExpressionASTNodeKind,
		MatchExpressionASTNodeKind,
		WhenExpressionASTNodeKind,
		TernaryExpressionASTNodeKind,
		OptionalChainingASTNodeKind,
		OptionalChainingCallASTNodeKind,
		TypeCastableQueryExpressionASTNodeKind,
		TypeCastExpressionASTNodeKind,
		RuntimeTypeCastExpressionASTNodeKind,
		BlockASTNodeKind,
		IfLetExpressionASTNodeKind,
		IfVarExpressionASTNodeKind,
		NullCoalesceExpressionASTNodeKind,
		BubbleValueToReturnASTNodeKind,
	}
)

type Location utils.Range[lexer.Position]

type ASTNode interface {
	Location() Location
	Group() ASTNodeGroup
	Kind() ASTNodeKind
}

type Statement interface {
	ASTNode

	statementNode()
}

type Declaration interface {
	ASTNode

	declarationNode()
}

type Definition interface {
	ASTNode

	definitionNode()
}

type Meta interface {
	ASTNode

	metaNode()
}

type Expression interface {
	ASTNode

	expressionNode()
}

type Literal interface {
	Expression

	literalNode()
}

type Type interface {
	ASTNode

	typeNode()
}

type (
	ImportStatementASTNodeKind struct {
		Loc  Location
		Name string
	}

	ConstDefinitionASTNodeKind struct {
		Loc   Location
		Name  string
		Value Expression
		Type  utils.Optional[Type]
	}

	VarDefinitionASTNodeKind struct {
		Loc   Location
		Name  string
		Value Expression
		Type  utils.Optional[Type]
	}

	LetDefinitionASTNodeKind struct {
		Loc   Location
		Name  string
		Value Expression
		Type  utils.Optional[Type]
	}

	VariableUsageExpresionASTNodeKind struct {
		Loc  Location
		Name string
	}

	BinaryExpressionASTNodeKind struct {
		Loc      Location
		Operator string
		Left     Expression
		Right    Expression
	}

	PostfixUnaryExpressionASTNodeKind struct {
		Loc      Location
		Operator string
		Left     Expression
	}

	PrefixUnaryExpressionASTNodeKind struct {
		Loc      Location
		Operator string
		Right    Expression
	}

	StructureDefFields []struct {
		Loc       Location
		Name      string
		IsMutable bool
		Type      Type
	}

	StructureDefinitionASTNodeKind struct {
		Loc    Location
		Name   string
		Fields StructureDefFields
	}

	ClassDefFields []struct {
		Loc       Location
		Name      string
		IsMutable bool
		Type      Type
	}

	ClassDefMethods []struct {
		Loc           Location
		Name          string
		ReturnType    Type
		ArgumentTypes []Type
		SelfType      utils.Optional[Type]
	}

	ClassDefConstructors []struct {
		Loc           Location
		Name          string
		MayReturnNull bool
		ArgumentTypes []Type
	}

	ClassDefinitionASTNodeKind struct {
		Loc          Location
		Name         string
		Fields       ClassDefFields
		Methods      ClassDefMethods
		Constructors ClassDefConstructors
	}

	FunctionDefinitionASTNodeKind struct{}

	MethodDefinitionASTNodeKind struct{}

	OperatorOverloadASTNodeKind struct{}

	AssignmentStatementASTNodeKind struct{}

	StructureInitilisationExpressionASTNodeKind struct{}

	StructureRefInitilisationExpressionASTNodeKind struct{}

	ImplicitReturnASTNodeKind struct{}

	ExplicitReturnASTNodeKind struct{}

	FunctionCallExpressionASTNodeKind struct{}

	MethodCallExpressionASTNodeKind struct{}

	MemberExpressionASTNodeKind struct{}

	ModulePathASTNodeKind struct{}

	LambdaExpressionASTNodeKind struct{}

	IfExpressionASTNodeKind struct{}

	IfStatementASTNodeKind struct{}

	SwitchStatementASTNodeKind struct{}

	MatchExpressionASTNodeKind struct{}

	WhenExpressionASTNodeKind struct{}

	InterfaceDefinitionASTNodeKind struct{}

	StringLiteralASTNodeKind struct{}

	ArrayLiteralASTNodeKind struct{}

	IntegerLiteralASTNodeKind struct{}

	DecimalLiteralASTNodeKind struct{}

	CStyleEnumDefinitionASTNodeKind struct{}

	SumTypeEnumDefinitionASTNodeKind struct{}

	NamespaceDefinitionASTNodeKind struct{}

	ExternalFnDeclarationASTNodeKind struct{}

	CStyleForLoopStatementASTNodeKind struct{}

	ForInLoopStatementASTNodeKind struct{}

	WhileLoopStatementASTNodeKind struct{}

	ForeverLoopStatementASTNodeKind struct{}

	TernaryExpressionASTNodeKind struct{}

	OptionalChainingASTNodeKind struct{}

	OptionalChainingCallASTNodeKind struct{}

	TypeCastableQueryExpressionASTNodeKind struct{}

	TypeCastExpressionASTNodeKind struct{}

	RuntimeTypeCastExpressionASTNodeKind struct{}

	InternalMacroDeclarationASTNodeKind struct{}

	CustomMacroDefinitionASTNodeKind struct{}

	MacroUsageASTNodeKind struct{}

	MacroVariableUsageASTNodeKind struct{}

	BlockASTNodeKind struct{}

	TupleDestructuringASTNodeKind struct{}

	ArrayCompTimeDestructuringASTNodeKind struct{}

	ArrayRuntimeDestructuringASTNodeKind struct{}

	StructOrClassDestructuringASTNodeKind struct{}

	ReferenceClassDestructuringASTNodeKind struct{}

	ConstraintASTNodeKind struct{}

	IfLetStatementASTNodeKind struct{}

	IfVarStatementASTNodeKind struct{}

	IfLetExpressionASTNodeKind struct{}

	IfVarExpressionASTNodeKind struct{}

	NullCoalesceExpressionASTNodeKind struct{}

	BubbleValueToReturnASTNodeKind struct{}
)
