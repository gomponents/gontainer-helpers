package container

type circularDeps struct {
	chain []string
}

func newCircularDeps() *circularDeps {
	return &circularDeps{
		chain: make([]string, 0),
	}
}

func (c *circularDeps) start(id string) []string {
	defer func() {
		c.chain = append(c.chain, id)
	}()
	for _, curr := range c.chain {
		if curr == id {
			return append(c.chain, id)
		}
	}
	return nil
}

func (c *circularDeps) stop() {
	c.chain = c.chain[:len(c.chain)-1]
}

func (c *circularDeps) isStopped() bool {
	return len(c.chain) == 0
}
