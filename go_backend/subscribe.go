package main

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func sub(client mqtt.Client) {
	topic := "topic/test/ac"
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("\nSubscribed to topic: %s", topic)
}

func subdb(client mqtt.Client) {
	topic := "topic/test/+"
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("\ndb client Subscribed to topic: %s", topic)
}
