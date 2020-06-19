package main

import (
	"flag"
	"log"

	"go.eth.moe/catbus-actuator-wakeonlan/config"
	"go.eth.moe/catbus-actuator-wakeonlan/mqtt"
	"go.eth.moe/catbus-actuator-wakeonlan/wakeonlan"
)

var (
	configPath = flag.String("config-path", "", "path to config")
)

func main() {
	flag.Parse()

	if *configPath == "" {
		log.Fatal("must set -config-path")
	}

	config, err := config.ParseFile(*configPath)
	if err != nil {
		log.Fatalf("could not parse config file: %q", err)
	}

	brokerOptions := mqtt.NewClientOptions()
	brokerOptions.AddBroker(config.Broker)
	brokerOptions.SetAutoReconnect(true)
	brokerOptions.SetConnectionLostHandler(func(_ mqtt.Client, err error) {
		log.Printf("disconnected from MQTT broker %s: %v", config.Broker, err)
	})
	brokerOptions.SetOnConnectHandler(func(broker mqtt.Client) {
		log.Printf("connected to MQTT broker %v", config.Broker)

		for topic := range config.MACsByTopic {
			token := broker.Subscribe(topic, mqtt.AtLeastOnce, func(_ mqtt.Client, msg mqtt.Message) {
				mac, ok := config.MACsByTopic[msg.Topic()]
				if !ok {
					return
				}
				if err := wakeonlan.Wake(mac); err != nil {
					log.Printf("could not send wake-on-lan: %v")
				}
			})
			if err := token.Error(); err != nil {
				log.Printf("could not subscribe to %q: %v", topic, err)

			}
		}
	})

	broker := mqtt.NewClient(brokerOptions)
	if token := broker.Connect(); token.Error() != nil {
		log.Fatalf("could not connect to MQTT broker: %v", token.Error())
	}

	select {}
}
