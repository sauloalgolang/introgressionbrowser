package main

import (
	"fmt"
)

import (
	"github.com/sauloalgolang/introgressionbrowser/web"
)

type WebCommand struct {
	Host string `long:"host" description:"Hostname" default:"127.0.0.1"`
	Port int    `long:"port" description:"Port" default:"8000"`
	Dir  string `long:"dir" description:"Directory to be served" default:"res/"`
}

func (w WebCommand) String() (res string) {
	res += fmt.Sprintf(" Host                   : %s\n", w.Host)
	res += fmt.Sprintf(" Port                   : %d\n", w.Port)
	res += fmt.Sprintf(" Dir                    : %s\n", w.Dir)
	return res
}

var webCommand WebCommand

func (x *WebCommand) Execute(args []string) error {
	fmt.Printf("Web\n")
	fmt.Println(x)

	web.NewWeb(x.Dir, x.Host, x.Port)

	return nil
}

func init() {
	parser.AddCommand(
		"web",
		"Start web interface",
		"Start web interface",
		&webCommand,
	)
}
