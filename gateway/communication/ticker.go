package communication

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"time"
)

const (
	topicTickerInterval  = 5 * time.Second
	reconnTickerInterval = 5 * time.Second
)

// A ticker periodically sends a message to each actuator using a given protocol
// Its Tick method must be called in a separate goroutine since it's blocking
type Ticker interface {
	Tick()
	Done()
}

// ConnCheckTicker periodically checks that the client is connected and tries to
// reconnect if it is not
type ConnCheckTicker struct {
	endTickChan chan bool
	ticker      *time.Ticker
	client      *mqtt.Client
}

func NewConnCheckTicker(client *mqtt.Client) *ConnCheckTicker {
	return &ConnCheckTicker{
		endTickChan: make(chan bool),
		client:      client,
	}
}

func (cct *ConnCheckTicker) Tick() {
	cct.ticker = time.NewTicker(reconnTickerInterval)
	for {
		select {
		case <-cct.endTickChan:
			return
		case <-cct.ticker.C:
			if !(*cct.client).IsConnected() {
				connectionToken := (*cct.client).Connect()

				if connectionToken.Wait() && connectionToken.Error() != nil {
					log.Println(connectionToken.Error())
				}
			}
		}
	}
}

func (cct *ConnCheckTicker) Done() {
	cct.ticker.Stop()
	cct.endTickChan <- true
}

// Ticker implementation for the MQTT protocol
type TopicTicker struct {
	endTickChan chan bool
	ticker      *time.Ticker
	topic       string
	client      *mqtt.Client
}

// Creates a new TopicTicker that periodically sends MQTT messages on the passed topic.
func NewTopicTicker(topic string, client *mqtt.Client) *TopicTicker {
	return &TopicTicker{
		endTickChan: make(chan bool),
		topic:       topic,
		client:      client,
	}
}

func (tt *TopicTicker) Tick() {
	tt.ticker = time.NewTicker(topicTickerInterval)
	for {
		select {
		case <-tt.endTickChan:
			return
		case <-tt.ticker.C:
			token := (*tt.client).Publish(tt.topic, 0, false, "1")
			if token.Wait() && token.Error() != nil {
				log.Println(token.Error())
			}
		}
	}
}

func (tt *TopicTicker) Done() {
	tt.ticker.Stop()
	tt.endTickChan <- true
}
