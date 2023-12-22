package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "strings"
)


func TestHandleGetTasks(t *testing.T) {
    req, err := http.NewRequest("GET", "/tasks", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(handleGetTasks)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }
}

func TestHandlePostTask(t *testing.T) {
    req, err := http.NewRequest("POST", "/task", strings.NewReader(`{"title": "Test Task", "description": "Test Description", "status": "Test Status"}`))
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(handlePostTask)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusCreated {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
    }
}

func TestHandlePutTask(t *testing.T) {
    req, err := http.NewRequest("PUT", "/task/1", strings.NewReader(`{"title": "Updated Task", "description": "Updated Description", "status": "Updated Status"}`))
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(handlePutTask)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }
}

func TestHandleDeleteTask(t *testing.T) {
    req, err := http.NewRequest("DELETE", "/task/1", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(handleDeleteTask)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }
}