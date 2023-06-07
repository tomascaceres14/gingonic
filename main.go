package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

type product struct {
	Name         string
	Price        int
	Stock        int
	Code         int
	Published    bool
	CreationDate string
}

func main() {
	router := gin.Default()

	data, err := os.ReadFile("./productos.json")
	if err != nil {
		log.Fatal(err)
	}

	var products []product
	if err := json.Unmarshal(data, &products); err != nil {
		log.Fatal(err)
	}

	fmt.Println(products)

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"products": products,
		})
	})

	router.Run()
}
