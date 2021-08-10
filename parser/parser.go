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
	token.LPAREN:   CALL,
}

// 前缀 和 中缀解析函数
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
	p.RegisterPrefix(token.TRUE, p.parseBoolean)
	p.RegisterPrefix(token.FALSE, p.parseBoolean)
	p.RegisterPrefix(token.LPAREN, p.parseGroupedExpression)

	p.RegisterPrefix(token.IF, p.parseIfExpression)
	p.RegisterPrefix(token.FUNCTION, p.parseFunctionLiteral)

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
	p.RegisterInfix(token.LPAREN, p.parseCallExpression) // call

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

// 解析语句
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
	p.nextToken() // 跳过 `=`

	stmt.Value = p.parseExpression(LOWEST)

	// 先跳过 expressions ,直到我们遇到 ;
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	defer untrace(trace("parseReturnStatement"))
	stmt := &ast.ReturnStatement{Token: p.curToken}
	// 解析 return <expression>;

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	blockStmt := &ast.BlockStatement{Token: p.curToken}

	p.nextToken() // skip {

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			blockStmt.Statements = append(blockStmt.Statements, stmt)
		}
		p.nextToken()
	}
	return blockStmt
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	defer untrace(trace("parseExpression"))
	// 流程就是使用递归方式构建多叉树 AST
	// 前缀
	prefix, ok := p.prefixParseFns[p.curToken.Type]
	if !ok {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	// 有前缀,可能是中缀的左值
	leftExp := prefix()

	// 中缀 key code
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		// 下一个token 不是`;` , 并且下一个token优先级大于传入参数优先级
		infix, ok := p.infixParseFns[p.peekToken.Type]
		// 不是中缀,跳出循环
		if !ok {
			//return leftExp
			break
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
	exp.Right = p.parseExpression(precedences)

	return exp
}

func (p *Parser) parseBoolean() ast.Expression {
	defer untrace(trace("parseBoolean"))
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	defer untrace(trace("parseGroupedExpression"))
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	defer untrace(trace("parseIfExpression"))
	exp := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) { // `if` (
		return nil
	}

	p.nextToken()                             // 跳过 (
	exp.Condition = p.parseExpression(LOWEST) // `if (` condition

	if !p.expectPeek(token.RPAREN) { // `if ( condition` )
		return nil
	}

	if !p.expectPeek(token.LBRACE) { // `if ( condition )` {
		return nil
	}

	exp.Consequence = p.parseBlockStatement() // `if ( condition )` { Consequence }

	if p.peekTokenIs(token.ELSE) { // else 部分
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		exp.Alternative = p.parseBlockStatement()
	}

	return exp
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	defer untrace(trace("parseFunctionLiteral"))
	exp := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) { // `fn` (
		return nil
	}
	exp.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) { // `fn ( params... )` {
		return nil
	}

	exp.Body = p.parseBlockStatement()
	return exp
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	defer untrace(trace("parseFunctionParameters"))
	// fn (a,b,c...)
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}
	p.nextToken()
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) { // 多个参数 (a,b,c)
		p.nextToken() // cur指向 `,`
		p.nextToken() // cur指向 `参数`

		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) { // 不是 ) 结束
		return nil
	}

	return identifiers
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	defer untrace(trace("parseCallExpression"))
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseCallParameters()
	return exp
}

func (p *Parser) parseCallParameters() []ast.Expression {
	defer untrace(trace("parseCallParameters"))
	args := []ast.Expression{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}
	p.nextToken()
	args = append(args, p.parseExpression(LOWEST)) // 添加表达式

	for p.peekTokenIs(token.COMMA) { // 多个参数 (a,b,c)
		p.nextToken() // cur指向 `,`
		p.nextToken() // cur指向 `参数`
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) { // 不是 ) 结束
		return nil
	}
	return args
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
	return p.precedence(p.peekToken.Type)
}

func (p *Parser) curPrecedence() int {
	return p.precedence(p.curToken.Type)
}

func (p *Parser) precedence(tk token.TokenType) int {
	if p, ok := precedences[tk]; ok {
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
