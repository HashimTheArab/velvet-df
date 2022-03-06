package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"
	_ "velvet/commands"
	"velvet/console"
	_ "velvet/db"
	"velvet/dfutils"
	_ "velvet/utils"
)

func main() {
	go log.Println(http.ListenAndServe("localhost:19133", nil))
	defer func() {
		if x := recover(); x != nil {
			if err := os.Mkdir("errors", os.ModePerm); err == nil || os.IsExist(err) {
				if file, err := os.Create("errors/" + time.Now().Format("Mon_Jan_1-03.04.05-MST_2006") + ".txt"); err == nil || os.IsExist(err) {
					_, _ = file.WriteString(fmt.Sprint(x))
				} else {
					fmt.Println("Failed making file, error: " + err.Error())
				}
			} else {
				fmt.Println("Failed making errors folder: " + err.Error())
			}
			panic(x)
		}
	}()

	console.StartNew()
	dfutils.StartServer()
}
