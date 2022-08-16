package contour

type config struct {
	loaded bool
	reload bool
}

func (c *config) ShouldReload() bool {
	return !c.loaded || c.reload
}
