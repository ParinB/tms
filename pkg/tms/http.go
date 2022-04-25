package tms

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

const (
	defaultLimit = 20
	defaultPage  = 0 //zero based pagination
)

//HTTPTaskManager contains http handler for task manager
type HTTPTaskManager struct {
	store    Storer
	validate *validator.Validate
}

//NewHTTPTaskManager returns a new HTTPTaskManager
func NewHTTPTaskManager(store Storer) *HTTPTaskManager {
	return &HTTPTaskManager{
		store:    store,
		validate: validator.New(),
	}
}
//CreateUserHandler creates a handler for create user http request
func (h *HTTPTaskManager) CreateUserHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(newResponseBody(err.Error()))
			return
		}
		if err := h.validate.Struct(req); err != nil {
			output := map[string]string{}
			for _, err := range err.(validator.ValidationErrors) {
				switch err.Field() {
				case "FirstName", "LastName":
					output[err.Field()] = "must be greater than 1 character and be alphabetic"
				case "Email":
					output[err.Field()] = "must be a valid email address"
				}
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&output)
			return
		}
		if err := h.store.CreateUser(r.Context(), req.FirstName, req.LastName, req.Email); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(newResponseBody(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(newResponseBody("user created successfully"))
	})
}
//CreateUserTask creates a task of a  particular user
func (h *HTTPTaskManager) CreateUserTaskHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID, _ := strconv.Atoi(vars["id"])
		var req CreateTaskRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(newResponseBody(err.Error()))
			return
		}

		output := map[string]string{}
		if err := h.validate.Struct(req); err != nil {
			for _, err := range err.(validator.ValidationErrors) {
				switch err.Field() {
				case "Description":
					output[err.Field()] = "must be greater than 1 character and be alphabetic"
				}
			}
		}

		err := h.store.CreateTask(r.Context(), req.Description, uint(userID), req.StartTime.t, req.EndTime.t, req.ReminderPeriod.t)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(newResponseBody(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(newResponseBody("task created successfully"))
	})
}
//GetTaskHandler returns a  task
func (h *HTTPTaskManager) GetTaskHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.Atoi(vars["id"])
		task, err := h.store.GetTask(r.Context(), uint(id))
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(newResponseBody("task not found"))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(newResponseBody(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&task)
	})
}
//GetUserHandler returns 1  user
func (h *HTTPTaskManager) GetUserHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.Atoi(vars["id"])
		user, err := h.store.GetUser(r.Context(), uint(id))
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(newResponseBody("user not found"))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(newResponseBody(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&user)
	})
}
//GetUsersHandler returns  users
func (h *HTTPTaskManager) GetUsersHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limit, page := extractPagination(r)
		offset := limit * page
		users, err := h.store.GetUsers(r.Context(), limit, offset)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(newResponseBody(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&users)
	})
}
//GetTasksHandler returns all tasks
func (h *HTTPTaskManager) GetTasksHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limit, page := extractPagination(r)
		offset := limit * page
		tasks, err := h.store.GetTasks(r.Context(), limit, offset)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(newResponseBody(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&tasks)
	})
}
//GetUserTasksHandler  returns  tasks belonging to a particular user
func (h *HTTPTaskManager) GetUserTasksHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID, _ := strconv.Atoi(vars["id"])

		limit, page := extractPagination(r)
		offset := limit * page
		tasks, err := h.store.GetUserTasks(r.Context(), uint(userID), limit, offset)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(newResponseBody(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&tasks)
	})
}
//DeleteUserHandler  deletes a specific user
func (h *HTTPTaskManager) DeleteUserHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.Atoi(vars["id"])
		err := h.store.DeleteUser(r.Context(), uint(id))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(newResponseBody(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(newResponseBody("user deleted successfully"))
	})
}
type responseBody struct {
	Message string `json:"message"`
}
func newResponseBody(msg string) *responseBody {
	return &responseBody{
		Message: msg,
	}
}
//extractPagination extracts pagination   parameters from  url
func extractPagination(r *http.Request) (uint, uint) {
	slimit := r.URL.Query().Get("limit")
	limit, _ := strconv.Atoi(slimit)
	if limit == 0 {
		limit = defaultLimit
	}

	spage := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(spage)
	return uint(limit), uint(page)
}
