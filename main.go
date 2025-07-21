package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"ljpprojects.org/sqopl/lexer"
	"ljpprojects.org/sqopl/parser"
)

func main() {
	log.SetOutput(io.Discard)

	file, err := os.Open("./test.sqopl")

	if err != nil {
		log.SetOutput(os.Stderr)
		log.Fatal(err)
	}

	lexer := lexer.NewLexer(file)
	parser := parser.NewParser(lexer)

	for {
		mnd, err := parser.ParseStatement()

		if err != nil {
			log.SetOutput(os.Stderr)
			log.Fatal(err)
		}

		nd, err := mnd.Value()

		if err != nil {
			break
		}

		fmt.Println(nd)
		fmt.Printf("%d-%d\n", nd.Location().Start, nd.Location().End)
	}
}
