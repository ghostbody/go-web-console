package exceptions

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RespError represents the error of a response
type RespError struct {
	status int
	err    error
}

func (re *RespError) Error() string {
	return re.err.Error()
}

// NewRespError creates a RespError with an error
func NewRespError(status int, err error) error {
	return &RespError{
		status: status,
		err:    err,
	}
}

// NewRespErrorWithStr creates a with string and formatters
func NewRespErrorWithStr(status int, msg string, a ...interface{}) error {
	return &RespError{
		status: status,
		err:    fmt.Errorf(msg, a...),
	}
}

// HandleErrorMiddleware is a gin middleware which handles common errors
func HandleErrorMiddleware(c *gin.Context) {
	c.Next()
	detectedErrors := c.Errors.ByType(gin.ErrorTypeAny)
	if len(detectedErrors) > 0 {
		err := detectedErrors[0].Err
		switch err.(type) {
		case *RespError:
			respErr := err.(*RespError)
			c.AbortWithStatusJSON(
				respErr.status,
				gin.H{
					"message": respErr.Error(),
				},
			)
		default:
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{
					"message": err.Error(),
				},
			)
		}
		return
	}
}
