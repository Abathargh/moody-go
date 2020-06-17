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
	serviceSubTopic       = 0
	actServerTopic        = 1
	modeSwitchTopic       = 0
	situationForwardTopic = 1

	//quiesce     = 50 // Client disconnect quiescence
	pingTimeout = 2 * time.Second
)

// MQTTClient implements the communication.Client ui
// and serves as the clients that receives the mqtt traffic from the
// WSAN in a similar way to BLE Centrals.
type MQTTClient struct {
	client mqtt.Client
	config MQTTConfig
}

// MQTTConfig stores the information contained in the mqtt section of the conf.json
type MQTTConfig struct {
	Host      string   `validate:"nonzero" json:"host"`
	Port      int      `validate:"nonzero,min=1,max=65536" json:"port"`
	DataTopic []string `validate:"nonzero,len=2" json:"dataTopic"` // 2 sub topic in the standard mqtt implementation
	PubTopics []string `validate:"nonzero,len=2" json:"pubTopics"` // 2 pub topic in the standard mqtt implementation
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
	opts.KeepAlive = 0
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		subscribing := true
		for subscribing {
			dataToken := client.Subscribe(c.config.DataTopic[serviceSubTopic], 0, dataCallback)
			if dataToken.Wait() && dataToken.Error() != nil {
				continue
			}
			actsToken := client.Subscribe(c.config.DataTopic[actServerTopic], 0, actCallback)
			if actsToken.Wait() && actsToken.Error() != nil {
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
	topic := fmt.Sprintf("%s/%s", c.config.PubTopics[modeSwitchTopic], group)
	token := c.client.Publish(topic, 0, false, rule)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (c *MQTTClient) Close() {
	log.Println("MQTT Client - Shutting down")
	if token := c.client.Unsubscribe(c.config.DataTopic[serviceSubTopic]); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	if token := c.client.Unsubscribe(c.config.DataTopic[actServerTopic]); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	// There's a bug that makes the disconnect hang, probably in the mqtt library.
	// Commenting the line below removes the bug but causes a broker-side client error
	// on the gateway closing.
	// c.client.Disconnect(quiesce)
	log.Println("MQTT Client - Stopped")
}

// The function that is called whenever a MQTT message is received on the
// sensor data topic.
func dataCallback(_ mqtt.Client, message mqtt.Message) {
	data := string(message.Payload())
	topicTokens := strings.Split(message.Topic(), "/")
	datatype := topicTokens[len(topicTokens)-1]
	DataHandler(datatype, data)
}

// The function that is called whenever a MQTT message is received on the
// actserver topic.
func actCallback(_ mqtt.Client, message mqtt.Message) {
	data := string(message.Payload())
	ActIPHandler(data)
}
