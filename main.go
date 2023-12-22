package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"

	_ "github.com/Atabek786/koznek/docs"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

var db *sql.DB
var dbMutex sync.Mutex

func main() {
	var err error
	connStr := "user=postgres dbname=koznek password=20050608 sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()

	http.HandleFunc("/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetTasks(w, r)
		case http.MethodPost:
			handlePostTask(w, r)
		case http.MethodPut:
			handlePutTask(w, r)
		case http.MethodDelete:
			handleDeleteTask(w, r)
		default:
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		}
	})

	http.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// @Summary Get all tasks
// @Description Get all tasks
// @Tags tasks
// @Accept  json
// @Produce  json
// @Success 200 {array} Task
// @Router /task [get]
func handleGetTasks(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT id, title, description, status FROM task")
    if err != nil {
        log.Println("Error retrieving tasks:", err)
        http.Error(w, "Failed to retrieve tasks", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    tasks := []Task{}
    for rows.Next() {
        var task Task
        err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status)
        if err != nil {
            log.Println("Error scanning task:", err)
            http.Error(w, "Failed to retrieve tasks", http.StatusInternalServerError)
            return
        }
        tasks = append(tasks, task)
    }

    if err = rows.Err(); err != nil {
        log.Println("Error iterating over rows:", err)
        http.Error(w, "Failed to retrieve tasks", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    err = json.NewEncoder(w).Encode(tasks)
    if err != nil {
        log.Println("Error encoding tasks:", err)
        http.Error(w, "Failed to encode tasks", http.StatusInternalServerError)
        return
    }
}

// @Summary Create a new task
// @Description Create a new task
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param id path int true "Task ID"
// @Success 201 {object} Task
// @Router /task [post]
func handlePostTask(w http.ResponseWriter, r *http.Request) {
    var newTask Task
    err := json.NewDecoder(r.Body).Decode(&newTask)
    if err != nil {
        log.Println("Error decoding request body:", err)
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if newTask.Title == "" || newTask.Description == "" {
        http.Error(w, "Title and Description are required fields", http.StatusBadRequest)
        return
    }

    dbMutex.Lock()
    defer dbMutex.Unlock()

    result, err := db.Exec("INSERT INTO task (id, title, description, status) VALUES ($1, $2, $3)", newTask.ID, newTask.Title, newTask.Description, newTask.Status)
    if err != nil {
        log.Println("Error inserting new task:", err)
        http.Error(w, "Failed to create task", http.StatusInternalServerError)
        return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        log.Println("Error getting rows affected:", err)
        http.Error(w, "Failed to create task", http.StatusInternalServerError)
        return
    }

    if rowsAffected == 0 {
        http.Error(w, "Failed to create task", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}

// @Summary Update on existing task
// @Description Update on existing task
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param id path int true "Task ID"
// @Success 200 {object} Task
// @Router /task/{id} [put]
func handlePutTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
    taskID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid task ID", http.StatusBadRequest)
        return
    }

    var updatedTask Task
    err = json.NewDecoder(r.Body).Decode(&updatedTask)
    if err != nil {
        log.Println("Error decoding request body:", err)
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if updatedTask.Title == "" || updatedTask.Description == "" {
        http.Error(w, "Title and Description are required fields", http.StatusBadRequest)
        return
    }

    dbMutex.Lock()
    defer dbMutex.Unlock()

    result, err := db.Exec("UPDATE task SET title = $2, description = $3, status = $4 WHERE id = $1", taskID, updatedTask.Title, updatedTask.Description, updatedTask.Status)
    if err != nil {
        log.Println("Error updating task:", err)
        http.Error(w, "Failed to update task", http.StatusInternalServerError)
        return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        log.Println("Error getting rows affected:", err)
        http.Error(w, "Failed to update task", http.StatusInternalServerError)
        return
    }

    if rowsAffected == 0 {
        http.Error(w, "Task not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

// @Summary Delete a task
// @Description Delete a task by ID
// @Tags tasks
// @Produce json
// @Param id path int true "Task ID"
// @Success 204
// @Router /task/{id} [delete]
func handleDeleteTask(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r) 
    taskID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid task ID", http.StatusBadRequest)
        return
    }

    dbMutex.Lock()
    defer dbMutex.Unlock()

    result, err := db.Exec("DELETE FROM task WHERE id = $1", taskID)
    if err != nil {
        log.Println("Error deleting task:", err)
        http.Error(w, "Failed to delete task", http.StatusInternalServerError)
        return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        log.Println("Error getting rows affected:", err)
        http.Error(w, "Failed to delete task", http.StatusInternalServerError)
        return
    }

    if rowsAffected == 0 {
        http.Error(w, "Task not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}