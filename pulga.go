package pulga

import (
  "fmt"
  "net/http"
  "net/url"
  "regexp"
  "strings"
)

type Handler struct {
  method string
  path   string
  regex  *regexp.Regexp
  http.HandlerFunc
}

type Router struct {
  handlers []*Handler
}

func New() *Router {
  return &Router{}
}

func routeRegex(path string) *regexp.Regexp {
  parts := strings.Split(path, "/")

  for i, p := range parts {
    if strings.HasPrefix(p, ":") {
      parts[i] = "(.*)"
    }
  }

  pattern := strings.Join(parts, "/")
  regex, err := regexp.Compile(pattern)
  if err != nil {
    panic(err)
  }

  return regex
}

func (r *Router) Add(method string, path string, handler http.HandlerFunc) {
  regex := routeRegex(path)
  r.handlers = append(r.handlers, &Handler{method, path, regex, handler})
}

func (r *Router) Get(path string, handler http.HandlerFunc) {
  r.Add("GET", path, handler)
}

func (r *Router) Post(path string, handler http.HandlerFunc) {
  r.Add("POST", path, handler)
}

func (r *Router) Put(path string, handler http.HandlerFunc) {
  r.Add("PUT", path, handler)
}

func (r *Router) Delete(path string, handler http.HandlerFunc) {
  r.Add("DELETE", path, handler)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
  for _, handler := range r.handlers {
    if params, valid := handler.validateRequest(req); valid {
      req.URL.RawQuery = params.Encode() + "&" + req.URL.RawQuery
      handler.ServeHTTP(w, req)
      return
    }
  }

  http.Error(w, "Not Found", http.StatusNotFound)
}

func (h *Handler) validateRequest(req *http.Request) (url.Values, bool) {
  if !h.regex.MatchString(req.URL.Path) || h.method != req.Method {
    return nil, false
  }

  matches := h.regex.FindStringSubmatch(req.URL.Path)[1:]

  if len(matches) == 0 {
    return nil, true
  }

  params := make(url.Values)
  keys := strings.Split(h.path, "/")

  j := 0
  for _, k := range keys {
    if strings.HasPrefix(k, ":") {
      params.Add(k[1:], matches[j])
      j++
    }
  }

  return params, true
}

func (r *Router) Listen(port int) {
  http.Handle("/", r)

  p := fmt.Sprintf(":%d", port)
  http.ListenAndServe(p, nil)
}
