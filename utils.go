package main

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func PrintCommand(cmd *exec.Cmd) {
	fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}

func AttachLinePrefixer(prefix string, cmd *exec.Cmd) {

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	outbr := bufio.NewReader(stdout)
	errbr := bufio.NewReader(stderr)

	printPrefixedLines := func(what string, br *bufio.Reader) {
		for {
			line, _, err := br.ReadLine()
			if err != nil {
				log.Fatalf("Mumble read error: %s\n", err)
				return
			}
			fmt.Println("[" + what + "] " + string(line))
		}
	}

	go printPrefixedLines(prefix, outbr)
	go printPrefixedLines(prefix, errbr)
}
