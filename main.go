package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Product struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
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
	productsGroup.GET("/search", Search)
	productsGroup.GET("/productparams", NewProduct)
	productsGroup.GET("/:id", SearchProduct)

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
func SearchProduct(c *gin.Context) {
	paramId := c.Param("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		log.Fatal(err)
	}
	for _, p := range products {
		if p.Id == id {
			c.JSON(http.StatusOK, p)
			return
		}
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

func NewProduct(c *gin.Context) {
	qId := c.Query("id")
	id, err := strconv.Atoi(qId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Id",
		})
		return
	}
	name := c.Query("name")
	qQuantity := c.Query("quantity")
	quantity, err := strconv.Atoi(qQuantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Quantity",
		})
		return
	}
	code_value := c.Query("code_value")
	qIs_published := c.Query("is_published")
	is_published, err := strconv.ParseBool(qIs_published)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Boolean",
		})
		return
	}
	expiration := c.Query("expiration")
	qPrice := c.Query("price")
	price, err := strconv.ParseFloat(qPrice, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Price",
		})
		return
	}

	product := Product{
		Id:          id,
		Name:        name,
		Quantity:    quantity,
		CodeValue:   code_value,
		IsPublished: is_published,
		Expiration:  expiration,
		Price:       price,
	}

	products := append(products, product)
	jsonList, err := json.Marshal(products)
	if err != nil {
		log.Fatal(err)
	}
	os.WriteFile("/products.json", jsonList, 066)
	c.JSON(http.StatusOK, product)
}
