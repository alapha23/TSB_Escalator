package main

import (
    "net/http"
    "fmt"
    "log"
    "context"
    "time"
    "strconv"

    "github.com/apache/pulsar/pulsar-client-go/pulsar"
    "github.com/urfave/negroni"
    //"github.com/davecgh/go-spew/spew"
)

func consume(batch int, topic string)pulsar.Message{
    // Instantiate a Pulsar client
    client, err := pulsar.NewClient(pulsar.ClientOptions{
            URL: "pulsar://localhost:6650",
    })

    if err != nil { log.Fatal(err) }

    // Use the client object to instantiate a consumer
    reader, err := client.CreateReader(pulsar.ReaderOptions{
        Topic:          topic,
        StartMessageID: pulsar.EarliestMessage,
	//ReceiverQueueSize: batch,
	//StartMessageID: pulsar.DeserializeMessageID(start_pos),
    })

    if err != nil { log.Fatal(err) }

    defer reader.Close()

    ctx := context.Background()
    // Listen on the topic with timeout 1 second
    var msg pulsar.Message
    c1 := make(chan string, 1)
    go func() {
        msg, err = reader.Next(ctx)
        if err != nil { log.Fatalf("Error reading from topic: %v", err) }
        c1 <- "Successly received a msg!\n"
    }()
    select {
    case res := <-c1:
        fmt.Println(res)
    case <-time.After(1 * time.Second):
        fmt.Println("timeout after 1 second\n")
	return nil
    }

    if err != nil { log.Fatal(err) }
    // Acknowledges a message to the Pulsar broker
    // consumer.Ack(msg)
    return msg
}

func myhandler(w http.ResponseWriter, r *http.Request) {
    // TODO: support for binlog, kafka, etc.
    // TODO: error checking 
    batch_s, ok := r.URL.Query()["batch"]
    if ok != true { log.Fatal("Error parsing batch") }
    batch_i, err := strconv.Atoi(batch_s[0])
    if err != nil { log.Fatal(err) }

/*    start_pos_s, ok := r.URL.Query()["start_pos"]
    if ok != true { log.Fatal("Error parsing start_pos") }
    start_pos_i, err := strconv.Atoi(start_pos_s[0])
    if err != nil { log.Fatal(err) }
    fmt.Fprintf(w, "start_pos:%d ", start_pos_i);
*/
    topic, ok := r.URL.Query()["topic"]
    if ok != true { log.Fatal("Error parsing topic") }
    // topic must be assigned

    for i:=0; i<batch_i; i++ {
        //LastSaveID := pulsar.DeserializeMessageID(pulsar.EarliestMessage)
        msg := consume(batch_i, topic[0])
        if(msg!=nil) {
            fmt.Fprintf(w, string(msg.Payload()))
        }else {
            fmt.Fprintf(w, "No log resides in pulsar\n")
        }
    }
}

func main() {

    mux := http.NewServeMux()
    mux.HandleFunc("/", myhandler)

    n := negroni.Classic() // Includes some default middlewares
    n.UseHandler(mux)

    http.ListenAndServe(":3000", n)
}
