package test

import (
	"fmt"
	"os"
	"strings"
)

type Flags struct {
	Gen			bool
}

func GetArgs() Flags {
	argv := os.Args
	argc := len(argv)

	res := Flags{}

	for argi := 1; argi < argc; argi++ {
		if flag, ok := strings.CutPrefix(argv[argi], "--"); ok {
			switch flag {
			case "gen":
				res.Gen = true
			default:
				fmt.Printf("Unknown flag '%s'\n", argv[argi])
			}
		} else {
			fmt.Printf("Unknown argument '%s'\n", argv[argi])
		}
	}

	return res
}
