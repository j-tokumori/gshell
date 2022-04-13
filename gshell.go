package gshell

import (
	"fmt"
	"strings"

	"github.com/peterh/liner"
)

type Shell struct {
	Client         *Client
	Scenario       interface{}
	ScenarioPlayer *ScenarioPlayer
}

func New(cfg Config) *Shell {
	return &Shell{
		Client:         NewClient(cfg.Host, cfg.IsSecure),
		ScenarioPlayer: NewScenarioPlayer(),
	}
}

func (s *Shell) Start() {
	s.bootstrap()

	line := liner.NewLiner()
	defer line.Close()

	// TODO: wait for signal and kill
	//sigs := make(chan os.Signal, 1)
	//signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	//go func() {
	//	sig := <-sigs
	//	fmt.Printf("sig: %d, exit...", sig)
	//	os.Exit(0)
	//}()

	for {
		l, err := line.Prompt("gshell> ")
		if err != nil {
			fmt.Println("error: ", err)
			continue
		}

		if s.exec(s.Client, l) {
			break
		}

		if strings.TrimSpace(l) != "" {
			line.AppendHistory(l)
		}
	}

	fmt.Println("exit...")
}

func (s *Shell) RegisterRPC(name string, f NewRPCFunc) {
	s.Client.rpcMap[name] = f
}

func (s *Shell) RegisterContext(f ContextFunc) {
	s.Client.ContextFunc = f
}

func (s *Shell) RegisterScenario(scenario interface{}) {
	s.Scenario = scenario
}

func (s *Shell) bootstrap() {
	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Println(rec)
			fmt.Println("failed bootstrap.")
		}
	}()

	s.ScenarioPlayer.Play(s.Client, s.Scenario, "Boot")
}

func (s *Shell) exec(c *Client, line string) bool {
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
		c.CallByJSON(second, []byte(args))
		c.PrintResponse(second)
	case "scenario", "s":
		//sce.Call(second)
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
