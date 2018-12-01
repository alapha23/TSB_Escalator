package main

import (
    "context"
    "log"
    //"fmt"

    "github.com/apache/pulsar/pulsar-client-go/pulsar"
    "github.com/davecgh/go-spew/spew"
)

func main() {
    // Instantiate a Pulsar client
    client, err := pulsar.NewClient(pulsar.ClientOptions{
            URL: "pulsar://localhost:6650",
    })

    if err != nil { log.Fatal(err) }

    // Use the client object to instantiate a consumer
    consumer, err := client.Subscribe(pulsar.ConsumerOptions{
        Topic:            "my-topic",
        SubscriptionName: "sub-1",
    })

    if err != nil { log.Fatal(err) }

    defer consumer.Close()

    ctx := context.Background()

    // Listen indefinitely on the topic
    for {
        msg, err := consumer.Receive(ctx)
        if err != nil { log.Fatal(err) }

        // Do something with the message
	// Acknowledges a message to the Pulsar broker
        // consumer.Ack(msg)
        spew.Printf("Message successfully received: %s\n", msg.Payload());
    }
}
