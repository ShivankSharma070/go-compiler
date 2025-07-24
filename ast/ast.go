package ast

import (
	"bytes"
	"strings"

	"github.com/ShivankSharma070/go-interpreter/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}
type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

// Token literal will return the literal value of token associated with a node
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		// First node's token literal value
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var buf bytes.Buffer
	for _, stmt := range p.Statements {
		buf.WriteString(stmt.String())
	}

	return buf.String()
}

// Let Statements
type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

// To implement statement & node interface
func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (ls *LetStatement) String() string {
	var buf bytes.Buffer

	buf.WriteString(ls.TokenLiteral() + " ")
	buf.WriteString(ls.Name.String())
	buf.WriteString(" = ")

	if ls.Value != nil {
		buf.WriteString(ls.Value.String())
	}

	buf.WriteString(";")
	return buf.String()
}

// Here identifier is implementing the expression interface. Why ??
// Although expression are those which produces some values,
// and identifier in let x = 5 does not produces a vlaue
// It is to keep the number of different nodes type minimum
// And a identifier will produce value in statement like let a = b; (b produces value)
type Identifier struct {
	Token token.Token
	Value string
}

// To implement expression & node interface
func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Value
}

// String Expression
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) String() string       { return sl.Token.Literal }
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }

// Return Statments
type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReturnStatement) String() string {
	var buf bytes.Buffer
	buf.WriteString(rs.TokenLiteral() + " ")
	buf.WriteString(rs.ReturnValue.String())

	buf.WriteString(";")
	return buf.String()
}

// Expression statement
// 5+10, (5*10)+5, foo(a,b) etc
// We are treating expression as statements because, we want to allow one line containing only expression as a statement
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	buf := bytes.Buffer{}
	buf.WriteString("(")
	buf.WriteString(pe.Operator)
	buf.WriteString(pe.Right.String())
	buf.WriteString(")")

	return buf.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var buf bytes.Buffer
	buf.WriteString("(")
	buf.WriteString(ie.Left.String())
	buf.WriteString(" " + ie.Operator + " ")
	buf.WriteString(ie.Right.String())
	buf.WriteString(")")
	return buf.String()
}

type BoolExpression struct {
	Token token.Token
	Value bool
}

func (be *BoolExpression) expressionNode()      {}
func (be *BoolExpression) TokenLiteral() string { return be.Token.Literal }
func (be *BoolExpression) String() string       { return be.Token.Literal }

type IfElseExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfElseExpression) expressionNode()      {}
func (ie *IfElseExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfElseExpression) String() string {
	var buf bytes.Buffer
	buf.WriteString("if ")
	buf.WriteString(ie.Condition.String())
	buf.WriteString(" ")
	buf.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		buf.WriteString("else ")
		buf.WriteString(ie.Alternative.String())
	}

	return buf.String()
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (be *BlockStatement) statementNode()       {}
func (be *BlockStatement) TokenLiteral() string { return be.Token.Literal }
func (be *BlockStatement) String() string {
	var buf bytes.Buffer
	for _, st := range be.Statements {
		buf.WriteString(st.String())
	}

	return buf.String()
}

type FunctionExpression struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fe *FunctionExpression) expressionNode()      {}
func (fe *FunctionExpression) TokenLiteral() string { return fe.Token.Literal }
func (fe *FunctionExpression) String() string {
	var buf bytes.Buffer

	parameters := []string{}
	for _, iden := range fe.Parameters {
		parameters = append(parameters, iden.String())
	}
	buf.WriteString(fe.TokenLiteral())
	buf.WriteString("(")
	buf.WriteString(strings.Join(parameters, ", "))
	buf.WriteString(")")
	buf.WriteString(fe.Body.String())
	return buf.String()
}

type CallExpression struct {
	Token    token.Token
	Function Expression
	Argument []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var buf bytes.Buffer
	arguments := []string{}
	for _, exp := range ce.Argument {
		arguments = append(arguments, exp.String())
	}
	buf.WriteString(ce.Function.String())
	buf.WriteString("(")
	buf.WriteString(strings.Join(arguments, ", "))
	buf.WriteString(")")
	return buf.String()
}

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}

type HashLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode()      {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String() string {
	var out bytes.Buffer
	pairs := []string{}
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}
