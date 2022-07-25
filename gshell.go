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
	Client   *Client
	Commands map[string]Command
}

func New(host string, opts ...Option) *Shell {
	options := &options{}
	for _, opt := range opts {
		opt.apply(options)
	}

	if options.rpcAliasHandler == nil {
		options.rpcAliasHandler = defaultRPCAliasHandler
	}

	s := &Shell{
		Client:   NewClient(host, options),
		Commands: make(map[string]Command, 0),
	}
	s.RegisterCommand([]string{"rpc", "r", "call"}, NewRPCCommand(options.rpcMap, options.rpcAliasHandler))
	if options.scenarioFactory != nil {
		scenario := options.scenarioFactory(s.Client)
		s.RegisterCommand([]string{"scenario", "s"}, &ScenarioCommand{NewScenarioPlayer(scenario)})
	}
	s.RegisterCommand([]string{"response"}, &ResponseCommand{})
	s.RegisterCommand([]string{"reply"}, &ReplyCommand{})
	s.RegisterCommand([]string{"Header"}, &HeaderCommand{})
	s.RegisterCommand([]string{"Trailer"}, &TrailerCommand{})
	s.RegisterCommand([]string{"sample"}, &SampleCommand{})
	s.RegisterCommand([]string{"list"}, &ListCommand{})
	s.RegisterCommand([]string{"trace"}, &TraceCommand{})
	s.RegisterCommand([]string{""}, &EmptyCommand{})
	s.RegisterCommand([]string{"help"}, &HelpCommand{})
	s.RegisterCommand([]string{"exit", "quit"}, &ExitCommand{})
	return s
}

func (s *Shell) Start() {
	s.bootstrap()

	line := liner.NewLiner()
	defer func() {
		if err := line.Close(); err != nil {
			panic(err)
		}
	}()

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
}

func (s *Shell) RegisterCommand(keys []string, cmd Command) {
	for _, key := range keys {
		s.Commands[key] = cmd
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

	s.scenario("Boot")
}

func (s *Shell) scenario(name string) {
	if cmd, ok := s.Commands["scenario"]; ok {
		cmd.Exec(s.Client, name)
	}
}

func (s *Shell) exec(c *Client, line string) bool {
	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Println(rec)
		}
	}()

	first, second, third := s.parse(line)

	if cmd, ok := s.Commands[first]; ok {
		return cmd.Exec(c, second, third)
	}

	fmt.Println("this command is not support.")

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
