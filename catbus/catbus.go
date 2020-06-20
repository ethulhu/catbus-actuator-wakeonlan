// SPDX-FileCopyrightText: 2020 Ethel Morgan
//
// SPDX-License-Identifier: MIT

// Package catbus is a convenience wrapper around MQTT for use with Catbus.
package catbus

import mqtt "github.com/eclipse/paho.mqtt.golang"

type (
	Message        = mqtt.Message
	MessageHandler = func(*Client, Message)

	Client struct {
		mqtt mqtt.Client
	}

	ClientOptions struct {
		DisconnectHandler func(*Client, error)
		ConnectHandler    func(*Client)
	}
)

const (
	atMostOnce byte = iota
	atLeastOnce
	exactlyOnce
)

const (
	Retain = true
)

func NewClient(brokerURI string, options ClientOptions) *Client {
	opts := mqtt.NewClientOptions()

	opts.AddBroker(brokerURI)
	opts.SetAutoReconnect(true)

	opts.SetConnectionLostHandler(func(c mqtt.Client, err error) {
		if options.DisconnectHandler != nil {
			options.DisconnectHandler(&Client{c}, err)
		}
	})
	opts.SetOnConnectHandler(func(c mqtt.Client) {
		if options.ConnectHandler != nil {
			options.ConnectHandler(&Client{c})
		}
	})

	return &Client{mqtt.NewClient(opts)}
}

// Subscribe subscribes to a Catbus MQTT topic.
func (c *Client) Subscribe(topic string, f MessageHandler) error {
	return c.mqtt.Subscribe(topic, atLeastOnce, func(c mqtt.Client, msg mqtt.Message) {
		f(&Client{c}, msg)
	}).Error()
}

// Connect connects to the Catbus MQTT broker and blocks forever.
func (c *Client) Connect() error {
	if err := c.mqtt.Connect().Error(); err != nil {
		return err
	}
	select {}
}
