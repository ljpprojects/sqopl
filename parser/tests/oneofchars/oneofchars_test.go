package oneofchars

import (
	"fmt"
	"os"
	"testing"

	"ljpprojects.org/sqopl/lexer"
	"ljpprojects.org/sqopl/parser"
)

func TestExpectOneOfChars(t *testing.T) {
	allowedChars := map[rune]*lexer.TokenGroup{
		'[': &lexer.TokenGroupingGroup,
		'{': &lexer.TokenGroupingGroup,
	}

	/**** Test the first working example ****/

	filename := "./works 1.sqopl"

	file, err := os.Open(filename)

	if err != nil {
		t.Fatalf("Error opening file '%s': %v\n", filename, err)
	}

	defer file.Close()

	flexer := lexer.NewLexer(file)
	fparser := parser.NewParser(flexer)

	tk, err := fparser.ExpectOneOfCharacters(allowedChars)

	if err != nil {
		switch err := err.(type) {
		case parser.ParseErrorExpectedOneOfCharacters:
			t.Fatalf("Error: expect failed when it should have succeeded.\nError: %s\n", err.Error())
		default:
			t.Fatalf("Expect failed (unexpected outcome) with error: %s.", err.Error())
		}
	}

	t.Logf("Expect succeeded; found token: %s", tk.ToDisplayString())

	file.Close()

	/**** Test the second working example ****/

	filename = "./works 2.sqopl"

	file, err = os.Open(filename)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file '%s': %v\n", filename, err)
		os.Exit(1)
	}

	flexer = lexer.NewLexer(file)
	fparser = parser.NewParser(flexer)

	tk, err = fparser.ExpectOneOfCharacters(allowedChars)

	if err != nil {
		switch err := err.(type) {
		case parser.ParseErrorExpectedOneOfCharacters:
			t.Fatalf("Error: expect failed when it should have succeeded.\nError: %s\n", err.Error())
		default:
			t.Fatalf("Expect failed (unexpected outcome) with error: %s.", err.Error())
		}
	}

	t.Logf("Expect succeeded; found token: %s", tk.ToDisplayString())

	file.Close()

	/**** Test the non-working example ****/

	filename = "./fails.sqopl"

	file, err = os.Open(filename)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file '%s': %v\n", filename, err)
		os.Exit(1)
	}

	flexer = lexer.NewLexer(file)
	fparser = parser.NewParser(flexer)

	tk, err = fparser.ExpectOneOfCharacters(allowedChars)

	if err != nil {
		switch err := err.(type) {
		case parser.ParseErrorExpectedOneOfCharacters:
			t.Logf("Expect failed (expected outcome) with error: %s.", err.Error())
			t.SkipNow()
		default:
			t.Fatalf("Expect failed (unexpected outcome) with error: %s.", err.Error())
		}
	}

	t.Fatalf("Error: expect succeeded when it should have failed; found token: %s", tk.ToDisplayString())

	file.Close()
}
