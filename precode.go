package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

func main() {
	r := chi.NewRouter()

	logger := log.New(os.Stdout, "logger: ", log.Ldate|log.Lmicroseconds)

	// getting all tasks from the map
	r.Get("/tasks", func(writer http.ResponseWriter, request *http.Request) {
		marshal, err := json.Marshal(tasks)
		if err != nil {
			logger.Print("Couldn't marshal the tasks map")
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		_, err = writer.Write(marshal)
		if err != nil {
			logger.Print("Error occurred during writing data to the response writer")
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		logger.Print("Getting all tasks")

		writer.WriteHeader(http.StatusOK)
	})

	// posting a new task to the map
	r.Post("/tasks", func(writer http.ResponseWriter, request *http.Request) {
		var buf bytes.Buffer
		_, err := buf.ReadFrom(request.Body)
		if err != nil {
			logger.Print("Couldn't read request's body content")
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		body := buf.Bytes()

		if len(body) == 0 {
			logger.Print("Empty body received")
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		task := &Task{}
		err = json.Unmarshal(body, task)
		if err != nil {
			logger.Print("Couldn't unmarshal received body")
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		_, ok := tasks[task.ID]
		if ok {
			logger.Printf("A task with the id %s already exists", task.ID)
			writer.WriteHeader(http.StatusConflict)
			return
		}

		tasks[task.ID] = *task
		logger.Print("New task successfully added")

		writer.WriteHeader(http.StatusCreated)
	})

	// gets a task by its id
	r.Get("/tasks/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id := chi.URLParam(request, "id")
		task, ok := tasks[id]
		if !ok {
			logger.Printf("Couldn't find a task with id equals %s\n", id)
			writer.WriteHeader(http.StatusNotFound)
			return
		}

		data, err := json.Marshal(task)
		if err != nil {
			logger.Print("Couldn't marshal a obtained task")
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		_, err = writer.Write(data)
		if err != nil {
			logger.Print("Couldn't write data to the response")
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	})

	// deletes a record from the map by its id
	r.Delete("/tasks/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id := chi.URLParam(request, "id")
		_, ok := tasks[id]
		if !ok {
			logger.Printf("Couldn't find a task with id equals %s\n", id)
			writer.WriteHeader(http.StatusNotFound)
			return
		}

		delete(tasks, id)
		writer.WriteHeader(http.StatusOK)
	})

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
