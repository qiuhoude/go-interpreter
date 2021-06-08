package parser

import (
	"fmt"
	"github.com/qiuhoude/go-interpreter/ast"
	"github.com/qiuhoude/go-interpreter/lexer"
	"github.com/qiuhoude/go-interpreter/token"
	"strconv"
)

const (
	// precedence operator
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <  >= <=
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

// precedence table , use in infix expression parse
var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.GT:       LESSGREATER,
	token.LT:       LESSGREATER,
	token.GEQ:      LESSGREATER,
	token.LEQ:      LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression // argument is left side
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token // cur point
	peekToken token.Token // next point

	prefixParseFns map[token.TokenType]prefixParseFn // 前缀解析方法
	infixParseFns  map[token.TokenType]infixParseFn  // 中缀解析方法
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:              l,
		errors:         []string{},
		prefixParseFns: map[token.TokenType]prefixParseFn{},
		infixParseFns:  map[token.TokenType]infixParseFn{},
	}

	p.RegisterPrefix(token.IDENT, p.parseIdentifier)
	p.RegisterPrefix(token.INT, p.parseIntegerLiteral)
	p.RegisterPrefix(token.BANG, p.parsePrefixExpression)
	p.RegisterPrefix(token.MINUS, p.parsePrefixExpression)
	p.RegisterPrefix(token.PLUS, p.parsePrefixExpression)

	p.RegisterInfix(token.EQ, p.parseInfixExpression)
	p.RegisterInfix(token.NOT_EQ, p.parseInfixExpression)
	p.RegisterInfix(token.GT, p.parseInfixExpression)
	p.RegisterInfix(token.LT, p.parseInfixExpression)
	p.RegisterInfix(token.GEQ, p.parseInfixExpression)
	p.RegisterInfix(token.LEQ, p.parseInfixExpression)
	p.RegisterInfix(token.PLUS, p.parseInfixExpression)
	p.RegisterInfix(token.MINUS, p.parseInfixExpression)
	p.RegisterInfix(token.SLASH, p.parseInfixExpression)
	p.RegisterInfix(token.ASTERISK, p.parseInfixExpression)

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) RegisterPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) RegisterInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	defer untrace(trace("parseExpressionStatement"))

	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// 先跳过 expressions ,直到我们遇到 ;
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	// 解析 return <expression>;

	p.nextToken()

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	defer untrace(trace("parseExpression"))
	prefix, ok := p.prefixParseFns[p.curToken.Type]
	if !ok {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()
	// key code
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix, ok := p.infixParseFns[p.peekToken.Type]
		if !ok {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	defer untrace(trace("parseIntegerLiteral"))
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	defer untrace(trace("parseIntegerLiteral"))
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit := &ast.IntegerLiteral{Token: p.curToken, Value: value}
	return lit
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	defer untrace(trace("parsePrefixExpression"))
	pe := &ast.PrefixExpression{Token: p.curToken, Operator: p.curToken.Literal}

	p.nextToken()
	pe.Right = p.parseExpression(PREFIX)
	return pe
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	defer untrace(trace("parseInfixExpression"))
	exp := &ast.InfixExpression{
		Token:    p.curToken,
		Left:     left,
		Operator: p.curToken.Literal,
	}

	precedences := p.curPrecedence()
	p.nextToken()
	if exp.Operator == "+" {
		exp.Right = p.parseExpression(precedences - 1)
	} else {
		exp.Right = p.parseExpression(precedences)
	}

	return exp
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken() // 进行下一步
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

/*
// pseudo-code



function parseLetStatement() {
advanceTokens()
identifier = parseIdentifier()
advanceTokens()
if currentToken() != EQUAL_TOKEN {
parseError("no equal sign!")
return null
}
advanceTokens()
value = parseExpression()
variableStatement = newVariableStatementASTNode()
variableStatement.identifier = identifier
variableStatement.value = value
return variableStatement
}
function parseIdentifier() {
identifier = newIdentifierASTNode()
identifier.token = currentToken()
return identifier
}
function parseExpression() {
if (currentToken() == INTEGER_TOKEN) {
if (nextToken() == PLUS_TOKEN) {
return parseOperatorExpression()
} else if (nextToken() == SEMICOLON_TOKEN) {
return parseIntegerLiteral()
}
} else if (currentToken() == LEFT_PAREN) {
return parseGroupedExpression()
}
// [...]
}
function parseOperatorExpression() {
operatorExpression = newOperatorExpression()
operatorExpression.left = parseIntegerLiteral()
operatorExpression.operator = currentToken()
operatorExpression.right = parseExpression()
return operatorExpression()
}
// [...]
*/
