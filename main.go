package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	_ "github.com/lib/pq"
)

type Message struct {
	Id      string
	Content string
	Sender  string
	Created time.Time
}

var connStr string

func main() {
	connStr = os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://ade:password@127.0.0.1:5432/chat?sslmode=disable"
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(logger.New())
	app.Static("/", "./scripts")
	app.Get("/", func(c *fiber.Ctx) error {
		username := ""
		if c.Cookies("username") != "" {
			username = c.Cookies("username")
		} else {

			username = generateUsername()
			cookie := new(fiber.Cookie)
			cookie.Name = "username"
			cookie.Value = username
			cookie.Expires = time.Now().Add(104 * 7 * 24 * time.Hour)
			c.Cookie(cookie)
		}

		return c.Render("index", fiber.Map{"User": username})
	})

	app.Get("/messages", func(c *fiber.Ctx) error {
		messages, e := getMessages(db)
		if e != nil {
			log.Println(e) // TODO: render an error dialog
		}

		return c.Render("messages", fiber.Map{"messages": messages})

	})

	app.Post("/new_message", func(c *fiber.Ctx) error {
		log.Println(c.FormValue("text"))
		content := c.FormValue("text")
		sender := c.Cookies("username")
		stmt := `INSERT into message("content","sender")
		VALUES($1,$2)
		`
		_, err = db.Exec(stmt, content, sender)
		if err != nil {
			log.Println(err)
		}
		messages, err := getMessages(db)
		if err != nil {
			log.Println(err)
		}
		return c.Render("messages", fiber.Map{"messages": messages})
	})

	app.Listen(":3000")
}

func getMessages(db *sql.DB) ([]Message, error) {

	stmt := `SELECT * FROM message`
	messages := make([]Message, 0)
	resultRows, e := db.Query(stmt)
	if e != nil {
		return nil, e

	}
	defer resultRows.Close()

	for resultRows.Next() {
		msg := Message{}
		e = resultRows.Scan(&msg.Id, &msg.Content, &msg.Sender, &msg.Created)
		if e != nil {
			log.Println(e) //TODO: Render an error dialog
			return nil, e
		}
		messages = append(messages, msg)
	}

	if err := resultRows.Err(); err != nil {
		log.Println(err)
		return nil, e
	}
	return messages, nil
}

func generateUsername() string {
	seed := rand.NewSource(time.Now().UnixNano())
	gen := rand.New(seed)

	username := fmt.Sprintf("anon%d%d%d%d%d%d%d", gen.Intn(10), gen.Intn(10), gen.Intn(10), gen.Intn(10), gen.Intn(10), gen.Intn(10), gen.Intn(10))
	return username

}
