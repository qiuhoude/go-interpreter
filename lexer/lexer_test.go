package lexer

import (
	"github.com/golang/mock/gomock"
	"github.com/prashantv/gostub"
	"github.com/qiuhoude/go-interpreter/lexer/mocks"
	"github.com/qiuhoude/go-interpreter/token"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
	"time"
)

var tokenTables = []token.Token{
	{token.LET, "let"},
	{token.IDENT, "five"},
	{token.ASSIGN, "="},
	{token.INT, "5"},
	{token.SEMICOLON, ";"},
	{token.LET, "let"},
	{token.IDENT, "ten"},
	{token.ASSIGN, "="},
	{token.INT, "10"},
	{token.SEMICOLON, ";"},
	{token.LET, "let"},
	{token.IDENT, "add"},
	{token.ASSIGN, "="},
	{token.FUNCTION, "fn"},
	{token.LPAREN, "("},
	{token.IDENT, "x"},
	{token.COMMA, ","},
	{token.IDENT, "y"},
	{token.RPAREN, ")"},
	{token.LBRACE, "{"},
	{token.IDENT, "x"},
	{token.PLUS, "+"},
	{token.IDENT, "y"},
	{token.SEMICOLON, ";"},
	{token.RBRACE, "}"},
	{token.SEMICOLON, ";"},
	{token.LET, "let"},
	{token.IDENT, "result"},
	{token.ASSIGN, "="},
	{token.IDENT, "add"},
	{token.LPAREN, "("},
	{token.IDENT, "five"},
	{token.COMMA, ","},
	{token.IDENT, "ten"},
	{token.RPAREN, ")"},
	{token.SEMICOLON, ";"},

	{token.BANG, "!"},
	{token.MINUS, "-"},
	{token.SLASH, "/"},
	{token.ASTERISK, "*"},
	{token.INT, "5"},
	{token.SEMICOLON, ";"},
	{token.INT, "5"},
	{token.LT, "<"},
	{token.INT, "10"},
	{token.GT, ">"},
	{token.INT, "5"},
	{token.SEMICOLON, ";"},

	{token.IF, "if"},
	{token.LPAREN, "("},
	{token.INT, "5"},
	{token.LT, "<"},
	{token.INT, "10"},
	{token.RPAREN, ")"},
	{token.LBRACE, "{"},
	{token.RETURN, "return"},
	{token.TRUE, "true"},
	{token.SEMICOLON, ";"},
	{token.RBRACE, "}"},
	{token.ELSE, "else"},
	{token.LBRACE, "{"},
	{token.RETURN, "return"},
	{token.FALSE, "false"},
	{token.SEMICOLON, ";"},
	{token.RBRACE, "}"},

	{token.INT, "10"},
	{token.EQ, "=="},
	{token.INT, "10"},
	{token.SEMICOLON, ";"},
	{token.INT, "10"},
	{token.NOT_EQ, "!="},
	{token.INT, "9"},
	{token.SEMICOLON, ";"},
	//5 <= 5 >= 5;
	{token.INT, "5"},
	{token.LEQ, "<="},
	{token.INT, "5"},
	{token.GEQ, ">="},
	{token.INT, "5"},
	{token.SEMICOLON, ";"},

	{token.STRING, "foobar"},
	{token.STRING, "foo bar"},

	// array
	{token.LBRACKET, "["},
	{token.INT, "1"},
	{token.COMMA, ","},
	{token.INT, "2"},
	{token.RBRACKET, "]"},
	{token.SEMICOLON, ";"},

	// hashtable
	{token.HASH, "hash"},
	{token.LBRACE, "{"},
	{token.STRING, "foo"},
	{token.COLON, ":"},
	{token.STRING, "bar"},
	{token.RBRACE, "}"},

	{token.EOF, ""},
}

