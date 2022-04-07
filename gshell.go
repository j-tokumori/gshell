package gshell

import (
	"fmt"
	"strings"

	"github.com/peterh/liner"
)

type Shell struct {
}

func New() *Shell {
	return &Shell{}
}

func (s *Shell) Start(cfg Config) {
	cli := NewClient("localhost:8080", false, cfg)
	sce := NewScenario(cli)

	//bootstrap(c, s)

	line := liner.NewLiner()
	defer line.Close()

	for {
		l, err := line.Prompt("gshell> ")
		if err != nil {
			fmt.Println("error: ", err)
			continue
		}

		if s.exec(cli, sce, l) {
			break
		}

		if strings.TrimSpace(l) != "" {
			line.AppendHistory(l)
		}
	}

	fmt.Println("exit...")
}

func (s *Shell) Register(rpc RPC) {
	// TODO:
}

func (s *Shell) bootstrap(cli *Client, sce *Scenario) {
	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Println(rec)
			fmt.Println("failed bootstrap.")
		}
	}()

	//sce.Boot()
}

func (s *Shell) exec(c *Client, sce *Scenario, line string) bool {
	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Println(rec)
		}
	}()

	fist, second, third := s.parse(line)
	switch fist {
	case "rpc", "r":
		args := s.parseArgs(third)
		c.CallWithJSON(second, []byte(args))
		c.PrintResponse(second)
	case "scenario", "s":
		sce.Call(second)
	case "response":
		if second == "" {
			c.PrintLastResponse()
		} else {
			c.PrintResponse(second)
		}
	case "reply":
		if second == "" {
			c.PrintLastReply()
		} else {
			c.PrintReply(second)
		}
	case "header":
		if second == "" {
			c.PrintLastHeader()
		} else {
			c.PrintHeader(second)
		}
	case "trailer":
		if second == "" {
			c.PrintLastTrailer()
		} else {
			c.PrintTrailer(second)
		}
	case "trace":
		if second == "" {
			c.PrintLastTraceURL()
		} else {
			c.PrintTraceURL(second)
		}
	case "sample":
		c.PrintSample(second)
	case "list":
		if second == "rpc" {
			c.PrintList(third)
		}
	case "":
	case "exit", "quit":
		return true
		//case "help":
		//default:
		//	help?
	default:
		s.help()
	}
	return false
}

func (s *Shell) help() {
	println("this command is not support.")
	println("todo: impl help... help you coming soon.")
}

func (s *Shell) parse(cmd string) (string, string, string) {
	arr := strings.Split(strings.TrimSpace(cmd), " ")
	switch len(arr) {
	case 0:
		return "", "", ""
	case 1:
		return arr[0], "", ""
	case 2:
		return arr[0], arr[1], ""
	default:
		return arr[0], arr[1], strings.Join(arr[2:], " ")
	}
}

// parseArgs 引数のパース
// {} を付け足して、key に "" を雑につけているだけの簡易処理
// value に , や : があったり、入れ子データに対応していなかったりするので、修正必須
func (s *Shell) parseArgs(str string) string {
	str = strings.TrimSpace(str)
	if str == "" {
		return "{}"
	}
	if str[0:1] == "{" { // 先頭が { なら生 json とみなす
		return str
	}
	args := strings.Split(str, ",")
	modArgs := make([]string, 0)
	for _, arg := range args {
		kv := strings.Split(arg, ":")
		if len(kv) != 2 {
			panic("TODO")
		}
		kv[0] = `"` + kv[0] + `"`
		modArgs = append(modArgs, strings.Join(kv[0:2], ":"))
	}
	return "{" + strings.Join(modArgs, ",") + "}"
}
