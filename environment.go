package main

import "fmt"

// Environment holds variable bindings
type Environment struct {
	values map[string]interface{}
	parent *Environment // For nested scopes later
}

func NewEnvironment() *Environment {
	return &Environment{
		values: make(map[string]interface{}),
		parent: nil,
	}
}

func NewEnclosedEnvironment(parent *Environment) *Environment {
	return &Environment{
		values: make(map[string]interface{}),
		parent: parent,
	}
}

func (e *Environment) Define(name string, value interface{}) {
	// we are not just defining a new variable but also redefining the variable when it is already defined
	e.values[name] = value
}

func (e *Environment) Get(name string) (interface{}, error) {
	if value, ok := e.values[name]; ok {
		return value, nil
	}

	if e.parent != nil {
		return e.parent.Get(name)
	}

	return nil, fmt.Errorf("undefined variable '%s'", name)
}

func (e *Environment) Assign(name string, value interface{}) error {
	if _, ok := e.values[name]; ok {
		e.values[name] = value
		return nil
	}

	if e.parent != nil {
		return e.parent.Assign(name, value)
	}

	return fmt.Errorf("undefined variable '%s'", name)
}
