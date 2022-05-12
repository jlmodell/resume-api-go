package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

const (
	ID string = "627b1d0ac160c4b9d29a6b30"
)

func main() {

	apiDetails := struct {
		AppName string `json:"appName"`
		Stack   string `json:"stack"`
		Version string `json:"version"`
		URL     string `json:"url"`
	}{
		AppName: "resume-api-go",
		Stack:   "Go, Fiber (v2), MongoDB",
		Version: "1.0.0",
		URL:     "https://github.com/jlmodell/resume-api-express",
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(apiDetails)
	})

	app.Get("/api/resume", func(c *fiber.Ctx) error {
		resume, err := getResumeById(ID)
		if err != nil {
			log.Fatal(err)
		}

		return c.JSON(resume)
	})

	app.Put("/api/resume/skills", func(c *fiber.Ctx) error {

		var body struct {
			Skill string `json:"skill" form:"skill" xml:"skill"`
		}
		if err := c.BodyParser(&body); err != nil {
			return err
		}

		if body.Skill == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "skill is required",
			})
		}

		if err := putSkillOnResume(ID, body.Skill); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"message": fmt.Sprintf("'%s' was added to skills", body.Skill),
		})
	})

	app.Listen(":8001")
}
