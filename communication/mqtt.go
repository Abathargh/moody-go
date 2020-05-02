package communication

import (
	"errors"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"moody-go/communication/handlers"
	"moody-go/models"
	"strings"
)

const (
	subTopicsCount = 2
	pubTopicsCount = 3

	greetTopic = 0
	dataTopic  = 1

	discoveryTopic        = 0
	ruleUpdateTopic       = 1
	situationForwardTopic = 2
)

var (
	client    mqtt.Client
	callbacks = []func(mqtt.Client, mqtt.Message){greetCallback, dataCallback}
)

// The function that is called whenever a MQTT message is received on the
// greet topic.
func greetCallback(_ mqtt.Client, message mqtt.Message) {
	node, err := models.NodeFromJson(message.Payload())
	if err != nil {
		log.Println(err.Error())
		return
	}
	handlers.GreetHandler(node)
}

// The function that is called whenever a MQTT message is received on the
// sensor data topic.
func dataCallback(_ mqtt.Client, message mqtt.Message) {
	data, err := models.DataFromJson(message.Payload())
	if err != nil {
		log.Println(err.Error())
		return
	}
	topicTokens := strings.Split(message.Topic(), "/")
	datatype := topicTokens[len(topicTokens)-1]
	handlers.DataHandler(datatype, data)
}

// MQTTClient implements the communication.Client interface
// and serves as the clients that receives the mqtt traffic from the
// WSAN in a similar way to BLE Centrals.
type MQTTClient struct {
	client mqtt.Client

	subTopics []string
	pubTopics []string
}

// Initializes the MQTTClient, for now we don't use a singleton in case in the future there's the need to
// use multiple clients to manage the traffic in a better way
func (c *MQTTClient) Init(host string, port int, hostId string, subTopics, pubTopics []string) error {

	if len(subTopics) != subTopicsCount || len(pubTopics) != pubTopicsCount {
		return errors.New("there's an error in your config file: wrong number of topics")
	}

	c.subTopics = subTopics
	c.pubTopics = pubTopics

	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%v:%v", host, port))
	opts.SetClientID(hostId)
	opts.OnConnect = func(client mqtt.Client) {
		subscribing := true
		for subscribing {
			for index := range subTopics {
				greetToken := client.Subscribe(subTopics[index], 0, callbacks[index])
				if greetToken.Wait() && greetToken.Error() != nil {
					continue
				}
			}
			subscribing = false
		}
	}

	c.client = mqtt.NewClient(opts)
	return nil
}

func (c *MQTTClient) Connect() error {
	connectionToken := client.Connect()

	if connectionToken.Wait() && connectionToken.Error() != nil {
		return connectionToken.Error()
	}

	return nil
}

func (c *MQTTClient) Discover() {

}
