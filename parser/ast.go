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
	MutableReferenceTypeASTNodeKind
	ImmutableReferenceTypeASTNodeKind
	RawPointerTypeASTNodeKind
	NamedTypeASTNodeKind
	UntaggedUnionTypeASTNodeKind
	NeverTypeASTNodeKind
	TableTypeASTNodeKind
	ArrayTypeASTNodeKind
	SliceTypeASTNodeKind
	ComputedVarDefinitionASTNodeKind
)

type ASTNodeGroup []ASTNodeKind

var (
	TypeASTNodeGroup ASTNodeGroup = ASTNodeGroup{
		MutableReferenceTypeASTNodeKind,
		ImmutableReferenceTypeASTNodeKind,
		RawPointerTypeASTNodeKind,
		NamedTypeASTNodeKind,
		UntaggedUnionTypeASTNodeKind,
		NeverTypeASTNodeKind,
		TableTypeASTNodeKind,
		ArrayTypeASTNodeKind,
		SliceTypeASTNodeKind,
	}

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
		ComputedVarDefinitionASTNodeKind,
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
	TypeGenericASTNode struct {
		ConformsTo []NamedTypeASTNode
	}

	RefType interface {
		Type

		Mutable() bool
		Escaping() bool
		InnerType() Type
	}

	MutableReference struct {
		Loc      Location
		escaping bool
		inner    Type
	}

	ImmutableReference struct {
		Loc      Location
		escaping bool
		inner    Type
	}

	RawPointer struct {
		Loc   Location
		inner Type
	}

	NamedTypeASTNode struct {
		Loc      Location
		Name     string
		Generics map[string]TypeGenericASTNode
	}

	UntaggedUnionTypeASTNode struct {
		Loc   Location
		Types []Type
	}

	NeverTypeASTNode struct {
		Loc Location
	}

	TableTypeASTNode struct {
		Loc Location

		KeyType   Type
		ValueType Type
	}

	ArrayTypeASTNode struct {
		Loc Location

		ValueType Type
		Length    uint64
	}

	SliceTypeASTNode struct {
		Loc Location

		ValueType Type
	}
)

func (ptr RawPointer) typeNode()               {}
func (ref MutableReference) typeNode()         {}
func (ref ImmutableReference) typeNode()       {}
func (typ NamedTypeASTNode) typeNode()         {}
func (typ UntaggedUnionTypeASTNode) typeNode() {}
func (typ NeverTypeASTNode) typeNode()         {}
func (typ TableTypeASTNode) typeNode()         {}
func (typ ArrayTypeASTNode) typeNode()         {}
func (typ SliceTypeASTNode) typeNode()         {}

func (ptr RawPointer) Mutable() bool         { return true }
func (ref MutableReference) Mutable() bool   { return true }
func (ref ImmutableReference) Mutable() bool { return false }

func (ptr RawPointer) Escaping() bool         { return true }
func (ref MutableReference) Escaping() bool   { return ref.escaping }
func (ref ImmutableReference) Escaping() bool { return ref.escaping }

func (ptr RawPointer) InnerType() Type         { return ptr.inner }
func (ref MutableReference) InnerType() Type   { return ref.inner }
func (ref ImmutableReference) InnerType() Type { return ref.inner }

func (ptr RawPointer) Location() Location               { return ptr.Loc }
func (ref MutableReference) Location() Location         { return ref.Loc }
func (ref ImmutableReference) Location() Location       { return ref.Loc }
func (typ NamedTypeASTNode) Location() Location         { return typ.Loc }
func (typ UntaggedUnionTypeASTNode) Location() Location { return typ.Loc }
func (typ NeverTypeASTNode) Location() Location         { return typ.Loc }
func (typ TableTypeASTNode) Location() Location         { return typ.Loc }
func (typ ArrayTypeASTNode) Location() Location         { return typ.Loc }
func (typ SliceTypeASTNode) Location() Location         { return typ.Loc }

