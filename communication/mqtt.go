package communication

import (
	"fmt"
	"github.com/Abathargh/moody-go/communication/handlers"
	"github.com/Abathargh/moody-go/models"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/validator.v2"
	"log"
	"strings"
)

const (
	clientId = "Moody-base"

	subTopicsCount = 2
	pubTopicsCount = 3

	greetTopic = 0
	dataTopic  = 1

	discoveryTopic        = 0
	ruleUpdateTopic       = 1
	situationForwardTopic = 2

	quiesce = 200 // Client disconect quiescence
)

var (
	callbacks = []func(mqtt.Client, mqtt.Message){greetCallback, dataCallback}
)

// MQTTClient implements the communication.Client ui
// and serves as the clients that receives the mqtt traffic from the
// WSAN in a similar way to BLE Centrals.
type MQTTClient struct {
	client mqtt.Client
	config MQTTConfig
}

type MQTTConfig struct {
	Host      string   `validate:"nonzero"`
	Port      int      `validate:"nonzero,min=1,max=65536"`
	SubTopics []string `validate:"nonzero,len=2"` // 2 sub topic defined in the standard mqtt implementatoin
	PubTopics []string `validate:"nonzero,len=3"` // 3 pub topic defined in the standard mqtt implementatoin
}

// Initializes the MQTTClient, for now we don't use a singleton in case in the future there's the need to
// use multiple clients to manage the traffic in a better way
func (c *MQTTClient) Init(conf interface{}) error {
	if err := mapstructure.Decode(conf, &(c.config)); err != nil {
		log.Println("wrong format for the mqtt section of the config file")
		return err
	}

	if err := validator.Validate(c.config); err != nil {
		log.Println("wrong values for one or more fields in the mqtt section of the config file")
		return err
	}

	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%v:%v", c.config.Host, c.config.Port))
	opts.SetClientID(clientId)
	opts.OnConnect = func(client mqtt.Client) {
		subscribing := true
		for subscribing {
			for index := range c.config.SubTopics {
				greetToken := client.Subscribe(c.config.SubTopics[index], 0, callbacks[index])
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
	connectionToken := c.client.Connect()

	if connectionToken.Wait() && connectionToken.Error() != nil {
		return connectionToken.Error()
	}
	return nil
}

func (c *MQTTClient) Discover() error {
	token := c.client.Publish(c.config.PubTopics[discoveryTopic], 0, false, clientId)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (c *MQTTClient) Forward(situation string) error {
	token := c.client.Publish(c.config.PubTopics[situationForwardTopic], 0, false, situation)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (c *MQTTClient) Update(group, rule string) error {
	topic := fmt.Sprintf("%s/%s", c.config.PubTopics[ruleUpdateTopic], group)
	token := c.client.Publish(topic, 0, false, rule)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (c *MQTTClient) Close() {
	c.client.Disconnect(quiesce)
}

// The function that is called whenever a MQTT message is received on the
// greet topic.
func greetCallback(_ mqtt.Client, message mqtt.Message) {
	node, err := models.NodeFromJson(message.Payload())
	if err != nil {
		log.Println("an error occurred while unmarshalling a greet packet")
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
		log.Println("an error occurred while unmarshalling a data packet")
		log.Println(err.Error())
		return
	}
	fmt.Println(data)
	topicTokens := strings.Split(message.Topic(), "/")
	datatype := topicTokens[len(topicTokens)-1]
	handlers.DataHandler(datatype, data)
}
