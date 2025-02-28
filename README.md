# go-rule-engine
An Extensible go rule engine written on top of grule


## Setting up direnv
- open your shell profile for mac its zshrc in the home folder
  - `sudo vim ~/.zshrc`
- Add below line in the file
  - `eval "$(direnv hook zsh)"`
- Reload the shell
  - `source ~/.zshrc`
- Navigte to the current project and run below command to enable direnv
  - `direnv allow`

## To generate the store code to fetch , save and remove processed events with sqlc
- `sqlc generate`
  - This will generate all code in `store` package from the schema and store sql file.

## Startup a local Postgres instance on docker
- Run below command from /go-rule-engine folder directly
```
  docker run --hostname=3482db53b646 \
  --name=postgres \
  --env PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/lib/postgresql/17/bin \
  --env GOSU_VERSION=1.17 \
  --env LANG=en_US.utf8 \
  --env PG_MAJOR=17 \
  --env PG_VERSION=17.4-1.pgdg120+2 \
  --env PGDATA=/var/lib/postgresql/data \
  --env POSTGRES_PASSWORD=dbpassword \
  --env POSTGRES_USER=dbuser \
  --env POSTGRES_DB=rule_engine \
  --volume $(pwd)/store/sqlc/store_schema.sql:/docker-entrypoint-initdb.d/store_schema.sql \
  --volume "$(pwd)/scripts/init.sh:/docker-entrypoint-initdb.d/init.sh" \
  --volume $(pwd)/pgdata:/var/lib/postgresql/data \
  --network bridge \
  -p 5432:5432 \
  --restart=no \
  --runtime=runc \
  -d postgres:latest
```

## TODO's
- Handle concurrency issue on state store 
  - Check if write fails how we can handle same using message bus commits
- How to handle event schema which would be needed for rule evaluation , 
     currently how does grule handles that.
