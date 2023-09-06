package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	gomatchmakek "github.com/LIOU2021/go-match-maker"
)

var myHub *gomatchmakek.Hub

type req struct {
	Data   interface{} `json:"data"`
	RoomId string      `json:"roomId"`
	Id     string      `json:"id"` // user id
}

func main() {
	go initMatchMaker()
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.POST("join", func(c *gin.Context) {
		r := &req{}

		if err := c.BindJSON(&r); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		testNewData := &gomatchmakek.Member{
			Data:   r.Data,
			RoomId: r.RoomId,
			Id:     r.Id,
		}
		go myHub.Join(testNewData)
		c.JSON(200, r)
	})

	r.POST("leave", func(c *gin.Context) {
		r := &req{}

		if err := c.BindJSON(&r); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		testNewData := &gomatchmakek.Member{
			Data:   r.Data,
			RoomId: r.RoomId,
			Id:     r.Id,
		}

		go myHub.Leave(testNewData)

		c.JSON(200, r)
	})

	srv := &http.Server{
		Addr:    ":8008",
		Handler: r,
	}

	ch := make(chan os.Signal, 1)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println(err)
		}
	}()

	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(c); err != nil {
		log.Println("srv.Shutdown:", err)
	}
	myHub.Close()
}

func initMatchMaker() {
	config := gomatchmakek.Config{
		RegisterBuff:   200,
		BroadcastBuff:  200,
		UnRegisterBuff: 200,
		Room:           []string{"a", "b", "c", "d"},
		HubName:        "go-match-maker",
		// Mode:           gomatchmakek.Debug,
		Mode: gomatchmakek.Release,
	}

	myHub = gomatchmakek.New(&config)

	myHub.Run()
}
