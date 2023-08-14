package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

func main() {
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

	app.Post("/new_message", func(c *fiber.Ctx) error {
		log.Println(c.FormValue("text"))
		return c.Render("messages", fiber.Map{})
	})

	err := godotenv.Load()
	if err != nil {
		log.Println("error, could not load envieonment variable")
	}
	db_url := os.Getenv("DB_URL")
	fmt.Println(db_url)
	app.Listen(":3000")
}

func generateUsername() string {
	seed := rand.NewSource(time.Now().UnixNano())
	gen := rand.New(seed)

	username := fmt.Sprintf("anon%d%d%d%d%d%d%d", gen.Intn(10), gen.Intn(10), gen.Intn(10), gen.Intn(10), gen.Intn(10), gen.Intn(10), gen.Intn(10))
	return username

}
