package main

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"strings"
)

func newMQ(amqpHost string) {
	if amqpHost == "" {
		return
	}

	conn, err := amqp.Dial(amqpHost)
	if err != nil {
		log.Fatal(err)
	}
	// defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	// defer ch.Close()

	err = ch.ExchangeDeclare(
		"predb",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	q, err := ch.QueueDeclare(
		"pre-releases",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.QueueBind(
		q.Name,
		"",
		"predb",
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"pre-api",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	go mqRun(msgs)
}

func mqRun(msgs <-chan amqp.Delivery) {
	for d := range msgs {
		log.Println(d.RoutingKey)
		switch d.RoutingKey {
		case "pre.insert":
			fallthrough
		case "pre.update":
			fallthrough
		case "pre.delete":
			var p preRow
			err := json.Unmarshal(d.Body, &p)
			if err != nil {
				log.Println(err)
				return
			}

			p.proc()
			if d.RoutingKey == "pre.delete" {
				_, err = p.deleteFromIndex()
				if err != nil {
					log.Println(err)
					return
				}
			} else {
				_, err = p.insertOrUpdateIndex()
				if err != nil {
					log.Println(err)
					return
				}
			}
			backendUpdates <- triggerAction{Action: strings.Split(d.RoutingKey, ".")[1], Row: &p}
			break
		case "nuke.insert":
			var n nuke
			err := json.Unmarshal(d.Body, &n)
			if err != nil {
				log.Println(err)
				return
			}

			p, err := getPre(sphinx, n.PreID, false)
			if err != nil {
				log.Println(err)
				return
			}

			p.setNuke(&n)
			backendUpdates <- triggerAction{Action: n.Type, Row: p}
			break
		}
	}
}
