package contour_test

import (
	"bytes"
	"testing"
)

func TestRenderWithFlash(t *testing.T) {
	engine, err := makeFSEngine()
	if err != nil {
		t.Fatalf("load: %v\n", err)
	}

	defer func() {
		if len(engine.Flashes.Get("info")) > 0 {
			t.Fatalf("success flash messages did not clear out")
		}
	}()

	engine.Flashes.Push("File has been saved")

	var buf bytes.Buffer
	engine.Render(&buf, "flash", map[string]interface{}{
		"Title": "Hello, Bob!",
	})
	expect := `<div class="alert">File has been saved</div>`
	result := trim(buf.String())
	if expect != result {
		t.Fatalf("Expected:\n%s\nResult:\n%s\n", expect, result)
	}

	buf.Reset()

	engine.Flashes.Push("Processing failed")

	engine.Render(&buf, "flash", map[string]interface{}{
		"Title": "Hello, Bob!",
	})
	expect = `<div class="alert">Processing failed</div>`
	result = trim(buf.String())
	if expect != result {
		t.Fatalf("Expected:\n%s\nResult:\n%s\n", expect, result)
	}
}

func TestRenderWithAllFlash(t *testing.T) {
	engine, err := makeFSEngine()
	if err != nil {
		t.Fatalf("load: %v\n", err)
	}

	defer func() {
		if len(engine.Flashes.All()) > 0 {
			t.Fatalf("flash messages did not clear out")
		}
	}()

	engine.Flashes.PushTo("success", "File has been saved")
	engine.Flashes.PushTo("error", "Processing failed")

	var buf bytes.Buffer
	engine.Render(&buf, "flash-all", map[string]interface{}{
		"Title": "Hello, Bob!",
	})
	expect := `<div class="alert">File has been saved</div><div class="alert">Processing failed</div>`
	result := trim(buf.String())
	if expect != result {
		t.Fatalf("Expected:\n%s\nResult:\n%s\n", expect, result)
	}
}
