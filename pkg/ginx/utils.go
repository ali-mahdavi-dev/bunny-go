package ginx

import (
	"encoding/json"
	"math"
	"net/http"
	"reflect"
	"strings"

	"github.com/ali-mahdavi-dev/bunny-go/internal/framwork/cerrors"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/spf13/cast"
)

// Get access token from header or query parameter
func GetToken(c *gin.Context) string {
	var token string
	auth := c.GetHeader("Authorization")
	prefix := "Bearer "

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

// Parse body json data to struct
func ParseJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		return cerrors.BadRequest("FailedParseJson", err.Error())
	}
	return nil
}

// Parse query parameter to struct
func ParseQuery(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		return cerrors.BadRequest("FailedParseQuery", err.Error())
	}
	return nil
}

// Parse Pagination query parameter to struct
func ParsePaginationQueryParam(c *gin.Context, obj *PaginationResult) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		return cerrors.BadRequest("FailedParseQuery", err.Error())
	}

	if obj.Limit < 1 {
		obj.Limit = 10
	}
	return nil
}

// Parse body form data to struct
func ParseForm(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindWith(obj, binding.Form); err != nil {
		return cerrors.BadRequest("FailedParseForm", err.Error())
	}
	return nil
}

// Response json data with status code
func ResJSON(c *gin.Context, status int, v interface{}) {
	buf, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	c.Set(ResBodyKey, buf)
	c.Data(status, "application/json; charset=utf-8", buf)
	c.Abort()
}

func ResSuccess(c *gin.Context, v interface{}) {
	ResJSON(c, http.StatusOK, ResponseResult{
		Success: true,
		Data:    v,
	})
}

func ResOK(c *gin.Context) {
	ResJSON(c, http.StatusOK, ResponseResult{
		Success: true,
	})
}

func CalculatePagination(total, limit, skip int64) (int64, int64) {
	pages := int64(math.Ceil(float64(total) / float64(limit)))
	page := (skip / limit) + 1
	return pages, page
}

func ResPage(c *gin.Context, v interface{}, pr *PaginationResult) {
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
	ResJSON(c, http.StatusOK, ResponseResult{
		Success: true,
		Data:    v,
		Total:   total,
		Page:    page,
		Pages:   pages,
	})
}

func ResError(c *gin.Context, err error, status ...int) {
	var ierr cerrors.Error
	if e, ok := err.(cerrors.Error); ok {
		ierr = e
	} else {
		ierr = cerrors.InternalServerError(cast.ToString(err))
	}

	code := int(ierr.Code())
	if len(status) > 0 {
		code = status[0]
	}

	ResJSON(c, code, ResponseResult{Error: cerrors.New(ierr.ID(), int32(code), ierr.Message(), ierr.Detail(), ierr.Status())})
}
