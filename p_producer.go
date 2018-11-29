producer, err := client.CreateProducer(pulsar.ProducerOptions{
    Topic: "my-topic",
})

if err != nil {
    log.Fatalf("Could not instantiate Pulsar producer: %v", err)
}

defer producer.Close()

msg := pulsar.ProducerMessage{
    Payload: []byte("Hello, Pulsar"),
}

if err := producer.Send(msg); err != nil {
    log.Fatalf("Producer could not send message: %v", err)
}
