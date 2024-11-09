package main

import (
	"fmt"
	"net/http"

	"example.com/simple-server/pkg/db"
	"example.com/simple-server/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// type albumBody struct {
// 	ID     string  `json:"id"`
// 	Title  string  `json:"title"`
// 	Artist string  `json:"artist"`
// 	Price  float64 `json:"price"`
// }

// Endpoints
func getAlbums(c *gin.Context, db *gorm.DB) {
	var albums []models.Album
	if result := db.Find(&albums); result.Error != nil {
		fmt.Println(result.Error)
	}
	c.IndentedJSON(http.StatusOK, albums)
}

func postAlbums(c *gin.Context, db *gorm.DB) {
	var newAlbum models.Album
	if result := db.Create(&newAlbum); result.Error != nil {
		fmt.Println(result.Error)
	}
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func getAlbumById(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var album models.Album
	if result := db.Find(&album, id); result.Error != nil {
		fmt.Println(result.Error)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
	} else {
		if album.ID == 0 {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
		} else {
			c.IndentedJSON(http.StatusOK, album)
		}
	}
}

func deleteAlbumById(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var album models.Album
	if result := db.Delete(&album, id); result.Error != nil {
		fmt.Println(result.Error)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
	} else {
		c.IndentedJSON(http.StatusOK, true)
	}
}

func patchAlbumById(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var data models.Body
	if err := c.BindJSON(&data); err != nil {
		c.IndentedJSON(http.StatusBadRequest, data)
	}
	var item models.Album
	if result := db.First(&item, id); result.Error != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
	} else {
		if item.ID == 0 {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
		} else {
			changes := map[string]interface{}{}
			if item.Title != data.Title {
				changes["Title"] = data.Title
			}
			if item.Artist != data.Artist {
				changes["Artist"] = data.Artist
			}
			if item.Price != data.Price {
				changes["Price"] = data.Price
			}
			if len(changes) != 0 {
				if err := db.Model(&item).Updates(changes).Error; err != nil {
					c.IndentedJSON(http.StatusOK, gin.H{"message": "Record didn't receive any changes"})
				} else {
					c.IndentedJSON(http.StatusOK, item)
				}
			} else {
				c.IndentedJSON(http.StatusOK, gin.H{"message": "Record didn't receive any changes"})
			}
		}
	}
}

func notImplemented(c *gin.Context) {
	c.IndentedJSON(http.StatusNotImplemented, gin.H{
		"error":   "Not implemented",
		"message": "The requested route does not exist",
	})
}

func main() {
	db := db.Init()
	router := gin.Default()

	router.GET("/albums", func(c *gin.Context) { getAlbums(c, db) })
	router.POST("/albums", func(c *gin.Context) { postAlbums(c, db) })
	router.GET("/albums/:id", func(c *gin.Context) { getAlbumById(c, db) })
	router.DELETE("/albums/:id", func(c *gin.Context) { deleteAlbumById(c, db) })
	router.PATCH("/albums/:id", func(c *gin.Context) { patchAlbumById(c, db) })
	router.NoRoute(notImplemented)

	router.Run("localhost:8080")
}
