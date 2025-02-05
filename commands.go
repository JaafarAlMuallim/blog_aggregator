package main

type commands struct {
	handler map[string]func(*state, command) error
}

type commandHandler interface {
	register(name string, f func(*state, command) error)
	run(state *state, cmd command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handler[name] = f
}

func (c *commands) run(state *state, cmd command) error {
	if handler, ok := c.handler[cmd.name]; ok {
		return handler(state, cmd)
	}
	return nil
}
