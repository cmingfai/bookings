package forms

import (
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

	r = httptest.NewRequest("POST", "/whatever", nil)
	r.PostForm = postedData

	form = New(r.PostForm)
	form.Required("a", "b", "c")

	if !form.Valid() {
		t.Error("form shows invalid when required fields exist")
	}

}

func TestForm_Has(t *testing.T) {
	// invalid case
	postedData := url.Values{}
	form := New(postedData)

	hasField := form.Has("a")
	if hasField {
		t.Error("should not have a field but the form says yes")
	}

	// valid case
	postedData = url.Values{}
	postedData.Add("a", "a")
	form = New(postedData)
	hasField = form.Has("a")

	if !hasField {
		t.Error("should have a field but the form says not")
	}
}

func TestForm_MinLength(t *testing.T) {
	// non existing field

	postedData := url.Values{}
	form := New(postedData)
	isMinLenSatisfied := form.MinLength("non_existing_field", 100)
	if isMinLenSatisfied {
		t.Error("form shows min length for non-existent field")
	}

	// invalid case
	postedData = url.Values{}
	postedData.Add("some_field", "ab")

	form = New(postedData)
	isMinLenSatisfied = form.MinLength("some_field", 10)
	if isMinLenSatisfied {
		t.Error("the field does not have enough length but it shows have")
	}

	isError := form.Errors.Get("some_field")
	if isError == "" {
		t.Error("should have an error, but did not get one")
	}

	// valid case
	postedData = url.Values{}
	postedData.Add("some_field", "1234567890")

	form = New(postedData)
	isMinLenSatisfied = form.MinLength("some_field", 10)
	if !isMinLenSatisfied {
		t.Error("the field does have enough length but it shows does not have")
	}

	isError = form.Errors.Get("some_field")
	if isError != "" {
		t.Error("should not have an error, but got one")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)
	isEmailValid := form.IsEmail("non_existing_field")
	if isEmailValid {
		t.Error("show email valid for an non existing field")
	}

	// valid case
	postedData = url.Values{}
	postedData.Add("email", "mingfai@live.com")

	form = New(postedData)
	isEmailValid = form.IsEmail("email")
	if !isEmailValid {
		t.Errorf("the email [%s] is valid but it says invalid", form.Get("email"))
	}

	// invalid case
	postedData = url.Values{}
	postedData.Add("email", "a")

	form = New(postedData)
	isEmailValid = form.IsEmail("email")
	if isEmailValid {
		t.Errorf("the email [%s] is invalid but it says valid", form.Get("email"))
	}
}
