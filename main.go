package main

import (
	"bufio"
	"log"
	"os"

	"ljpprojects.org/sqopl/lexer"
)

func main() {
	file, err := os.Open("./test.sqopl")

	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(file)
	lexer := lexer.NewLexer(reader)

	for {
		ntk, err := lexer.NextToken()

		if err != nil {
			log.Fatal(err)
		}

		tk, err := ntk.Value()

		if err != nil {
			log.Fatal(err)
		}

		log.Println(tk)
	}
}
