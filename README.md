# Gator - RSS Feed Aggregator

Gator is a command-line RSS feed aggregator built with Go and PostgreSQL. It allows you to manage RSS feeds, subscribe to channels, and browse aggregated content from a single interface.

> This project is part of the [Boot.dev](https://boot.dev) guided curriculum for learning backend development.

## Prerequisites

Before you can run gator, you'll need to have the following installed:

- **Go** (1.26.3 or higher) - [Download Go](https://golang.org/dl/)
- **PostgreSQL** - either:
  - Install PostgreSQL locally - [Download PostgreSQL](https://www.postgresql.org/download/)
  - Or use Docker and Docker Compose (see [Setting up PostgreSQL with Docker](#setting-up-postgresql-with-docker) below)

## Installation

Install the gator CLI using `go install`:

```bash
go install github.com/LuisCabantac/gator@latest
```

This will build and install the `gator` binary to your `$GOPATH/bin` directory (typically `~/go/bin`). Make sure this directory is in your `$PATH`.

## Setting up PostgreSQL with Docker

If you don't have PostgreSQL installed locally, you can run it in Docker using the included `docker-compose.yaml`:

```bash
docker-compose up -d
```

This will start a PostgreSQL container on port `5433` with the following credentials:
- **Username**: `gator`
- **Password**: `postgres`
- **Database**: `gator`

Use the connection string `postgres://gator:postgres@localhost:5433/gator` in your config file (note the port is `5433`).

To stop the database:

```bash
docker-compose down
```

## Configuration

Gator requires a configuration file to connect to your PostgreSQL database. 

### Setting up the config file

Create a `.gatorconfig.json` file in your home directory with the following format:

```json
{
  "db_url": "postgres://user:password@localhost:5432/gator",
  "current_user_name": ""
}
```

Replace `user`, `password`, and database name as needed for your PostgreSQL setup. If you're using Docker, use `postgres://gator:postgres@localhost:5433/gator`. The `current_user_name` field will be updated automatically as you log in.

## Running the Program

Once configured, you can run gator commands using:

```bash
gator <command> [args...]
```

### Available Commands

Here are some of the key commands you can use:

- **`register <username>`** - Create a new user account
- **`login <username>`** - Log in as an existing user
- **`addfeed <feed_url> <feed_name>`** - Add an RSS feed to track
- **`feeds`** - List all available feeds
- **`follow <feed_name>`** - Subscribe to a feed
- **`unfollow <feed_name>`** - Unsubscribe from a feed
- **`following`** - List feeds you're currently following
- **`browse [limit]`** - View posts from feeds you follow
- **`agg <duration>`** - Continuously aggregate feeds (e.g., `agg 1m` for 1 minute interval)
- **`users`** - List all registered users
- **`reset`** - Clear all data from the database

## Example Workflow

```bash
# Register a new user
gator register john

# Log in
gator login john

# Add a feed
gator addfeed "https://example.com/feed.xml" "Example Feed"

# Follow the feed
gator follow "Example Feed"

# Browse posts
gator browse 10

# Start aggregating feeds every 30 seconds
gator agg 30s
```

## Development

To build gator from source:

```bash
git clone https://github.com/LuisCabantac/gator.git
cd gator
go build
```

This will create a `gator` executable in the current directory.
