// See https://gobyexample.com/command-line-subcommands.
package main

import (
	"fmt"
	"os"

	"github.com/jtoar/rwtools/fw"
	"github.com/jtoar/rwtools/gh"
	"github.com/jtoar/rwtools/prj"
	"github.com/jtoar/rwtools/renovate"
)

func main() {
	if len(os.Args) < 2 {
		printErrMsg()
		os.Exit(1)
	}

	switch os.Args[1] {

	case "smoke-test":
		fmt.Println("ok")

	case "framework":
		if len(os.Args) < 3 {
			fw.PrintErrMsg()
			os.Exit(1)
		}

		switch os.Args[2] {

		case "clean":
			fw.Clean()

		default:
			fw.PrintErrMsg()
			os.Exit(1)
		}

	case "project":
		if len(os.Args) < 3 {
			prj.PrintErrMsg()
			os.Exit(1)
		}

		switch os.Args[2] {

		case "clean":
			prj.Clean()

		default:
			prj.PrintErrMsg()
			os.Exit(1)
		}

	case "renovate":
		if len(os.Args) < 3 {
			renovate.PrintErrMsg()
			os.Exit(1)
		}

		switch os.Args[2] {

		case "open":
			renovate.Open()

		case "update":
			renovate.Update()

		default:
			renovate.PrintErrMsg()
			os.Exit(1)
		}

	case "github":
		if len(os.Args) < 3 {
			gh.PrintErrMsg()
			os.Exit(1)
		}

		switch os.Args[2] {

		case "cache-clean":
			gh.CacheClean()

		default:
			gh.PrintErrMsg()
			os.Exit(1)
		}

	default:
		printErrMsg()
		os.Exit(1)
	}
}

func printErrMsg() {
	fmt.Println("Expected one of: framework, github project, renovate")
}
