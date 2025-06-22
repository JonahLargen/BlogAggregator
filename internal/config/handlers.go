package config

import "fmt"

func handlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("login command requires a username argument")
	}
	userName := cmd.Args[0]
	if err := s.Config.SetUser(userName); err != nil {
		return err
	}
	fmt.Printf("Logged in as user %s\n", userName)
	return nil
}
