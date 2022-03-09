package cli

import (
	"bufio"
	. "file-server/structs"
	. "file-server/utils"
	"fmt"
	"os"
	"strings"
)

type SubjectCli interface {
	Register(observer CliObserver)
	Unregister(observer CliObserver)
	NotifyAll(options ...string)
}

type Cli struct {
	Observers []CliObserver
	Source    string
	Active    chan bool
}

func NewCli(source string) *Cli {
	return &Cli{
		Source: source,
		Active: make(chan bool),
	}
}

func (cli *Cli) Start() {
	for {
		cli.listenInput()
	}
}

func (cli *Cli) listenInput() {
	for {
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
		cli.handleInput(input)
	}
}

/*
subscribe channel:name
unsubscribe channel:name
send channel:name file:path

start server
stop server
*/

func (cli *Cli) handleInput(input string) {
	input = SingleSpacePattern.ReplaceAllString(strings.TrimSpace(input), " ")

	if cli.Source == "server" {
		switch input {
		case START:
			cli.NotifyAll(START)
		case STOP:
			cli.NotifyAll(STOP)
			cli.Active <- false
		case HELP:
			cli.help()
		default:
			cli.invalid()
		}
		return
	}

	options := strings.Split(input, " ")
	if len(options) < 1 {
		return
	}
	action := strings.ToLower(options[0])

	if RegexSubscribe.MatchString(input) || RegexUnsubscribe.MatchString(input) {
		cli.NotifyAll(action, cli.value(options[1]))
		return
	}

	if RegexSend.MatchString(input) {
		cli.NotifyAll(action, cli.value(options[1]), cli.value(options[2]))
		return
	}

	switch action {
	case EXIT:
		cli.Active <- false
	case HELP:
		cli.help()
	default:
		cli.invalid()
	}

}

func (cli *Cli) value(param string) string {
	data := strings.SplitN(param, ":", 2)
	return data[1]
}

func (cli *Cli) help() {
	PrintHelp("----subscribe channel:name----",
		"----unsubscribe channel:name----",
		"----send channel:name file:path----")
}

func (cli *Cli) invalid() {
	PrintError("Invalid parameters", "Run command help")
}

func (cli *Cli) Register(observer CliObserver) {
	cli.Observers = append(cli.Observers, observer)
}

func (cli *Cli) Unregister(observer CliObserver) {
	for i, obs := range cli.Observers {
		if obs.Identifier() == observer.Identifier() {
			cli.Observers = append(cli.Observers[:i], cli.Observers[i+1:]...)
		}
	}
}

func (cli *Cli) NotifyAll(options ...string) {
	for _, obs := range cli.Observers {
		if obs.Identifier() == cli.Source {
			obs.OnEntry(options)
		}
	}
}
