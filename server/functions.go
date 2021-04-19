package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
)

func authorize(ctx *fiber.Ctx) error {
	authorization := ctx.Get("Authorization")
	if authorization == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "No authorization token provided.")
	}

	// _ctx := context.Background()

	// token, err := cdnAuth.VerifyIDToken(_ctx, authorization)
	// if err != nil {
	// 	return fiber.NewError(fiber.StatusUnauthorized, "Invalid authorization token provided.")
	// }

	// ctx.Locals("user", token.UID)

	if authorization != cdnConfig.Authorization {
		return fiber.NewError(fiber.StatusUnauthorized, "No authorization token provided.")
	}

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

func toJSON(s interface{}) []byte {
	body, _ := json.Marshal(s)
	return body
}

func Set(s []string) []string {
	seen := make(map[string]struct{}, len(s))
	j := 0
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		s[j] = v
		j++
	}
	return s[:j]
}

func roundFloat64(num float64) string {
	return fmt.Sprintf("%.2f", num)
}

func getFileSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%vB", size)
	} else if size < 1048576 {
		num := float64(size / 1024)
		return fmt.Sprintf("%vKiB", roundFloat64(num))
	} else if size < 1073741824 {
		num := float64(size / 1048576)
		return fmt.Sprintf("%vMiB", roundFloat64(num))
	} else {
		num := float64(size / 1073741824)
		return fmt.Sprintf("%vGiB", roundFloat64(num))
	}
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

type JSONModifier struct {
	Map map[string]interface{}
}

func NewJSONModifer(data interface{}) *JSONModifier {
	jsonMap := new(map[string]interface{})
	json.Unmarshal(toJSON(data), jsonMap)

	return &JSONModifier{
		Map: *jsonMap,
	}
}

func (m *JSONModifier) AddField(key string, value interface{}) {
	m.Map[key] = value
}

func (m *JSONModifier) ToJSON() []byte {
	return toJSON(m.Map)
}
