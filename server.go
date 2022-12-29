package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/chewwey/assessment/user"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
)

func main() {
	db := user.InitDB()

	e := echo.New()
	h := user.Handler{DB: db}
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Logger.SetLevel(log.INFO)
	e.POST("/expenses", h.CreateExpensesHandler)
	e.GET("/expenses/:id", h.GetExpensesByIdHandler)
	e.PUT("/expenses/:id", h.UpdateByIdHandler)
	e.GET("/expenses", h.GetAllUserHandler)

	go func() {
		if err := e.Start(os.Getenv("PORT")); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
