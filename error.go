package working_with_users

import (
	"github.com/gin-gonic/gin"
	"github.com/siruspen/logrus"
)

type ModifyError struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Details []Detail `json:"details"`
}

type Detail struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func NewModifyError(code string, message string, details []Detail) ModifyError {
	logrus.Error(message)
	return ModifyError{code, message, details}
}

func RunModifyError(c *gin.Context, statusCode int, error ModifyError) {
	c.AbortWithStatusJSON(statusCode, error)
}
