package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/E-Kelly2018/AirportBnB/internal/config"
	"github.com/E-Kelly2018/AirportBnB/internal/forms"
	"github.com/E-Kelly2018/AirportBnB/internal/models"
	"github.com/E-Kelly2018/AirportBnB/internal/render"
	"log"
	"net/http"
)

const key = "reservation"

//Repo the repository used by the handles
var Repo *Repository

//Repository is the repository typer
type Repository struct {
	App *config.AppConfig
}

//NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

//NewHandlers sets the repository for the handles
func NewHandlers(r *Repository) {
	Repo = r
}

//Home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIp := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIp)

	render.RenderTemplate(w, r, "home.page.html", &models.TemplateData{})
}

//About page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	//preform so logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello Template"
	remoteIp := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIp
	//Send data to template
	render.RenderTemplate(w, r, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data[key] = emptyReservation
	render.RenderTemplate(w, r, "make-reservation.page.html", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email", "phone")
	form.MinLength("first_name", 3, r)
	form.IsEmail("email")
	if !form.Valid() {
		data := make(map[string]interface{})
		data[key] = reservation

		render.RenderTemplate(w, r, "make-reservation.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	m.App.Session.Put(r.Context(), key, reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), key).(models.Reservation)
	if !ok {
		log.Println("Cannot get item from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	m.App.Session.Remove(r.Context(), key)
	data := make(map[string]interface{})
	data[key] = reservation
	render.RenderTemplate(w, r, "reservation-summary.page.html", &models.TemplateData{
		Data: data,
	})
}

//Generals renders the generals page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "generals.page.html", &models.TemplateData{})
}

//Majors renders the majors page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "majors.page.html", &models.TemplateData{})
}

//Availability renders the mAvailability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.html", &models.TemplateData{})
}

//PostAvailability renders the mAvailability page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start_date")
	end := r.Form.Get("end_date")

	w.Write([]byte(fmt.Sprintf("Start date is %s and end date is %s", start, end)))

}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Meaagse string `json:"message"`
}

//AvailabilityJson handles request for availability and sends JSON response
func (m *Repository) AvailabilityJson(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Meaagse: "Available!",
	}
	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		log.Println(err)
	}
	log.Println(string(out))
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

//Contact renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.html", &models.TemplateData{})
}
