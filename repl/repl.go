package repl

import (
	"bufio"
	"io"
	"fmt"

	"monkey/lexer"
	"monkey/token"
)

const PROMPT = "🐒 "

func Start(reader io.Reader, writer io.Writer) {
	scanner := bufio.NewScanner(reader)

	for {
		fmt.Printf(PROMPT)

		if !scanner.Scan() {
			return
		}

		line := scanner.Text()
		lex  := lexer.New(line)

		for currToken := lex.NextToken(); currToken.Type != token.EOF; currToken = lex.NextToken() {
			// %v -> default value
			// +  -> field names added (if value is struct). otherwise, ignored
			fmt.Printf("%+v\n", currToken)
		}
	}
}
