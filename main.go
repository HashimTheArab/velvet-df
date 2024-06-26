package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"
	_ "velvet/commands"
	"velvet/console"
	_ "velvet/db"
	_ "velvet/utils"
)

func main() {
	gob.Register(struct{}{})
	go func() {
		log.Println(http.ListenAndServe(":19133", nil))
	}()
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
	startServer()
}
