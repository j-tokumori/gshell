package gshell

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/peterh/liner"
)

type Shell struct {
	Client         *Client
	Commands       map[string]CommandFunc
	Scenario       interface{}
	ScenarioPlayer *ScenarioPlayer
}

type Config struct {
	Host     string
	IsSecure bool

	ContextFunc ContextFunc
	ErrorFunc   ErrorFunc

	Scenario interface{}
}

func New(cfg Config) *Shell {
	return &Shell{
		Client:         NewClient(cfg.Host, cfg.IsSecure, cfg.ContextFunc, cfg.ErrorFunc),
		Commands:       make(map[string]CommandFunc, 0),
		Scenario:       cfg.Scenario,
		ScenarioPlayer: NewScenarioPlayer(),
	}
}

func (s *Shell) Start() {
	s.RegisterCommand(RPCCommand, "rpc", "r")

	s.bootstrap()

	line := liner.NewLiner()
	defer line.Close()

	// TODO: wait for signal and kill
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Printf("\nsig: %d, exit...", sig)
		os.Exit(0)
	}()

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

func (s *Shell) RegisterCommand(f CommandFunc, keys ...string) {
	for _, key := range keys {
		s.Commands[key] = f
	}
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

	first, second, third := s.parse(line)

	if f, ok := s.Commands[first]; ok {
		f(c, second, third)
		return false
	}

	switch first {
	//case "rpc", "r":
	//	args := s.parseArgs(third)
	//	c.CallByJSON(second, []byte(args))
	//	c.PrintResponse(second)
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
