import (
    "log"
    "runtime"

    "github.com/apache/pulsar/pulsar-client-go/pulsar"
)

func main() {
    client, err := pulsar.NewClient(pulsar.ClientOptions{
        URL: "pulsar://localhost:6650",
        OperationTimeoutSeconds: 5,
        MessageListenerThreads: runtime.NumCPU(),
    })

    if err != nil {
        log.Fatalf("Could not instantiate Pulsar client: %v", err)
    }
}
