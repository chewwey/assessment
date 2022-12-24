package main

import (
	"log"
	"os"

	"github.com/chewwey/assessment/user"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func main() {
	user.InitDB()

	e := echo.New()

	e.POST("/expenses", user.CreateExpensesHandler)
	e.GET("/expenses/:id", user.GetExpensesById)
	e.PUT("expenses/:id", user.UpdateByIdHandler)

	log.Fatal(e.Start(os.Getenv("PORT")))
}
