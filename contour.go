// Package contour is a simple and fast template engine for Go.
//
// It features layouts, simple partial includes and flash messages.
// Template files can be loaded at runtime, or from a go:embed directive.
package contour

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
)

// layoutKey defines a method to embed a template within a layout
const layoutKey = "body"

// Engine defines the template engine
type Engine struct {
	cfg       *config
	dir       string
	fs        http.FileSystem
	ext       string
	layout    string
	funcLock  sync.RWMutex
	funcMap   map[string]interface{}
	Flashes   *Flash
	Templates *template.Template
}

// New returns a new template engine with a default configuration
func New(dir, ext string) *Engine {
	engine := &Engine{
		cfg:     &config{},
		dir:     dir,
		ext:     ext,
		layout:  layoutKey,
		funcMap: make(map[string]interface{}),
		Flashes: NewFlash(),
	}

	engine.AddFunc(engine.layout, func() error {
		return errors.New("include called unexpectedly")
	})

	engine.AddFunc("flash", func(key string) []string {
		if len(key) > 0 {
			return engine.Flashes.Get(key)
		}

		return engine.Flashes.All()
	})

	return engine
}

// NewFS returns a template engine with a default configuration and given file system
func NewFS(fs http.FileSystem, ext string) *Engine {
	e := New("/", ext)
	e.fs = fs
	return e
}

// Layout sets a new method to embed a template within a layout
func (e *Engine) Layout(key string) *Engine {
	e.layout = key
	return e
}

// Reload if set to true the templates are reloading on each render,
// use it when you're in development and you don't want to restart
// the application when you edit a template file.
func (e *Engine) Reload(enabled bool) *Engine {
	e.cfg.reload = enabled
	return e
}

// AddFunc adds the function to the template's function map.
// It is legal to overwrite elements of the default actions
func (e *Engine) AddFunc(name string, fn interface{}) *Engine {
	e.funcLock.Lock()
	e.funcMap[name] = fn
	e.funcLock.Unlock()
	return e
}

// AddFuncMap adds the functions from a map to the template's function map.
// It is legal to overwrite elements of the default actions
func (e *Engine) AddFuncMap(m map[string]interface{}) *Engine {
	e.funcLock.Lock()
	for name, fn := range m {
		e.funcMap[name] = fn
	}
	e.funcLock.Unlock()
	return e
}

// Load parses the templates to the engine.
func (e *Engine) Load() error {
	if e.cfg.loaded {
		return nil
	}

	e.funcLock.Lock()
	defer e.funcLock.Unlock()

	e.Templates = template.New(e.dir)
	e.Templates.Funcs(e.funcMap)

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip file if it's a directory or has no file info
		if info == nil || info.IsDir() {
			return nil
		}

		// Skip file if it does not equal the given template extension
		if filepath.Ext(path) != e.ext {
			return nil
		}

		// Get the relative file path
		rel, err := filepath.Rel(e.dir, path)
		if err != nil {
			return err
		}

		name := filepath.ToSlash(rel)
		name = strings.TrimSuffix(name, e.ext)
		buf, err := readFile(e.fs, path)
		if err != nil {
			return err
		}

		// Create new template associated with the current one
		// This enable use to invoke other templates {{ template .. }}
		_, err = e.Templates.New(name).Parse(string(buf))
		if err != nil {
			return err
		}

		return err
	}

	// set config to loaded so that all templates have been parsed
	e.cfg.loaded = true

	if e.fs != nil {
		info, err := stat(e.fs, e.dir)
		if err != nil {
			return walkFn(e.dir, nil, err)
		}
		return walk(e.fs, e.dir, info, walkFn)
	}

	return filepath.Walk(e.dir, walkFn)
}

// Render renders a template with given name using provided data and layout template
func (e *Engine) Render(w io.Writer, name string, data interface{}, layout ...string) error {
	if e.cfg.ShouldReload() {
		if e.cfg.reload {
			e.cfg.loaded = false
		}

		if err := e.Load(); err != nil {
			return err
		}
	}

	tmpl := e.Templates.Lookup(name)
	if tmpl == nil {
		return fmt.Errorf("render: template not found: %s", name)
	}

	// Ensure flashes are cleared after rendering
	defer e.Flashes.Clear()

	// Got the layout
	if len(layout) > 0 && layout[0] != "" {
		lay := e.Templates.Lookup(layout[0])
		if lay == nil {
			return fmt.Errorf("render: layout not found: %s", layout[0])
		}

		e.funcLock.Lock()
		defer e.funcLock.Unlock()

		lay.Funcs(map[string]interface{}{
			e.layout: func() string {
				if err := tmpl.Execute(w, data); err != nil {
					return err.Error()
				}

				return ""
			},
		})

		return lay.Execute(w, data)
	}

	return tmpl.Execute(w, data)
}
