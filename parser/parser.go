package parser

import (
	"fmt"
	"github.com/qiuhoude/go-interpreter/ast"
	"github.com/qiuhoude/go-interpreter/lexer"
	"github.com/qiuhoude/go-interpreter/token"
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.NextToken()
	p.NextToken()
	return p
}

func (p *Parser) NextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
		//case token.IF:

	}
	return nil
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.NextToken()
	}
	return program
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
		p.NextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	// 解析 return <expression>;

	p.NextToken()

	for !p.curTokenIs(token.SEMICOLON) {
		p.NextToken()
	}

	return stmt
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.NextToken() // 进行下一步
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
