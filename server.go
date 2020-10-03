package main

import (
  "strconv"
	"github.com/gofiber/fiber"
)

type Task struct {
  Id uint64 `json:"id"`
	Title string `json:"title"`
}

var id uint64

func main() {
	app := fiber.New()

  var tasks []Task

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(tasks)
  })

  app.Post("/", func(c *fiber.Ctx) error {
    task := new(Task)

    if err := c.BodyParser(task); err != nil {
      return err
    }

    id++
    task.Id = id
    tasks = append(tasks, *task)

    return c.JSON(task)
  })

  app.Get("/:id", func (c *fiber.Ctx) error {
    taskId, err := strconv.ParseUint(c.Params("id"), 10, 64)

    if err != nil {
      return err
    }

    for _, s := range tasks {
      if s.Id == taskId {
        return c.JSON(s)
      }
    }

    return c.JSON(nil)
  })

  app.Put("/:id", func(c *fiber.Ctx) error {
    t := new(Task)
    taskId, err := strconv.ParseUint(c.Params("id"), 10, 64)

    if err != nil {
      return err
    }

    if err := c.BodyParser(t); err != nil {
      return err
    }

    for i, s := range tasks {
      if s.Id == taskId {
        t.Id = s.Id
        tasks[i] = *t
      }
    }

    return c.JSON(t)
  })

  app.Delete("/:id", func(c *fiber.Ctx) error {
    taskId, err := strconv.ParseUint(c.Params("id"), 10, 64)

    if err != nil {
      return err
    }

    var t Task
    for i, s := range tasks {
      if s.Id == taskId {
        t = s
        tasks = append(tasks[:i], tasks[i + 1:]...)
      }
    }

    return c.JSON(t)
  })

	app.Listen(":8000")
}

