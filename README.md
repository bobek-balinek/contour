# contour

Tiny library for handling templating in a Go web app, batteries included.

## Features

- Layouts
- Flash messages
- Auto-reloading (useful in development)

## Usage

Initialize and integrate **contour** like so:

```
package main

import (
  "net/http"
  "github.com/bobek-balinek/contour"
)

var tmpl *contour.Engine

func init() {
  tmpl = contour.New("./views", ".html")
}

func main() {
  mux := http.NewServeMux()
	mux.HandleFunc("/", index)

	port := ":3000"
	log.Println("Listening on port ", port)
	http.ListenAndServe(port, mux)
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
  w.WriteHeader(http.StatusOK)

  tmpl.Render(w, "index", map[string]interface{}{
		"Title": "Hello World",
	}, "layouts/default")

	if _, err := w.Write([]byte(out)); err != nil {
    w.Write([]byte(err.Error()))
  }
}
```
