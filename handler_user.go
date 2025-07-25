package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/h0dy/blog-aggregator/internal/database"
)

// handlerLogin func is for login cmd, it sets the given user to the config and check if the user is in the database
func handlerLogin(st *state, cmd command) error {
	if len(cmd.Arg) < 1 {
		return fmt.Errorf("\nusage: %s <name>", cmd.Name)
	}
	username := cmd.Arg[0]

	user, err := st.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("\ncouldn't get the user: %v", err) 
	}

	if err := st.cfg.SetUser(user.Name); err != nil {
		return err
	}
	fmt.Println("User switched/set successfully:")
	logUser(user)
	return nil
}

// handlerRegister func is for registering the given user to database and sets the user to the config
func handlerRegister(st *state, cmd command) error {
	if cmd.Arg == nil {
		return fmt.Errorf("\nusage: %v <name>", cmd.Name)
	}

	username := cmd.Arg[0]
	user, err := st.db.CreateUser(context.Background(), database.CreateUserParams{
		ID: 	   uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      username,
	})
	if err != nil {
		return fmt.Errorf("\ncouldn't create user: %v", err)
	}

	if err := st.cfg.SetUser(user.Name); err != nil {
		return err
	}
	fmt.Println("User created successfully:")
	logUser(user)
	return nil
}

// handlerListUsers func print all the users in the console
func handlerListUsers(st *state, cmd command) error {
	users, err := st.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("\ncouldn't get the user: %v", err)
	}
	for _, user := range users {
		if st.cfg.CurrentUsername == user.Name {
			fmt.Printf("* %v (current)\n", user.Name)
			continue
		}
		fmt.Printf("* %v\n", user.Name)
	}
	return nil
}

func logUser(user database.User) {
	fmt.Printf("USER ID: %v\nUSER Name:%v\n", user.ID, user.Name)
}
// handlerLogUser func log current user
func handlerLogUser(st *state, cmd command, user database.User)error {
	logUser(user)
	return nil
}