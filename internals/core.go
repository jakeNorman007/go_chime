package internals

import (
  "log"
  "sync"
  "bytes"
  "html/template"
)

type Message struct {
  ClientId        string
  MessageContent  string
}

type WebSocketMessage struct {
  WsMessageContent string       `json:"text"`
  Headers          interface{}  `json:"HEADERS"`
}

type Core struct {
  sync.RWMutex

  clients     map[*Client]bool
  broadcast   chan  *Message
  register    chan  *Client
  unregister  chan  *Client 
  messages    []*Message
}

func NewCore() *Core {
  return &Core {
    clients:     make(map[*Client]bool),
    broadcast:   make(chan *Message),
    register:    make(chan *Client),
    unregister:  make(chan *Client),
  }
}

func (c *Core) Run() {
  for {
    select {
    case client := <- c.register:
      c.Lock()
      c.clients[client] = true
      c.Unlock()

      for _, msg := range c.messages {
        client.send <- GetMessageTemplate(msg)
      }

    case client := <- c.unregister:
      c.Lock()
      if _, ok := c.clients[client]; ok {
        close(client.send)
        delete(c.clients, client)
      }
      c.Unlock()

    case msg := <- c.broadcast:
      c.RLock()
      c.messages = append(c.messages, msg)

      for client := range c.clients {
        select {
        case client.send <- GetMessageTemplate(msg):
        default:
          close(client.send)
          delete(c.clients, client)
        }
      }
      c.RUnlock()
    }
  }
}

func GetMessageTemplate(msg *Message) []byte {
  tmpl, err := template.ParseFiles("templates/message.html")
  if err != nil {
    log.Fatalf("Template parsing: %s", err)
  }

  var renderedMessage bytes.Buffer

  data := struct {
    Username string
    Content  string
  }{
    Username: msg.ClientId,
    Content:  msg.MessageContent,
  }

  err = tmpl.Execute(&renderedMessage, data)
  if err != nil {
    log.Fatalf("Template execution: %s", err)
  }

  return renderedMessage.Bytes()
}
