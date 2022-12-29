package user

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h *Handler) UpdateByIdHandler(c echo.Context) error {
	e := Expenses{}
	err := c.Bind(&e)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	stmt, err := h.DB.Prepare("UPDATE expenses SET title = $1, amount = $2, note = $3, tags = $4 WHERE id=$5;")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	res, err := stmt.Exec(e.Title, e.Amount, e.Note, pq.Array(e.Tag), c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	count, err := res.RowsAffected()
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	if count == 0 {
		return c.JSON(http.StatusNotFound, Err{Message: "expense not found" + err.Error()})
	}

	e.ID = id

	return c.JSON(http.StatusOK, e)
}
