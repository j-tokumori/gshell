package gshell

import (
	"strings"
)

type CommandFunc func(*Client, string, string)

func RPCCommand(c *Client, arg1, arg2 string) {
	parsed := parseArgs(arg2)
	c.CallByJSON(arg1, []byte(parsed))
	c.PrintResponse(arg1)
}

// parseArgs 引数のパース
// {} を付け足して、key に "" を雑につけているだけの簡易処理
// value に , や : があったり、入れ子データに対応していなかったりするので、修正必須
func parseArgs(str string) string {
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
