package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/h0dy/blog-aggregator/internal/config"
	"github.com/h0dy/blog-aggregator/internal/database"
	_ "github.com/lib/pq"
)

// state struct holds a pointer to a config and a pointer to a database
type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	// ReadConfigFile() read the config file in home directory
	cfg, err := config.ReadConfigFile()
	if err != nil {
		log.Fatalf("error reading config file: %v\n", err)
	}

	// open connection to the local database
	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		log.Fatalf("error in connecting to the database (PostgreSQL): %v\n", err)
	}
	defer db.Close()
	//dbQueries holds the queries that are generated by sqlc
	dbQueries := database.New(db)
	
	// userState holds the application state for a user, including config and database access
	userState := &state{
		cfg: &cfg,
		db: dbQueries,
	}

	// commands holds a map of available command names to their associated handler functions
	commands := commands{
		cmds: make(map[string]func(*state, command) error),
	}
	
	// registering all the commands to their associated handler
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	commands.register("users", handlerListUsers)
	commands.register("agg", handlerFeed)
	commands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	commands.register("feeds", handlerListFeeds)
	commands.register("follow", middlewareLoggedIn(handlerFollowFeed))
	commands.register("following", middlewareLoggedIn(handlerFollowingFeeds))
	commands.register("unfollow", middlewareLoggedIn(handlerUnfollowFeed))
	commands.register("current-user", middlewareLoggedIn(handlerLogUser))
	commands.register("browse", middlewareLoggedIn(handlerPosts))
	
	// ensure the user provide a command 
	if len(os.Args) < 2 { 
		log.Fatal("Usage: cli <command> [args...]")
	}

	cmdName := os.Args[1]
	cmdArg := os.Args[2:]

	// run the command 
	if err = commands.run(userState, command{
		Name:cmdName,
		Arg: cmdArg,
	}); err != nil {
		log.Fatal(err)
	}
}

// middlewareLoggedIn is a middleware that checks if the current user is logged in before executing certain handlers    
func middlewareLoggedIn(handler func(st *state, cmd command, user database.User) error) func(*state, command) error {
	
	return func(st *state, cmd command)  error {
		user, err := st.db.GetUser(context.Background(), st.cfg.CurrentUsername)
		if err != nil {
			return err
		}
		return handler(st, cmd, user)
	}
}