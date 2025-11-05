package httputils

import (
	"encoding/json"
	"math"
	"reflect"
	"strings"

	"shikposh-backend/pkg/framework/cerrors"

	"github.com/gofiber/fiber/v3"
	"github.com/spf13/cast"
)

// Token
func GetToken(c fiber.Ctx) string {
	auth := c.Get("Authorization")
	prefix := "Bearer "
	token := ""

	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	} else {
		token = auth
	}

	if token == "" {
		token = c.Query("accessToken")
	}

	return token
}

// Parsing
func ParseJSON(c fiber.Ctx, obj interface{}) error {
	if err := c.Bind().Body(obj); err != nil {
		return cerrors.BadRequest("FailedParseJson", err.Error())
	}
	return nil
}

func ParseQuery(c fiber.Ctx, obj interface{}) error {
	if err := c.Bind().Query(obj); err != nil {
		return cerrors.BadRequest("FailedParseQuery", err.Error())
	}
	return nil
}

func ParsePaginationQueryParam(c fiber.Ctx, obj *PaginationResult) error {
	if err := c.Bind().Query(obj); err != nil {
		return cerrors.BadRequest("FailedParseQuery", err.Error())
	}
	if obj.Limit < 1 {
		obj.Limit = 10
	}
	return nil
}

func ParseForm(c fiber.Ctx, obj interface{}) error {
	if err := c.Bind().Body(obj); err != nil {
		return cerrors.BadRequest("FailedParseForm", err.Error())
	}
	return nil
}

// Responses
func ResJSON(c fiber.Ctx, status int, v interface{}) error {
	buf, err := json.Marshal(v)
	if err != nil {
		return err
	}
	c.Set(ResBodyKey, string(buf))
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	c.Status(status)
	return c.Send(buf)
}

func ResSuccess(c fiber.Ctx, v interface{}) error {
	return ResJSON(c, fiber.StatusOK, ResponseResult{
		Success: true,
		Data:    v,
	})
}

func ResOK(c fiber.Ctx) error {
	return ResJSON(c, fiber.StatusOK, ResponseResult{
		Success: true,
	})
}

func CalculatePagination(total, limit, skip int64) (int64, int64) {
	if limit <= 0 {
		limit = 1 // Prevent division by zero
	}
	pages := int64(math.Ceil(float64(total) / float64(limit)))
	page := (skip / limit) + 1
	return pages, page
}

func ResPage(c fiber.Ctx, v interface{}, pr *PaginationResult) error {
	var total, pages, page int64
	if pr != nil {
		total = pr.Total
		pages, page = CalculatePagination(total, pr.Limit, pr.Skip)
	}
	if page < 1 {
		page = 1
	}
	if pages < 1 {
		pages = 1
	}
	if page > pages {
		page = pages
	}

	reflectValue := reflect.Indirect(reflect.ValueOf(v))
	if reflectValue.IsNil() {
		v = make([]interface{}, 0)
	}

	return ResJSON(c, fiber.StatusOK, ResponseResult{
		Success: true,
		Data:    v,
		Total:   total,
		Page:    page,
		Pages:   pages,
	})
}

func ResError(c fiber.Ctx, err error) error {
	var ierr cerrors.Error
	if e, ok := err.(cerrors.Error); ok {
		ierr = e
	} else {
		ierr = cerrors.InternalServerError(cast.ToString(err))
	}

	code := int(ierr.Code())

	return ResJSON(c, code, ResponseResult{Error: ierr})
}
