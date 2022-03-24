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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbclient *mongo.Client
var err error
var clients []mqtt.Client
var responseDataId string

type MessageTopic struct {
	Topic string
}

type messageReceived struct {
	Name        string
	Status      string
	Topic       string
	TimeStampFE string
}

type messageSent struct {
	Id          string
	Name        string
	Status      string
	Topic       string
	TimeStampFE string
	TimeStampBE string
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("\nReceived message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var messagePubHandlerForDB mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Println("\ndb handler called")

	// fmt.Println(msg.Payload())
	// fmt.Println(string(msg.Payload()))

	var finaljson data
	err := json.Unmarshal(msg.Payload(), &finaljson)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("finaljson=")
	fmt.Println(finaljson)

	doc := bson.D{{"Name", finaljson.Name}, {"Status", finaljson.Status}, {"Topic", finaljson.Topic}, {"TimeStampFE", finaljson.TimeStampFE}}

	dbObjId := insertIntoDB(dbclient, doc)
	responseDataId = dbObjId
	fmt.Printf("\n inserted into db: %s\n", doc)

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
	dbclient, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://cherindranath:cherry2580@cluster0.s1pf5.mongodb.net/publish_msg1.msg1?retryWrites=true&w=majority"))
	if err != nil {
		fmt.Println("this is connecting to mongodb error")
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
	router.HandleFunc("/api/fetch", getTopicHistory)
	handler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe(":8000", handler))

	defer func() {
		for i := 0; i < no_devices; i++ {
			clients[i].Disconnect(250)
		}
	}()
}

func getTopicHistory(writer http.ResponseWriter, request *http.Request) {
	var topicName MessageTopic
	err := json.NewDecoder(request.Body).Decode(&topicName)
	if err != nil {
		fmt.Println("---" + err.Error() + "---")
		return
	}
	fmt.Print("\ntopic Name=")
	fmt.Println(topicName)

	docs := get_topic_history(dbclient, topicName.Topic)
	fmt.Println(docs[0])

	type msgTopicSend struct {
		Topic   string
		Updates []primitive.M
	}

	var msgTopicSend1 msgTopicSend
	msgTopicSend1.Topic = topicName.Topic
	msgTopicSend1.Updates = docs

	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(msgTopicSend1)

	return
}

func publishData(writer http.ResponseWriter, request *http.Request) {
	var jsonData messageReceived
	err := json.NewDecoder(request.Body).Decode(&jsonData)
	if err != nil {
		fmt.Println(err.Error() + "---")
		return
	}
	fmt.Println("")
	fmt.Println(jsonData)
	publish(clients[0], jsonData)

	var senddata messageSent
	senddata.Id = responseDataId
	senddata.Name = jsonData.Name
	senddata.Status = jsonData.Status
	senddata.Topic = jsonData.Topic
	senddata.TimeStampFE = jsonData.TimeStampFE

	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(senddata)
	return
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

func insertIntoDB(dbclient *mongo.Client, data bson.D) (dbid string) {
	collection := dbclient.Database("publish_msg").Collection("msg")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		log.Fatal(err)
	}
	id := res.InsertedID
	fmt.Println(id)
	stringObjectID := id.(primitive.ObjectID).Hex()
	return stringObjectID
}

func get_topic_history(dbclient *mongo.Client, topic string) (docs []bson.M) {
	collection := dbclient.Database("publish_msg").Collection("msg")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filterCursor, err := collection.Find(ctx, bson.M{"Topic": topic})
	if err != nil {
		log.Fatal(err)
	}
	var updates []bson.M
	err = filterCursor.All(ctx, &updates)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(updates)

	return updates
}
