package interpreter

import "fmt"

func (ip *Interpreter) load(name string) (StackValue, error) {
	var value StackValue
	value, ok := ip.memory[name]
	if !ok {
		return value, fmt.Errorf("failed to load(%s). value not found.", name)
	}

	s1, ok := value.String()
	if ok {
		ip.runtimev("loaded \"%s\"\n", s1)
		return value, nil
	}

	n1, ok := value.Int()
	if ok {
		ip.runtimev("loaded %d\n", n1)
		return value, nil
	}

	return value, fmt.Errorf("failed to load(%s). value found, but unknown type.", name)
}

func (ip *Interpreter) store(name string, value StackValue) error {
	s1, okString := value.String()
	n1, okInt := value.Int()
	if okString {
		ip.runtimev("stored \"%s\"\n", s1)
	} else if okInt {
		ip.runtimev("stored %d\n", n1)
	} else {
		return fmt.Errorf("failed to store(%s). value given, but unknown type.", name)
	}

	ip.memory[name] = value
	return nil
}
