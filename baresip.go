package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"text/template"
)

var stdin io.WriteCloser

func WriteDTMF(chars string) {
	io.WriteString(stdin, chars)
}

func writeBaresipConfigFile(config BaresipConfig) {
	configTemplate, err := template.ParseFiles("templates/baresip_config")
	if err != nil {
		log.Fatalf("loading template failed with %s", err)
	}
	fmt.Println(configTemplate)
	f, err := os.Create(path.Join(config.StoragePath, "config"))
	if err != nil {
		log.Fatalf("Error while creating file: %s", err)
	}
	err = configTemplate.Execute(f, config)
	if err != nil {
		log.Print("execute: ", err)
		return
	}
	f.Close()
}

type BaresipConfig struct {
	Server      string
	Port        uint
	Username    string
	Password    string
	CallNumber  string
	JackName    string
	StoragePath string
}

func StartBaresip(conf BaresipConfig) *exec.Cmd {
	writeBaresipConfigFile(conf)

	cmd := exec.Command("baresip",
		"-f", conf.StoragePath,
		// TODO: escape strings
		"-e", "/uanew <sip:"+conf.Username+"@"+conf.Server+":"+string(conf.Port)+">;auth_pass="+conf.Password+";regint=0",
		"-e", "/dial "+conf.CallNumber)

	stdin2, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalf("cmd.StdinPipe() failed with %s", err)
	}
	stdin = stdin2
	log.Println(stdin)

	AttachLinePrefixer("baresip", cmd)
	PrintCommand(cmd)

	err2 := cmd.Start()
	if err2 != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err2)
	}

	return cmd
}
