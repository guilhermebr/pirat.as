package shortener

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/asaskevich/govalidator"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

// Response struct to http request
type Response struct {
	Status string      `json:"status,omitempty"`
	Data   interface{} `json:"data,omitempty"`
	Error  interface{} `json:"error,omitempty"`
}

// Encode handler
func Encode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	longurl := r.URL.Query().Get("url")

	if longurl == "" {
		errorResponse(w, r, errors.New("url is needed"))
		return
	}

	validURL := govalidator.IsURL(longurl)

	if !validURL {
		errorResponse(w, r, errors.New("invalid url"))
		return
	}

	u, _ := url.Parse(longurl)
	if u.Scheme == "" {
		longurl = "http://" + longurl
	}

	// Search if url exist.
	s := &Shortener{}
	s.LongURL = longurl
	if err := s.SearchByURL(); err != nil {
		errorResponse(w, r, err)
		return
	}

	if s.ShortURL != "" {
		// Return success
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(Response{
			Status: "success",
			Data:   "http://pirat.as/" + s.ShortURL,
		})
		return
	}

	// Insert if new.
	err := s.Insert()

	if err != nil {
		errorResponse(w, r, err)
		return
	}

	// Return success
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(Response{
		Status: "success",
		Data:   "http://pirat.as/" + s.ShortURL,
	})
}

// Redir handler
func Redir(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shorturl := vars["key"]

	s := &Shortener{}
	id, err := shortToID(shorturl)
	if err != nil {
		errorResponse(w, r, err)
		return
	}
	s.ID = id
	err = s.Read()
	if err != nil {
		errorResponse(w, r, err)
		return
	}

	// redirect to LongURL.
	http.Redirect(w, r, s.LongURL, http.StatusMovedPermanently)

	// update stats.
	s.Views = s.Views + 1
	s.LastView = time.Now()
	s.Update()
}

func errorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Println(err)
	w.WriteHeader(422)
	json.NewEncoder(w).Encode(Response{
		Status: "error",
		Error: bson.M{
			"code":    422,
			"message": err.Error(),
		},
	})
}
