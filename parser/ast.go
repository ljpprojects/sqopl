package parser

import (
	"ljpprojects.org/sqopl/lexer"
	"ljpprojects.org/sqopl/utils"
)

// An enum representing the different kinds of AST nodes.
type ASTNodeKind uint8

const (
	// The AST node kind representing the import statement.
	ImportStatementASTNodeKind ASTNodeKind = iota

	// The AST node kind representing the compile-time-constant definition
	ConstDefinitionASTNodeKind

	// The AST node kind representing the mutable variable definition.
	VarDefinitionASTNodeKind

	// The AST node kind representing the immutable variable definition.
	LetDefinitionASTNodeKind

	// The AST node kind representing the identifier literal.
	IdentifierLiteralASTNodeKind

	// The AST node kind representing the binary expression.
	BinaryExpressionASTNodeKind

	// The AST node kind representing the postfix unary expression.
	PostfixUnaryExpressionASTNodeKind

	// The AST node kind representing the prefix unary expression.
	PrefixUnaryExpressionASTNodeKind

	// The AST node kind representing the structure definition.
	StructureDefinitionASTNodeKind

	// The AST node kind representing the class definition.
	ClassDefinitionASTNodeKind

	// The AST node kind representing the function definition.
	FunctionDefinitionASTNodeKind

	// The AST node kind representing the method definition.
	MethodDefinitionASTNodeKind

	// The AST node kind representing the operator definition.
	OperatorOverloadASTNodeKind

	// The AST node kind representing the assignment statement.
	AssignmentStatementASTNodeKind

	// The AST node kind representing the structure initialisation expression.
	StructureInitilisationExpressionASTNodeKind

	// The AST node kind representing the structure allocation (structure reference initialisation) expression.
	StructureRefInitilisationExpressionASTNodeKind

	// The AST node kind representing the implicit return statement.
	ImplicitReturnASTNodeKind

	// The AST node kind representing the explicit return statement.
	ExplicitReturnASTNodeKind

	// The AST node kind representing the function call expression.
	FunctionCallExpressionASTNodeKind

	// The AST node kind representing the method call expression.
	MethodCallExpressionASTNodeKind

	// The AST node kind representing the member expression.
	MemberExpressionASTNodeKind

	// The AST node kind representing the module path expression.
	ModulePathASTNodeKind

	// The AST node kind representing the lambda expression.
	LambdaExpressionASTNodeKind

	// The AST node kind representing the if expression (an if with an else clause).
	IfExpressionASTNodeKind

	// The AST node kind representing the if statement (an if with no else clause).
	IfStatementASTNodeKind

	// The AST node kind representing the switch statement.
	SwitchStatementASTNodeKind

	// The AST node kind representing the match expression.
	MatchExpressionASTNodeKind

	// The AST node kind representing the when expression.
	WhenExpressionASTNodeKind

	// The AST node kind representing the interface definition.
	InterfaceDefinitionASTNodeKind

	// The AST node kind representing the string literal.
	StringLiteralASTNodeKind

	// The AST node kind representing the array literal.
	ArrayLiteralASTNodeKind

	// The AST node kind representing the integer literal.
	IntegerLiteralASTNodeKind

	// The AST node kind representing the decimal literal.
	DecimalLiteralASTNodeKind

	// The AST node kind representing the C-style enum definition.
	CStyleEnumDefinitionASTNodeKind

	// The AST node kind representing the sum type (enum) definition.
	SumTypeEnumDefinitionASTNodeKind

	// The AST node kind representing the namespace definition.
	NamespaceDefinitionASTNodeKind

	// The AST node kind representing the external function declaration.
	ExternalFnDeclarationASTNodeKind

	// The AST node kind representing the C-style for loop (the for loop form with the definition, condition, and incrementor statements) statement.
	CStyleForLoopStatementASTNodeKind

	// The AST node kind representing the for-in (the iterative for loop) statement.
	ForInLoopStatementASTNodeKind

	// The AST node kind representing the while loop statement.
	WhileLoopStatementASTNodeKind

	// The AST node kind representing the infinite loop statement.
	ForeverLoopStatementASTNodeKind

	// The AST node kind representing the ternary expression.
	TernaryExpressionASTNodeKind

	// The AST node kind representing the optional chaining expression.
	OptionalChainingASTNodeKind

	// The AST node kind representing the is-type-castable-to query expression.
	TypeCastableQueryExpressionASTNodeKind

	// The AST node kind representing the type cast expression.
	TypeCastExpressionASTNodeKind

	// The AST node kind representing the at-runtime type cast expression.
	RuntimeTypeCastExpressionASTNodeKind

	// The AST node kind representing the internal macro declaration.
	InternalMacroDeclarationASTNodeKind

	// The AST node kind representing the custom macro definition.
	CustomMacroDefinitionASTNodeKind

	// The AST node kind representing the macro usage meta-node.
	MacroUsageASTNodeKind

	// The AST node kind representing the macro variable usage meta-node.
	MacroVariableUsageASTNodeKind

	// The AST node kind representing the block expression.
	BlockASTNodeKind

	// The AST node kind representing the tuple destructuring segment-node.
	TupleDestructuringASTNodeKind

	// The AST node kind representing the array (of known length) destructuring segment-node.
	ArrayCompTimeDestructuringASTNodeKind

	// The AST node kind representing the array (of unknown length) destructuring segment-node.
	ArrayRuntimeDestructuringASTNodeKind

	// The AST node kind representing the structure/class destructuring segment-node.
	StructOrClassDestructuringASTNodeKind

	// The AST node kind representing the reference destructuring segment-node.
	ReferenceDestructuringASTNodeKind

	// The AST node kind representing the constraint segment-node.
	ConstraintASTNodeKind

	// The AST node kind representing the if-let statement (an if-let with no else clause).
	IfLetStatementASTNodeKind

	// The AST node kind representing the if-var statement (an if-var with no else clause).
	IfVarStatementASTNodeKind

	// The AST node kind representing the if-let expression (an if-let with an else clause).
	IfLetExpressionASTNodeKind

	// The AST node kind representing the if-var expression (an if-var with an else clause).
	IfVarExpressionASTNodeKind

	// The AST node kind representing the null coalesce expression.
	NullCoalesceExpressionASTNodeKind

	// The AST node kind representing the buuble value to return expression.
	BubbleValueToReturnASTNodeKind

	// The AST node kind representing the mutable reference type.
	MutableReferenceTypeASTNodeKind

	// The AST node kind representing the immutable reference type.
	ImmutableReferenceTypeASTNodeKind

	// The AST node kind representing the raw pointer type.
	RawPointerTypeASTNodeKind

	// The AST node kind representing the named type.
	NamedTypeASTNodeKind

	// The AST node kind representing the untagged union type.
	UntaggedUnionTypeASTNodeKind

	// The AST node kind representing the never (!) type.
	NeverTypeASTNodeKind

	// The AST node kind representing the table type.
	TableTypeASTNodeKind

	// The AST node kind representing the array (fixed length) type.
	ArrayTypeASTNodeKind

	// The AST node kind representing the slice (dynamic-length array, a reference) type.
	SliceTypeASTNodeKind

	// The AST node kind representing the tuple type.
	TupleTypeASTNodeKind

	// The AST node kind representing the computed variable definition.
	ComputedVarDefinitionASTNodeKind
)

