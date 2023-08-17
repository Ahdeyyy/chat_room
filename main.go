package main

import (
	"database/sql"
	"fmt"

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

type MessageData struct {
	Messages []Message
	Sender   string
}

func (m Message) IsSender(username string) bool {
	return m.Sender == username
}

func (m Message) DateString() string {
	server_time := time.Now().Local()
	date := m.Created.In(server_time.Location())

	if date.Day() == time.Now().Day() {
		return date.Format("15:04")
	}
	if date.Year() == time.Now().Year() {
		return date.Format("02 Jan 15:04")
	}
	return date.Format("02 Jan 2006 15:04")
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
	defer db.Close()

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
			return c.Render("error", fiber.Map{"message": "An error occured while fetching messages"})
		}

		data := MessageData{Messages: messages, Sender: c.Cookies("username")}

		return c.Render("messages", fiber.Map{"data": data})

	})

	app.Post("/new_message", func(c *fiber.Ctx) error {
		// log.Println(c.FormValue("text"))
		content := c.FormValue("text")
		sender := c.Cookies("username")

		if content != "" && sender != "" {

			stmt := `INSERT into message("content","sender")
		VALUES($1,$2)
		`
			_, err = db.Exec(stmt, content, sender)
			if err != nil {
				return c.Render("error", fiber.Map{"message": "An error occured while sending message"})

			}
		}
		messages, err := getMessages(db)
		if err != nil {
			return c.Render("error", fiber.Map{"message": "An error occured while fetching messages"})

		}

		data := MessageData{Messages: messages, Sender: c.Cookies("username")}
		return c.Render("messages", fiber.Map{"data": data})
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
			return nil, e
		}
		messages = append(messages, msg)
	}

	if err := resultRows.Err(); err != nil {

		return nil, e
	}
	return messages, nil
}

func generateUsername() string {
	seed := rand.NewSource(time.Now().UnixNano())
	gen := rand.New(seed)

	username := fmt.Sprintf("user%d%d%d%d%d%d%d", gen.Intn(10), gen.Intn(10), gen.Intn(10), gen.Intn(10), gen.Intn(10), gen.Intn(10), gen.Intn(10))
	return username

}
