package lexer

import (
	"github.com/qiuhoude/go-interpreter/token"
)

/*
脚本代码目前只支持 ASCII
number类型只支持 Integer
*/

// 词法分析器
type ILexer interface {
	NextToken() token.Token
}

type Lexer struct {
	input        string
	position     int  // 当前的位置
	readPosition int  // 当前读到的位置
	ch           byte // 当前char
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l

}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // 0 -> ASCII code is NUL
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	// 跳过空格
	l.skipWhitespace()

	switch l.ch {
	case '=': // = , ==
		if l.peekChar() == '=' { // ==
			preCh := l.ch
			l.readChar()
			tok = makeStrCharToken(token.EQ, string(preCh)+string(l.ch))
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!': // !, !=
		if l.peekChar() == '=' { // !=
			preCh := l.ch
			l.readChar()
			tok = makeStrCharToken(token.NOT_EQ, string(preCh)+string(l.ch))
		} else { // !
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		if l.peekChar() == '=' { // <=
			preCh := l.ch
			l.readChar()
			tok = makeStrCharToken(token.LEQ, string(preCh)+string(l.ch))
		} else {
			tok = newToken(token.LT, l.ch)
		}
	case '>':
		if l.peekChar() == '=' { // >=
			preCh := l.ch
			l.readChar()
			tok = makeStrCharToken(token.GEQ, string(preCh)+string(l.ch))
		} else {
			tok = newToken(token.GT, l.ch)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) { // 字母
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok // readIdentifier()里面已经 调用了 l.readChar() 所以要return
		} else if isDigit(l.ch) { // 数字
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else { // 非法字符
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool { // [a-z|A-Z|_]
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func makeStrCharToken(tokenType token.TokenType, str string) token.Token {
	return token.Token{Type: tokenType, Literal: str}
}
