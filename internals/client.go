package internals

import (
  "log"
  "time"
  "net/http"
  "encoding/json"
  "github.com/gorilla/websocket"
  "github.com/jakeNorman007/go_chime/auth/users"
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
  writeWaitTime = 10 * time.Second
  pongWait = 60 * time.Second
  pingPeriod = (pongWait * 9) / 10
  maxMessageSize = 512
)

func ServeWebSocket(core *Core, w http.ResponseWriter, r *http.Request) {
  token, err := r.Cookie("jwt")
  if err != nil {
    log.Println("Error in retrieving JWT: ", err)
    http.Error(w, "Unauthorized", http.StatusUnauthorized)
    return
  }
  
  username, err := users.ExtractUsernameFromToken(token.Value)
  if err != nil {
    log.Println("Error in decoding JWT: ", err)
    http.Error(w, "Unauthorized", http.StatusUnauthorized)
    return
  }

  connection, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    log.Println("Error upgrading to websocket: ", err)
    return
  }

  client := &Client {
    id: username,
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

  cl.connection.SetReadLimit(maxMessageSize)
  cl.connection.SetReadDeadline(time.Now().Add(pongWait))
  cl.connection.SetPongHandler(func(appData string) error {
    cl.connection.SetReadDeadline(time.Now().Add(pongWait))
    return nil
  })

  for {
    _, text, err := cl.connection.ReadMessage()
    if err != nil {
      if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
        log.Printf("Unexpected close error: %v", err)
      }
      break
    }

    var wsMessage WebSocketMessage
    if err := json.Unmarshal(text, &wsMessage); err != nil {
      log.Printf("Error decoding websocket message: %v", err)
      continue
    }

    msg := &Message {
      ClientId: cl.id,
      MessageContent: wsMessage.WsMessageContent,
    }

    cl.core.broadcast <- msg
  }
}

func (cl *Client) readFromWebSocketConnection() {
  ticker := time.NewTicker(pingPeriod)

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