// Formats an ASTNodeKind into a string for printing out.
func (k ASTNodeKind) ToDisplayString() string {
	switch k {
	case ImportStatementASTNodeKind:
		return "Kind(Import Statement)"
	case ConstDefinitionASTNodeKind:
		return "Kind(Const Definition)"
	case VarDefinitionASTNodeKind:
		return "Kind(Var Definition)"
	case LetDefinitionASTNodeKind:
		return "Kind(Let Definition)"
	case IdentifierLiteralASTNodeKind:
		return "Kind(Identifier Literal)"
	case BinaryExpressionASTNodeKind:
		return "Kind(Binary Expression)"
	case PostfixUnaryExpressionASTNodeKind:
		return "Kind(Postfix Unary Expression)"
	case PrefixUnaryExpressionASTNodeKind:
		return "Kind(Prefix Unary Expression)"
	case StructureDefinitionASTNodeKind:
		return "Kind(Structure Definition)"
	case ClassDefinitionASTNodeKind:
		return "Kind(Class Definition)"
	case FunctionDefinitionASTNodeKind:
		return "Kind(Function Definition)"
	case MethodDefinitionASTNodeKind:
		return "Kind(Method Definition)"
	case OperatorOverloadASTNodeKind:
		return "Kind(Operator Overload Definition)"
	case AssignmentStatementASTNodeKind:
		return "Kind(Assignment Statement)"
	case StructureInitilisationExpressionASTNodeKind:
		return "Kind(Structure Initialisation Expression)"
	case StructureRefInitilisationExpressionASTNodeKind:
		return "Kind(Structure Ref Initialisation Expression)"
	case ImplicitReturnASTNodeKind:
		return "Kind(Implicit Return Expression)"
	case ExplicitReturnASTNodeKind:
		return "Kind(Explicit Return Expression)"
	case FunctionCallExpressionASTNodeKind:
		return "Kind(Function Call Expression)"
	case MethodCallExpressionASTNodeKind:
		return "Kind(Method Call Expression)"
	case MemberExpressionASTNodeKind:
		return "Kind(Member Expression)"
	case ModulePathASTNodeKind:
		return "Kind(Module Path Expression)"
	case LambdaExpressionASTNodeKind:
		return "Kind(Lambda Expression)"
	case IfExpressionASTNodeKind:
		return "Kind(If Expression)"
	case IfStatementASTNodeKind:
		return "Kind(If Statement)"
	case SwitchStatementASTNodeKind:
		return "Kind(Switch Statement)"
	case MatchExpressionASTNodeKind:
		return "Kind(Match Expression)"
	case WhenExpressionASTNodeKind:
		return "Kind(When Expression)"
	case InterfaceDefinitionASTNodeKind:
		return "Kind(Interface Defintion)"
	case StringLiteralASTNodeKind:
		return "Kind(String Literal)"
	case ArrayLiteralASTNodeKind:
		return "Kind(Array Literal)"
	case IntegerLiteralASTNodeKind:
		return "Kind(Integer Literal)"
	case DecimalLiteralASTNodeKind:
		return "Kind(Decimal Literal)"
	case CStyleEnumDefinitionASTNodeKind:
		return "Kind(C-Style Enum Defintion)"
	case SumTypeEnumDefinitionASTNodeKind:
		return "Kind(Sum Type Enum Defintion)"
	case NamespaceDefinitionASTNodeKind:
		return "Kind(Namespace Defintion)"
	case ExternalFnDeclarationASTNodeKind:
		return "Kind(External Function Declaration)"
	case CStyleForLoopStatementASTNodeKind:
		return "Kind(C-Style For Loop Statement)"
	case ForInLoopStatementASTNodeKind:
		return "Kind(For-In Loop Statement)"
	case WhileLoopStatementASTNodeKind:
		return "Kind(While Loop Statement)"
	case ForeverLoopStatementASTNodeKind:
		return "Kind(Forever Loop Statement)"
	case TernaryExpressionASTNodeKind:
		return "Kind(Ternary Expression)"
	case OptionalChainingASTNodeKind:
		return "Kind(Optional Chaining Expression)"
	case TypeCastableQueryExpressionASTNodeKind:
		return "Kind(Type-Is-Castable-To Query Expression)"
	case TypeCastExpressionASTNodeKind:
		return "Kind(Type Cast Expression)"
	case RuntimeTypeCastExpressionASTNodeKind:
		return "Kind(Runtime Type Cast Expression)"
	case InternalMacroDeclarationASTNodeKind:
		return "Kind(Internal MAcro Declaration)"
	case CustomMacroDefinitionASTNodeKind:
		return "Kind(Custom Macro Definition)"
	case MacroUsageASTNodeKind:
		return "Kind(Macro Usage Node)"
	case MacroVariableUsageASTNodeKind:
		return "Kind(Macro Variable Usage Node)"
	case BlockASTNodeKind:
		return "Kind(Block Expression)"
	case TupleDestructuringASTNodeKind:
		return "Kind(Tuple Destructuring Segment)"
	case ArrayCompTimeDestructuringASTNodeKind:
		return "Kind(Array Known-Length Destructuring Segment)"
	case ArrayRuntimeDestructuringASTNodeKind:
		return "Kind(Array Unknown-Length Destructuring Segment)"
	case StructOrClassDestructuringASTNodeKind:
		return "Kind(Struct/Class Destructuring Segment)"
	case ReferenceDestructuringASTNodeKind:
		return "Kind(Reference Destructuring Segment)"
	case ConstraintASTNodeKind:
		return "Kind(Constraint Segment)"
	case IfLetStatementASTNodeKind:
		return "Kind(If-Let Statement)"
	case IfVarStatementASTNodeKind:
		return "Kind(If-Var Statement)"
	case IfLetExpressionASTNodeKind:
		return "Kind(If-Let Expression)"
	case IfVarExpressionASTNodeKind:
		return "Kind(If-Var Expression)"
	case NullCoalesceExpressionASTNodeKind:
		return "Kind(Null Coalesce Expression)"
	case BubbleValueToReturnASTNodeKind:
		return "Kind(Bubble To Return Expression)"
	case MutableReferenceTypeASTNodeKind:
		return "Kind(Mutable Reference Type)"
	case ImmutableReferenceTypeASTNodeKind:
		return "Kind(Immutable Reference Type)"
	case RawPointerTypeASTNodeKind:
		return "Kind(Raw Pointer Type)"
	case NamedTypeASTNodeKind:
		return "Kind(Named Type)"
	case UntaggedUnionTypeASTNodeKind:
		return "Kind(Untagged Union Type)"
	case NeverTypeASTNodeKind:
		return "Kind(Never Type)"
	case TableTypeASTNodeKind:
		return "Kind(Table Type)"
	case ArrayTypeASTNodeKind:
		return "Kind(Array Type)"
	case SliceTypeASTNodeKind:
		return "Kind(Slice Type)"
	case TupleTypeASTNodeKind:
		return "Kind(Tuple Type)"
	case ComputedVarDefinitionASTNodeKind:
		return "Kind(Computed Var Definition)"
	}

	return "Unknown"
}

// Identical to lexer.TokenGroup, just with ASTNodeKinds.
type ASTNodeGroup []ASTNodeKind

