package gshell

import (
	"reflect"
)

type ScenarioPlayer struct {
}

func NewScenarioPlayer() *ScenarioPlayer {
	return &ScenarioPlayer{}
}

func (s *ScenarioPlayer) Play(client *Client, scenario interface{}, name string) {
	v := reflect.ValueOf(scenario)
	m := v.MethodByName(name)
	if m.Kind() != reflect.Func {
		panic(m.Kind())
	}
	t := m.Type()
	argv := make([]reflect.Value, t.NumIn())
	argv[0] = reflect.ValueOf(client)
	result := m.Call(argv)
	if len(result) != 0 {
		panic(result) // TODO
	}
}
