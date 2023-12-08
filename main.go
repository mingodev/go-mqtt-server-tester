package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
)

var mqttMessagePublishedHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var mqttConnectionHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connection to MQTT server successful.")
}

var mqttConnectionLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection to MQTT server lost: %v", err)
}

func publishMqttTestMessages(client mqtt.Client, topic string, amount int, delay time.Duration) {
	for i := 0; i < amount; i++ {
		text := fmt.Sprintf("Message %d", i+1)
		token := client.Publish(topic, 0, false, text)
		token.Wait()
		time.Sleep(delay)
	}
}

func subscribeToMqttTopic(client mqtt.Client, topic string) {
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s\n", topic)
}

func main() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal("Error: .env file not found. Aborting test.")
	}

	broker := os.Getenv("MQTT_SERVER_ADDRESS")
	port, _ := strconv.Atoi(os.Getenv("MQTT_SERVER_PORT"))

	fmt.Println(port)

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID(os.Getenv("MQTT_CLIENT_ID"))
	opts.SetUsername(os.Getenv("MQTT_CLIENT_USER"))
	opts.SetPassword(os.Getenv("MQTT_CLIENT_PASSWORD"))
	opts.SetDefaultPublishHandler(mqttMessagePublishedHandler)
	opts.OnConnect = mqttConnectionHandler
	opts.OnConnectionLost = mqttConnectionLostHandler

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Subscribe to topic
	topic := os.Getenv("MQTT_SERVER_TOPIC")
	subscribeToMqttTopic(client, topic)

	// Send test messages to topic
	publishMqttTestMessages(client, topic, 10, time.Duration(3*time.Second))

	client.Disconnect(250)
}
