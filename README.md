# BlogAggregator

Long Running, Persistent RSS Blog Aggregator

## Requirements

- **Go** (>= 1.20)
- **Postgres** (running instance, for persistent storage)

## Installing the gator CLI

You can install the `gator` CLI by running:

```bash
go install github.com/JonahLargen/BlogAggregator/cmd/gator@latest
```

This will place the compiled binary in your `$GOPATH/bin` or `$HOME/go/bin` directory. Make sure this path is in your `$PATH`.

## Setting up the configuration

Before running `gator`, you need to create a configuration file in your home directory:

- File: `~/.gatorconfig.json`

Example contents:

```json
{
  "db_url": "postgres://youruser:yourpassword@localhost:5432/yourdb?sslmode=disable",
  "current_user_name": ""
}
```

Replace the `db_url` with your own Postgres connection string.

## Running the program

**For production and normal usage, always use the `gator` binary:**

```bash
gator <command> [args]
```

For development purposes, you can use:

```bash
go run .
```

## Available Commands

Here are some of the commands you can run with the `gator` CLI:

- `register <username>`  
  Register a new user.
- `login <username>`  
  Log in as an existing user.
- `reset`  
  Reset the database (dangerous, wipes data!).
- `users`
  List all users.
- `agg <time_between_requests>`  
  Start aggregating feeds.
- `addfeed <name> <url>`  
  Add a new RSS feed.
- `feeds`  
  List available feeds.
- `follow <feed_url>`  
  Follow a feed.
- `following`  
  Show feeds you're following.
- `unfollow <feed_url>`  
  Unfollow a feed.
- `browse [limit]`  
  Browse your aggregated posts (optionally limit number of posts shown).

### Example usage

Register (automatically logs you in):

```bash
gator register alice
```

Add a feed (you automatically follow any feed you add):

```bash
gator addfeed "MyBlog" "https://myblog.com/rss"
```

Aggregate feeds every 60 seconds:

```bash
gator agg 60s
```

Browse your most recent posts:

```bash
gator browse 10
```

## Notes

- You must have Postgres running and accessible at the `db_url` provided in your config.
- The config file must be present and valid before running most commands.
- The Go toolchain is **not required** to run the binary after installation (`go install` produces a statically compiled binary).

MIT License © 2025 Jonah Largen