# Blog Aggregator

## Ref

`postgres@15` and `go1.23.4` were used to create the blog aggregator

## Installation

you can use `go install <repo_link>` to install the blog aggregator

## Usage

Set up a config file in the root directory of your machine

```json
{
  "db_url": "postgres://postgres:postgres@localhost:5432/blog_agg?sslmode=disable"
}
```

Start by using adding a user `register <name>`, `login <name>` will change the user
if there are multiple users in the database. Names are unique, case sensitive.
`addFeed <name> <url>` will add a feed to the database.
`following` will show the list of feeds that the user is following.
`follow <url>` will follow a feed.
`unfollow <url>` will unfollow a feed.
`agg <time:1m30s>` will start the aggregator with the time between requests.
`browse <limit>` will show the recent `limit` posts.
