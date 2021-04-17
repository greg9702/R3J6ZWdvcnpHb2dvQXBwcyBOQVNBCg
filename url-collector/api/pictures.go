package api

import (
	"net/http"
	"url-collector/url-collector/fetcher"
	"url-collector/url-collector/models"

	"github.com/gin-gonic/gin"
)

type picturesController struct {
	dataFetcher fetcher.FetcherInterface
}

func NewPicturesController(fetcher fetcher.FetcherInterface) *picturesController {
	p := picturesController{fetcher}
	return &p
}

func (pc *picturesController) GetImages(ctx *gin.Context) {

	pictures := &models.PicturesToBeFetched{}

	err := ctx.Bind(pictures)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = pictures.Validate()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	urlList, err := pc.dataFetcher.FetchData(pictures)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"urls": urlList,
	})
}
