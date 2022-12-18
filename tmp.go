package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	//"strings"
)

func executeCommand(name string, args []string) {
	cmd := exec.Command(name, args...)
	//cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("in all caps: %q\n", out.String())
}

func main() {
	command := "uname"
	args := []string{"-a"}
	executeCommand(command, args)
}
