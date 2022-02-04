package console

import (
	"bufio"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

type Reader struct {
	sender *CommandSender
}

func StartNew() {
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{ForceColors: true}
	log.Level = logrus.InfoLevel
	c := Reader{sender: &CommandSender{log}}
	c.Start()
}

// Start starts the CommandSender command reader. This should typically be called at the start of the program.
func (r *Reader) Start() {
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for {
			scanner.Scan()
			if text := strings.ToValidUTF8(strings.TrimSpace(scanner.Text()), ""); len(text) != 0 {
				args := strings.Split(text, " ")
				name := args[0]
				if command, ok := cmd.ByAlias(name); ok {
					command.Execute(strings.Join(args[1:], " "), r.sender)
				} else {
					output := &cmd.Output{}
					output.Errorf("Unknown command: %v", name)
					r.sender.SendCommandOutput(output)
				}
			}
		}
	}()
}
