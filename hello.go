package hello

import (
  "fmt"
  "bufio"
  "io/ioutil"
  "net"
  "net/http"
  "appengine"
  "appengine/socket"
  "appengine/urlfetch"
)

func init() {
    http.HandleFunc("/", handler)
}

func urlfetchBased(req *http.Request) {
  ctx := appengine.NewContext(req)
  client := urlfetch.Client(ctx)
  req, err := http.NewRequest("GET", "https://gae-node-firebase.firebaseio.com/temp.json", nil)
  //req.Header.Add("Accept", "text/event-stream")
  resp, err := client.Do(req)
  if err != nil {
    fmt.Printf("Error: %s", err)
  }

  reader := bufio.NewReader(resp.Body)
  for {
    var bytes []byte
    _, err := reader.Read(bytes)
    if err != nil {
      break
    }
    ctx.Errorf("%s\n", bytes)
  }
}

func printRespStreaming(ctx appengine.Context, resp* http.Response) {
  reader := bufio.NewReader(resp.Body)
  for {
    line, err := reader.ReadBytes('\n')
    if err != nil {
      ctx.Errorf("Error print streaming failed: %s\n", err)
      ctx.Errorf("No more output\n")
      //return
    }
    ctx.Errorf("%s\n", line)
  }
}

func printResp(ctx appengine.Context, resp* http.Response) {
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    ctx.Errorf("Error print failed: %s\n", err)
    ctx.Errorf("No more output\n")
    return
  }
  ctx.Errorf("%s", body)
}

func getSocketClient(ctx appengine.Context) (* http.Client) {
  return &http.Client{
    Transport: &http.Transport{
      Dial: func(network, addr string) (net.Conn, error) {
        ctx.Errorf("args: %s %s", network, addr)
        return socket.Dial(ctx, network, addr)
      },
    },
  }
}

func getUrlFetchClient(ctx appengine.Context) (* http.Client) {
  return &http.Client{
    Transport: &urlfetch.Transport{Context: ctx},
  }
}

func firebaseSubscribe(req* http.Request) {
  ctx := appengine.NewContext(req)
  ctx.Errorf("Created request")

  req, _ = http.NewRequest("GET",
                              "https://gae-node-firebase.firebaseio.com/temp.json",
                              nil)

  req.Header.Add("Accept", "text/event-stream")

  client := getSocketClient(ctx)

  go doRequestAndGetResult(ctx, client, req)
}

func doRequestAndGetResult(ctx appengine.Context,
                           client *http.Client,
                           req *http.Request) {
  ctx.Errorf("Do request")
  resp, err := client.Do(req)

  if err != nil {
    ctx.Errorf("Error Do failed: %s %i", err.Error(),
               http.StatusInternalServerError)
    return
  }
  printRespStreaming(ctx, resp)
}

func handler(w http.ResponseWriter, r *http.Request) {
  firebaseSubscribe(r)
  fmt.Fprint(w, "Hello, world!")
}
