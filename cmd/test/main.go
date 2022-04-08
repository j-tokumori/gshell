package main

import (
	"github.com/j-tokumori/gshell"
	"github.com/j-tokumori/gshell/cmd/test/api"
)

func main() {
	s := gshell.New(gshell.Config{Host: "localhost:8080", IsSecure: false})
	register(s)
	s.Start()
}

// 以下手書きのシナリオ

type Scenario struct {
	e *gshell.Client
}

func (s *Scenario) Play() {
	s.e.Call(&Register{})
}

// 手書きデフォルト値
func (r *Register) Default(e *gshell.Client) *api.RegisterArgs {
	return &api.RegisterArgs{
		Country:        "jp",
		Platform:       "apple",
		PlatformUserId: "test",
		DeviceName:     "gshell",
	}
}
