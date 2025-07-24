package token

type TokenType string

const (
	EOF     = "EOF"
	ELLEGAL = "ELLEGAL"

	// Identifier and Literals
	IDEN = "IDEN" // Variable names
	INT  = "INT"
	STRING = "STRING"

	// Operators
	ASSIGN  = "="
	PLUS    = "+"
	ASTERIK = "*"
	MINUS   = "-"
	BANG    = "!"
	SLASH   = "/"
	GT      = ">"
	LT      = "<"

	EQ     = "=="
	NOT_EQ = "!="

	//Delimeters
	SEMICOLON = ";"
	COLON = ":"
	COMMA     = ","

	// Braces
	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

type Token struct {
	Type    TokenType
	Literal string
}

var keyword = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func LookUpIden(iden string) TokenType {
	if tokType, ok := keyword[iden]; ok {
		return tokType
	}
	return IDEN
}

func NewToken(t TokenType, char byte) Token {
	return Token{t, string(char)}
}
