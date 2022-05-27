package gshell

import (
	"reflect"
	"strconv"
)

type ScenarioPlayer struct {
	Scenario interface{}
}

type ScenarioFactory func(*Client) interface{}

func NewScenarioPlayer(scenario interface{}) *ScenarioPlayer {
	return &ScenarioPlayer{Scenario: scenario}
}

func (s *ScenarioPlayer) Play(_ *Client, name string, args ...string) {
	v := reflect.ValueOf(s.Scenario)
	m := v.MethodByName(name)
	if m.Kind() != reflect.Func {
		panic(m.Kind())
	}
	t := m.Type()
	argv := make([]reflect.Value, t.NumIn())

	// 引数を埋める
	// int, string のみ対応、埋めれない場合 zero値で埋める
	for i := 0; i < t.NumIn(); i++ {
		it := t.In(i)
		switch it.Kind() {
		case reflect.Int:
			var ii int
			if len(args) > i {
				ii, _ = strconv.Atoi(args[i])
			}
			argv[i] = reflect.ValueOf(ii)
		case reflect.String:
			var s string
			if len(args) > i {
				s = args[i]
			}
			argv[i] = reflect.ValueOf(s)
		default:
			panic("unexpected scenario args type.")
		}
	}

	result := m.Call(argv)
	if len(result) != 0 {
		panic(result) // TODO
	}
}
