package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Product struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   int     `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

var products []Product

func main() {
	loadProducts("products.json")
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })

	productsGroup := router.Group("/products")
	{
		productsGroup.GET("/search", Search)
	}

	router.Run()
}

func loadProducts(path string) {
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal([]byte(file), &products); err != nil {
		panic(err)
	}
}

func Search(c *gin.Context) {

	query := c.Query("priceGt")
	priceGt, err := strconv.ParseFloat(query, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid price",
		})
		return
	}

	list := []Product{}
	for _, v := range products {
		if v.Price > priceGt {
			list = append(list, v)
		}
	}

	c.JSON(http.StatusOK, list)
}
