package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func publish(client mqtt.Client, data Message) {
	num := 1
	for i := 0; i < num; i++ {
		//text := bson.D{{Key: string(data.Name), Value: string(data.Msg)}}
		//token := client.Publish("topic/test/ac", 0, false, text)
		Name, Msg := string(data.Name), string(data.Msg)
		text := "Name:" + string(Name) + "\n" + "Msg:" + string(Msg)
		token := client.Publish("topic/test/ac", 0, false, text)
		token.Wait()
		//time.Sleep(time.Second)
	}
}