func (ptr RawPointer) Kind() ASTNodeKind               { return RawPointerTypeASTNodeKind }
func (ref MutableReference) Kind() ASTNodeKind         { return MutableReferenceTypeASTNodeKind }
func (ref ImmutableReference) Kind() ASTNodeKind       { return ImmutableReferenceTypeASTNodeKind }
func (typ NamedTypeASTNode) Kind() ASTNodeKind         { return NamedTypeASTNodeKind }
func (typ UntaggedUnionTypeASTNode) Kind() ASTNodeKind { return UntaggedUnionTypeASTNodeKind }
func (typ NeverTypeASTNode) Kind() ASTNodeKind         { return NeverTypeASTNodeKind }
func (typ TableTypeASTNode) Kind() ASTNodeKind         { return TableTypeASTNodeKind }
func (typ ArrayTypeASTNode) Kind() ASTNodeKind         { return ArrayTypeASTNodeKind }
func (typ SliceTypeASTNode) Kind() ASTNodeKind         { return SliceTypeASTNodeKind }

func (ptr RawPointer) Group() ASTNodeGroup               { return TypeASTNodeGroup }
func (ref MutableReference) Group() ASTNodeGroup         { return TypeASTNodeGroup }
func (ref ImmutableReference) Group() ASTNodeGroup       { return TypeASTNodeGroup }
func (typ NamedTypeASTNode) Group() ASTNodeGroup         { return TypeASTNodeGroup }
func (typ UntaggedUnionTypeASTNode) Group() ASTNodeGroup { return TypeASTNodeGroup }
func (typ NeverTypeASTNode) Group() ASTNodeGroup         { return TypeASTNodeGroup }
func (typ TableTypeASTNode) Group() ASTNodeGroup         { return TypeASTNodeGroup }
func (typ ArrayTypeASTNode) Group() ASTNodeGroup         { return TypeASTNodeGroup }
func (typ SliceTypeASTNode) Group() ASTNodeGroup         { return TypeASTNodeGroup }

