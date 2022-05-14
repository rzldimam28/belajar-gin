package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Articles struct {
	gorm.Model
	Title string `gorm:"type:varchar(40)" json:"title"`
	Slug string `gorm:"unique_index" json:"slug"`
	Description string `gorm:"type:text" json:"description"`
}

var db *gorm.DB

func main() {

	var err error
	db, err = gorm.Open("mysql", "root:leesrcyng__@/belajar_gin?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("Can not connect to database")
	}
	defer db.Close()

	db.AutoMigrate(&Articles{})

	router := gin.Default()
	
	router.GET("/", GetHome)

	apiV1 := router.Group("/api/v1/") 
	{
		articles := apiV1.Group("/articles/")
		{
			articles.GET("/", GetArticles)
			articles.GET("/:slug", GetArticle)
			articles.POST("/", PostArticle)
		}
	}

	log.Fatal(router.Run("localhost:8080"))

}

func GetHome(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": "Hello World",
	})
}

func GetArticles(ctx *gin.Context) {
	articles := []Articles{}
	db.Find(&articles)
	
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"status": "OK",
		"data": articles,
	})
}

func GetArticle(ctx *gin.Context) {
	slug := ctx.Param("slug")
	article := Articles{}

	if db.First(&article, "slug = ?", slug).RecordNotFound() {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"status": "Data Not Found",
			"data": nil,
		})
		ctx.Abort()
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 200,
			"status": "OK",
			"data": article,
		})
	}
}

func PostArticle(ctx *gin.Context) {
	article := Articles{
		Title : ctx.PostForm("title"),
		Slug : slug.Make(ctx.PostForm("title")),
		Description : ctx.PostForm("description"),
	}
	db.Create(&article)
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"status": "OK",
		"data": "Success Create Article",
	})
}