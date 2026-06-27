package lib

import (
	"fmt"
	"strings"
)

type Plugins struct {
	Cmd      []string
	Tags     string
	IsOwner  bool
	NoPrefix bool
	Run      func(conn IClient, m M)
}
type item struct {
	Cmd string
}

var list []Plugins

func AddPlugins(plug *Plugins) {
	list = append(list, *plug)
}
func GetList() []Plugins {
	return list
}

func GetMenu() string {
	var str string
	var tags map[string][]item
	str += "*LIST COMMAND*\n\n"
	for _, cmd := range list {
		if tags == nil {
			tags = make(map[string][]item)
		}
		tg := strings.ToUpper(cmd.Tags)
		if _, ok := tags[tg]; !ok {
			tags[tg] = []item{}
		}
		tags[tg] = append(tags[tg], item{Cmd: cmd.Cmd[0]})
	}
	for key := range tags {
		count := 1
		str += fmt.Sprintf("┏━❰ *%s* ❱\n", key)
		for _, e := range tags[key] {
			str += fmt.Sprintf("┃➣ %d. %s\n", count, e.Cmd)
			count++
		}
		str += "┗━━━━━━━━━━━⦿\n"

	}
	return str
}
