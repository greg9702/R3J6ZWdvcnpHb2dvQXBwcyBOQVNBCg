package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type picturesController struct {
	// TODO
}

func NewPicturesController() *picturesController {
	p := picturesController{}
	return &p
}

func (pc *picturesController) GetImages(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"text": "hello",
	})
}
