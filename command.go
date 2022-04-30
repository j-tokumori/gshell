package gshell

import (
	"strings"
)

type Command interface {
	Exec(*Client, ...string) bool
}

type RPCCommand struct {
}

func (c *RPCCommand) Exec(client *Client, args ...string) bool {
	parsed := c.parseArgs(args[1])
	client.CallByJSON(args[0], []byte(parsed))
	client.PrintResponse(args[0])
	return false
}

// parseArgs 引数のパース
// {} を付け足して、key に "" を雑につけているだけの簡易処理
// value に , や : があったり、入れ子データに対応していなかったりするので、修正必須
func (c *RPCCommand) parseArgs(str string) string {
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

type ScenarioCommand struct {
	Scenario       interface{}
	ScenarioPlayer *ScenarioPlayer
}

func (c *ScenarioCommand) Exec(client *Client, args ...string) bool {
	c.ScenarioPlayer.Play(client, c.Scenario, args[0])
	return false
}
