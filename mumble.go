package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path"
	"text/template"
)

type MumbleConfig struct {
	Username    string
	Host        string
	Port        uint
	Channel     string
	JackName    string
	StoragePath string
}

func writeMumbleConfigFile(config MumbleConfig) {
	configTemplate, err := template.ParseFiles("templates/Mumble.conf")
	if err != nil {
		log.Fatalf("loading template failed with %s", err)
	}
	fmt.Println(configTemplate)

	err = os.MkdirAll(path.Join(config.StoragePath, "Mumble"), 0700)
	if err != nil {
		log.Fatal(err)
	}

	f, err2 := os.Create(path.Join(config.StoragePath, "Mumble/Mumble.conf"))
	if err2 != nil {
		log.Fatalf("Error while creating file: %s", err)
	}
	err = configTemplate.Execute(f, config)
	if err != nil {
		log.Print("execute: ", err)
		return
	}
	f.Close()
}

// RunMumble :
func StartMumble(config MumbleConfig) *exec.Cmd {

	u := &url.URL{
		Scheme: "mumble",
		User:   url.User(config.Username),
		Host:   config.Host,
		Path:   config.Channel,
	}

	writeMumbleConfigFile(config)

	cmd := exec.Command("mumble", "-m", "-jn", config.JackName, u.String())
	// cmd.Stderr = os.Stderr

	cmd.Env = append(os.Environ(),
		// hides the window
		"QT_QPA_PLATFORM=vnc",

		// overwrite config & data directories
		"XDG_CONFIG_HOME="+config.StoragePath,
		"XDG_DATA_HOME="+config.StoragePath,
	)

	AttachLinePrefixer("mumble", cmd)
	PrintCommand(cmd)

	err2 := cmd.Start()
	if err2 != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err2)
	}

	return cmd
}
