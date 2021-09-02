package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	port := flag.Int("port", 9000, "HTTP listen port")
	flag.Parse()

	db, err := newInventory()
	if err != nil {
		log.Fatalf("Error creating sample inventory: %v", err)
	}

	r := gin.Default()
	r.POST("/inventory", newAddInventoryItem(db))
	r.GET("/inventory", newSearchInventoryHandler(db))
	addr := fmt.Sprintf(":%d", *port)
	log.Fatal(r.Run(addr))
}

func newAddInventoryItem(db *inventory) gin.HandlerFunc {
	return func(c *gin.Context) {
		var item InventoryItem
		if err := c.BindJSON(&item); err != nil {
			log.Printf("Data binding error: %v", err)
			// BindJSON will implicitly return HTTP 400 on failure
			return
		}
		if err := db.add(item); err != nil {
			log.Printf("Failed to add item: %v", err)
			c.AbortWithStatus(http.StatusConflict)
			return
		}
		c.Status(http.StatusCreated)
	}
}

func newSearchInventoryHandler(db *inventory) gin.HandlerFunc {
	return func(c *gin.Context) {
		search := c.Query("searchString")

		skipVal := c.Query("skip")
		skip, err := strconv.Atoi(skipVal)
		if err != nil {
			log.Printf("Invalid skip parameter: %v", err)
			skip = 0
		}

		limitVal := c.Query("limit")
		limit, err := strconv.Atoi(limitVal)
		if err != nil {
			log.Printf("Invalid limit parameter: %v", err)
			limit = 0
		}

		log.Printf("Searching inventory for %q, skipping %d, limiting to %d", search, skip, limit)
		res := db.search(search, skip, limit)
		c.JSON(http.StatusOK, res)
	}
}
