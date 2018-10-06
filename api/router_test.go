package api

import (
	"bytes"
	"github.com/gorilla/mux"
	"github.com/limjcst/riviere/listener"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewRouter(t *testing.T) {
	router := NewRouter("")
	err := router.Walk(func(route *mux.Route, router *mux.Router,
		ancestors []*mux.Route) (err error) {
		pathTemplate, err := route.GetPathTemplate()

		if err == nil {
			methods, err := route.GetMethods()
			if err == nil {
				t.Logf("ROUTE: %s %s", pathTemplate, strings.Join(methods, ","))
			}
		}
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestSpec(t *testing.T) {
	req, err := http.NewRequest("GET", "/spec", nil)
	if err != nil {
		t.Fatal(err)
	}
	handler := http.HandlerFunc(GetSpecEndpoint)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
}

func CheckEditTunnel(t *testing.T, method string, handler http.HandlerFunc,
	data io.Reader, targetStatus int) {
	req, err := http.NewRequest(method, "/tunnel", data)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != targetStatus {
		t.Errorf("%s /tunnel returned wrong status code: got %v want %v",
			method, status, targetStatus)
	}
}

func TestAddTunnel(t *testing.T) {
	GlobalPool = listener.NewPool("127.0.0.1")
	defer GlobalPool.Close()
	body := `{"forward_address": "127.0.0.1", "forward_port": 80}`
	handler := http.HandlerFunc(AddTunnelEndpoint)

	// Invalid data
	data := bytes.NewBufferString(body)
	CheckEditTunnel(t, "POST", handler, data, http.StatusBadRequest)

	body = `{"port": 10000,"forward_address": "127.0.0.1","forward_port": 80}`
	// Add for the first time
	data = bytes.NewBufferString(body)
	CheckEditTunnel(t, "POST", handler, data, http.StatusCreated)

	// Add for the second time
	data = bytes.NewBufferString(body)
	CheckEditTunnel(t, "POST", handler, data, http.StatusConflict)
}

func TestDeleteTunnel(t *testing.T) {
	GlobalPool = listener.NewPool("127.0.0.1")
	defer GlobalPool.Close()
	body := `{}`
	handler := http.HandlerFunc(DeleteTunnelEndpoint)

	// Invalid data
	data := bytes.NewBufferString(body)
	CheckEditTunnel(t, "DELETE", handler, data, http.StatusBadRequest)

	body = `{"port": 10000,"forward_address": "127.0.0.1","forward_port": 80}`
	// Add
	data = bytes.NewBufferString(body)
	CheckEditTunnel(t, "POST", http.HandlerFunc(AddTunnelEndpoint),
		data, http.StatusCreated)

	// Delete
	body = `{"port": 10000}`
	data = bytes.NewBufferString(body)
	CheckEditTunnel(t, "DELETE", handler, data, http.StatusAccepted)

	// Delete again
	data = bytes.NewBufferString(body)
	CheckEditTunnel(t, "DELETE", handler, data, http.StatusNotFound)
}
