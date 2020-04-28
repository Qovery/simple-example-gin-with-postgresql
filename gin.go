package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"io/ioutil"
	"log"
	"net/http"
	"qovery-gin-postgresql/db"
	"time"
)

type JokeAPIResponseBody struct {
	CreatedAt  string        `json:"created_at"`
	Value      string        `json:"value"`
}

type Joke struct {
	Content string    `json:"content"`
	AddedAt time.Time `json:"added_at"`
}

func main() {
	r := gin.Default()

	postgres := db.Connect()

	CreateJokeTable(postgres)

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

		joke := Joke{
			Content: jokeRes.Value,
			AddedAt: time.Now(),
		}

		jokes := []Joke{joke}

		insert, err := postgres.Model(&jokes).Insert(&jokes)

		if err != nil {
			log.Printf("Error while inserting a new Joke to the database, reason: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Something went wrong",
			})
			return
		}

		print(insert.RowsAffected())

		if err != nil {
			log.Printf("Error while a new Joke, reason: %v\n", err)
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

func CreateJokeTable(db *pg.DB) error {
	opts := orm.CreateTableOptions{
		IfNotExists: true,
	}

	createError := db.CreateTable(&Joke{}, &opts)

	if createError != nil {
		log.Printf("Error while creating Jokes table, reason: %v\n", createError)
		return createError
	}

	log.Printf("Jokes table created")
	return nil
}
