package config

import "fmt"

type commands struct {
	commands map[string]func(*State, Command) error
}

func NewCommands() *commands {
	c := &commands{}
	c.register("login", handlerLogin)
	c.register("register", handlerRegister)
	c.register("reset", handlerReset)
	c.register("users", handlerListUsers)
	c.register("agg", handlerAgg)
	c.register("addfeed", handlerAddFeed)
	c.register("feeds", handlerFeeds)
	c.register("follow", handlerFollow)
	c.register("following", handlerFollowing)
	return c
}

func (c *commands) Run(s *State, cmd Command) error {
	handler, exists := c.commands[cmd.Name]
	if !exists {
		return fmt.Errorf("command %s not found", cmd.Name)
	}
	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*State, Command) error) {
	if c.commands == nil {
		c.commands = make(map[string]func(*State, Command) error)
	}
	if _, exists := c.commands[name]; exists {
		panic(fmt.Sprintf("command %s already registered", name))
	}
	c.commands[name] = f
}
