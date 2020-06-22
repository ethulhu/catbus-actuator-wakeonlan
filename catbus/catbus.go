// SPDX-FileCopyrightText: 2020 Ethel Morgan
//
// SPDX-License-Identifier: MIT

// Package catbus is a convenience wrapper around MQTT for use with Catbus.
package catbus

import (
	"math/rand"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type (
	Message        = mqtt.Message
	MessageHandler = func(*Client, Message)

	Client struct {
		mqtt mqtt.Client

		rebroadcastMu      sync.Mutex
		rebroadcastByTopic map[string]*time.Timer
		rebroadcastPeriod  time.Duration
		rebroadcastJitter  time.Duration
	}


	ClientOptions struct {
		DisconnectHandler func(*Client, error)
		ConnectHandler    func(*Client)

		// Rebroadcast previously seen values every RebroadcastPeriod ± [0,RebroadcastJitter).
		RebroadcastPeriod   time.Duration
		RebroadcastJitter   time.Duration

		// RebroadcastDefaults are optional values to seed rebroadcasting if no prior values are seen.
		// E.g. unless we've been told otherwise, assume a device is off.
		RebroadcastDefaults map[string][]byte
	}

	// Retention is whether or not the MQTT broker should retain the message.
	Retention bool
)

const (
	atMostOnce byte = iota
	atLeastOnce
	exactlyOnce
)

const (
	Retain     = Retention(true)
	DontRetain = Retention(false)
)

const (
	DefaultRebroadcastPeriod = 1 * time.Minute
	DefaultRebroadcastJitter = 15 * time.Second
)

func NewClient(brokerURI string, options ClientOptions) *Client {
	client := &Client{
		rebroadcastByTopic: map[string]*time.Timer{},
		rebroadcastPeriod:  DefaultRebroadcastPeriod,
		rebroadcastJitter:  DefaultRebroadcastJitter,
	}

	if options.RebroadcastPeriod != 0 {
		client.rebroadcastPeriod = options.RebroadcastPeriod
	}
	if options.RebroadcastJitter != 0 {
		client.rebroadcastJitter = options.RebroadcastJitter
	}
	for topic, payload := range options.RebroadcastDefaults {
		// TODO: Allow users to set retention?
		client.rebroadcastLater(topic, DontRetain, payload)
	}

	mqttOpts := mqtt.NewClientOptions()
	mqttOpts.AddBroker(brokerURI)
	mqttOpts.SetAutoReconnect(true)
	mqttOpts.SetConnectionLostHandler(func(c mqtt.Client, err error) {
		if options.DisconnectHandler != nil {
			options.DisconnectHandler(client, err)
		}
	})
	mqttOpts.SetOnConnectHandler(func(c mqtt.Client) {
		if options.ConnectHandler != nil {
			options.ConnectHandler(client)
		}
	})
	client.mqtt = mqtt.NewClient(mqttOpts)

	return client
}

// Connect connects to the Catbus MQTT broker and blocks forever.
func (c *Client) Connect() error {
	if err := c.mqtt.Connect().Error(); err != nil {
		return err
	}
	select {}
}

// Subscribe subscribes to a Catbus MQTT topic.
func (c *Client) Subscribe(topic string, f MessageHandler) error {
	return c.mqtt.Subscribe(topic, atLeastOnce, func(_ mqtt.Client, msg mqtt.Message) {
		c.rebroadcastLater(msg.Topic(), Retention(msg.Retained()), msg.Payload())

		f(c, msg)
	}).Error()
}

// Publish publishes to a Catbus MQTT topic.
func (c *Client) Publish(topic string, retention Retention, payload []byte) error {
	c.rebroadcastLater(topic, retention, payload)

	return c.mqtt.Publish(topic, atLeastOnce, bool(retention), payload).Error()
}

func (c *Client) rebroadcastLater(topic string, retention Retention, payload []byte) {
	c.rebroadcastMu.Lock()
	defer c.rebroadcastMu.Unlock()

	if timer := c.rebroadcastByTopic[topic]; timer != nil {
		_ = timer.Stop()
	}
	c.rebroadcastByTopic[topic] = time.AfterFunc(c.rebroadcastDuration(), func() {
		_ = c.Publish(topic, retention, payload)
	})
}
func (c *Client) rebroadcastDuration() time.Duration {
	jitter := time.Duration(rand.Intn(int(c.rebroadcastJitter)))
	if rand.Intn(1) == 0 {
		return c.rebroadcastPeriod + jitter
	}
	return c.rebroadcastPeriod - jitter
}
