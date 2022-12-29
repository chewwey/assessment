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
	user.InitDB()

	e := echo.New()
	h := handler{}
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Logger.SetLevel(log.INFO)
	e.POST("/expenses", user.CreateExpensesHandler)
	e.GET("/expenses/:id", user.GetExpensesByIdHandler)
	e.PUT("/expenses/:id", user.UpdateByIdHandler)
	e.GET("/expenses", user.GetAllUserHandler)

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
