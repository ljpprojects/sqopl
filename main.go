package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"ljpprojects.org/sqopl/lexer"
	"ljpprojects.org/sqopl/parser"
)

var (
	inputFile  = flag.String("file", "", "Input .sqopl file to parse")
	inputFileF = flag.String("f", "", "Input .sqopl file to parse (shorthand)")
	outFile    = flag.String("out", "", "Input .sqopl file to parse")
	outFileO   = flag.String("o", "", "Input .sqopl file to parse (shorthand)")
	verbose    = flag.Bool("verbose", false, "Enable verbose output")
	verboseV   = flag.Bool("v", false, "Enable verbose output (shorthand)")
	quiet      = flag.Bool("quiet", false, "Suppress all output except errors")
	quietQ     = flag.Bool("q", false, "Suppress all output except errors (shorthand)")
	help       = flag.Bool("help", false, "Show help message")
	helpH      = flag.Bool("h", false, "Show help message (shorthand)")
	version    = flag.Bool("version", false, "Show version information")
)

const (
	VERSION = "1.0.0"
	USAGE   = `SQOPL compiler

USAGE:
    sqopl [OPTIONS] [FILE]

ARGUMENTS:
    FILE                    Input .sqopl file to compile

OPTIONS:
    -f, --file=FILE         Input .sqopl file
    -o, --out=FILE          Output file
    -v, --verbose           Enable verbose output
    -q, --quiet             Suppress all output except errors
    -h, --help              Show this help message
        --version           Show version information
`
)

func init() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, USAGE)
	}
}

func main() {
	flag.Parse()

	if *help || *helpH {
		flag.Usage()
		os.Exit(0)
	}

	if *version {
		fmt.Printf("SQOPL version %s\n", VERSION)
		os.Exit(0)
	}

	var filename string
	var isVerbose bool
	var isQuiet bool

	if *inputFile != "" {
		filename = *inputFile
	} else if *inputFileF != "" {
		filename = *inputFileF
	} else if flag.NArg() > 0 {
		filename = flag.Arg(0)
	}

	isVerbose = *verbose || *verboseV

	isQuiet = *quiet || *quietQ

	if isVerbose && isQuiet {
		fmt.Fprintf(os.Stderr, "Error: --verbose and --quiet flags cannot be used together\n")
		os.Exit(1)
	}

	if filename == "" {
		fmt.Fprint(os.Stderr, "Must specify an input file.\n")
		os.Exit(1)
	}

	if !strings.HasSuffix(filename, ".sqopl") {
		fmt.Fprintf(os.Stderr, "Warning: File '%s' does not have .sqopl extension\n", filename)
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: File '%s' does not exist\n", filename)
		os.Exit(1)
	}

	if isQuiet || !isVerbose {
		log.SetOutput(io.Discard)
	} else {
		log.SetOutput(os.Stderr)
		log.SetPrefix("[DEBUG] ")
	}

	if isVerbose {
		fmt.Fprintf(os.Stderr, "Compiling file %s\n", filename)
	}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file '%s': %v\n", filename, err)
		os.Exit(1)
	}
	defer file.Close()

	lexer := lexer.NewLexer(file)
	parser := parser.NewParser(lexer)

	for {
		mnd, err := parser.ParseStatement()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Parse error: %v\n", err)
			os.Exit(1)
		}

		nd, err := mnd.Value()
		if err != nil {
			break
		}

		fmt.Println(nd)
	}
}
