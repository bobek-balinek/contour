package contour_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/bobek-balinek/contour"
)

func trim(str string) string {
	trimmed := strings.TrimSpace(regexp.MustCompile(`\s+`).ReplaceAllString(str, " "))
	trimmed = strings.Replace(trimmed, " <", "<", -1)
	trimmed = strings.Replace(trimmed, "> ", ">", -1)
	return trimmed
}

func makeEngine() (*contour.Engine, error) {
	engine := contour.New("./views", ".html")
	engine.AddFunc("isSignedIn", func(user string) bool {
		return user == "bob"
	})

	if err := engine.Load(); err != nil {
		return nil, err
	}

	return engine, nil
}

func makeFSEngine() (*contour.Engine, error) {
	engine := contour.NewFS(http.Dir("./views"), ".html")
	engine.AddFunc("isSignedIn", func(user string) bool {
		return user == "bob"
	})

	if err := engine.Load(); err != nil {
		return nil, err
	}

	return engine, nil
}

func TestRenderWithPartisl(t *testing.T) {
	engine, err := makeEngine()
	if err != nil {
		t.Fatalf("load: %v\n", err)
	}

	var buf bytes.Buffer
	engine.Render(&buf, "index", map[string]interface{}{
		"Title": "Hello, Bob!",
	})
	expect := `<h2>Header</h2><h1>Hello, Bob!</h1><h4>Footer</h4>`
	result := trim(buf.String())
	if expect != result {
		t.Fatalf("Expected:\n%s\nResult:\n%s\n", expect, result)
	}
}

func TestRenderWithFunc(t *testing.T) {
	engine, err := makeEngine()
	if err != nil {
		t.Fatalf("load: %v\n", err)
	}

	var buf bytes.Buffer
	engine.Render(&buf, "signed-in", map[string]interface{}{
		"Title": "Dashboard",
		"User":  "bob",
	})
	expect := `Hello Bob!`
	result := trim(buf.String())
	if expect != result {
		t.Fatalf("Expected:\n%s\nResult:\n%s\n", expect, result)
	}
}

func TestRender(t *testing.T) {
	engine, err := makeEngine()
	if err != nil {
		t.Fatalf("load: %v\n", err)
	}

	var buf bytes.Buffer
	engine.Render(&buf, "errors/404", map[string]interface{}{
		"Error": "404 Not Found",
	})
	expect := `<h1>404 Not Found</h1>`
	result := trim(buf.String())
	if expect != result {
		t.Fatalf("Expected:\n%s\nResult:\n%s\n", expect, result)
	}
}

func TestRenderWithLayout(t *testing.T) {
	engine, err := makeEngine()
	if err != nil {
		t.Fatalf("load: %v\n", err)
	}

	var buf bytes.Buffer
	engine.Render(&buf, "index", map[string]interface{}{
		"Title": "Application",
	}, "layouts/app")
	expect := `<!DOCTYPE html><html lang="en"><head><title>Application</title></head><body><h2>Header</h2><h1>Application</h1><h4>Footer</h4></body></html>`
	result := trim(buf.String())
	if expect != result {
		t.Fatalf("Expected:\n%s\nResult:\n%s\n", expect, result)
	}
}

func TestRenderWithFileSystem(t *testing.T) {
	engine, err := makeFSEngine()
	if err != nil {
		t.Fatalf("load: %v\n", err)
	}

	var buf bytes.Buffer
	engine.Render(&buf, "index", map[string]interface{}{
		"Title": "Application",
	}, "layouts/app")
	expect := `<!DOCTYPE html><html lang="en"><head><title>Application</title></head><body><h2>Header</h2><h1>Application</h1><h4>Footer</h4></body></html>`
	result := trim(buf.String())
	if expect != result {
		t.Fatalf("Expected:\n%s\nResult:\n%s\n", expect, result)
	}
}

func TestReload(t *testing.T) {
	engine, err := makeFSEngine()
	if err != nil {
		t.Fatalf("load: %v\n", err)
	}

	engine.Reload(true) // Optional. Default: false

	engine.AddFunc("isAdmin", func(user string) bool {
		return user == "admin"
	})
	if err := engine.Load(); err != nil {
		t.Fatalf("load: %v\n", err)
	}

	if err := ioutil.WriteFile("./views/reload.html", []byte("after reload\n"), 0644); err != nil {
		t.Fatalf("write file: %v\n", err)
	}
	defer func() {
		if err := ioutil.WriteFile("./views/reload.html", []byte("before reload\n"), 0644); err != nil {
			t.Fatalf("write file: %v\n", err)
		}
	}()

	engine.Load()

	var buf bytes.Buffer
	engine.Render(&buf, "reload", nil)
	expect := "after reload"
	result := trim(buf.String())
	if expect != result {
		t.Fatalf("Expected:\n%s\nResult:\n%s\n", expect, result)
	}
}
