package main

import (
    "net/http"
    "fmt"
    "log"
    "context"
    "time"

    "github.com/apache/pulsar/pulsar-client-go/pulsar"
    "github.com/urfave/negroni"
    //"github.com/davecgh/go-spew/spew"
)

func consume()string{
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
    // Listen on the topic
    var msg pulsar.Message
    c1 := make(chan string, 1)
    go func() {
        msg, err = consumer.Receive(ctx)
        c1 <- "Successly received a msg!\n"
    }()
    select {
    case res := <-c1:
        fmt.Println(res)
    case <-time.After(1 * time.Second):
        fmt.Println("timeout after 1 second\n")
	return ""
    }

    // Listen on the topic
    if err != nil { log.Fatal(err) }
    // Acknowledges a message to the Pulsar broker
    consumer.Ack(msg)
    // spew.Printf("Message successfully received: %s\n", msg);
    return string(msg.Payload())
}

func main() {

    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
        fmt.Fprintf(w, "Welcome to the home page!\n")
        s := consume();
	if(s!="") {
            fmt.Fprintf(w, s);
        }else {
            fmt.Fprintf(w, "No logs reside in pulsar\n");
	}
    })

    n := negroni.Classic() // Includes some default middlewares
    n.UseHandler(mux)

    http.ListenAndServe(":3000", n)
}
