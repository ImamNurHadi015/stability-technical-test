package handlers

import (
	"strconv"

	"stability-test-task-api/models"
	"stability-test-task-api/store"

	"github.com/gofiber/fiber/v2"
)

func GetTasks(c *fiber.Ctx) error {
	tasks := store.GetAllTasks()
	return c.JSON(tasks)
}

func GetTask(c *fiber.Ctx) error {
	idParam := c.Params("id")

	id, _ := strconv.Atoi(idParam)

	task := store.GetTaskByID(id)

	if task == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "task not found",
		})
	}

	return c.JSON(task)
}

func CreateTask(c *fiber.Ctx) error {
	var task models.Task

	if err := c.BodyParser(&task); err != nil {
		return err
	}

	   // Validasi title agar tidak boleh kosong (bisa dengan 400 status code)
	   if task.Title == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "title is required",
        })
    }

	//Menambahkan pointer untuk mengikuti perubahan
	store.AddTask(&task)

	return c.JSON(task)
}

func DeleteTask(c *fiber.Ctx) error {
	idParam := c.Params("id")

	id, _ := strconv.Atoi(idParam)

	//Memeriksa terlebih dulu apakah task ada
    task := store.GetTaskByID(id)
    if task == nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "task not found",
        })
    }
	store.DeleteTask(id)

	return c.Status(200).JSON(fiber.Map{
        "message": "deleted",
    })
}
