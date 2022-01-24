package console

import (
	"bufio"
	"fmt"
	"github.com/df-mc/dragonfly/server/cmd"
	"os"
	"strings"
)

// StartReader starts the CommandSender command reader. This should typically be called at the start of the program.
func StartReader() {
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for {
			scanner.Scan()
			text := strings.Split(scanner.Text(), " ")
			if command, ok := cmd.ByAlias(text[0]); ok {
				command.Execute(strings.Join(text[1:], " "), sender)
			} else {
				fmt.Println("Unknown Command")
			}
		}
	}()
}
