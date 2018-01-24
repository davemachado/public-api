package main

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var (
	myHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("bar"))
	})
	myHandlerWithError = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
	})
)

func TestNoConfig(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/should/be/stdout/", nil)
	req.RemoteAddr = "111.222.333.444"
	myHandler.ServeHTTP(res, req)

	expect(t, res.Code, http.StatusOK)
	expect(t, res.Body.String(), "bar")
}

func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected [%v] (type %v) - Got [%v] (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}
