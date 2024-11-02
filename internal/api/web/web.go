package apiweb

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/qreator/worker-pool/pkg/web"
)

type workersPool interface {
	Delete(id []int)
	Add(n int)
	Alive() []int
	Work(jobs []string)
}

type WebWorkersAPI struct {
	workers workersPool
}

func NewWebWorkersAPI(workers workersPool) *WebWorkersAPI {
	return &WebWorkersAPI{
		workers: workers,
	}
}

func (wp *WebWorkersAPI) Register(router *mux.Router) {
	router.HandleFunc("/delete", wp.delete).Methods(http.MethodDelete)
	router.HandleFunc("/alive", wp.alive).Methods(http.MethodGet)
	router.HandleFunc("/add", wp.add).Methods(http.MethodPost)
	router.HandleFunc("/work", wp.work).Methods(http.MethodPost)
}


func (wp *WebWorkersAPI) delete(w http.ResponseWriter, r *http.Request) {
	msg := web.NewLogMsg(r.URL.Path, r.Method)

	body, err := web.ReadRequestBody(r)
	if err != nil {
		msg.Set(err.Error(), http.StatusUnsupportedMediaType)
		web.WriteError(w, msg)
		return
	}

	idsStr := strings.Split(string(body), ",")
	ids := make([]int, len(idsStr))

	for _, idStr := range idsStr {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			msg.Set(err.Error(), http.StatusBadRequest)
			web.WriteError(w, msg)
			return
		}
		ids = append(ids, id)
	}

	wp.workers.Delete(ids)

	msg.Set("success", http.StatusOK)
	web.WriteData(w, msg, map[string]interface{}{"status": "ok"})
}


func (wp *WebWorkersAPI) add(w http.ResponseWriter, r *http.Request) {
	msg := web.NewLogMsg(r.URL.Path, r.Method)

	body, err := web.ReadRequestBody(r)
	if err != nil {
		msg.Set(err.Error(), http.StatusUnsupportedMediaType)
		web.WriteError(w, msg)
		return
	}

	n, err := strconv.Atoi(string(body))
	if err != nil {
		msg.Set(err.Error(), http.StatusBadRequest)
		web.WriteError(w, msg)
		return
	}

	wp.workers.Add(n)

	msg.Set("success", http.StatusOK)
	web.WriteData(w, msg, map[string]interface{}{"status": "ok"})
}


func (wp *WebWorkersAPI) alive(w http.ResponseWriter, r *http.Request) {
	msg := web.NewLogMsg(r.URL.Path, r.Method)

	data := wp.workers.Alive()

	msg.Set("success", http.StatusOK)
	web.WriteData(w, msg, map[string]interface{}{"ids": data})
}


func (wp *WebWorkersAPI) work(w http.ResponseWriter, r *http.Request) {
	msg := web.NewLogMsg(r.URL.Path, http.MethodPost)

	body, err := web.ReadRequestBody(r)
	if err != nil {
		msg.Set(err.Error(), http.StatusUnsupportedMediaType)
		web.WriteError(w, msg)
		return
	}

	jobs := strings.Split(string(body), ",")

	wp.workers.Work(jobs)

	msg.Set("success", http.StatusOK)
	web.WriteData(w, msg, map[string]interface{}{"message": "ok"})
}
