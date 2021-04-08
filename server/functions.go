package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/fatih/structs"
	"github.com/gofiber/fiber/v2"
)

func authorize(ctx *fiber.Ctx) error {
	authorization := ctx.Get("Authorization")
	if authorization == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "No authorization token provided.")
	}

	_ctx := context.Background()

	token, err := cdnAuth.VerifyIDToken(_ctx, authorization)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid authorization token provided.")
	}

	ctx.Locals("user", token.UID)

	return ctx.Next()
}

func dummyMiddleware(ctx *fiber.Ctx) error {
	ctx.Locals("user", "solaris")
	return ctx.Next()
}

func contains(values []string, value string) bool {
	for _, elem := range values {
		if elem == value {
			return true
		}
	}

	return false
}

func indexOf(values []string, value string) int {
	for index, elem := range values {
		if elem == value {
			return index
		}
	}

	return -1
}

func addField(s []byte, k string, v interface{}) []byte {
	dummyMap := new(map[string]interface{})
	_ = json.Unmarshal(s, dummyMap)
	(*dummyMap)[k] = v
	return toJSON(dummyMap)
}

func toJSON(s interface{}) []byte {
	body, _ := json.Marshal(s)
	return body
}

// combines the object, latter fields take priority of existing fields
func combine(main map[string]interface{}, rest ...map[string]interface{}) map[string]interface{} {
	for _, m := range rest {
		structs.FillMap(m, main)
	}

	return main
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randSeq(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

type JSONResponse struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func NewResponse(code int, message string) *JSONResponse {
	return &JSONResponse{
		Message: message,
		Code:    code,
	}
}

func NewResponseByError(code int, err error) *JSONResponse {
	return &JSONResponse{
		Message: err.Error(),
		Code:    code,
	}
}

func (response *JSONResponse) SetSuccess(success bool) {
	response.Success = success
}

func (response *JSONResponse) SetCode(code int) {
	response.Code = code
}

func (response *JSONResponse) SetData(data interface{}) {
	response.Data = data
}

func (response *JSONResponse) ToJSON() []byte {
	return toJSON(response)
}
