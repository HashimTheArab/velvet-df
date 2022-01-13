package main

import (
	_ "velvet/commands"
	_ "velvet/db"
	"velvet/dfutils"
	_ "velvet/utils"
)

func main() {
	dfutils.StartServer()
}
