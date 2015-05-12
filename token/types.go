package token

type Type int

const (
	InvalidType Type = 1 << iota

	KeywordType    // keyword, e.g. "function", "end"
	LiteralType    // literal, e.g. 2.34, "a string", false
	MarkerType     // marker for tables, groupings, etc; e.g. {, (
	OperatorType   // operator, e.g. +, ==, &
	IdentifierType // identifier, e.g. the function: print
	Significant    = KeywordType | LiteralType | MarkerType | IdentifierType

	CommentType
	WhitespaceType
)

var tokenTypes = map[Token]Type{
	EOF:   InvalidType,
	Error: InvalidType,
	Space: WhitespaceType,

	FunctionName: IdentifierType,
	TableBegin:   MarkerType,
	TableEnd:     MarkerType,

	Comma:        MarkerType,
	StatementEnd: MarkerType,

	If:       KeywordType,
	Else:     KeywordType,
	ElseIf:   KeywordType,
	For:      KeywordType,
	Function: KeywordType,
	End:      KeywordType,
	While:    KeywordType,
	Break:    KeywordType,
	Do:       KeywordType,
	Goto:     KeywordType,
	In:       KeywordType,
	Local:    KeywordType,
	Until:    KeywordType,
	Then:     KeywordType,
	Repeat:   KeywordType,
	Return:   KeywordType,

	OpenParen:  MarkerType,
	CloseParen: MarkerType,

	Comment: CommentType,

	Nil:            IdentifierType,
	StringLiteral:  LiteralType,
	NumberLiteral:  LiteralType,
	BooleanLiteral: LiteralType,

	Identifier: IdentifierType,

	AssignmentOperator:    OperatorType,
	AdditionOperator:      OperatorType,
	SubtractionOperator:   OperatorType,
	DivisionOperator:      OperatorType,
	MultOperator:          OperatorType,
	ConcatenationOperator: OperatorType,
	VarargsOperator:       OperatorType,
	LengthOperator:        OperatorType,
	ObjectOperator:        OperatorType,
	MethodOperator:        OperatorType,
	//ComparisonOperator:      OperatorType,
	EqualityOperator:    OperatorType,
	NotEqualityOperator: OperatorType,

	NotOperator: OperatorType,
	AndOperator: OperatorType,
	OrOperator:  OperatorType,

	TableLookupOperatorLeft:  MarkerType,
	TableLookupOperatorRight: MarkerType,
	BitwiseShiftOperator:     OperatorType,
	BitwiseAndOperator:       OperatorType,
	BitwiseXorOperator:       OperatorType,
	BitwiseOrOperator:        OperatorType,
	BitwiseNotOperator:       OperatorType,

	Require: KeywordType,
}
