package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	mqttBroker = flag.String("mqttBroker", "", "MQTT broker URI (mandatory). E.g.:192.168.1.1:1883")
	topic      = flag.String("topic", "", "Topic where hub-ctrl messages will be received")
)

type Message struct {
	hub   int
	port  int
	power bool
}

func connectMQTT() (mqtt.Client, error) {
	opts := mqtt.NewClientOptions().AddBroker("tcp://" + *mqttBroker)

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("%s", token.Error())
	}

	return client, nil
}

func mqttCallback(client mqtt.Client, msq mqtt.Message) {

}

func main() {

	//Check command line parameters
	if *mqttBroker == "" || *topic == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	//Channel used to block while receiving messages
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	clientMQTT, err := connectMQTT()
	if err != nil {
		log.Fatalf("Error connecting to MQTT broker: %s", err)
	}

	if token := clientMQTT.Subscribe(*topic, 0, mqttCallback); token.Wait() && token.Error() != nil {
		log.Fatalf("Error subscribing to topic %s : %s", *topic, err)
	}

	<-c

}
