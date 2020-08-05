package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/xthexder/go-jack"
)

func main() {
	mumbleHost := flag.String("mumble-host", "localhost", "the server to connect to")
	mumblePort := flag.Uint("mumble-port", 64738, "port of the mumble server")
	mumbleUsername := flag.String("mumble-user", "sip_bridge", "name the mumble client uses")
	mumbleChannel := flag.String("mumble-channel", "", "channel the client joins to")

	sipServer := flag.String("sip-host", "localhost", "sip server IP or Domain")
	sipPort := flag.Uint("sip-port", 5060, "SIP Port")
	sipUsername := flag.String("sip-username", "", "SIP Username")
	sipPassword := flag.String("sip-password", "", "SIP Password")
	sipCallNumber := flag.String("call-number", "", "Phonenumber which gets called")

	dtmfTones := flag.String("dtmf", "", "Digits to enter after call initiates")
	dtmfWait := flag.Uint("dtmf-wait", 10, "Seconds to wait until DTMF tones are played")

	// mumblePort := flag.
	flag.Parse()

	fmt.Println(mumbleHost)

	StartMumble(MumbleConfig{
		Username:    *mumbleUsername,
		Host:        *mumbleHost,
		Port:        *mumblePort,
		Channel:     *mumbleChannel,
		StoragePath: "/tmp/mumble-sip",
		JackName:    "ms1_mumble",
	})

	StartBaresip(BaresipConfig{
		Server:      *sipServer,
		Port:        *sipPort,
		Username:    *sipUsername,
		Password:    *sipPassword,
		CallNumber:  *sipCallNumber,
		StoragePath: "/tmp/mumble-sip",
		JackName:    "ms1_sip",
	})

	client, _ := jack.ClientOpen("Example Client", jack.NoStartServer)
	if client == nil {
		log.Fatal("Could not connect to jack server.")
	}
	defer client.Close()

	// lets wait 2 seconds before we try to find the mumble jack ports
	time.Sleep(2 * time.Second)

	// get mumble ports
	mumbleSource := client.GetPortByName("ms1_mumble:output_1")
	for mumbleSource == nil {
		log.Println("jack: could not get mumble source.")
		log.Println("Maybe mumble is still starting? try again in 2s")
		time.Sleep(2 * time.Second)
		mumbleSource = client.GetPortByName("ms1_mumble:output_1")
	}

	mumbleSink := client.GetPortByName("ms1_mumble:input")
	if mumbleSink == nil {
		log.Fatalln("jack: could not get mumble sink")
	}

	// get baresip ports
	baresipSource := client.GetPortByName("baresip:output_1")
	for baresipSource == nil {
		log.Fatalln("jack: could not get baresip source")
		log.Println("Maybe baresip is still starting? try again in 2s")
		time.Sleep(2 * time.Second)
		baresipSource = client.GetPortByName("baresip:output_1")
	}

	baresipSink := client.GetPortByName("baresip-01:input_1")
	if baresipSink == nil {
		log.Fatalln("jack: could not get baresip sink")
	}

	// connect ports
	client.ConnectPorts(mumbleSource, baresipSink)
	client.ConnectPorts(baresipSource, mumbleSink)

	if *dtmfTones != "" {
		time.Sleep(time.Duration(*dtmfWait) * time.Second)
		fmt.Println("send code: " + *dtmfTones)
		WriteDTMF(*dtmfTones)
	}
	select {}
}
