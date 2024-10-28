package main

import (
	"encoding/json"
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
)

type Article struct {
    Title string `json:"title"`
}

type Category struct {
    Articles []Article `json:"articles"`
}

type ResponseData struct {
    Data []Category `json:"data"`
}

func main() {
    app := fiber.New()

    app.Get("/", func(c *fiber.Ctx) error {
        client := resty.New()

        // Main API call
        resp, err := client.R().
            SetHeader("Content-Type", "application/json").
            Get("http://localhost:1337/api/categories?populate[articles][populate]=main_img")

        if err != nil {
            log.Fatalf("Error while making request: %v", err)
            return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
        }

        // Parsing JSON response to decode into struct
        var responseData ResponseData
        err = json.Unmarshal(resp.Body(), &responseData)
        if err != nil {
            log.Fatalf("Error while parsing response: %v", err)
            return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
        }

        titles := make([]string, 0)
        for _, category := range responseData.Data {
            for _, article := range category.Articles {
                titles = append(titles, article.Title)
            }
        }

        return c.JSON(titles)
    })

    app.Listen(":3000")
}
