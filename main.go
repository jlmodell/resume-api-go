package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
		URL:     "https://github.com/jlmodell/resume-api-go",
	}

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost, https://odellmay.com",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

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

	app.Put("/api/resume/:field", func(c *fiber.Ctx) error {
		field := c.Params("field")
		fields := []string{
			"skills", "projects", "links", "certifications",
		}
		if !stringSliceContainsString(fields, field) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":  "invalid field",
				"fields": fields,
			})
		}

		var body struct {
			Add string `json:"add" form:"add" xml:"add"`
		}
		if err := c.BodyParser(&body); err != nil {
			return err
		}

		if body.Add == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "'add' field is required",
			})
		}

		fieldValue, err := putAdditionalItemInFieldSlice(ID, body.Add, field)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"message": fmt.Sprintf("'%s' was added to %s", body.Add, field),
			field:     fieldValue,
		})
	})

	app.Delete("/api/resume/:field", func(c *fiber.Ctx) error {
		field := c.Params("field")
		fields := []string{
			"skills", "projects", "links", "certifications",
		}
		if !stringSliceContainsString(fields, field) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":  "invalid field",
				"fields": fields,
			})
		}

		var body struct {
			Remove string `json:"remove" form:"remove" xml:"remove"`
		}
		if err := c.BodyParser(&body); err != nil {
			return err
		}

		if body.Remove == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "'remove' field is required",
			})
		}

		fieldValue, err := delItemInFieldSlice(ID, body.Remove, field)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"message": fmt.Sprintf("'%s' was removed from %s", body.Remove, field),
			field:     fieldValue,
		})
	})

	app.Listen(":8001")
}
