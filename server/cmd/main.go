package main

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

/**
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hallo from golang")
}

func main() {
	fmt.Println("Started server.. ")
	var httpPort int = 8080
	fmt.Println("Listening on port :" + strconv.Itoa(httpPort))
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":"+strconv.Itoa(httpPort), nil)
	if err != nil {
		panic(err)
	}
}

*/

func main() {

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	var httpPort int = 8080
	fmt.Println("Listening on port :" + strconv.Itoa(httpPort))

	r.Run(":" + strconv.Itoa(httpPort)) // listen and serve on 0.0.0.0:8080
}
