package token

// Type is token type
type Type string

const (
	// ILLEGAL is unknown or error
	ILLEGAL = "ILLEGAL"
	// EOF is end of file
	EOF = "EOF"

	// IDENT is identifier (x, y,,,)
	IDENT = "IDENT"
	// INT is 1, 2,,,
	INT = "INT"

	// Operater

	// ASSIGN is =
	ASSIGN = "="
	// PLUS is +
	PLUS = "+"
	// MINUS is -
	MINUS = "-"
	// BANG is !
	BANG = "!"
	// ASTERISK is *
	ASTERISK = "*"
	// SLASH is /
	SLASH = "/"
	// LT is <
	LT = "<"
	// GT is >
	GT = ">"
	// EQ is ==
	EQ = "=="
	// NOTEQ is !=
	NOTEQ = "!="

	// Delimiter

	// COMMA is ,
	COMMA = ","
	// SEMICOLON is ;
	SEMICOLON = ";"
	// LPAREN is (
	LPAREN = "("
	// RPAREN is )
	RPAREN = ")"
	// LBRACE is {
	LBRACE = "{"
	// RBRACE is }
	RBRACE = "}"

	// Keyword

	// FUNCTION is function()
	FUNCTION = "FUNCTION"
	// LET is let
	LET = "LET"
	// TRUE is true
	TRUE = "true"
	// FALSE is false
	FALSE = "false"
	// IF is if
	IF = "if"
	// ELSE is else
	ELSE = "else"
	// RETURN is return
	RETURN = "return"
)

// Token is single token
type Token struct {
	Type    Type
	Literal string
}

var keywords = map[string]Type{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// LookupIdent return keyword or ident
func LookupIdent(ident string) Type {
	if tp, ok := keywords[ident]; ok {
		return tp
	}

	return IDENT
}
