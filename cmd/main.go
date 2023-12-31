package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	gomatchmaker "github.com/LIOU2021/go-match-maker"
)

var myHub *gomatchmaker.Hub

type req struct {
	Data   interface{} `json:"data" binding:"required"` //require
	RoomId string      `json:"roomId" binding:"required"`
	Id     string      `json:"id" binding:"required"` // user id
}

var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

func init() {
	gomatchmaker.RegisterRedisClient(rdb)

	config := gomatchmaker.Config{
		RegisterBuff:   200,
		BroadcastBuff:  200,
		UnRegisterBuff: 200,
		Room:           []string{"a", "b", "c", "d"},
		HubName:        "go-match-maker",
		// Mode:           gomatchmaker.Debug,
		Mode:           gomatchmaker.Release,
		Interval:       time.Millisecond * 200,
		RegisterFunc:   func(m gomatchmaker.Member) {},
		UnRegisterFunc: func(m gomatchmaker.Member) {},
	}

	myHub = gomatchmaker.New(&config)

	go myHub.Run()
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(
		gin.Recovery(),
		// gin.Logger(),
	)

	r.GET("members", func(c *gin.Context) {
		c.JSON(200, myHub.GetMembers())
	})

	r.POST("join", func(c *gin.Context) {
		r := &req{}

		if err := c.ShouldBindJSON(&r); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		testNewData := &gomatchmaker.Member{
			Data:   r.Data,
			RoomId: r.RoomId,
			Id:     r.Id,
		}
		go myHub.Join(testNewData)
		c.JSON(200, r)
	})

	r.POST("leave", func(c *gin.Context) {
		r := &req{}

		if err := c.ShouldBindJSON(&r); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		testNewData := &gomatchmaker.Member{
			Data:   r.Data,
			RoomId: r.RoomId,
			Id:     r.Id,
		}

		go myHub.Leave(testNewData)

		c.JSON(200, r)
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	ch := make(chan os.Signal, 1)

	go func() {
		for ms := range myHub.Notification() {
			fmt.Println("撮合成功： ", ms[0], ms[1])
		}
	}()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println(err)
		}
	}()

	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	fmt.Println("shutdown...")
	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(c); err != nil {
		fmt.Println("srv.Shutdown:", err)
	}
	myHub.Close()
}
