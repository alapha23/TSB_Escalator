package main

import (
    "net/http"
    "fmt"
    "strings"
    "log"
    "context"

    "github.com/apache/pulsar/pulsar-client-go/pulsar"
    "github.com/urfave/negroni"
    //"github.com/davecgh/go-spew/spew"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
    // writes something on http://localhost:8080/something
    message := r.URL.Path
    message = strings.TrimPrefix(message, "/")
    message = "Hello " + message

    w.Write([]byte(message))
}

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
    msg, err := consumer.Receive(ctx)

    if err != nil { log.Fatal(err) }
    // Do something with the message
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
    fmt.Fprintf(w, s);
    s = consume();
    fmt.Fprintf(w, s);
  })

  n := negroni.Classic() // Includes some default middlewares
  n.UseHandler(mux)

  http.ListenAndServe(":3000", n)
}

/*
func main() {
    // a basic web server
    http.HandleFunc("/", sayHello)
    if err := http.ListenAndServe(":8080", nil); err != nil {
        panic(err)
    }
}
*/
