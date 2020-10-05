package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

var db *sql.DB

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "root"
	dbname   = "tasks"
)

type Task struct {
	Id    uint64 `json:"id"`
	Title string `json:"title"`
}

type Tasks struct {
	Tasks []Task `json:"tasks"`
}

func Connect() error {
	var err error
	db, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname))

	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := Connect(); err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		rows, err := db.Query("SELECT id, title FROM tasks")

		if err != nil {
			return c.Status(500).JSON(err.Error())
		}

		defer rows.Close()
		result := Tasks{}

		for rows.Next() {
			task := Task{}

			if err := rows.Scan(&task.Id, &task.Title); err != nil {
				return err
			}

			result.Tasks = append(result.Tasks, task)
		}

		return c.JSON(result)
	})

	app.Post("/", func(c *fiber.Ctx) error {
		task := new(Task)

		if err := c.BodyParser(task); err != nil {
			return c.Status(400).JSON(err.Error())
		}

		rows, err := db.Query("INSERT INTO tasks (title) VALUES ($1) RETURNING *", task.Title)

		if err != nil {
			return err
		}

		defer rows.Close()

		result := Task{}
		for rows.Next() {
			if err := rows.Scan(&result.Id, &result.Title); err != nil {
				return err
			}
		}
		return c.JSON(result)
	})

	app.Get("/:id", func(c *fiber.Ctx) error {
		taskId, err := strconv.ParseUint(c.Params("id"), 10, 64)
		task := new(Task)

		if err != nil {
			return err
		}

		rows, err := db.Query("SELECT id, title FROM tasks WHERE id=$1", taskId)

		if err != nil {
			return c.Status(500).JSON(err.Error())
		}

		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(&task.Id, &task.Title); err != nil {
				return err
			}
		}

		return c.JSON(task)
	})

	app.Put("/:id", func(c *fiber.Ctx) error {
		taskId, err := strconv.ParseUint(c.Params("id"), 10, 64)
		task := new(Task)

		if err != nil {
			return err
		}

		if err := c.BodyParser(task); err != nil {
			return c.Status(400).JSON(err.Error())
		}

		_, err = db.Query("UPDATE tasks SET title=$1 WHERE id=$2", task.Title, taskId)

		if err != nil {
			return err
		}

		task.Id = taskId
		return c.Status(201).JSON(task)
	})

	app.Delete("/:id", func(c *fiber.Ctx) error {
		taskId, err := strconv.ParseUint(c.Params("id"), 10, 64)

		if err != nil {
			return err
		}

		_, err = db.Query("DELETE FROM tasks WHERE id=$1", taskId)

		if err != nil {
			return err
		}

		return c.JSON(taskId)
	})

	log.Fatal(app.Listen(":8000"))
}
