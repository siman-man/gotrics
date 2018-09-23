package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"strconv"

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
			f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)

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

func report(data []gotrics.GoMetrics) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Function Name", "Method Length", "Parameter Count", "Nesting Level", "ABC Size"})

	for _, m := range data {
		table.Append([]string{
			m.Name,
			strconv.Itoa(m.Length),
			strconv.Itoa(m.Count),
			strconv.Itoa(m.Level),
			strconv.FormatFloat(m.ABCSize, 'f', -1, 64),
		})
	}

	table.Render()
}
