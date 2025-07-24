package parser

import (
	"fmt"
	"strconv"

	"github.com/ShivankSharma070/go-interpreter/ast"
	"github.com/ShivankSharma070/go-interpreter/lexer"
	"github.com/ShivankSharma070/go-interpreter/token"
)

// Constants for deciding precedence of operators with parsing them as expression
const (
	_ = iota
	LOWEST
	EQUALS      // == or !=
	LESSGREATER // > or <
	SUM         // + or -
	PRODUCT     // * or /
	PREFIX      // -X or !X
	CALL        // myfunction(x)
	INDEX 		// Array index
)

var precedence = map[token.TokenType]int{
	token.EQ:      EQUALS,
	token.NOT_EQ:  EQUALS,
	token.LT:      LESSGREATER,
	token.GT:      LESSGREATER,
	token.PLUS:    SUM,
	token.MINUS:   SUM,
	token.SLASH:   PRODUCT,
	token.ASTERIK: PRODUCT,
	token.LPAREN:  CALL,
	token.LBRACKET : INDEX, 
}

type Parser struct {
	l            *lexer.Lexer
	currentToken token.Token
	peekToken    token.Token
	errors       []string

	// Maps to associate a token with a parser function
	prefixParserMap map[token.TokenType]prefixParserFunc
	infixParserMap  map[token.TokenType]infixParserFunc
}

type (
	prefixParserFunc func() ast.Expression               // For token found in prefix position
	infixParserFunc  func(ast.Expression) ast.Expression // For token found in infix position
)

func (p *Parser) Errors() []string {
	return p.errors
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	// Prefix Functions
	p.prefixParserMap = map[token.TokenType]prefixParserFunc{}
	p.registerPrefix(token.IDEN, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBooleanExpression)
	p.registerPrefix(token.FALSE, p.parseBooleanExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfElseExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionExpression)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.LBRACE, p.parseHashLiteral)

	// Infix Functions
	p.infixParserMap = map[token.TokenType]infixParserFunc{}
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERIK, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)

	// Read Two tokens so that currentToken and peekToken are set
	p.nextToken()
	p.nextToken()
	return p
}

// Reading next token from our lexer
func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.currentToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// Parsing return statements
func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{Token: p.currentToken}

	p.nextToken()
	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.isPeekToken(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// Parsing let Statements
func (p *Parser) parseLetStatement() ast.Statement {
	stmt := &ast.LetStatement{Token: p.currentToken}

	if !p.expectPeek(token.IDEN) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.isPeekToken(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{Token: p.currentToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.isPeekToken(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParserMap[p.currentToken.Type]
	if prefix == nil {
		p.errors = append(p.errors, fmt.Sprintf("No prefix parsing function found for %s", p.currentToken.Literal))
		return nil
	}

	leftExp := prefix()

	for !p.isPeekToken(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParserMap[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

// ===================================================================

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.currentToken}

	value, err := strconv.ParseInt(lit.TokenLiteral(), 0, 64)
	if err != nil {
		err := fmt.Sprintf("Could not parse %q as integer", lit.TokenLiteral())
		p.errors = append(p.errors, err)
	}

	lit.Value = value
	return lit
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	exp := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}

	p.nextToken()

	exp.Right = p.parseExpression(PREFIX)

	return exp
}

func (p *Parser) parseBooleanExpression() ast.Expression {
	return &ast.BoolExpression{Token: p.currentToken, Value: p.isCurToken(token.TRUE)}
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	exp := &ast.InfixExpression{
		Token:    p.currentToken,
		Left:     left,
		Operator: p.currentToken.Literal,
	}

	precedence := p.currPrecedence()
	p.nextToken()
	exp.Right = p.parseExpression(precedence)

	return exp
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseIfElseExpression() ast.Expression {
	exp := &ast.IfElseExpression{Token: p.currentToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	exp.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	exp.Consequence = p.parseBlockExpression()

	if p.isPeekToken(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		exp.Alternative = p.parseBlockExpression()
	}

	return exp
}

func (p *Parser) parseFunctionExpression() ast.Expression {
	function := &ast.FunctionExpression{Token: p.currentToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	function.Parameters = p.parseFunctionParameters()
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	function.Body = p.parseBlockExpression()

	return function
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifier := []*ast.Identifier{}

	if p.isPeekToken(token.RPAREN) {
		p.nextToken()
		return identifier
	}

	p.nextToken()
	iden := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	identifier = append(identifier, iden)

	for p.isPeekToken(token.COMMA) {
		p.nextToken()
		p.nextToken()
		iden := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
		identifier = append(identifier, iden)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifier
}

func (p *Parser) parseBlockExpression() *ast.BlockStatement {
	blockStmt := &ast.BlockStatement{Token: p.currentToken}
	blockStmt.Statements = []ast.Statement{}

	p.nextToken()

	for !p.isCurToken(token.RBRACE) && !p.isCurToken(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			blockStmt.Statements = append(blockStmt.Statements, stmt)
		}
		p.nextToken()
	}

	return blockStmt
}

func (p *Parser) parseCallExpression(exp ast.Expression) ast.Expression {
	callExp := &ast.CallExpression{Token: p.currentToken, Function: exp}
	callExp.Argument = p.parseCallArgument()
	return callExp
}

func (p *Parser) parseCallArgument() []ast.Expression {
	args := []ast.Expression{}
	args = p.parseExpressionList(token.RPAREN)
	return args
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.currentToken}
	array.Elements = p.parseExpressionList(token.RBRACKET)
	return array
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.isPeekToken(end) {
		p.nextToken()
		return list
	}
	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.isPeekToken(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression{
	exp := &ast.IndexExpression{Token: p.currentToken, Left:left}
	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
}

func (p *Parser) parseHashLiteral() ast.Expression{
	hash := &ast.HashLiteral{Token : p.currentToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.isPeekToken(token.RBRACE){
		p.nextToken()
		key := p.parseExpression(LOWEST)
		if !p.expectPeek(token.COLON) {
			return nil
		}
		p.nextToken()
		value := p.parseExpression(LOWEST)

		hash.Pairs[key] = value

		if !p.isPeekToken(token.RBRACE) && !p.expectPeek(token.COMMA){
			return nil
		}
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}
	return hash
}

// expectPeek function reads next token only if the next token is what we expect
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.isPeekToken(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) isCurToken(t token.TokenType) bool {
	return p.currentToken.Type == t
}

func (p *Parser) isPeekToken(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peekError(t token.TokenType) {
	err := fmt.Sprintf("Expected next token to be %s , got %s", t, p.peekToken.Type)
	p.errors = append(p.errors, err)
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedence[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) currPrecedence() int {
	if p, ok := precedence[p.currentToken.Type]; ok {
		return p
	}

	return LOWEST
}

// Helper function to asscoiate parser fucntion for a token in prefix position
func (p *Parser) registerPrefix(tokentype token.TokenType, fn prefixParserFunc) {
	p.prefixParserMap[tokentype] = fn
}

// Helper function to asscoiate parser fucntion for a token in prefix position
func (p *Parser) registerInfix(tokentype token.TokenType, fn infixParserFunc) {
	p.infixParserMap[tokentype] = fn
}
