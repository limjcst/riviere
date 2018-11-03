package api

import (
	"bytes"
	"errors"
	"github.com/gorilla/mux"
	"github.com/limjcst/riviere/config"
	"github.com/limjcst/riviere/listener"
	"github.com/limjcst/riviere/models"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockDB struct {
	ready bool
}

func (db *MockDB) NewTunnel(tunnel *models.Tunnel) error {
	if !db.ready {
		return errors.New("port occupied")
	}
	return nil
}

func (db *MockDB) DeleteTunnel(tunnel *models.Tunnel) (int64, error) {
	if !db.ready {
		return 0, errors.New("tunnel doesn't exist")
	}
	return 1, nil
}

func (db *MockDB) ListTunnel() ([]*models.Tunnel, error) {
	tunnels := make([]*models.Tunnel, 0)
	if !db.ready {
		return tunnels, errors.New("Failed to fetch tunnels")
	}
	tunnels = append(tunnels, &models.Tunnel{})
	return tunnels, nil
}

const testDBDriver = "sqlite3"
const testDBSourceName = "file:test.db?cache=shared&mode=memory"

func GetRouter() (router *mux.Router) {
	return
}

func TestNewRouter(t *testing.T) {
	router := NewRouter(new(config.Config), new(MockDB))
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
	ctx := &ContextInjector{&MockDB{false}, &config.Config{Spec: ""}}
	handler := http.HandlerFunc(ctx.GetSpecEndpoint)
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
	ctx := &ContextInjector{&MockDB{false}, nil}
	handler := http.HandlerFunc(ctx.AddTunnelEndpoint)

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
	ctx := &ContextInjector{&MockDB{false}, nil}
	handler := http.HandlerFunc(ctx.DeleteTunnelEndpoint)

	// Invalid data
	data := bytes.NewBufferString(body)
	CheckEditTunnel(t, "DELETE", handler, data, http.StatusBadRequest)

	body = `{"port": 10000,"forward_address": "127.0.0.1","forward_port": 80}`
	// Add
	data = bytes.NewBufferString(body)
	CheckEditTunnel(t, "POST", http.HandlerFunc(ctx.AddTunnelEndpoint),
		data, http.StatusCreated)

	// Delete
	body = `{"port": 10000}`
	data = bytes.NewBufferString(body)
	CheckEditTunnel(t, "DELETE", handler, data, http.StatusAccepted)

	// Delete again
	data = bytes.NewBufferString(body)
	CheckEditTunnel(t, "DELETE", handler, data, http.StatusNotFound)
}

func TestListTunnel(t *testing.T) {
	ctx := &ContextInjector{&MockDB{false}, nil}
	handler := http.HandlerFunc(ctx.ListTunnelEndpoint)

	// DB fails
	CheckEditTunnel(t, "GET", handler, nil, http.StatusInternalServerError)
	// DB is ready
	ctx = &ContextInjector{&MockDB{true}, nil}
	handler = http.HandlerFunc(ctx.ListTunnelEndpoint)
	CheckEditTunnel(t, "GET", handler, nil, http.StatusOK)
}
