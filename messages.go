package main

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"strings"
)

func newMQ(amqpHost string) {
	if amqpHost != "" {
		conn, err := amqp.Dial(amqpHost)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		ch, err := conn.Channel()
		if err != nil {
			log.Fatal(err)
		}
		defer ch.Close()

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

		go func() {
			for d := range msgs {
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
		}()
	}
}
