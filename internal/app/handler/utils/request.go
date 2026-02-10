package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vukieuhaihoa/bookmark-management/pkg/common"
)

func bindInputFromBodyWithoutValidation[T any](c *gin.Context) (*T, error) {
	input := new(T)
	if err := c.ShouldBindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, common.InputFieldError(err))
		c.Abort()
		return nil, err
	}

	return input, nil
}
