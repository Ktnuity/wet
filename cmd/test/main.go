package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/ktnuity/wet/internal/test"
	"github.com/ktnuity/wet/internal/util/testutil"
)

func main() {
	args := test.GetArgs()
	tests := test.GetTests()
	match, err := test.LoadTest()
	if err != nil {
		fmt.Printf("Failed to load test: %v\n", err)
		os.Exit(1)
	}

	result := make([]string, 0, 8)
	live := len(match) > 0 && !args.Gen

	executable := "./wet"
	testExec := "./test"
	if runtime.GOOS == "windows" {
		executable = ".\\wet.exe"
		testExec = ".\\test.exe"
	}

	if live || args.Gen {
		printf := func(format string, args...any) {
			result = append(result, fmt.Sprintf(format, args...))
		}

		for idx, step := range tests {
			stepName := step 
			if strings.HasPrefix(step, "--") {
				stepName = fmt.Sprintf("%s %s", executable, step)
			} else if len(step) == 0 {
				stepName = executable
			}
			printf("Test %d | %s", idx, stepName)

			testutil.Printf("Testing %s <lightgray>|<reset> ", stepName)

			res, ok := test.RunFile(step, args.Forward)
			result = append(result, res...)

			if !ok {
				testutil.Printf("<red>fail<reset>\n")
				fmt.Printf("%s\n", strings.Join(res, "\n"))
				os.Exit(1)
			}

			if live {
				disc := testutil.FindDiscrepency(match, result)
				if len(disc) > 0 {
					testutil.Printf("<red>fail<reset>\n")
					fmt.Printf("%s\n", strings.Join(disc, "\n"))
					os.Exit(1)
				}

				testutil.Printf("<green>ok<reset>\n")
			} else {
				testutil.Printf("<cyan>saving...<reset>\n")
			}
		}
	} else {
		fmt.Printf("No tests found.\nGenerate tests with %s --gen\n", testExec)
		os.Exit(1)
	}

	fmt.Printf("Test Done!\n")

	if args.Gen {
		testFile, err := test.GetTestFile()
		if err != nil {
			fmt.Printf("Failed to get test file for gen saving: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Saving test...\n")
		test := strings.Join(result, "\n")
		os.WriteFile(testFile, []byte(test), 0644)
	}
}

