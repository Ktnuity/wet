package test

import (
	"fmt"
	"os"
	"strings"
)

type Args struct {
	Gen			bool
	Forward		[]string
}

func GetArgs() Args {
	argv := os.Args
	argc := len(argv)

	res := Args{}

	for argi := 1; argi < argc; argi++ {
		if argv[argi] == "--" {
			res.Forward = argv[argi+1:]
			break
		}

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
