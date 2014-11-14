package pulga

import (
  "fmt"
  "net/http"
  "net/http/httptest"
  "testing"
)

var basicHandler = func(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hi. Pulga here!")
  w.WriteHeader(http.StatusOK)
}

func TestRouteWithParams(t *testing.T) {
  router := &Router{}
  router.Get("/", basicHandler)
  router.Get("/posts/:category/:slug", basicHandler)

  r, _ := http.NewRequest("GET", "/posts/golang/pulga-router?repo=github", nil)
  w := httptest.NewRecorder()
  router.ServeHTTP(w, r)

  category := r.URL.Query().Get("category")
  slug := r.URL.Query().Get("slug")
  repo := r.URL.Query().Get("repo")

  if category != "golang" {
    t.Errorf("The category must be [golang]. Got [%s].", category)
  }

  if slug != "pulga-router" {
    t.Errorf("The category must be [pulga-router]. Got [%s].", slug)
  }

  if repo != "github" {
    t.Errorf("The category must be [github]. Got [%s].", repo)
  }
}

func TestRouteNotFound(t *testing.T) {
  handler := &Router{}

  r, _ := http.NewRequest("GET", "/", nil)
  w := httptest.NewRecorder()
  handler.ServeHTTP(w, r)

  if w.Code != 404 {
    t.Errorf("Status code must be [404]; Got [%v]", w.Code)
  }
}

func BenchmarkRouteWithoutParams(b *testing.B) {
  handler := &Router{}
  handler.Get("/", basicHandler)

  for i := 0; i < b.N; i++ {
    r, _ := http.NewRequest("GET", "/", nil)
    w := httptest.NewRecorder()
    handler.ServeHTTP(w, r)
  }
}

func BenchmarkRouteWithParams(b *testing.B) {
  handler := &Router{}
  handler.Get("/posts/:category/:slug", basicHandler)

  for i := 0; i < b.N; i++ {
    r, _ := http.NewRequest("GET", "/posts/ruby/slow-language", nil)
    w := httptest.NewRecorder()
    handler.ServeHTTP(w, r)
  }
}

func BenchmarkBareRouter(b *testing.B) {
  r, _ := http.NewRequest("GET", "/", nil)
  w := httptest.NewRecorder()
  mux := http.NewServeMux()
  mux.HandleFunc("/", basicHandler)

  for i := 0; i < b.N; i++ {
    mux.ServeHTTP(w, r)
  }
}
