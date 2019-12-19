package vanilla

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	// time allowed to write a message to the peer.
	writeWait = 30 * time.Second

	// time allowed to read the next pong message from the peer.
	pongWait = 2 * time.Second

	// send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}

	// channel reading error
	errChannelReading = errors.New("channel reading error")
	// writer beign error
	errWriterBegin = errors.New("write begin error")
	// writing error
	errWriting = errors.New("writing error")

	// websocket upgrader
	ug = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type client struct {
	// The websocket connection.
	ws *websocket.Conn
	// Buffered channel of outbound messages.
	send chan []byte
}

func getWSHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ws, err := ug.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"reason": "websocket upgrade failed"})
			return
		}

		wsc := &client{ws: ws, send: make(chan []byte, 256)}
		defer close(wsc.send)

		go wsc.writePump()
		go wsc.readPump()
	}
}

func (c *client) readPump() {
	defer func() {
		// c.hub.unregister <- c
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		msg = bytes.TrimSpace(bytes.Replace(msg, newline, space, -1))
		log.Println("ws message <", msg)
	}
}

func (c *client) writePump() {
	// ping ticker
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()

	for {
		select {
		// receive the sending channel message
		case msg, ok := <-c.send:
			if !ok {
				c.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.ws.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Println("failed to create writer")
				return
			}
			w.Write(msg)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				log.Println("failed to close writer")
				return
			}
		case <-ticker.C:
			c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("failed to write ping")
				return
			}
		}
	}
}
