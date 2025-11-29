package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	SERVER = "http://127.0.0.1:8080"
)

type TicketEntry struct {
	ID        int     `json:"id"`
	PersonID  int     `json:"person_id"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"created_at"`
	CalledAt  *string `json:"called_at"`
}

type QueueNextResponse struct {
	Available bool         `json:"available"`
	Entry     *TicketEntry `json:"entry"`
}

type CallNextResponse struct {
	Entry *TicketEntry `json:"entry"`
}

type StatusRequest struct {
	Status string `json:"status"`
}

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// API для получения следующего в очереди
	router.GET("/api/next", getNext)

	// API для вызова следующего
	router.POST("/api/call-next", callNext)

	// API для отметки как обслуженного
	router.POST("/api/done/:ticket_id", markDone)

	if err := router.Run("0.0.0.0:3000"); err != nil {
		panic(err)
	}
}

func getNext(c *gin.Context) {
	resp, err := http.Get(fmt.Sprintf("%s/queue/next", SERVER))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"error": "Server error"})
		return
	}

	body, _ := io.ReadAll(resp.Body)
	var queueResp QueueNextResponse
	json.Unmarshal(body, &queueResp)

	c.JSON(http.StatusOK, queueResp)
}

func callNext(c *gin.Context) {
	// Вызываем на сервере
	serverResp, err := http.Post(fmt.Sprintf("%s/queue/call_next", SERVER), "application/json", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer serverResp.Body.Close()

	if serverResp.StatusCode == http.StatusNoContent {
		c.JSON(http.StatusOK, gin.H{"empty": true})
		return
	}

	if serverResp.StatusCode != http.StatusOK {
		c.JSON(serverResp.StatusCode, gin.H{"error": "Server error"})
		return
	}

	body, _ := io.ReadAll(serverResp.Body)
	var callResp CallNextResponse
	json.Unmarshal(body, &callResp)

	c.JSON(http.StatusOK, callResp)
}

func markDone(c *gin.Context) {
	ticketID := c.Param("ticket_id")

	statusReq := StatusRequest{Status: "done"}
	payloadBytes, _ := json.Marshal(statusReq)

	resp, err := http.Post(
		fmt.Sprintf("%s/tickets/%s/status", SERVER, ticketID),
		"application/json",
		io.NopCloser(bytes.NewReader(payloadBytes)),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"error": "Server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
