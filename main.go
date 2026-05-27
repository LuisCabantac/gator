package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/LuisCabantac/gator/internal/config"
	"github.com/LuisCabantac/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	dbURL := cfg.DBURL

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	programState := &state{
		db:  dbQueries,
		cfg: &cfg,
	}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handleUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.run(programState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}
