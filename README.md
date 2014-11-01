# Pulga

Pulga is a simple HTTP Routing API implemented for Go programming language.

_This package was implemented for knowledge purpose. I do not recommend use it in production environments._

### Install

    go get github.com/andrielfn/pulga

### Usage

```go
func postHandler(w http.ResponseWriter, r *http.Request) {
  params := r.URL.Query()
  category := params.Get("category")
  slug := params.Get("slug")
  fmt.Fprintf(w, "Category: %s\nSlug: %s", category, slug)
}

func main() {
  r := pulga.New()

  r.Get("/posts/:category/:slug", postHandler)

  r.Listen(4000)
}
```

### Route examples

```go
r.Get("/any/route/:params")
r.Post("/any/route/:params")
r.Put("/any/route/:params")
r.Delete("/any/route/:params")
```
