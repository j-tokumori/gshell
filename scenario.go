package gshell

import (
	"reflect"
)

// ScenarioPlayer TODO: Client に Play 機能持たせていいか..
type ScenarioPlayer struct {
	Scenario interface{}
}

func NewScenarioPlayer(scenario interface{}) *ScenarioPlayer {
	return &ScenarioPlayer{scenario}
}

func (s *ScenarioPlayer) Play(client *Client, name string) {
	// TODO: scenario nil チェック
	v := reflect.ValueOf(s.Scenario)
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