var input = `
let five = 5;
let ten = 10;
let add = fn(x, y) {
	x + y;
};
let result = add(five, ten);

!-/*5;
5 < 10 > 5;

if (5 < 10) {
return true;
} else {
return false;
}

10 == 10;
10 != 9;

5 <= 5 >= 5;

"foobar"
"foo bar"
[1, 2];
hash{"foo": "bar"}
`

// mock?????????Lexer
func newMockILexer(t *testing.T) (l ILexer, deferFn func()) {
	ctrl := gomock.NewController(t)
	// ????????????Mock??????
	mockLexer := mocks.NewMockILexer(ctrl)

	//MockILexer???????????? ???????????? NextToken() ???????????????????????????
	for i := range tokenTables {
		mockLexer.EXPECT().NextToken().Return(tokenTables[i])
	}
	l = mockLexer
	// go 1.14+ ??????mockgen 1.5.0+, ????????????ctrl.Finish(),?????????????????? cleanup
	deferFn = func() {
		ctrl.Finish()
	}
	return
}

// ?????????Lexer
func newLexer() (l ILexer, deferFn func()) {
	l = New(input)
	deferFn = func() {}
	return
}

func newILexer(t *testing.T) (ILexer, func()) {
	stubs := gostub.New()
	stubs.SetEnv("GO_MOCK_TEST", "0")
	defer stubs.Reset()

	env := os.Getenv("GO_MOCK_TEST")
	if env == "1" {
		return newMockILexer(t)
	}
	return newLexer()
}

func TestNextToken(t *testing.T) {
	l, fn := newILexer(t)
	defer fn()

	for i, tt := range tokenTables {
		tok := l.NextToken()
		if tok.Type != tt.Type {
			t.Fatalf("tokenTables[%d] - tokentype wrong. literal=%v expected=%q, got=%q",
				i, tt.Literal, tt.Type, tok.Type)
		}
		if tok.Literal != tt.Literal {
			t.Fatalf("tokenTables[%d] - literal wrong. expected=%q, got=%q",
				i, tt.Type, tok.Type)
		}
	}
}

// ===== GoConvey????????? ====

func TestStringSliceEqual(t *testing.T) {
	// ??????????????????????????????Convey????????????????????????
	// ??????????????????string????????????????????????
	// ???????????????????????????????????????????????????*testing.T??????
	// ?????????????????????????????????????????????????????????????????????????????????????????????
	// So???????????????????????? , ???????????????????????????, ????????????????????????????????????
	Convey("TestStringSliceEqual", t, func() {

		// ?????? Convey
		Convey("TestStringSliceEqual should return true when a != nil  && b != nil", func() {
			a := []string{"hello", "goconvey"}
			b := []string{"hello", "goconvey"}
			So(stringSliceEqual(a, b), ShouldBeTrue)

		})

		Convey("TestStringSliceEqual should return true when a ???= nil  && b ???= nil", func() {
			So(stringSliceEqual(nil, nil), ShouldBeTrue)
		})

		Convey("TestStringSliceEqual should return false when a ???= nil  && b != nil", func() {
			a := []string(nil)
			b := []string{}
			So(stringSliceEqual(a, b), ShouldBeFalse)
		})

		Convey("TestStringSliceEqual should return false when a != nil  && b != nil", func() {
			a := []string{"hello", "world"}
			b := []string{"hello", "goconvey"}
			So(stringSliceEqual(a, b), ShouldBeFalse)
		})
	})

}

func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// ===== GoStub =====
func TestStubDemo(t *testing.T) {
	Convey("TestStubDemo", t, func() {
		Convey("TestStubGlobalVariable", func() {
			stubs := gostub.Stub(&input, `let a = 1`)
			defer stubs.Reset()

			So(input, ShouldEqual, `let a = 1`)
		})

		Convey("TestStubFunc", func() {
			var timeNow = time.Now
			stubs := gostub.Stub(&timeNow, func() time.Time {
				return time.Date(2015, 6, 1, 0, 0, 0, 0, time.UTC)
			})
			defer stubs.Reset()

			So(1, ShouldEqual, timeNow().Day())
		})
	})
}
