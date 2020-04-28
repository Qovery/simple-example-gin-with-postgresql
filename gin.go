package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"qovery-gin-postgresql/db"
	"time"
)

type JokeAPIResponseBody struct {
	CreatedAt string `json:"created_at"`
	Value     string `json:"value"`
}

func main() {
	r := gin.Default()

	postgres := db.Connect()

	r.GET("/", func(c *gin.Context) {
		resp, err := http.Get("https://api.chucknorris.io/jokes/random")

		if err != nil {
			log.Printf("Error while getting a new Joke, reason: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Something went wrong",
			})
			return
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		var jokeRes JokeAPIResponseBody
		json.Unmarshal(body, &jokeRes)

		print(jokeRes.Value)

		joke := db.Joke{
			Content: jokeRes.Value,
			AddedAt: time.Now(),
		}

		jokes := []db.Joke{joke}

		_, err = postgres.Model(&jokes).Insert(&jokes)

		if err != nil {
			log.Printf("Error while inserting a new Joke to the database, reason: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Something went wrong",
			})
			return
		}

		err = postgres.Model(&jokes).Select()

		if err != nil {
			log.Printf("Error while getting all Jokes, reason: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Something went wrong :(",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Jokes from PostgreSQL",
			"data":    jokes,
		})
		return
	})

	r.Run()
}
