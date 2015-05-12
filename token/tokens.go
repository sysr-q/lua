package token

import (
	"sort"
	"strconv"
)

type Token int

const (
	EOF Token = iota
	Error
	Space

	FunctionName
	TableBegin
	TableEnd

	Comma
	StatementEnd

	If
	Else
	ElseIf
	For
	Function
	End
	While
	Break
	Do
	Goto
	In
	Local
	Until
	Then
	Repeat
	Return

	OpenParen
	CloseParen

	Comment

	Nil
	StringLiteral
	NumberLiteral
	BooleanLiteral

	Identifier

	AssignmentOperator
	AdditionOperator
	SubtractionOperator
	DivisionOperator
	MultOperator
	ConcatenationOperator
	VarargsOperator
	LengthOperator
	ObjectOperator
	MethodOperator
	ComparisonOperator
	EqualityOperator
	NotEqualityOperator

	NotOperator
	AndOperator
	OrOperator

	TableLookupOperatorLeft
	TableLookupOperatorRight
	BitwiseShiftOperator
	BitwiseAndOperator
	BitwiseXorOperator
	BitwiseOrOperator
	BitwiseNotOperator

	Require
)

var tokens = []string{
	EOF:   "EOF",
	Space: "space",

	FunctionName: "function name",
	TableBegin:   "table begin",
	TableEnd:     "table end",

	Comma:        "comma",
	StatementEnd: ";",

	If:       "if",
	Else:     "else",
	ElseIf:   "elseif",
	For:      "for",
	Function: "function",
	End:      "end",
	While:    "while",
	Break:    "break",
	Do:       "do",
	Goto:     "goto",
	In:       "in",
	Local:    "local",
	Until:    "until",
	Then:     "then",
	Repeat:   "repeat",
	Return:   "return",

	OpenParen:  "open-paren",
	CloseParen: "close-paren",

	Comment: "comment", // "--[[ --]]",

	Nil:            "nil",
	StringLiteral:  "string-literal",
	NumberLiteral:  "number-literal",
	BooleanLiteral: "bool-literal",

	Identifier: "identifier",

	AssignmentOperator:    "=",
	AdditionOperator:      "+",
	SubtractionOperator:   "-",
	DivisionOperator:      "/|//",
	MultOperator:          "*%",
	ConcatenationOperator: "..",
	VarargsOperator:       "...",
	LengthOperator:        "#",
	ObjectOperator:        ".",
	MethodOperator:        ":",
	// ComparisonOperator: ">=><=<",
	EqualityOperator:    "==",
	NotEqualityOperator: "~=",

	NotOperator: "not",
	AndOperator: "and",
	OrOperator:  "or",

	TableLookupOperatorLeft:  "[",
	TableLookupOperatorRight: "]",
	BitwiseShiftOperator:     "<<>>",
	BitwiseAndOperator:       "&",
	BitwiseXorOperator:       "^",
	BitwiseOrOperator:        "|",
	BitwiseNotOperator:       "~",

	Require: "require",
}

var TokenList []string

func init() {
	TokenList = make([]string, len(TokenMap))
	i := 0
	for token := range TokenMap {
		TokenList[i] = token
		i += 1
	}
	sort.Sort(sort.Reverse(sort.StringSlice(TokenList)))
}

// TokenMap maps source code string tokens to  types when strings can
// be represented directly. Not all  types will be represented here.
var TokenMap = map[string]Token{
	"{": TableBegin,
	"}": TableEnd,

	",": Comma,
	";": StatementEnd,

	"if":       If,
	"else":     Else,
	"elseif":   ElseIf,
	"for":      For,
	"function": Function,
	"end":      End,
	"while":    While,
	"break":    Break,
	"do":       Do,
	"goto":     Goto,
	"in":       In,
	"local":    Local,
	"until":    Until,
	"then":     Then,
	"repeat":   Repeat,
	"return":   Return,

	"(": OpenParen,
	")": CloseParen,

	"--":   Comment,
	"--[[": Comment,
	"--]]": Comment,

	"nil":   Nil,
	"true":  BooleanLiteral,
	"false": BooleanLiteral,

	"=":   AssignmentOperator,
	"+":   AdditionOperator,
	"-":   SubtractionOperator,
	"/":   DivisionOperator,
	"//":  DivisionOperator,
	"*":   MultOperator,
	"%":   MultOperator,
	"..":  ConcatenationOperator,
	"...": VarargsOperator,
	"#":   LengthOperator,
	".":   ObjectOperator,
	":":   MethodOperator,
	">=":  ComparisonOperator,
	">":   ComparisonOperator,
	"<=":  ComparisonOperator,
	"<":   ComparisonOperator,
	"==":  EqualityOperator,
	"~=":  NotEqualityOperator,

	"not": NotOperator,
	"and": AndOperator,
	"or":  OrOperator,

	"[":  TableLookupOperatorLeft,
	"]":  TableLookupOperatorRight,
	"<<": BitwiseShiftOperator,
	">>": BitwiseShiftOperator,
	"&":  BitwiseAndOperator,
	"^":  BitwiseXorOperator,
	"|":  BitwiseOrOperator,
	"~":  BitwiseNotOperator,

	"require": Require,
}

func (i Token) String() string {
	TypeName := tokens[i]
	if len(TypeName) == 0 {
		return strconv.Itoa(int(i))
	}
	return TypeName
}

func (i Token) IsType(t Type) bool {
	return i.Type()&t != 0
}

func (i Token) Type() Type {
	// tokenTypes is in types.go
	typ, ok := tokenTypes[i]
	if !ok {
		panic("token without type")
	}
	return typ
}
