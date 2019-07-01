package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	go nginx()
	go cybermonday()
	select {}
}

func nginx() {
	for {
		var cmd = exec.Command("nginx", "-g", "daemon off;")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		var err = cmd.Run()
		if err != nil {
			log.Panic(err)
		}
	}
}

func cybermonday() {
	for {
		var cmd = exec.Command("cybermonday")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		var err = cmd.Run()
		if err != nil {
			log.Panic(err)
		}
	}
}