var (
	// The AST node group containing all type nodes.
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
		TupleTypeASTNodeKind,
	}

	// The AST node group containing all meta nodes.
	MetaASTNodeGroup ASTNodeGroup = ASTNodeGroup{
		MacroUsageASTNodeKind,
		MacroVariableUsageASTNodeKind,
	}

	// The AST node group containing all segment nodes.
	SegmentASTNodeGroup ASTNodeGroup = ASTNodeGroup{
		TupleDestructuringASTNodeKind,
		ArrayCompTimeDestructuringASTNodeKind,
		ArrayRuntimeDestructuringASTNodeKind,
		StructOrClassDestructuringASTNodeKind,
		ReferenceDestructuringASTNodeKind,
		ConstraintASTNodeKind,
		BlockASTNodeKind,
	}

	// The AST node group containing all statement nodes.
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

	// The AST node group containing all definition nodes.
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

	// The AST node group containing all declaration nodes.
	DeclarationASTNodeGroup ASTNodeGroup = ASTNodeGroup{
		ExternalFnDeclarationASTNodeKind,
		InternalMacroDeclarationASTNodeKind,
	}

	// The AST node group containing all literal nodes.
	LiteralASTNodeGroup ASTNodeGroup = ASTNodeGroup{
		StringLiteralASTNodeKind,
		ArrayLiteralASTNodeKind,
		IntegerLiteralASTNodeKind,
		DecimalLiteralASTNodeKind,
		IdentifierLiteralASTNodeKind,
	}

	// The AST node group containing all expression nodes.
	ExpressionASTNodeGroup ASTNodeGroup = ASTNodeGroup{
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

// The interface which all valid AST nodes inheret from.
type ASTNode interface {
	// Gets the lexer.Location (alias for utils.Range[lexer.Position]) of the ASTNode
	Location() lexer.Location

	// Gets the group the ASTNode belongs to.
	Group() ASTNodeGroup

	// Gets the kind of ASTNode.
	// The ASTNodeKind is guarenteed to be unique for different kinds of ASTNodes.
	Kind() ASTNodeKind
}

// The interface which all statements inherit from.
type Statement interface {
	ASTNode

	// Stub to make compiler interpret this as a different ASTNode
	statementNode()
}

// The interface which all declarations and definitions (they are different in SQOPL) inherit from.
type Declaration interface {
	Statement

	// Stub to make compiler interpret this as a different ASTNode
	declarationNode()
}

// The interface which all definitions inherit from.
type Definition interface {
	Declaration

	// Stub to make compiler interpret this as a different ASTNode
	definitionNode()
}

// The interface which all 'meta' nodes (e.g. macros) inherit from.
type Meta interface {
	ASTNode

	// Stub to make compiler interpret this as a different ASTNode
	metaNode()
}

// The interface which all ecpressions inherit from.
// Expressions are a specialised subset of statements that return values.
type Expression interface {
	Statement

	// Stub to make compiler interpret this as a different ASTNode
	expressionNode()
}

// The interface which all literals inherit from.
type Literal interface {
	Expression

	// Stub to make compiler interpret this as a different ASTNode
	literalNode()
}

// The interface which all identifiers inherit from.
// An identifier is any literal that can identify data, not just the raw literal.
type Identifier interface {
	Literal

	// Stub to make compiler interpret this as a different ASTNode
	identifierNode()
}

// The interface which all components of other ASTNodes inherit from.
type Component interface {
	ASTNode

	// Stub to make compiler interpret this as a different ASTNode
	componentNode()
}

// The interface which all destructuring components inherit from.
type Destructure interface {
	Component

	// Stub to make compiler interpret this as a different ASTNode
	destructureNode()
}

// The interface which all constraint components inherit from.
type Constraint interface {
	Component

	// Stub to make compiler interpret this as a different ASTNode
	constraintNode()
}

// The interface which all types inherit from.
type Type interface {
	ASTNode

	// Stub to make compiler interpret this as a different ASTNode
	typeNode()
}

type (
	// A component used in declarations/types to specify compile-time generics.
	TypeGenericASTNode struct {
		ConformsTo []NamedTypeASTNode
	}

	// The interface which all reference types inherit from.
	RefType interface {
		Type

		Mutable() bool
		Escaping() bool
		InnerType() Type
		IsDyn() bool
	}

	// The mutable reference type.
	MutableReference struct {
		Loc        lexer.Location
		IsEscaping bool
		Inner      Type
		IsDynamic  bool
	}

	// The immutable reference type.
	ImmutableReference struct {
		Loc        lexer.Location
		IsEscaping bool
		Inner      Type
		IsDynamic  bool
	}

	// The slice reference type.
	// The only reference type that cannot be dynamic.
	SliceTypeASTNode struct {
		Loc        lexer.Location
		ValueType  Type
		IsMutable  bool
		IsEscaping bool
	}

	// The raw pointer type.
	RawPointer struct {
		Loc       lexer.Location
		Inner     Type
		IsDynamic bool
	}

	// Any type represented with an identifier.
	// Types as parts of modules are mangled.
	NamedTypeASTNode struct {
		Loc      lexer.Location
		Name     string
		Generics map[string]TypeGenericASTNode
	}

	// The untagged union type.
	UntaggedUnionTypeASTNode struct {
		Loc   lexer.Location
		Types []Type
	}

	// The never type, which is primarily returned by functions which never return.
	NeverTypeASTNode struct {
		Loc lexer.Location
	}

	// The table type.
	TableTypeASTNode struct {
		Loc lexer.Location

		KeyType   Type
		ValueType Type
	}

	// The array type. Unlike a slice, this is not a reference type because it cannot grow or shrink.
	ArrayTypeASTNode struct {
		Loc lexer.Location

		ValueType Type
		Length    uint64
	}

	// The tuple type.
	// Although identical to the untagged union in structure, all elements specified in ValueTypes are present, not just one.
	TupleTypeASTNode struct {
		Loc        lexer.Location
		ValueTypes []Type
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
func (typ TupleTypeASTNode) typeNode()         {}

func (ptr RawPointer) Mutable() bool         { return true }
func (ref MutableReference) Mutable() bool   { return true }
func (slc SliceTypeASTNode) Mutable() bool   { return slc.IsMutable }
func (ref ImmutableReference) Mutable() bool { return false }

func (ptr RawPointer) Escaping() bool         { return true }
func (ref MutableReference) Escaping() bool   { return ref.IsEscaping }
func (slc SliceTypeASTNode) Escaping() bool   { return slc.IsEscaping }
func (ref ImmutableReference) Escaping() bool { return ref.IsEscaping }

func (ptr RawPointer) InnerType() Type         { return ptr.Inner }
func (ref MutableReference) InnerType() Type   { return ref.Inner }
func (slc SliceTypeASTNode) InnerType() Type   { return slc.ValueType }
func (ref ImmutableReference) InnerType() Type { return ref.Inner }

func (ptr RawPointer) IsDyn() bool         { return ptr.IsDynamic }
func (ref MutableReference) IsDyn() bool   { return ref.IsDynamic }
func (slc SliceTypeASTNode) IsDyn() bool   { return false }
func (ref ImmutableReference) IsDyn() bool { return ref.IsDynamic }

func (ptr RawPointer) Location() lexer.Location               { return ptr.Loc }
func (ref MutableReference) Location() lexer.Location         { return ref.Loc }
func (ref ImmutableReference) Location() lexer.Location       { return ref.Loc }
func (typ NamedTypeASTNode) Location() lexer.Location         { return typ.Loc }
func (typ UntaggedUnionTypeASTNode) Location() lexer.Location { return typ.Loc }
func (typ NeverTypeASTNode) Location() lexer.Location         { return typ.Loc }
func (typ TableTypeASTNode) Location() lexer.Location         { return typ.Loc }
func (typ ArrayTypeASTNode) Location() lexer.Location         { return typ.Loc }
func (typ SliceTypeASTNode) Location() lexer.Location         { return typ.Loc }
func (typ TupleTypeASTNode) Location() lexer.Location         { return typ.Loc }

func (ptr RawPointer) Kind() ASTNodeKind               { return RawPointerTypeASTNodeKind }
func (ref MutableReference) Kind() ASTNodeKind         { return MutableReferenceTypeASTNodeKind }
func (ref ImmutableReference) Kind() ASTNodeKind       { return ImmutableReferenceTypeASTNodeKind }
func (typ NamedTypeASTNode) Kind() ASTNodeKind         { return NamedTypeASTNodeKind }
func (typ UntaggedUnionTypeASTNode) Kind() ASTNodeKind { return UntaggedUnionTypeASTNodeKind }
func (typ NeverTypeASTNode) Kind() ASTNodeKind         { return NeverTypeASTNodeKind }
func (typ TableTypeASTNode) Kind() ASTNodeKind         { return TableTypeASTNodeKind }
func (typ ArrayTypeASTNode) Kind() ASTNodeKind         { return ArrayTypeASTNodeKind }
func (typ SliceTypeASTNode) Kind() ASTNodeKind         { return SliceTypeASTNodeKind }
func (typ TupleTypeASTNode) Kind() ASTNodeKind         { return TupleTypeASTNodeKind }

func (ptr RawPointer) Group() ASTNodeGroup               { return TypeASTNodeGroup }
func (ref MutableReference) Group() ASTNodeGroup         { return TypeASTNodeGroup }
func (ref ImmutableReference) Group() ASTNodeGroup       { return TypeASTNodeGroup }
func (typ NamedTypeASTNode) Group() ASTNodeGroup         { return TypeASTNodeGroup }
func (typ UntaggedUnionTypeASTNode) Group() ASTNodeGroup { return TypeASTNodeGroup }
func (typ NeverTypeASTNode) Group() ASTNodeGroup         { return TypeASTNodeGroup }
func (typ TableTypeASTNode) Group() ASTNodeGroup         { return TypeASTNodeGroup }
func (typ ArrayTypeASTNode) Group() ASTNodeGroup         { return TypeASTNodeGroup }
func (typ SliceTypeASTNode) Group() ASTNodeGroup         { return TypeASTNodeGroup }
func (typ TupleTypeASTNode) Group() ASTNodeGroup         { return TypeASTNodeGroup }

type (
	ImportStatementASTNode struct {
		Loc  lexer.Location
		Path []string
	}

	ConstDefinitionASTNode struct {
		Loc   lexer.Location
		Name  string
		Value Expression
		Type  utils.Optional[Type]
	}

	VarDefinitionASTNode struct {
		Loc   lexer.Location
		Name  string
		Value Expression
		Type  utils.Optional[Type]
	}

	LetDefinitionASTNode struct {
		Loc   lexer.Location
		Name  string
		Value Expression
		Type  utils.Optional[Type]
	}

	IdentifierLiteralASTNode struct {
		Loc  lexer.Location
		Name string
	}

	BinaryExpressionASTNode struct {
		Loc      lexer.Location
		Operator string
		Left     Expression
		Right    Expression
	}

	PostfixUnaryExpressionASTNode struct {
		Loc      lexer.Location
		Operator string
		Left     Expression
	}

	PrefixUnaryExpressionASTNode struct {
		Loc      lexer.Location
		Operator string
		Right    Expression
	}

	StructureDefFields []struct {
		Loc       lexer.Location
		Name      string
		IsMutable bool
		Type      Type
	}

	StructureDefinitionASTNode struct {
		Loc    lexer.Location
		Name   string
		Fields StructureDefFields
	}

	ClassDefFields []struct {
		Loc       lexer.Location
		Name      string
		IsMutable bool
		Type      Type
	}

	ClassDefMethods []struct {
		Loc           lexer.Location
		Name          string
		ReturnType    Type
		ArgumentTypes []Type
		SelfType      utils.Optional[Type]
	}

	ClassDefConstructors []struct {
		Loc           lexer.Location
		Name          string
		MayReturnNull bool
		ArgumentTypes []Type
	}

	ClassDefinitionASTNode struct {
		Loc          lexer.Location
		Name         string
		Fields       ClassDefFields
		Methods      ClassDefMethods
		Constructors ClassDefConstructors
	}

	FunctionDefinitionASTNode struct {
		Loc        lexer.Location
		Name       string
		ReturnType Type
		Parameters map[string]Type
		Generics   map[string]TypeGenericASTNode
		Body       BlockASTNode
	}

	MethodDefinitionASTNode struct {
		Loc         lexer.Location
		Name        string
		ReturnType  Type
		Parameters  map[string]Type
		ContextType Type
		Generics    map[string]TypeGenericASTNode
		Body        BlockASTNode
	}

	OperatorOverloadASTNode struct {
		Loc           lexer.Location
		Operator      string
		ReturnType    Type
		ContextType   Type
		RightHandType Type
		Body          BlockASTNode
	}

	AssignmentStatementASTNode struct {
		Loc   lexer.Location
		Left  Expression
		Right Expression
	}

	StructureInitilisationExpressionASTNode struct {
		Loc         lexer.Location
		StructName  string
		FieldValues map[string]Expression
	}

	StructureRefInitilisationExpressionASTNode struct {
		Loc         lexer.Location
		StructName  string
		FieldValues map[string]Expression
		RefType     RefType
	}

	ImplicitReturnASTNode struct {
		Loc   lexer.Location
		Value Expression
	}

	ExplicitReturnASTNode struct {
		Loc   lexer.Location
		Value Expression
	}

	FunctionCallExpressionASTNode struct {
		Loc        lexer.Location
		Name       string
		Parameters map[string]Expression
		Generics   map[string]TypeGenericASTNode
	}

	MethodCallExpressionASTNode struct {
		Loc        lexer.Location
		Context    Type
		Name       string
		Parameters map[string]Expression
		Generics   map[string]TypeGenericASTNode
	}

	MemberExpressionASTNode struct {
		Loc      lexer.Location
		Segments []Expression
	}

	ModulePathASTNode struct {
		Loc      lexer.Location
		Segments []Expression
	}

	LambdaExpressionASTNode struct {
		Loc        lexer.Location
		Parameters map[string]Type
		ReturnType Type
		Body       BlockASTNode
	}

	IfExpressionASTNode struct {
		Loc          lexer.Location
		Condition    Expression
		Body         BlockASTNode
		FallbackBody BlockASTNode
	}

	IfStatementASTNode struct {
		Loc       lexer.Location
		Condition Expression
		Body      BlockASTNode
	}

	SwitchStatementCase struct {
		Loc  lexer.Location
		Body BlockASTNode
	}

	SwitchStatementASTNode struct {
		Loc   lexer.Location
		Value Expression
		Cases map[Expression]SwitchStatementCase
	}

	MatchExpressionCase struct {
		Loc  lexer.Location
		Body BlockASTNode
	}

	MatchOrWhenExpressionCase struct {
		Loc  lexer.Location
		Body BlockASTNode
	}

	MatchExpressionASTNode struct {
		Loc          lexer.Location
		Value        Expression
		Cases        map[Constraint]MatchOrWhenExpressionCase
		FallbackCase MatchOrWhenExpressionCase
	}

	WhenExpressionASTNode struct {
		Loc          lexer.Location
		Cases        map[Expression]MatchOrWhenExpressionCase
		FallbackCase MatchOrWhenExpressionCase
	}

	GetterMethodDecl struct {
		Loc      lexer.Location
		SelfType RefType
	}

	SetterMethodDecl struct {
		Loc      lexer.Location
		SelfType RefType
	}

	InterfaceDefField struct {
		Loc            lexer.Location
		Type           Type
		ComputedGetter utils.Optional[GetterMethodDecl]
		ComputedSetter utils.Optional[SetterMethodDecl]
	}

	InterfaceDefMethod struct {
		Loc         lexer.Location
		ReturnType  Type
		Parameters  map[string]Type
		ContextType Type
		Generics    map[string]TypeGenericASTNode
	}

	InterfaceDefinitionASTNode struct {
		Loc      lexer.Location
		Extends  []NamedTypeASTNode
		Fields   map[string]InterfaceDefField
		Methods  map[string]InterfaceDefMethod
		Generics map[string]TypeGenericASTNode
	}

	StringLiteralASTNode struct {
		Loc    lexer.Location
		String string
	}

	ArrayLiteralASTNode struct {
		Loc    lexer.Location
		Values []Expression
	}

	IntegerLiteralASTNode struct {
		Loc   lexer.Location
		Value int64
	}

	DecimalLiteralASTNode struct {
		Loc   lexer.Location
		Value float64
	}

	CStyleEnumDefinitionASTNode struct {
		Loc         lexer.Location
		Name        string
		Variants    []string
		OrdinalType Type
	}

	SumTypeEnumVariant struct {
		Loc        lexer.Location
		Parameters map[string]Type
	}

	SumTypeEnumDefinitionASTNode struct {
		Loc      lexer.Location
		Name     string
		Variants map[string]SumTypeEnumVariant
	}

	NamespaceDefinitionASTNode struct {
		Loc          lexer.Location
		Name         string
		Declarations []Declaration
	}

	ExternalFnDeclarationASTNode struct {
		Loc        lexer.Location
		Name       string
		ReturnType Type
		Parameters map[string]Type
	}

	CStyleForLoopStatementASTNode struct {
		Loc            lexer.Location
		Initialisation utils.Optional[VarDefinitionASTNode]
		Check          utils.Optional[Expression]
		Increment      utils.Optional[Statement]
		Body           BlockASTNode
	}

	ForInLoopStatementASTNode struct {
		Loc      lexer.Location
		Variable string
		Iterator Expression
		Body     BlockASTNode
	}

	WhileLoopStatementASTNode struct {
		Loc       lexer.Location
		Condition Expression
		Body      BlockASTNode
	}

	ForeverLoopStatementASTNode struct {
		Loc  lexer.Location
		Body BlockASTNode
	}

	TernaryExpressionASTNode struct {
		Loc           lexer.Location
		Condition     Expression
		SuccessValue  Expression
		FallbackValue Expression
	}

	OptionalChainingASTNode struct {
		Loc   lexer.Location
		Chain []Expression
	}

	TypeCastableQueryExpressionASTNode struct {
		Loc   lexer.Location
		Value Expression
		Type  Type
	}

	TypeCastExpressionASTNode struct {
		Loc   lexer.Location
		Value Expression
		Type  Type
	}

	RuntimeTypeCastExpressionASTNode struct {
		Loc   lexer.Location
		Value Expression
		Type  Type
	}

	BlockASTNode struct {
		Loc  lexer.Location
		Code []ASTNode
	}

	DestructedElement struct {
		Loc       lexer.Location
		Mutable   bool
		ValueType Type
	}

	TupleDestructuringASTNode struct {
		Loc      lexer.Location
		Elements map[string]DestructedElement
	}

	ArrayCompTimeDestructuringASTNode struct {
		Loc      lexer.Location
		Elements map[string]DestructedElement
	}

	ArrayRuntimeDestructuringASTNode struct {
		Loc      lexer.Location
		Elements map[string]DestructedElement
	}

	StructOrClassDestructuringASTNode struct {
		Loc      lexer.Location
		Elements map[string]DestructedElement
	}

	ReferenceDestructuringASTNode struct {
		Loc           lexer.Location
		ReferenceType RefType
		Destructuring Destructure
	}

	ConstraintDestructedElement struct {
		Loc       lexer.Location
		Mutable   bool
		ValueType Type
		AliasName utils.Optional[string]
	}

	ConstraintASTNode struct {
		Loc          lexer.Location
		Variant      Identifier
		WhereClauses []Expression
	}

	IfLetStatementASTNode struct {
		Loc   lexer.Location
		Name  string
		Value Expression
		Body  BlockASTNode
	}

	IfVarStatementASTNode struct {
		Loc   lexer.Location
		Name  string
		Value Expression
		Body  BlockASTNode
	}

	IfLetExpressionASTNode struct {
		Loc          lexer.Location
		Name         string
		Value        Expression
		Body         BlockASTNode
		FallbackBody BlockASTNode
	}

	IfVarExpressionASTNode struct {
		Loc          lexer.Location
		Name         string
		Value        Expression
		Body         BlockASTNode
		FallbackBody BlockASTNode
	}

	NullCoalesceExpressionASTNode struct {
		Loc           lexer.Location
		Value         Expression
		FallbackValue Expression
	}

	BubbleValueToReturnASTNode struct {
		Loc   lexer.Location
		value Expression
	}

	GetterMethodDef struct {
		Loc      lexer.Location
		Body     BlockASTNode
		SelfType RefType
	}

	SetterMethodDef struct {
		Loc             lexer.Location
		Body            BlockASTNode
		SelfType        RefType
		NewValueArgName string
	}

	ComputedVarDefinitionASTNode struct {
		Loc       lexer.Location
		Name      string
		Getter    GetterMethodDef
		Setter    utils.Optional[SetterMethodDef]
		ValueType Type
	}
)

/*
ImportStatementASTNode
ConstDefinitionASTNode
VarDefinitionASTNode
LetDefinitionASTNode
VariableUsageExpresionASTNode
BinaryExpressionASTNode
PostfixUnaryExpressionASTNode
PrefixUnaryExpressionASTNode
StructureDefinitionASTNode
ClassDefinitionASTNode
FunctionDefinitionASTNode
MethodDefinitionASTNode
OperatorOverloadASTNode
AssignmentStatementASTNode
StructureInitilisationExpressionASTNode
StructureRefInitilisationExpressionASTNode
ImplicitReturnASTNode
ExplicitReturnASTNode
FunctionCallExpressionASTNode
MethodCallExpressionASTNode
MemberExpressionASTNode
ModulePathASTNode
LambdaExpressionASTNode
IfExpressionASTNode
IfStatementASTNode
SwitchStatementASTNode
MatchExpressionASTNode
WhenExpressionASTNode
InterfaceDefinitionASTNode
StringLiteralASTNode
ArrayLiteralASTNode
IntegerLiteralASTNode
DecimalLiteralASTNode
CStyleEnumDefinitionASTNode
SumTypeEnumDefinitionASTNode
NamespaceDefinitionASTNode
ExternalFnDeclarationASTNode
CStyleForLoopStatementASTNode
ForInLoopStatementASTNode
WhileLoopStatementASTNode
ForeverLoopStatementASTNode
TernaryExpressionASTNode
OptionalChainingASTNode
TypeCastableQueryExpressionASTNode
TypeCastExpressionASTNode
RuntimeTypeCastExpressionASTNode
InternalMacroDeclarationASTNode
CustomMacroDefinitionASTNode
MacroUsageASTNode
MacroVariableUsageASTNode
BlockASTNode
TupleDestructuringASTNode
ArrayCompTimeDestructuringASTNode
ArrayRuntimeDestructuringASTNode
StructOrClassDestructuringASTNode
ReferenceDestructuringASTNode
ConstraintASTNode
IfLetStatementASTNode
IfVarStatementASTNode
IfLetExpressionASTNode
IfVarExpressionASTNode
NullCoalesceExpressionASTNode
BubbleValueToReturnASTNode

func (node) Location() lexer.Location { return node.Loc }
*/

func (node ImportStatementASTNode) Location() lexer.Location                     { return node.Loc }
func (node ConstDefinitionASTNode) Location() lexer.Location                     { return node.Loc }
func (node VarDefinitionASTNode) Location() lexer.Location                       { return node.Loc }
func (node LetDefinitionASTNode) Location() lexer.Location                       { return node.Loc }
func (node IdentifierLiteralASTNode) Location() lexer.Location                   { return node.Loc }
func (node BinaryExpressionASTNode) Location() lexer.Location                    { return node.Loc }
func (node PostfixUnaryExpressionASTNode) Location() lexer.Location              { return node.Loc }
func (node PrefixUnaryExpressionASTNode) Location() lexer.Location               { return node.Loc }
func (node StructureDefinitionASTNode) Location() lexer.Location                 { return node.Loc }
func (node ClassDefinitionASTNode) Location() lexer.Location                     { return node.Loc }
func (node FunctionDefinitionASTNode) Location() lexer.Location                  { return node.Loc }
func (node MethodDefinitionASTNode) Location() lexer.Location                    { return node.Loc }
func (node OperatorOverloadASTNode) Location() lexer.Location                    { return node.Loc }
func (node AssignmentStatementASTNode) Location() lexer.Location                 { return node.Loc }
func (node StructureInitilisationExpressionASTNode) Location() lexer.Location    { return node.Loc }
func (node StructureRefInitilisationExpressionASTNode) Location() lexer.Location { return node.Loc }
func (node ImplicitReturnASTNode) Location() lexer.Location                      { return node.Loc }
func (node ExplicitReturnASTNode) Location() lexer.Location                      { return node.Loc }
func (node FunctionCallExpressionASTNode) Location() lexer.Location              { return node.Loc }
func (node MethodCallExpressionASTNode) Location() lexer.Location                { return node.Loc }
func (node MemberExpressionASTNode) Location() lexer.Location                    { return node.Loc }
func (node ModulePathASTNode) Location() lexer.Location                          { return node.Loc }
func (node LambdaExpressionASTNode) Location() lexer.Location                    { return node.Loc }
func (node IfExpressionASTNode) Location() lexer.Location                        { return node.Loc }
func (node IfStatementASTNode) Location() lexer.Location                         { return node.Loc }
func (node SwitchStatementASTNode) Location() lexer.Location                     { return node.Loc }
func (node MatchExpressionASTNode) Location() lexer.Location                     { return node.Loc }
func (node WhenExpressionASTNode) Location() lexer.Location                      { return node.Loc }
func (node InterfaceDefinitionASTNode) Location() lexer.Location                 { return node.Loc }
func (node StringLiteralASTNode) Location() lexer.Location                       { return node.Loc }
func (node ArrayLiteralASTNode) Location() lexer.Location                        { return node.Loc }
func (node IntegerLiteralASTNode) Location() lexer.Location                      { return node.Loc }
func (node DecimalLiteralASTNode) Location() lexer.Location                      { return node.Loc }
func (node CStyleEnumDefinitionASTNode) Location() lexer.Location                { return node.Loc }
func (node SumTypeEnumDefinitionASTNode) Location() lexer.Location               { return node.Loc }
func (node NamespaceDefinitionASTNode) Location() lexer.Location                 { return node.Loc }
func (node ExternalFnDeclarationASTNode) Location() lexer.Location               { return node.Loc }
func (node CStyleForLoopStatementASTNode) Location() lexer.Location              { return node.Loc }
func (node ForInLoopStatementASTNode) Location() lexer.Location                  { return node.Loc }
func (node WhileLoopStatementASTNode) Location() lexer.Location                  { return node.Loc }
func (node ForeverLoopStatementASTNode) Location() lexer.Location                { return node.Loc }
func (node TernaryExpressionASTNode) Location() lexer.Location                   { return node.Loc }
func (node OptionalChainingASTNode) Location() lexer.Location                    { return node.Loc }
func (node TypeCastableQueryExpressionASTNode) Location() lexer.Location         { return node.Loc }
func (node TypeCastExpressionASTNode) Location() lexer.Location                  { return node.Loc }
func (node RuntimeTypeCastExpressionASTNode) Location() lexer.Location           { return node.Loc }
func (node BlockASTNode) Location() lexer.Location                               { return node.Loc }
func (node TupleDestructuringASTNode) Location() lexer.Location                  { return node.Loc }
func (node ArrayCompTimeDestructuringASTNode) Location() lexer.Location          { return node.Loc }
func (node ArrayRuntimeDestructuringASTNode) Location() lexer.Location           { return node.Loc }
func (node StructOrClassDestructuringASTNode) Location() lexer.Location          { return node.Loc }
func (node ReferenceDestructuringASTNode) Location() lexer.Location              { return node.Loc }
func (node ConstraintASTNode) Location() lexer.Location                          { return node.Loc }
func (node IfLetStatementASTNode) Location() lexer.Location                      { return node.Loc }
func (node IfVarStatementASTNode) Location() lexer.Location                      { return node.Loc }
func (node IfLetExpressionASTNode) Location() lexer.Location                     { return node.Loc }
func (node IfVarExpressionASTNode) Location() lexer.Location                     { return node.Loc }
func (node NullCoalesceExpressionASTNode) Location() lexer.Location              { return node.Loc }
func (node BubbleValueToReturnASTNode) Location() lexer.Location                 { return node.Loc }

func (node ImportStatementASTNode) Kind() ASTNodeKind   { return ImportStatementASTNodeKind }
func (node ConstDefinitionASTNode) Kind() ASTNodeKind   { return ConstDefinitionASTNodeKind }
func (node VarDefinitionASTNode) Kind() ASTNodeKind     { return VarDefinitionASTNodeKind }
func (node LetDefinitionASTNode) Kind() ASTNodeKind     { return LetDefinitionASTNodeKind }
func (node IdentifierLiteralASTNode) Kind() ASTNodeKind { return IdentifierLiteralASTNodeKind }
func (node BinaryExpressionASTNode) Kind() ASTNodeKind  { return BinaryExpressionASTNodeKind }
func (node PostfixUnaryExpressionASTNode) Kind() ASTNodeKind {
	return PostfixUnaryExpressionASTNodeKind
}
func (node PrefixUnaryExpressionASTNode) Kind() ASTNodeKind { return PrefixUnaryExpressionASTNodeKind }
func (node StructureDefinitionASTNode) Kind() ASTNodeKind   { return StructureDefinitionASTNodeKind }
func (node ClassDefinitionASTNode) Kind() ASTNodeKind       { return ClassDefinitionASTNodeKind }
func (node FunctionDefinitionASTNode) Kind() ASTNodeKind    { return FunctionDefinitionASTNodeKind }
func (node MethodDefinitionASTNode) Kind() ASTNodeKind      { return MethodDefinitionASTNodeKind }
func (node OperatorOverloadASTNode) Kind() ASTNodeKind      { return OperatorOverloadASTNodeKind }
func (node AssignmentStatementASTNode) Kind() ASTNodeKind   { return AssignmentStatementASTNodeKind }
func (node StructureInitilisationExpressionASTNode) Kind() ASTNodeKind {
	return StructureInitilisationExpressionASTNodeKind
}
func (node StructureRefInitilisationExpressionASTNode) Kind() ASTNodeKind {
	return StructureRefInitilisationExpressionASTNodeKind
}
func (node ImplicitReturnASTNode) Kind() ASTNodeKind { return ImplicitReturnASTNodeKind }
func (node ExplicitReturnASTNode) Kind() ASTNodeKind { return ExplicitReturnASTNodeKind }
func (node FunctionCallExpressionASTNode) Kind() ASTNodeKind {
	return FunctionCallExpressionASTNodeKind
}
func (node MethodCallExpressionASTNode) Kind() ASTNodeKind  { return MethodCallExpressionASTNodeKind }
func (node MemberExpressionASTNode) Kind() ASTNodeKind      { return MemberExpressionASTNodeKind }
func (node ModulePathASTNode) Kind() ASTNodeKind            { return ModulePathASTNodeKind }
func (node LambdaExpressionASTNode) Kind() ASTNodeKind      { return LambdaExpressionASTNodeKind }
func (node IfExpressionASTNode) Kind() ASTNodeKind          { return IfExpressionASTNodeKind }
func (node IfStatementASTNode) Kind() ASTNodeKind           { return IfStatementASTNodeKind }
func (node SwitchStatementASTNode) Kind() ASTNodeKind       { return SwitchStatementASTNodeKind }
func (node MatchExpressionASTNode) Kind() ASTNodeKind       { return MatchExpressionASTNodeKind }
func (node WhenExpressionASTNode) Kind() ASTNodeKind        { return WhenExpressionASTNodeKind }
func (node InterfaceDefinitionASTNode) Kind() ASTNodeKind   { return InterfaceDefinitionASTNodeKind }
func (node StringLiteralASTNode) Kind() ASTNodeKind         { return StringLiteralASTNodeKind }
func (node ArrayLiteralASTNode) Kind() ASTNodeKind          { return ArrayLiteralASTNodeKind }
func (node IntegerLiteralASTNode) Kind() ASTNodeKind        { return IntegerLiteralASTNodeKind }
func (node DecimalLiteralASTNode) Kind() ASTNodeKind        { return DecimalLiteralASTNodeKind }
func (node CStyleEnumDefinitionASTNode) Kind() ASTNodeKind  { return CStyleEnumDefinitionASTNodeKind }
func (node SumTypeEnumDefinitionASTNode) Kind() ASTNodeKind { return SumTypeEnumDefinitionASTNodeKind }
func (node NamespaceDefinitionASTNode) Kind() ASTNodeKind   { return NamespaceDefinitionASTNodeKind }
func (node ExternalFnDeclarationASTNode) Kind() ASTNodeKind { return ExternalFnDeclarationASTNodeKind }
func (node CStyleForLoopStatementASTNode) Kind() ASTNodeKind {
	return CStyleForLoopStatementASTNodeKind
}
func (node ForInLoopStatementASTNode) Kind() ASTNodeKind   { return ForInLoopStatementASTNodeKind }
func (node WhileLoopStatementASTNode) Kind() ASTNodeKind   { return WhileLoopStatementASTNodeKind }
func (node ForeverLoopStatementASTNode) Kind() ASTNodeKind { return ForeverLoopStatementASTNodeKind }
func (node TernaryExpressionASTNode) Kind() ASTNodeKind    { return TernaryExpressionASTNodeKind }
func (node OptionalChainingASTNode) Kind() ASTNodeKind     { return OptionalChainingASTNodeKind }
func (node TypeCastableQueryExpressionASTNode) Kind() ASTNodeKind {
	return TypeCastableQueryExpressionASTNodeKind
}
func (node TypeCastExpressionASTNode) Kind() ASTNodeKind { return TypeCastExpressionASTNodeKind }
func (node RuntimeTypeCastExpressionASTNode) Kind() ASTNodeKind {
	return RuntimeTypeCastExpressionASTNodeKind
}
func (node BlockASTNode) Kind() ASTNodeKind              { return BlockASTNodeKind }
func (node TupleDestructuringASTNode) Kind() ASTNodeKind { return TupleDestructuringASTNodeKind }
func (node ArrayCompTimeDestructuringASTNode) Kind() ASTNodeKind {
	return ArrayCompTimeDestructuringASTNodeKind
}
func (node ArrayRuntimeDestructuringASTNode) Kind() ASTNodeKind {
	return ArrayRuntimeDestructuringASTNodeKind
}
func (node StructOrClassDestructuringASTNode) Kind() ASTNodeKind {
	return StructOrClassDestructuringASTNodeKind
}
func (node ReferenceDestructuringASTNode) Kind() ASTNodeKind {
	return ReferenceDestructuringASTNodeKind
}
func (node ConstraintASTNode) Kind() ASTNodeKind      { return ConstraintASTNodeKind }
func (node IfLetStatementASTNode) Kind() ASTNodeKind  { return IfLetStatementASTNodeKind }
func (node IfVarStatementASTNode) Kind() ASTNodeKind  { return IfVarStatementASTNodeKind }
func (node IfLetExpressionASTNode) Kind() ASTNodeKind { return IfLetExpressionASTNodeKind }
func (node IfVarExpressionASTNode) Kind() ASTNodeKind { return IfVarExpressionASTNodeKind }
func (node NullCoalesceExpressionASTNode) Kind() ASTNodeKind {
	return NullCoalesceExpressionASTNodeKind
}
func (node BubbleValueToReturnASTNode) Kind() ASTNodeKind { return BubbleValueToReturnASTNodeKind }

func (node ImportStatementASTNode) Group() ASTNodeGroup        { return StatementASTNodeGroup }
func (node IfLetStatementASTNode) Group() ASTNodeGroup         { return StatementASTNodeGroup }
func (node IfVarStatementASTNode) Group() ASTNodeGroup         { return StatementASTNodeGroup }
func (node AssignmentStatementASTNode) Group() ASTNodeGroup    { return StatementASTNodeGroup }
func (node ImplicitReturnASTNode) Group() ASTNodeGroup         { return StatementASTNodeGroup }
func (node ExplicitReturnASTNode) Group() ASTNodeGroup         { return StatementASTNodeGroup }
func (node IfStatementASTNode) Group() ASTNodeGroup            { return StatementASTNodeGroup }
func (node SwitchStatementASTNode) Group() ASTNodeGroup        { return StatementASTNodeGroup }
func (node CStyleForLoopStatementASTNode) Group() ASTNodeGroup { return StatementASTNodeGroup }
func (node ForInLoopStatementASTNode) Group() ASTNodeGroup     { return StatementASTNodeGroup }
func (node WhileLoopStatementASTNode) Group() ASTNodeGroup     { return StatementASTNodeGroup }
func (node ForeverLoopStatementASTNode) Group() ASTNodeGroup   { return StatementASTNodeGroup }

func (node ImportStatementASTNode) statementNode()        {}
func (node IfLetStatementASTNode) statementNode()         {}
func (node IfVarStatementASTNode) statementNode()         {}
func (node AssignmentStatementASTNode) statementNode()    {}
func (node ImplicitReturnASTNode) statementNode()         {}
func (node ExplicitReturnASTNode) statementNode()         {}
func (node IfStatementASTNode) statementNode()            {}
func (node SwitchStatementASTNode) statementNode()        {}
func (node CStyleForLoopStatementASTNode) statementNode() {}
func (node ForInLoopStatementASTNode) statementNode()     {}
func (node WhileLoopStatementASTNode) statementNode()     {}
func (node ForeverLoopStatementASTNode) statementNode()   {}

func (node ConstDefinitionASTNode) Group() ASTNodeGroup       { return DefinitionASTNodeGroup }
func (node VarDefinitionASTNode) Group() ASTNodeGroup         { return DefinitionASTNodeGroup }
func (node LetDefinitionASTNode) Group() ASTNodeGroup         { return DefinitionASTNodeGroup }
func (node StructureDefinitionASTNode) Group() ASTNodeGroup   { return DefinitionASTNodeGroup }
func (node ClassDefinitionASTNode) Group() ASTNodeGroup       { return DefinitionASTNodeGroup }
func (node FunctionDefinitionASTNode) Group() ASTNodeGroup    { return DefinitionASTNodeGroup }
func (node MethodDefinitionASTNode) Group() ASTNodeGroup      { return DefinitionASTNodeGroup }
func (node OperatorOverloadASTNode) Group() ASTNodeGroup      { return DefinitionASTNodeGroup }
func (node InterfaceDefinitionASTNode) Group() ASTNodeGroup   { return DefinitionASTNodeGroup }
func (node CStyleEnumDefinitionASTNode) Group() ASTNodeGroup  { return DefinitionASTNodeGroup }
func (node SumTypeEnumDefinitionASTNode) Group() ASTNodeGroup { return DefinitionASTNodeGroup }
func (node NamespaceDefinitionASTNode) Group() ASTNodeGroup   { return DefinitionASTNodeGroup }

func (node ConstDefinitionASTNode) statementNode()         {}
func (node VarDefinitionASTNode) statementNode()           {}
func (node LetDefinitionASTNode) statementNode()           {}
func (node StructureDefinitionASTNode) statementNode()     {}
func (node ClassDefinitionASTNode) statementNode()         {}
func (node FunctionDefinitionASTNode) statementNode()      {}
func (node MethodDefinitionASTNode) statementNode()        {}
func (node OperatorOverloadASTNode) statementNode()        {}
func (node InterfaceDefinitionASTNode) statementNode()     {}
func (node CStyleEnumDefinitionASTNode) statementNode()    {}
func (node SumTypeEnumDefinitionASTNode) statementNode()   {}
func (node NamespaceDefinitionASTNode) statementNode()     {}
func (node ConstDefinitionASTNode) declarationNode()       {}
func (node VarDefinitionASTNode) declarationNode()         {}
func (node LetDefinitionASTNode) declarationNode()         {}
func (node StructureDefinitionASTNode) declarationNode()   {}
func (node ClassDefinitionASTNode) declarationNode()       {}
func (node FunctionDefinitionASTNode) declarationNode()    {}
func (node MethodDefinitionASTNode) declarationNode()      {}
func (node OperatorOverloadASTNode) declarationNode()      {}
func (node InterfaceDefinitionASTNode) declarationNode()   {}
func (node CStyleEnumDefinitionASTNode) declarationNode()  {}
func (node SumTypeEnumDefinitionASTNode) declarationNode() {}
func (node NamespaceDefinitionASTNode) declarationNode()   {}
func (node ConstDefinitionASTNode) definitionNode()        {}
func (node VarDefinitionASTNode) definitionNode()          {}
func (node LetDefinitionASTNode) definitionNode()          {}
func (node StructureDefinitionASTNode) definitionNode()    {}
func (node ClassDefinitionASTNode) definitionNode()        {}
func (node FunctionDefinitionASTNode) definitionNode()     {}
func (node MethodDefinitionASTNode) definitionNode()       {}
func (node OperatorOverloadASTNode) definitionNode()       {}
func (node InterfaceDefinitionASTNode) definitionNode()    {}
func (node CStyleEnumDefinitionASTNode) definitionNode()   {}
func (node SumTypeEnumDefinitionASTNode) definitionNode()  {}
func (node NamespaceDefinitionASTNode) definitionNode()    {}

func (node StringLiteralASTNode) Group() ASTNodeGroup     { return LiteralASTNodeGroup }
func (node ArrayLiteralASTNode) Group() ASTNodeGroup      { return LiteralASTNodeGroup }
func (node IntegerLiteralASTNode) Group() ASTNodeGroup    { return LiteralASTNodeGroup }
func (node DecimalLiteralASTNode) Group() ASTNodeGroup    { return LiteralASTNodeGroup }
func (node IdentifierLiteralASTNode) Group() ASTNodeGroup { return LiteralASTNodeGroup }

func (node StringLiteralASTNode) statementNode()      {}
func (node ArrayLiteralASTNode) statementNode()       {}
func (node IntegerLiteralASTNode) statementNode()     {}
func (node DecimalLiteralASTNode) statementNode()     {}
func (node IdentifierLiteralASTNode) statementNode()  {}
func (node StringLiteralASTNode) expressionNode()     {}
func (node ArrayLiteralASTNode) expressionNode()      {}
func (node IntegerLiteralASTNode) expressionNode()    {}
func (node DecimalLiteralASTNode) expressionNode()    {}
func (node IdentifierLiteralASTNode) expressionNode() {}
func (node StringLiteralASTNode) literalNode()        {}
func (node ArrayLiteralASTNode) literalNode()         {}
func (node IntegerLiteralASTNode) literalNode()       {}
func (node DecimalLiteralASTNode) literalNode()       {}
func (node IdentifierLiteralASTNode) literalNode()    {}

func (node ExternalFnDeclarationASTNode) Group() ASTNodeGroup { return DeclarationASTNodeGroup }
func (node ExternalFnDeclarationASTNode) statementNode()      {}
func (node ExternalFnDeclarationASTNode) declarationNode()    {}

func (node BlockASTNode) Group() ASTNodeGroup                      { return SegmentASTNodeGroup }
func (node TupleDestructuringASTNode) Group() ASTNodeGroup         { return SegmentASTNodeGroup }
func (node ArrayCompTimeDestructuringASTNode) Group() ASTNodeGroup { return SegmentASTNodeGroup }
func (node ArrayRuntimeDestructuringASTNode) Group() ASTNodeGroup  { return SegmentASTNodeGroup }
func (node StructOrClassDestructuringASTNode) Group() ASTNodeGroup { return SegmentASTNodeGroup }
func (node ReferenceDestructuringASTNode) Group() ASTNodeGroup     { return SegmentASTNodeGroup }
func (node ConstraintASTNode) Group() ASTNodeGroup                 { return SegmentASTNodeGroup }

func (node BlockASTNode) componentNode()                      {}
func (node TupleDestructuringASTNode) componentNode()         {}
func (node ArrayCompTimeDestructuringASTNode) componentNode() {}
func (node ArrayRuntimeDestructuringASTNode) componentNode()  {}
func (node StructOrClassDestructuringASTNode) componentNode() {}
func (node ReferenceDestructuringASTNode) componentNode()     {}
func (node ConstraintASTNode) componentNode()                 {}

func (node IfLetExpressionASTNode) Group() ASTNodeGroup        { return ExpressionASTNodeGroup }
func (node IfVarExpressionASTNode) Group() ASTNodeGroup        { return ExpressionASTNodeGroup }
func (node NullCoalesceExpressionASTNode) Group() ASTNodeGroup { return ExpressionASTNodeGroup }
func (node BubbleValueToReturnASTNode) Group() ASTNodeGroup    { return ExpressionASTNodeGroup }
func (node BinaryExpressionASTNode) Group() ASTNodeGroup       { return ExpressionASTNodeGroup }
func (node PostfixUnaryExpressionASTNode) Group() ASTNodeGroup { return ExpressionASTNodeGroup }
func (node PrefixUnaryExpressionASTNode) Group() ASTNodeGroup  { return ExpressionASTNodeGroup }
func (node StructureInitilisationExpressionASTNode) Group() ASTNodeGroup {
	return ExpressionASTNodeGroup
}
func (node StructureRefInitilisationExpressionASTNode) Group() ASTNodeGroup {
	return ExpressionASTNodeGroup
}
func (node FunctionCallExpressionASTNode) Group() ASTNodeGroup      { return ExpressionASTNodeGroup }
func (node MethodCallExpressionASTNode) Group() ASTNodeGroup        { return ExpressionASTNodeGroup }
func (node MemberExpressionASTNode) Group() ASTNodeGroup            { return ExpressionASTNodeGroup }
func (node ModulePathASTNode) Group() ASTNodeGroup                  { return ExpressionASTNodeGroup }
func (node LambdaExpressionASTNode) Group() ASTNodeGroup            { return ExpressionASTNodeGroup }
func (node IfExpressionASTNode) Group() ASTNodeGroup                { return ExpressionASTNodeGroup }
func (node MatchExpressionASTNode) Group() ASTNodeGroup             { return ExpressionASTNodeGroup }
func (node WhenExpressionASTNode) Group() ASTNodeGroup              { return ExpressionASTNodeGroup }
func (node TernaryExpressionASTNode) Group() ASTNodeGroup           { return ExpressionASTNodeGroup }
func (node OptionalChainingASTNode) Group() ASTNodeGroup            { return ExpressionASTNodeGroup }
func (node TypeCastableQueryExpressionASTNode) Group() ASTNodeGroup { return ExpressionASTNodeGroup }
func (node TypeCastExpressionASTNode) Group() ASTNodeGroup          { return ExpressionASTNodeGroup }

func (node IfLetExpressionASTNode) statementNode()                      {}
func (node IfVarExpressionASTNode) statementNode()                      {}
func (node NullCoalesceExpressionASTNode) statementNode()               {}
func (node BubbleValueToReturnASTNode) statementNode()                  {}
func (node BinaryExpressionASTNode) statementNode()                     {}
func (node PostfixUnaryExpressionASTNode) statementNode()               {}
func (node PrefixUnaryExpressionASTNode) statementNode()                {}
func (node StructureInitilisationExpressionASTNode) statementNode()     {}
func (node StructureRefInitilisationExpressionASTNode) statementNode()  {}
func (node FunctionCallExpressionASTNode) statementNode()               {}
func (node MethodCallExpressionASTNode) statementNode()                 {}
func (node MemberExpressionASTNode) statementNode()                     {}
func (node ModulePathASTNode) statementNode()                           {}
func (node LambdaExpressionASTNode) statementNode()                     {}
func (node IfExpressionASTNode) statementNode()                         {}
func (node MatchExpressionASTNode) statementNode()                      {}
func (node WhenExpressionASTNode) statementNode()                       {}
func (node TernaryExpressionASTNode) statementNode()                    {}
func (node OptionalChainingASTNode) statementNode()                     {}
func (node TypeCastableQueryExpressionASTNode) statementNode()          {}
func (node TypeCastExpressionASTNode) statementNode()                   {}
func (node IfLetExpressionASTNode) expressionNode()                     {}
func (node IfVarExpressionASTNode) expressionNode()                     {}
func (node NullCoalesceExpressionASTNode) expressionNode()              {}
func (node BubbleValueToReturnASTNode) expressionNode()                 {}
func (node BinaryExpressionASTNode) expressionNode()                    {}
func (node PostfixUnaryExpressionASTNode) expressionNode()              {}
func (node PrefixUnaryExpressionASTNode) expressionNode()               {}
func (node StructureInitilisationExpressionASTNode) expressionNode()    {}
func (node StructureRefInitilisationExpressionASTNode) expressionNode() {}
func (node FunctionCallExpressionASTNode) expressionNode()              {}
func (node MethodCallExpressionASTNode) expressionNode()                {}
func (node MemberExpressionASTNode) expressionNode()                    {}
func (node ModulePathASTNode) expressionNode()                          {}
func (node LambdaExpressionASTNode) expressionNode()                    {}
func (node IfExpressionASTNode) expressionNode()                        {}
func (node MatchExpressionASTNode) expressionNode()                     {}
func (node WhenExpressionASTNode) expressionNode()                      {}
func (node TernaryExpressionASTNode) expressionNode()                   {}
func (node OptionalChainingASTNode) expressionNode()                    {}
func (node TypeCastableQueryExpressionASTNode) expressionNode()         {}
func (node TypeCastExpressionASTNode) expressionNode()                  {}
