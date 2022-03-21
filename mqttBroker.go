package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbclient *mongo.Client
var err error
var clients []mqtt.Client

type Message struct {
	Name  string
	Topic string
	Msg   string
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("\nReceived message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var messagePubHandlerForDB mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Println("\ndb handler called")
	//call for db storage of messages
	dbVal := bson.D{{"msg", string(msg.Payload())}, {"topic", string(msg.Topic())}}
	insertIntoDB(dbclient, dbVal)
	fmt.Printf("\n inserted into db: %s\n", dbVal)
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("\nConnect lost: %v", err)
}

func main() {
	//connection with database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dbclient, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://cherindranath:cherry2580@cluster0.s1pf5.mongodb.net/publish_msg.msg?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = dbclient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	//create n devices
	no_devices := 3
	clients = create_n_devices(no_devices)
	//create a devices wildcard that listens to all msgs and store them in Database
	clientdb := create_db_device()

	sub(clients[1])
	sub(clients[2])
	subdb(clientdb)
	// publish(clients[0], "hi")

	//backend router
	router := mux.NewRouter()
	router.HandleFunc("/api/publish", publishData)
	handler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe(":8000", handler))

	defer func() {
		for i := 0; i < no_devices; i++ {
			clients[i].Disconnect(250)
		}
	}()
}

func publishData(writer http.ResponseWriter, request *http.Request) {
	var jsonData Message
	err := json.NewDecoder(request.Body).Decode(&jsonData)
	if err != nil {
		fmt.Println(err.Error() + "---")
		return
	}
	fmt.Println("")
	fmt.Println(jsonData)
	publish(clients[0], jsonData)
}

func create_db_device() (device mqtt.Client) {
	var broker = "broker.emqx.io"
	var port = 1883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("go_mqtt_client" + "db")
	opts.SetUsername("emqx" + "db")
	opts.SetPassword("public" + "db")
	opts.SetDefaultPublishHandler(messagePubHandlerForDB)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	clientdb := mqtt.NewClient(opts)
	if token := clientdb.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return clientdb
}

func create_n_devices(n int) (clientsSlice []mqtt.Client) {
	var broker = "broker.emqx.io"
	var port = 1883

	no_devices := n
	var clients []mqtt.Client
	for i := 0; i < no_devices; i++ {
		opts := mqtt.NewClientOptions()
		opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
		opts.SetClientID("go_mqtt_client" + strconv.Itoa(i))
		opts.SetUsername("emqx" + strconv.Itoa(i))
		opts.SetPassword("public" + strconv.Itoa(i))
		opts.SetDefaultPublishHandler(messagePubHandler)
		opts.OnConnect = connectHandler
		opts.OnConnectionLost = connectLostHandler

		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		clients = append(clients, client)
	}
	return clients
}

func insertIntoDB(dbclient *mongo.Client, data bson.D) {
	collection := dbclient.Database("publish_msg").Collection("msg")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		log.Fatal(err)
	}
	id := res.InsertedID
	fmt.Println(id)
}
