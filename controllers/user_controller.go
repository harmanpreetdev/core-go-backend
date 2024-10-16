package controllers

import (
	"context"
	"core_two_go/models"
	"core_two_go/services"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	switch r.Method {
	case "POST":
		h.createUser(w, r.WithContext(ctx))
	case "GET":
		h.getUser(w, r.WithContext(ctx))
	case "PUT":
		h.updateUser(w, r.WithContext(ctx))
	case "DELETE":
		h.deleteUser(w, r.WithContext(ctx))
	default:
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("Error decoding user: %v", err)
		writeErrorResponse(w, http.StatusBadRequest, "Invalid input")
		return
	}

	if err := h.userService.CreateUser(ctx, &user); err != nil {
		log.Printf("Error Creating user: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	log.Printf("User created: %v", user)
	writeSuccessResponse(w, http.StatusCreated, user)
}

func (h *UserHandler) getUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	ids, ok := r.URL.Query()["id"]
	if !ok || len(ids) == 0 {
		writeErrorResponse(w, http.StatusBadRequest, "Missing id")
		return
	}

	id, err := strconv.Atoi(ids[0])
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid id format")
	}

	user, err := h.userService.GetUser(ctx, id)
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		writeErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	log.Printf("User retrieved: %v", user)
	writeSuccessResponse(w, http.StatusOK, user)
}

func (h *UserHandler) updateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	ids, ok := r.URL.Query()["id"]
	if !ok || len(ids) == 0 {
		writeErrorResponse(w, http.StatusBadRequest, "Missing id")
		return
	}

	id, err := strconv.Atoi(ids[0])
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid id format")
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("Error decoding user: %v", err)
		writeErrorResponse(w, http.StatusBadRequest, "Invalid input")
		return
	}

	user.ID = id

	if err := h.userService.UpdateUser(ctx, &user); err != nil {
		log.Printf("Error updating user: %v", err)
		writeErrorResponse(w, http.StatusBadRequest, "Invalid input")
		return
	}

	log.Printf("User updated: %v", user)
	writeSuccessResponse(w, http.StatusOK, user)
}

func (h *UserHandler) deleteUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	ids, ok := r.URL.Query()["id"]
	if !ok || len(ids) == 0 {
		writeErrorResponse(w, http.StatusBadRequest, "Missing id")
		return
	}

	id, err := strconv.Atoi(ids[0])
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid id format")
		return
	}

	if err := h.userService.DeleteUser(ctx, id); err != nil {
		log.Printf("Error deleting user: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	log.Printf("User deleted with id: %d", id)
	writeSuccessResponse(w, http.StatusNoContent, nil)
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func writeSuccessResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
