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

type Purchase struct {
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	TotalPrice  float64 `json:"total_price"`
}

var products []Product

func main() {
	loadProducts("products.json")
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })

	productsGroup := router.Group("/products")
	productsGroup.GET("/search", PriceGtThan)
	productsGroup.GET("/productparams", createProduct)
	productsGroup.GET("/:id", FindById)
	productsGroup.GET("/searchbyquantity", SearchByQuantityRange)
	productsGroup.GET("/buy", PurchaseDetail)

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

func FindById(c *gin.Context) {
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

func PriceGtThan(c *gin.Context) {

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

func createProduct(c *gin.Context) {
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
	f, err := os.Create("products.json")
	if err != nil {
		log.Fatal(err)
	}

	f.Write(jsonList)
	c.JSON(http.StatusOK, product)
}

func SearchByQuantityRange(c *gin.Context) {
	qFrom := c.Query("from")
	from, err := strconv.Atoi(qFrom)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid 'from' Quantity",
		})
		return
	}

	qTo := c.Query("to")
	to, err := strconv.Atoi(qTo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid 'to' Quantity",
		})
		return
	}

	list := []Product{}
	for _, p := range products {
		if p.Quantity >= from && p.Quantity <= to {
			list = append(list, p)
		}
	}

	c.JSON(http.StatusOK, list)
}

func PurchaseDetail(c *gin.Context) {
	qCodeValue := c.Query("code_value")

	qQuantity := c.Query("quantity")
	quantity, err := strconv.Atoi(qQuantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid quantity",
		})
		return
	}

	for _, p := range products {
		if p.CodeValue == qCodeValue {
			totalPrice := p.Price * float64(quantity)
			purchase := Purchase{
				ProductName: p.Name,
				Quantity:    quantity,
				TotalPrice:  totalPrice,
			}
			c.JSON(http.StatusOK, purchase)
			return
		}
	}
}
