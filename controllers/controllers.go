package controllers

import (
	"net/http"
	"xml_comparator/models"

	"github.com/gin-gonic/gin"
)

//XML ...
var XML *models.XMLTree

//GetXML ...
func GetXML(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"rss_json": *XML})
}