type (
	ImportStatementASTNode struct {
		Loc  Location
		Name string
	}

	ConstDefinitionASTNode struct {
		Loc   Location
		Name  string
		Value Expression
		Type  utils.Optional[Type]
	}

	VarDefinitionASTNode struct {
		Loc   Location
		Name  string
		Value Expression
		Type  utils.Optional[Type]
	}

	LetDefinitionASTNode struct {
		Loc   Location
		Name  string
		Value Expression
		Type  utils.Optional[Type]
	}

	VariableUsageExpresionASTNode struct {
		Loc  Location
		Name string
	}

	BinaryExpressionASTNode struct {
		Loc      Location
		Operator string
		Left     Expression
		Right    Expression
	}

	PostfixUnaryExpressionASTNode struct {
		Loc      Location
		Operator string
		Left     Expression
	}

	PrefixUnaryExpressionASTNode struct {
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

	StructureDefinitionASTNode struct {
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

	ClassDefinitionASTNode struct {
		Loc          Location
		Name         string
		Fields       ClassDefFields
		Methods      ClassDefMethods
		Constructors ClassDefConstructors
	}

	FunctionDefinitionASTNode struct {
		Loc        Location
		Name       string
		ReturnType Type
		Parameters map[string]Type
		Generics   map[string]TypeGenericASTNode
		Body       []ASTNode
	}

	MethodDefinitionASTNode struct {
		Loc         Location
		Name        string
		ReturnType  Type
		Parameters  map[string]Type
		ContextType Type
		Generics    map[string]TypeGenericASTNode
		Body        []ASTNode
	}

	OperatorOverloadASTNode struct {
		Loc           Location
		Operator      string
		ReturnType    Type
		ContextType   Type
		RightHandType Type
		Body          []ASTNode
	}

	AssignmentStatementASTNode struct {
		Loc   Location
		Left  Expression
		Right Expression
	}

	StructureInitilisationExpressionASTNode struct {
		Loc         Location
		StructName  string
		FieldValues map[string]Expression
	}

	StructureRefInitilisationExpressionASTNode struct {
		Loc         Location
		StructName  string
		FieldValues map[string]Expression
		RefType     RefType
	}

	ImplicitReturnASTNode struct {
		Loc   Location
		Value Expression
	}

	ExplicitReturnASTNode struct {
		Loc   Location
		Value Expression
	}

	FunctionCallExpressionASTNode struct {
		Loc        Location
		Name       string
		Parameters map[string]Expression
		Generics   map[string]TypeGenericASTNode
	}

	MethodCallExpressionASTNode struct {
		Loc        Location
		Context    Type
		Name       string
		Parameters map[string]Expression
		Generics   map[string]TypeGenericASTNode
	}

	MemberExpressionASTNode struct {
		Loc      Location
		Segments []Expression
	}

	ModulePathASTNode struct {
		Loc      Location
		Segments []Expression
	}

	LambdaExpressionASTNode struct {
		Loc        Location
		Parameters map[string]Type
		ReturnType Type
		Body       []ASTNode
	}

	IfExpressionASTNode struct {
		Loc          Location
		Condition    Expression
		Body         []ASTNode
		FallbackBody []ASTNode
	}

	IfStatementASTNode struct {
		Loc       Location
		Condition Expression
		Body      []ASTNode
	}

	SwitchStatementCase struct {
		Loc  Location
		Body []ASTNode
	}

	SwitchStatementASTNode struct {
		Loc   Location
		Value Expression
		Cases map[Expression]SwitchStatementCase
	}

	MatchExpressionCase struct {
		Loc  Location
		Code []ASTNode
	}

	MatchOrWhenExpressionCase struct {
		Loc  Location
		Code []ASTNode
	}

	MatchExpressionASTNode struct {
		Loc          Location
		Value        Expression
		Cases        map[TupleDestructuringASTNode]MatchOrWhenExpressionCase
		FallbackCase MatchOrWhenExpressionCase
	}

	WhenExpressionASTNode struct {
		Loc          Location
		Cases        map[Expression]MatchOrWhenExpressionCase
		FallbackCase MatchOrWhenExpressionCase
	}

	InterfaceDefField struct {
		Loc        Location
		IsComputed bool
	}

	InterfaceDefinitionASTNode struct {
		Loc       Location
		SubTypeOf NamedTypeASTNode
		Fields    map[string]InterfaceDefField
	}

	StringLiteralASTNode struct{}

	ArrayLiteralASTNode struct{}

	IntegerLiteralASTNode struct{}

	DecimalLiteralASTNode struct{}

	CStyleEnumDefinitionASTNode struct{}

	SumTypeEnumDefinitionASTNode struct{}

	NamespaceDefinitionASTNode struct{}

	ExternalFnDeclarationASTNode struct{}

	CStyleForLoopStatementASTNode struct{}

	ForInLoopStatementASTNode struct{}

	WhileLoopStatementASTNode struct{}

	ForeverLoopStatementASTNode struct{}

	TernaryExpressionASTNode struct{}

	OptionalChainingASTNode struct{}

	OptionalChainingCallASTNode struct{}

	TypeCastableQueryExpressionASTNode struct{}

	TypeCastExpressionASTNode struct{}

	RuntimeTypeCastExpressionASTNode struct{}

	InternalMacroDeclarationASTNode struct{}

	CustomMacroDefinitionASTNode struct{}

	MacroUsageASTNode struct{}

	MacroVariableUsageASTNode struct{}

	BlockASTNode struct{}

	TupleDestructuringASTNode struct{}

	ArrayCompTimeDestructuringASTNode struct{}

	ArrayRuntimeDestructuringASTNode struct{}

	StructOrClassDestructuringASTNode struct{}

	ReferenceClassDestructuringASTNode struct{}

	ConstraintASTNode struct{}

	IfLetStatementASTNode struct{}

	IfVarStatementASTNode struct{}

	IfLetExpressionASTNode struct{}

	IfVarExpressionASTNode struct{}

	NullCoalesceExpressionASTNode struct{}

	BubbleValueToReturnASTNode struct{}
)
