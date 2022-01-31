package models

import "github.com/E-Kelly2018/AirportBnB/internal/forms"

//TemplateData holds data sent from handler to templates
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[int]int
	FloatMap  map[float32]float32
	Data      map[string]interface{}
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
	Form      *forms.Form
}
