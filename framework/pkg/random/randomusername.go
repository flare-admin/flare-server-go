package random

import (
	"math/rand"
	"strings"
)

// 一些预定义的前缀、中间名和后缀
var prefixes = []string{"user", "guest", "member", "player"}
var names = []string{"alpha", "beta", "gamma", "delta", "epsilon", "zeta", "eta", "theta"}
var suffixes = []string{"01", "99", "123", "xyz", "abc", "007", "789"}

// GenerateName 生成一个随机用户名
func GenerateName() string {
	prefix := prefixes[rand.Intn(len(prefixes))]
	name := names[rand.Intn(len(names))]
	suffix := suffixes[rand.Intn(len(suffixes))]
	return strings.Join([]string{prefix, name, suffix}, "_")
}
