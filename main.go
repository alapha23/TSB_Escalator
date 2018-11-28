package main

import (
    "net/http"
    "strings"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
    // writes something on http://localhost:8080/something
    message := r.URL.Path
    message = strings.TrimPrefix(message, "/")
    message = "Hello " + message

    w.Write([]byte(message))
}

func main() {
    http.HandleFunc("/", sayHello)
    if err := http.ListenAndServe(":8080", nil); err != nil {
        panic(err)
    }
}
