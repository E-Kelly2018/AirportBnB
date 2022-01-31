package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/url", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have recived valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/url", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	r, _ = http.NewRequest("POST", "/url", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("Shows does not have required fileds when it does")
	}
}

func TestForm_Has(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	has := form.Has("whatever")
	if has {
		t.Error("form shows has field when it does not")
	}

	postedData = url.Values{}
	postedData.Add("a", "a")
	form = New(postedData)
	has = form.Has("a")
	if !has {
		t.Error("shows for does not have field when it should")
	}

}

func TestForm_MinLength(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	form.MinLength("whatever", 10)
	if form.Valid() {
		t.Error("form shows min length fornot existent field")
	}

	isError := form.Errors.Get("whatever")
	if isError == "" {
		t.Error("Should have an error but did not get one")
	}

	postedData = url.Values{}
	postedData.Add("a", "some value")
	form = New(postedData)
	form.MinLength("a", 100)
	if form.Valid() {
		t.Error("Shows minlength of 100 met when data is shorter ")
	}

	postedData = url.Values{}
	postedData.Add("a", "abc123")
	form = New(postedData)
	form.MinLength("a", 1)
	if !form.Valid() {
		t.Error("Shows minlength of 1 is not met when it is")
	}

	isError = form.Errors.Get("a")
	if isError != "" {
		t.Error("Should not have an error but got one")
	}

}

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	form.IsEmail("x")
	if form.Valid() {
		t.Error("form shows valid email for non existent field")
	}

	postedData = url.Values{}
	postedData.Add("email", "eoin@gmail.com")
	form = New(postedData)
	form.IsEmail("email")
	if !form.Valid() {
		t.Error("got invalid email when we should not have ")

	}

	postedData = url.Values{}
	postedData.Add("email", "x")
	form = New(postedData)
	form.IsEmail("email")
	if form.Valid() {
		t.Error("got valid for invalid email address")

	}
}
