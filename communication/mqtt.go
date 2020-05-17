package communication

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/validator.v2"
	"log"
	"strings"
	"time"
)

const (
	clientId = "Moody-base"

	discoveryTopic        = 0
	ruleUpdateTopic       = 1
	situationForwardTopic = 2

	quiesce     = 200 // Client disconnect quiescence
	pingTimeout = 2 * time.Second
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
	DataTopic string   `validate:"nonzero"`       // 2 sub topic defined in the standard mqtt implementation
	PubTopics []string `validate:"nonzero,len=3"` // 3 pub topic defined in the standard mqtt implementation
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

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%v:%v", c.config.Host, c.config.Port))
	opts.SetPingTimeout(pingTimeout)
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		subscribing := true
		for subscribing {
			dataToken := client.Subscribe(c.config.DataTopic, 0, dataCallback)
			if dataToken.Wait() && dataToken.Error() != nil {
				continue
			}
			subscribing = false
		}
	})

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
	log.Println("MQTT Client - Shutting down")
	c.client.Unsubscribe(c.config.DataTopic)
	c.client.Disconnect(quiesce)
}

// The function that is called whenever a MQTT message is received on the
// sensor data topic.
func dataCallback(_ mqtt.Client, message mqtt.Message) {
	data := string(message.Payload())
	topicTokens := strings.Split(message.Topic(), "/")
	datatype := topicTokens[len(topicTokens)-1]
	DataHandler(datatype, data)
}
