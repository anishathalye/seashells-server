package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"github.com/anishathalye/seashells-server/datamanager"
	"time"
)

const minPasswordLength = 10

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsHandler(w http.ResponseWriter, r *http.Request, manager *datamanager.DataManager) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("failed to set websocket upgrade: %v", err)
		return
	}
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	defer conn.Close()
	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pingTimeout))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pingPeriod))
		return nil
	})

	t, msg, err := conn.ReadMessage()
	if err != nil {
		log.Printf("error while reading from websocket: %v", err)
		return
	}
	if t != websocket.TextMessage {
		return
	}
	id := string(msg)
	sess := manager.Get(id)
	if sess == nil {
		// this shouldn't happen unless the user got really unlucky
		// with the timing or if someone is up to no good
		return
	}
	data, done := sess.Subscribe()
	defer done()

	for {
		select {
		case message, ok := <-data:
			conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			if !ok {
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				return
			}
			w.Write(message)
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func runWeb(manager *datamanager.DataManager) {
	r := gin.Default()
	r.Use(gin.Logger())
	r.Static("/static", "static")
	r.StaticFile("/favicon.ico", "resources/favicon.ico")
	r.LoadHTMLGlob("templates/*.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	r.GET("/v/:id", func(c *gin.Context) {
		id := c.Param("id")
		if manager.Get(id) == nil {
			c.HTML(http.StatusNotFound, "oops.html", gin.H{
				"message": "Session not found.",
			})
			return
		}
		c.HTML(http.StatusOK, "terminal.html", gin.H{
			"id": id,
		})
	})

	r.GET("/ws", func(c *gin.Context) {
		wsHandler(c.Writer, c.Request, manager)
	})

	// plaintext
	r.GET("/p/:id", func(c *gin.Context) {
		sess := manager.Get(c.Param("id"))
		if sess == nil {
			c.String(http.StatusNotFound, "404: Not Found")
			return
		}
		c.Data(http.StatusOK, "text/plain; charset=utf-8", sess.Dump())
	})

	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "oops.html", gin.H{
			"message": "Page not found.",
		})
	})

	password := os.Getenv("ADMIN_PASSWORD")
	if len(password) > minPasswordLength {
		attachAdmin(r, password, manager)
	}

	r.Run()
}

func attachAdmin(base *gin.Engine, password string, manager *datamanager.DataManager) {
	admin := base.Group("/inspect", gin.BasicAuth(gin.Accounts{
		"admin": password,
	}))
	admin.GET("/", func(c *gin.Context) {
		var lines []string
		for _, sess := range manager.All() {
			lines = append(lines, fmt.Sprintf("%s", sess.String()))
		}
		c.HTML(http.StatusOK, "admin.html", gin.H{
			"sessions": lines,
		})
	})
}
