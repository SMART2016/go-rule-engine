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
- 