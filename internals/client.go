package internals

import (
  "log"
  "time"
  "bytes"
  "net/http"
  "encoding/json"
  "github.com/google/uuid"
  "github.com/gorilla/websocket"
)

type Client struct {
  id          string
  core        *Core
  connection  *websocket.Conn
  send        chan []byte
}

var upgrader = websocket.Upgrader {
  ReadBufferSize: 1024,
  WriteBufferSize: 1024,
}

const (
  writeWaitTime = 10 * time.Second // gives 10 seconds to write to peer
  pongWait = 60 * time.Second //time that is allowed to recieve next message from peer
  pingPeriod = (pongWait * 9) / 10 // sends "pings" to peer, must be less time than pongWait
  maxMessageSize = 512 // maximum message size allowed from the peer (bytes)
)

func ServeWebSocket(core *Core, w http.ResponseWriter, r *http.Request) {
  connection, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    log.Println(err)
    return
  }

  // for now client id is just going to be a uuid
  id := uuid.New()
  
  client := &Client {
    id: id.String(),
    core: core,
    connection: connection,
    send: make(chan []byte, 256),
  }

  client.core.register <- client

  go client.readFromCore()
  go client.readFromWebSocketConnection()
}

func (cl *Client) readFromCore() {
  defer func() {
    cl.connection.Close()
    cl.core.unregister <- cl
  }()

  cl.connection.SetReadLimit(maxMessageSize) // sets max byte size you can read from a peer
  cl.connection.SetReadDeadline(time.Now().Add(pongWait)) // throws error after 60s of not reading a sent message
  cl.connection.SetPongHandler(func(appData string) error { // appData is the data for the pong application data recieved from the peer
    cl.connection.SetReadDeadline(time.Now().Add(pongWait))
    return nil
  })

  for {
    _, text, err := cl.connection.ReadMessage()
    log.Printf("value: %v", string(text));

    // client application doesn't normally cause an unexpected closing of websocket connection. ".CloseGoingAway" (1001) is the expected closing
    // from the client (closing browser, navigating away, etc.) ".CloseAbnormalClosure" (1006) means connection was terminated by browser or
    // client unexpectedly.
    if err != nil {
      if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
        log.Printf("Error: %v", err)
      }
      break
    }

    msg := &WebSocketMessage {}
    reader := bytes.NewReader(text)
    decoder := json.NewDecoder(reader)
    err = decoder.Decode(msg)
    if err != nil {
      log.Printf("Error: %v", err)
    }

    cl.core.broadcast <- &Message {
      ClientId: cl.id,
      MessageContent: msg.WsMessageContent,
    }
  }
}

func (cl *Client) readFromWebSocketConnection() {
  ticker := time.NewTicker(pingPeriod)

  /// defers closing the connection until something goes wrong
  defer func() {
    ticker.Stop()
    cl.connection.Close()
  }()

  for {
    select {
    case msg, ok := <- cl.send:
      cl.connection.SetWriteDeadline(time.Now().Add(writeWaitTime))
      if !ok {
          cl.connection.WriteMessage(websocket.CloseMessage, []byte{})
          return
        }

      wr, err := cl.connection.NextWriter(websocket.TextMessage)
      if err != nil {
        return
      }

      wr.Write(msg)

      // adding queued messages to chat
      num := len(cl.send)
      for i := 0; i < num; i++ {
        wr.Write(msg)
      }

      if err := wr.Close(); err != nil {
        return
      }
    case <- ticker.C:
      cl.connection.SetWriteDeadline(time.Now().Add(writeWaitTime))
      if err := cl.connection.WriteMessage(websocket.PingMessage, nil); err != nil {
        return
      }
    }
  }
}
