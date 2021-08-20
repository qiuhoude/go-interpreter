package ast

import (
	"bytes"
	"fmt"
	"github.com/qiuhoude/go-interpreter/token"
	"strings"
)

type Node interface {
	fmt.Stringer
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// ----- implementation------

// Program Node is root node
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}
func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// ================== statement ======================
// let <identifier> = <expression>;
type LetStatement struct {
	Token token.Token //the token.LET
	Name  *Identifier
	Value Expression
}

func (l *LetStatement) TokenLiteral() string { return l.Token.Literal }
func (l *LetStatement) statementNode()       {}
func (l *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(l.TokenLiteral() + " ")
	out.WriteString(l.Name.String())
	out.WriteString(" = ")
	if l.Value != nil {
		out.WriteString(l.Value.String())
	}
	out.WriteString(";")

	return out.String()
}

// ReturnStatement
// return <expression>;
type ReturnStatement struct {
	Token token.Token // the token.RETURN
	Value Expression
}

func (r *ReturnStatement) TokenLiteral() string { return r.Token.Literal }
func (r *ReturnStatement) statementNode()       {}
func (r *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(r.TokenLiteral() + " ")
	if r.Value != nil {
		out.WriteString(r.Value.String())
	}
	out.WriteString(";")

	return out.String()
}

// ExpressionStatement
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// 块语句
type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// ================== expression ======================
// PrefixExpression
// <prefix operator><expression>; eg -5;
type PrefixExpression struct {
	Token    token.Token //the token.BANG or token.MINUS or token.PLUS?
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

// InfixExpression
// <expression> <infix operator> <expression>
type InfixExpression struct {
	Token    token.Token // +,-,/,*,<,>,==,!=
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

// IfExpression
// if (<condition>) <consequence> else <alternative>
type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (i *IfExpression) expressionNode()      {}
func (i *IfExpression) TokenLiteral() string { return i.Token.Literal }
func (i *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(i.Condition.String())
	out.WriteString(" ")
	out.WriteString(i.Consequence.String())
	if i.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(i.Alternative.String())
	}
	return out.String()
}

// fn <parameters> <block statement>, fn(a,b){return a + b;}
type FunctionLiteral struct {
	Token      token.Token // token.FUNCTION
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fn *FunctionLiteral) expressionNode()      {}
func (fn *FunctionLiteral) TokenLiteral() string { return fn.Token.Literal }
func (fn *FunctionLiteral) String() string {
	var out bytes.Buffer

	var params []string
	for _, p := range fn.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(fn.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fn.Body.String())

	return out.String()
}

// 调用表达式
// 可以分查两部分identifier和参数部分中间通过 ( 分割, `(` 注册成 infixFn
// <expression>(<comma separated expressions>) , fn(x, y) { x + y; }(2, 3), add(2, 3), add(2 + 2, 3 * 3 * 3)
type CallExpression struct {
	Token     token.Token  // The '(' token
	Function  Expression   // Identifier or FunctionLiteral ,eg add(1,2), add
	Arguments []Expression // eg add(1,2), 1,2
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	var params []string
	for _, p := range ce.Arguments {
		params = append(params, p.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")

	return out.String()
}

// BlockExpression 语句块表达式, 单纯的{}语句表达式
type BlockExpression struct {
	Token token.Token // the { token
	Body  *BlockStatement
}

func (be *BlockExpression) expressionNode()      {}
func (be *BlockExpression) TokenLiteral() string { return be.Token.Literal }
func (be *BlockExpression) String() string {
	var out bytes.Buffer
	out.WriteString("{")
	out.WriteString(be.Body.String())
	out.WriteString("}")
	return out.String()
}

// 分配表达式 identifier = expression
type AssignExpression struct {
	Name  *Identifier
	Value Expression
}

func (ae *AssignExpression) expressionNode()      {}
func (ae *AssignExpression) TokenLiteral() string { return ae.Name.Token.Literal }
func (ae *AssignExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ae.Name.String())
	out.WriteString(" = ")
	if ae.Value != nil {
		out.WriteString(ae.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// ==================== 叶子节点 ==================
// IdentifierExpression
type Identifier struct {
	Token token.Token //the token.IDENT
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// IntegerExpression
type IntegerLiteral struct {
	Token token.Token //the token.INT
	Value int64
}

func (i *IntegerLiteral) expressionNode()      {}
func (i *IntegerLiteral) TokenLiteral() string { return i.Token.Literal }
func (i *IntegerLiteral) String() string       { return i.Token.Literal }

// BooleanExpression
type Boolean struct {
	Token token.Token //the token.TRUE, token.FALSE
	Value bool
}

func (i *Boolean) expressionNode()      {}
func (i *Boolean) TokenLiteral() string { return i.Token.Literal }
func (i *Boolean) String() string       { return i.Token.Literal }

// string
type StringLiteral struct {
	Token token.Token // token.STRING
	Value string
}

func (s *StringLiteral) expressionNode()      {}
func (s *StringLiteral) TokenLiteral() string { return s.Token.Literal }
func (s *StringLiteral) String() string       { return s.Token.Literal }
