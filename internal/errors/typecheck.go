package errors

import (
	"fmt"
	"os"
	"strings"

	"github.com/ktnuity/wet/internal/types"
)

type PreparedTypeCheck struct {
	operator		*types.Word
	info			[]string
}

func BadTypeCheck(operator *types.Word, info...string) {
	fmt.Printf("Type Checker Failed!\n")
	fmt.Printf("Operator: %s\n", operator.UnwrapName())

	traceLines := trace(operator)

	lines := make([]string, 0, len(info) + len(traceLines) + 1)

	lines = append(lines, traceLines...)
	lines = append(lines, "")
	lines = append(lines, info...)

	for _, line := range lines {
		fmt.Printf(" %s\n", line)
	}

	os.Exit(1)
}

func PrepareTypeCheck(operator *types.Word, info...string) *PreparedTypeCheck {
	return &PreparedTypeCheck{
		operator, info,
	}
}

func (p *PreparedTypeCheck) With(info...string) *PreparedTypeCheck {
	p.info = append(p.info, info...)

	return p
}

func (p *PreparedTypeCheck) Throw(info...string) {
	infos := []string{}
	infos = append(infos, p.info...)
	infos = append(infos, info...)
	BadTypeCheck(p.operator, infos...)
}

func (p *PreparedTypeCheck) Empty() {
	p.Throw("Stack is empty.")
}

func (p *PreparedTypeCheck) Stack(expect, actual int) {
	if expect <= 1 {
		p.Empty()
	}

	p.Throw(fmt.Sprintf("Stack size is %d. %d is required.", actual, expect))
}

func (p *PreparedTypeCheck) GetValue(index int) {
	names := []string{
		"first",
		"second",
		"third",
		"fourth",
		"fifth",
		"sixth",
		"seventh",
		"eighth",
		"nineth",
		"tenth",
	}

	if index < 0 {
		p.Throw("Failed to get value.")
	}

	p.Throw(fmt.Sprintf("Failed to get %s value.", names[index]))
}

func (p *PreparedTypeCheck) GetNameValue(index int, name string) {
	names := []string{
		"first",
		"second",
		"third",
		"fourth",
		"fifth",
		"sixth",
		"seventh",
		"eighth",
		"nineth",
		"tenth",
	}

	if index < 0 {
		p.Throw(fmt.Sprintf("Failed to get %s.", name))
	}

	p.Throw(fmt.Sprintf("Failed to get %s %s.", names[index], name))
}

func (p *PreparedTypeCheck) GetValueType(index int) {
	names := []string{
		"first",
		"second",
		"third",
		"fourth",
		"fifth",
		"sixth",
		"seventh",
		"eighth",
		"nineth",
		"tenth",
	}

	if index < 0 {
		p.Throw("Failed to get value type.")
	}

	p.Throw(fmt.Sprintf("Failed to get %s value type.", names[index]))
}

func (p *PreparedTypeCheck) GetNameValueType(index int, name string) {
	names := []string{
		"first",
		"second",
		"third",
		"fourth",
		"fifth",
		"sixth",
		"seventh",
		"eighth",
		"nineth",
		"tenth",
	}

	if index < 0 {
		p.Throw(fmt.Sprintf("Failed to get %s type.", name))
	}

	p.Throw(fmt.Sprintf("Failed to get %s %s type.", names[index], name))
}

func (p *PreparedTypeCheck) ExpectType(got string, expects...string) {
	expect := strings.Join(expects, " or ")

	p.Throw(fmt.Sprintf("Expected %s, got %s.", expect, got))
}

func (p *PreparedTypeCheck) ExpectNameType(name, got string, expects...string) {
	expect := strings.Join(expects, " or ")

	p.Throw(fmt.Sprintf("Expected %s %s, got %s.", expect, name, got))
}

func (p *PreparedTypeCheck) ExpectNameTypeCause(got *types.Word, expects...string) {
	expect := strings.Join(expects, " or ")

	p.Throw(fmt.Sprintf("Expected %s cause, got %v.", expect, got))
}

func (p *PreparedTypeCheck) UnexpectedType(index int, got string) {
	names := []string{
		"first",
		"second",
		"third",
		"fourth",
		"fifth",
		"sixth",
		"seventh",
		"eighth",
		"nineth",
		"tenth",
	}

	if index < 0 {
		p.Throw(fmt.Sprintf("Unexpected type %s.", got))
	}

	p.Throw(fmt.Sprintf("Unexpected %s type %s.", names[index], got))
}

func (p *PreparedTypeCheck) ConnectedTokenError(name string, err error) {
	p.Throw(fmt.Sprintf("Connected keyword '%s' failed.", name), fmt.Sprintf("Error: %v", err))
}

func (p *PreparedTypeCheck) CallLeadingError(word *types.Word, err error) {
	p.withTrace(word, fmt.Sprintf("Error leading into %s call: %v", word.UnwrapName(), err), "Referencing:")
}

func (p *PreparedTypeCheck) CallLeavingError(name string, err error) {
	p.Throw(fmt.Sprintf("Error leaving '%s' call: %v", name, err))
}

func (p *PreparedTypeCheck) withTrace(word *types.Word, info...string) {
	result := make([]string, 0, 8)
	result = append(result, info...)
	result = append(result, trace(word)...)
	p.Throw(result...)
}

func trace(word *types.Word) []string {
	result := make([]string, 0, 8)

	result = append(result, fmt.Sprintf("at %s:%d:%d:", word.LineFile(), word.LineNumber(), word.EntryNumber()))

	snippet := word.UnwrapLine().Parent.Parent

	for snippet != nil {
		result = append(result, fmt.Sprintf("at %s:%d:", snippet.Snippet.Name, snippet.Line))
		snippet = snippet.Snippet.Parent
	}

	return result
}
