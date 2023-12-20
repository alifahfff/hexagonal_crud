package utils

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strings"
)

// Response is standard api response model.
type Response struct {
	Status     int         `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

// Pagination is pagination response model.
type Pagination struct {
	Total       int `json:"total"`
	Limit       int `json:"limit"`
	CurrentPage int `json:"current_page"`
	LastPage    int `json:"last_page"`
}

// ResponseWithJSON to write response with JSON format.
func ResponseWithJSON(c *fiber.Ctx, code int, data interface{}, err error, pagination ...*Pagination) {
	r := Response{
		Status:  code,
		Message: strings.ToLower(http.StatusText(code)),
	}
	if len(pagination) > 0 && pagination[0] != nil {
		r.Pagination = pagination[0]
		if r.Pagination.CurrentPage <= 0 {
			r.Pagination.CurrentPage = 1
		}
		r.Pagination.LastPage = RoundUp(float64(r.Pagination.Total) / float64(r.Pagination.Limit))
		if r.Pagination.LastPage <= 0 {
			r.Pagination.LastPage = 1
		}
	}

	r.Data = data
	if err != nil {
		r.Message = err.Error()
	}

	// Set response header.
	c.Accepts("application/json")

	_ = c.JSON(r)
}
