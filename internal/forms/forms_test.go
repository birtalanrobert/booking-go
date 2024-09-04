package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()

	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")

	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	r, _ = http.NewRequest("POST", "/whatever", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")

	if !form.Valid() {
		t.Error("form shows invalid when does it should be valid")
	}
}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	has := form.Has("a")

	if has {
		t.Error("Has should return false when field is not in the body")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")

	r, _ = http.NewRequest("POST", "/whatever", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	has = form.Has("a")

	if !has {
		t.Error("Has should return true when field is in the body")
	}
}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	postedData := url.Values{}
	postedData.Add("a", "a")

	form.MinLength("a", 3)

	if form.Valid() {
		t.Error("shows minlegnth of 3 met when data is shorter")
	}

	isError := form.Errors.Get("a")

	if isError == "" {
		t.Error("should have an error, but did not get one")
	}

	postedData = url.Values{}
	postedData.Add("aa", "aaaaaa")

	form = New(postedData)
	form.MinLength("aa", 3)

	if !form.Valid() {
		t.Error("shows minlength of 3 is not met when it is")
	}

	isError = form.Errors.Get("aa")

	if isError != "" {
		t.Error("should not have an error, but got one")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	postedData.Add("a", "a")

	form.IsEmail("a")

	if form.Valid() {
		t.Error("Has should return false when field is not in the body")
	}

	postedData = url.Values{}
	postedData.Add("a", "aaa@gmail.com")

	form = New(postedData)
	form.IsEmail("a")

	if !form.Valid() {
		t.Error("Has should return true when field is in the body")
	}
}