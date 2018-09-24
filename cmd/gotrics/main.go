package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"

	"github.com/siman-man/gotrics"
)

var (
	// gotrics options
	format = flag.String("f", "", "set output format")
)

var (
	exitCode = 0
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: gotrics [flags] [path ...]\n")
	flag.PrintDefaults()
	os.Exit(0)
}

func main() {
	gotricsMain()
	os.Exit(exitCode)
}

func gotricsMain() {
	flag.Usage = usage
	flag.Parse()

	for i := 0; i < flag.NArg(); i++ {
		path := flag.Arg(i)

		switch dir, err := os.Stat(path); {
		case err != nil:
			exitCode = 2
			return
		case dir.IsDir():
			exitCode = 2
			fmt.Println("Not support directory path.")
			return
		default:
			fset := token.NewFileSet()
			f, err := parseSourceCode(fset, path)

			if err != nil {
				exitCode = 2
				return
			}

			result := gotrics.Analyze(fset, f)

			if *format == "json" {
				str, err := json.Marshal(result)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(string(str))

			} else {
				report(result)
			}
		}
	}
}

func parseSourceCode(fset *token.FileSet, path string) (f *ast.File, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	src, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	f, err = parser.ParseFile(fset, path, src, parser.ParseComments)

	if err == nil || !strings.Contains(err.Error(), "expected 'package'") {
		return
	}

	psrc := append([]byte("package p;"), src...)
	f, err = parser.ParseFile(fset, path, psrc, parser.ParseComments)

	return
}

func report(data []gotrics.GoMetrics) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Function Name", "Method Length", "Parameter Count", "Nesting Level", "ABC Size"})

	for _, m := range data {
		table.Append([]string{
			m.Name,
			strconv.Itoa(m.MethodLength),
			strconv.Itoa(m.ParameterCount),
			strconv.Itoa(m.NestingLevel),
			strconv.FormatFloat(m.ABCSize, 'f', -1, 64),
		})
	}

	table.Render()
}
