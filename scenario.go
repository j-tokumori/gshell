package gshell

import (
	"reflect"
)

type Scenario struct {
	e *Client
}

func NewScenario(e *Client) *Scenario {
	return &Scenario{e: e}
}

func (s *Scenario) Call(scenario string) {
	v := reflect.ValueOf(s)
	m := v.MethodByName(scenario)
	if m.Kind() != reflect.Func {
		panic(m.Kind())
	}
	t := m.Type()
	argv := make([]reflect.Value, t.NumIn())
	result := m.Call(argv)
	if len(result) != 0 {
		panic(result) // TODO
	}
}
