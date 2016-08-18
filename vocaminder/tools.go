package vocaminder

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response is following the [JSend convention](https://labs.omniti.com/labs/jsend)
type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

func sendFailResponse(c *gin.Context, message string) {

	r := &Response{
		Status: "fail",
		Data: map[string]string{
			"message": message,
		},
	}
	jsonResponse, _ := json.Marshal(r)

	c.String(http.StatusBadRequest, string(jsonResponse))
}

func sendSuccessResponse(c *gin.Context, data interface{}) {

	r := &Response{
		Status: "success",
		Data:   data,
	}
	jsonResponse, _ := json.Marshal(r)

	c.String(http.StatusOK, string(jsonResponse))
}
