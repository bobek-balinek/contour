package contour

// Flash defines a flash message container
type Flash struct {
	msgs map[string][]string
}

// NewFlash returns a new flash message container
func NewFlash() *Flash {
	return &Flash{msgs: make(map[string][]string)}
}

// PushTo pushes a new message onto the stack with given type i.e. `info`, `error`, or `success`
func (f *Flash) PushTo(key, msg string) {
	f.msgs[key] = append(f.msgs[key], msg)
}

// Push pushes a new message onto the stack with default type `info`
func (f *Flash) Push(msg string) {
	f.PushTo("info", msg)
}

// Get returns all messages of given type
func (f *Flash) Get(key string) []string {
	val := f.msgs[key]

	delete(f.msgs, key)

	return val
}

// Clear clears all messages
func (f *Flash) Clear() {
	f.msgs = make(map[string][]string, 0)
}

// All returns all messages regardless of their type
func (f *Flash) All() []string {
	var msgs []string
	for _, v := range f.msgs {
		msgs = append(msgs, v...)
	}

	f.Clear()

	return msgs
}
