// SPDX-FileCopyrightText: 2020 Ethel Morgan
//
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"flag"
	"log"

	"go.eth.moe/catbus-wakeonlan/catbus"
	"go.eth.moe/catbus-wakeonlan/config"
	"go.eth.moe/catbus-wakeonlan/logger"
	"go.eth.moe/catbus-wakeonlan/wakeonlan"
)

var (
	configPath = flag.String("config-path", "", "path to config")
)

func main() {
	flag.Parse()

	if *configPath == "" {
		log.Fatal("must set -config-path")
	}

	log, _ := logger.FromContext(context.Background())

	config, err := config.ParseFile(*configPath)
	if err != nil {
		log.AddField("config-path", *configPath)
		log.WithError(err).Fatal("could not parse config file")
	}

	log.AddField("broker-uri", config.BrokerURI)

	catbusOptions := catbus.ClientOptions{
		DisconnectHandler: func(_ *catbus.Client, err error) {
			log := log
			if err != nil {
				log = log.WithError(err)
			}
			log.Error("disconnected from MQTT broker")
		},
		ConnectHandler: func(client *catbus.Client) {
			log.Info("connected to MQTT broker")

			for topic := range config.MACsByTopic {
				err := client.Subscribe(topic, func(_ *catbus.Client, msg catbus.Message) {
					if string(msg.Payload()) != "on" {
						return
					}
					mac, ok := config.MACsByTopic[msg.Topic()]
					if !ok {
						return
					}

					log.AddField("mac", mac)
					log.AddField("topic", topic)
					if err := wakeonlan.Wake(mac); err != nil {
						log.WithError(err).Error("could not send wake-on-lan packet")
						return
					}
					log.Info("sent wake-on-lan packet")
				})
				if err != nil {
					log := log.WithError(err)
					log.AddField("topic", topic)
					log.Error("could not subscribe to MQTT topic")
				}
			}
		},
	}

	catbusOptions.DefaultPayloadByTopic = map[string][]byte{}
	for topic := range config.MACsByTopic {
		catbusOptions.DefaultPayloadByTopic[topic] = []byte("off")
	}

	catbus := catbus.NewClient(config.BrokerURI, catbusOptions)

	if err := catbus.Connect(); err != nil {
		log.WithError(err).Fatal("could not connect to MQTT broker")
	}
}
