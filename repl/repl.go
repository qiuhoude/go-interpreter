package repl

import (
	"bufio"
	"fmt"
	"github.com/qiuhoude/go-interpreter/lexer"
	"github.com/qiuhoude/go-interpreter/token"
	"io"
)

/*
Read Eval Print Loop = repl
*/

const PROMPT = ">>"

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		_, _ = fmt.Fprintln(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			_, _ = fmt.Fprintf(out, "%+v\n", tok)

		}
	}
}
