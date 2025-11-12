package main

import (
	"fmt"

	"github.com/ktnuity/wet/internal/app"
	"github.com/ktnuity/wet/internal/source"
	"github.com/ktnuity/wet/internal/util"
)

func main() {
	args, err := util.GetCommandArguments()
	if err != nil {
		if v, ok := err.(*util.ArgError); ok {
			fmt.Printf("%s\n", v.Message)
		} else {
			fmt.Printf("Error: %v\n", err)
		}

		return
	}

	src, exit := source.Load(args)
	defer exit()

	if src == nil {
		fmt.Printf("No source\n")
		return
	}

	err = app.EntryPoint(src, args)
	util.ExitWithError(err, util.AsRef("Runtime failure"))
}

