package main

import (
	"encoding/json"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type data struct {
	Name        string `json:"name"`
	Status      string `json:"status"`
	Topic       string `json:"topic"`
	TimeStampFE string `json:"timeStampFE"`
}

func publish(client mqtt.Client, data messageReceived) {
	num := 1
	for i := 0; i < num; i++ {
		name, status, topic, timeStampFE := string(data.Name), string(data.Status), string(data.Topic), string(data.TimeStampFE)
		//text := "Name:" + string(name) + "\n" + "Status:" + string(status) + "\n" + "Topic:" + string(topic) + "\n" + "Time:" + string(timeStampFE)
		//data := { item: "canvas", qty: 100, tags: ["cotton"], size: { h: 28, w: 35.5, uom: "cm" } }
		//token := client.Publish("topic/test/ac", 0, false, text)

		//text1 := "{\"Bird\":10,\"Cat\":\"Fuzzy\"}"
		// text := "{\"Name\":%s,\"Status\":\"%s\",\"Topic\":%s,\"TimeStampFE\":\"%s\"}", (name, status, topic, timeStampFE)
		// token := client.Publish("topic/test/ac", 0, false, text)
		token := client.Publish("topic/test/ac", 0, false, createJSON(name, status, topic, timeStampFE))
		token.Wait()
	}
}

func createJSON(name string, status string, topic string, timeStampFE string) (jsondata []byte) {
	entry := data{Name: name, Status: status, Topic: topic, TimeStampFE: timeStampFE}
	fmt.Println("entry = %+v\n", entry)

	finalBson, err := json.Marshal(entry)
	if err != nil {
		panic(err)
	}
	return finalBson
}
