package handlers

import (
	"currency-checker/internal/service"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/mail"
)

type RespondBouncedEmails struct {
	FailedEmails []string `json:"failed_emails"`
}

type HTTPHandler struct {
	emailServ *service.EmailService
	log       logrus.FieldLogger
}

func (h *HTTPHandler) NewGorillaMux(router *mux.Router) *mux.Router {
	router.HandleFunc("/rate", h.GetRate).Methods(http.MethodGet)
	router.HandleFunc("/sendEmails", h.SendEmails).Methods(http.MethodPost)
	router.HandleFunc("/subscribe", h.Subscribe).Methods(http.MethodPost)
	return router
}

func NewHandler(emailService *service.EmailService, logger logrus.FieldLogger) *HTTPHandler {
	return &HTTPHandler{
		emailServ: emailService,
		log:       logger,
	}
}

func (h *HTTPHandler) GetRate(w http.ResponseWriter, r *http.Request) {

	rate, err := h.emailServ.Rate()
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		h.log.Error(err)
		w.WriteHeader(400)
		marshal, err := json.Marshal("can't get bitcoin price: " + err.Error())
		if err != nil {
			http.Error(w, "can't get bitcoin price: "+err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(marshal)
		if err != nil {
			http.Error(w, "can't get bitcoin price: "+err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	marshal, _ := json.Marshal(rate)
	_, err = w.Write(marshal)
	if err != nil {
		h.log.Error(err.Error())
		http.Error(w, "can't get bitcoin price: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
}

func (h *HTTPHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	headerContentType := r.Header.Get("Content-Type")
	if headerContentType != "application/x-www-form-urlencoded" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "parsing failed", http.StatusInternalServerError)
		return
	}

	data := r.FormValue("email")
	_, err = mail.ParseAddress(data)
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		marshal, err := json.Marshal("invalid email")
		if err != nil {
			http.Error(w, "marshalling failed", http.StatusInternalServerError)
			return
		}
		_, err = w.Write(marshal)
		if err != nil {
			http.Error(w, "writing failed", http.StatusInternalServerError)
			return
		}
		return
	}
	marshal, _ := json.Marshal("email is already in the list\n")
	err = h.emailServ.Subscribe(data, h.log)
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(409)
		_, err := w.Write(marshal)
		if err != nil {
			http.Error(w, "writing failed", http.StatusInternalServerError)
			return
		}
		return
	}
}

func (h *HTTPHandler) SendEmails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	failedEmails := h.emailServ.SendEmails(h.log)
	if len(failedEmails) > 0 {
		err := json.NewEncoder(w).Encode(failedEmails)
		if err != nil {
			http.Error(w, "can't return broken emails: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

}
