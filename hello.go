package hello

import (
  "bufio"
  "fmt"
  "net/http"
)

func init() {
    http.HandleFunc("/", handler)
    firebaseSubscribe()
}

func firebaseSubscribe() {
  req, err := http.NewRequest("GET",
                              "https://gae-node-firebase.firebaseio.com/temp.json",
                              nil)
  req.Header.Add("Accept", "text/event-stream")
  client := &http.Client{}
  resp, err := client.Do(req)

  if err != nil {
    fmt.Printf("Error: %s", err)
  }

  reader := bufio.NewReader(resp.Body)
  for {
    line, err := reader.ReadBytes('\n')
    if err != nil {
      fmt.Println("Got an error!")
    }
    fmt.Printf("%s\n", line)
  }
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hello, world!")
}
