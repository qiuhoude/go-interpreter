package repl

import (
	"bufio"
	"fmt"
	"github.com/qiuhoude/go-interpreter/evaluator"
	"github.com/qiuhoude/go-interpreter/lexer"
	"github.com/qiuhoude/go-interpreter/object"
	"github.com/qiuhoude/go-interpreter/parser"
	"io"
)

/*
Read Eval Print Loop = repl
*/

const PROMPT = ">>"

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Println(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)

		p := parser.New(l)

		program := p.ParseProgram() // AST
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}
		//_, _ = fmt.Fprintf(out, "%s\n", program.String())

		//for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		//	_, _ = fmt.Fprintf(out, "%+v\n", tok)
		//}

		evaluated := evaluator.Eval(program, object.GlobalEnv())
		if evaluated != nil {
			_, _ = fmt.Fprintf(out, "%s\n", evaluated.Inspect())
		}
	}
}
func printParserErrors(out io.Writer, errors []string) {
	_, _ = fmt.Fprint(out, " parser errors:\n")
	for _, msg := range errors {
		_, _ = fmt.Fprintf(out, "\t%v\n", msg)
	}
}
