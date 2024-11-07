package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// Add some data to use at the beginning
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

// Endpoints

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}
func postAlbums(c *gin.Context) {
	var newAlbum album

	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}
func getAlbumById(c *gin.Context) {
	id := c.Param("id")

	for _, a := range albums {
		if a.ID == id {

			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}
func deleteAlbumById(c *gin.Context) {
	id := c.Param("id")

	for i, item := range albums {
		if item.ID == id {
			albums = append(albums[:i], albums[i+1:]...)
			c.IndentedJSON(http.StatusOK, item)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}
func patchAlbumById(c *gin.Context) {
	id := c.Param("id")
	var newAlbum album

	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	for i, item := range albums {
		if item.ID == id {
			if newAlbum.Artist != "" && newAlbum.Artist != item.Artist {
				albums[i].Artist = newAlbum.Artist
			}
			if newAlbum.Title != "" && newAlbum.Title != item.Title {
				albums[i].Title = newAlbum.Artist
			}
			if newAlbum.Price != item.Price {
				albums[i].Price = newAlbum.Price
			}
			c.IndentedJSON(http.StatusOK, albums[i])
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func notImplemented(c *gin.Context) {
	c.IndentedJSON(http.StatusNotImplemented, gin.H{
		"error":   "Not implemented",
		"message": "The requested route does not exist",
	})
}

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.POST("/albums", postAlbums)
	router.GET("/albums/:id", getAlbumById)
	router.DELETE("/albums/:id", deleteAlbumById)
	router.PATCH("/albums/:id", patchAlbumById)
	router.NoRoute(notImplemented)

	router.Run("localhost:8080")
}
