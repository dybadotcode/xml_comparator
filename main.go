package main

import (
	"xml_comparator/configs"
	"xml_comparator/controllers"
	"xml_comparator/models"

	"github.com/gin-gonic/gin"
)

func main() {
	configs.ReadConfig("config.json")
	xml1 := models.ReadXMLFromFile("xmls/test.xml")
	rss1 := models.RssFromXML(xml1)
	xml2 := models.ReadXMLFromFile("xmls/test update.xml")
	rss2 := models.RssFromXML(xml2)
	models.RssCompare(*rss1, rss2)

	controllers.XML = rss2.RSSXMLTree
	route := gin.Default()
	route.GET("/xml", controllers.GetXML)
	route.Run(configs.Configs.HTTP)
}
