package utils

import (
	"net/http"

	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"

	"treafik-api/pkg/logger"
)

// ErrResponse defines the return messages when an error occurred.
// Reference will be omitted if it does not exist.
// swagger:model
type ErrResponse struct {
	// Code defines the business error code.
	Code int `json:"code"`

	// Message contains the detail of this message.
	// This message is suitable to be exposed to external
	Message string `json:"message"`

	// Reference returns the reference document which maybe useful to solve this error.
	Reference string `json:"reference,omitempty"`
}

type SuccessResponse struct {
	// Code defines the business error code.
	Code int `json:"code"`
	// Message contains the detail of this message.
	// This message is suitable to be exposed to external
	Message string `json:"message"`

	Result interface{} `json:"result"`
}

func WriteResponse(c *gin.Context, err error, data interface{}) {
	if err != nil {
		logger.Errorf("%#+v", err)
		coder := errors.ParseCoder(err)
		c.JSON(coder.HTTPStatus(), ErrResponse{
			Code:      coder.Code(),
			Message:   coder.String(),
			Reference: coder.Reference(),
		})
		return
	}
	result := SuccessResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Result:  data,
	}
	c.JSON(http.StatusOK, result)
}

func WriteSuccessResponse(c *gin.Context, data interface{}) {
	WriteResponse(c, nil, data)
}
