# API

## Building

You'll need to add [templ-main](https://templ.guide/quick-start/installation) to this folder (easiest method is probably to just download the zip from github, extract, and then drop into the root of your local version of this api project).

The API uses `templ`, which requires a build step.

To generate Go code from the templ files, run: `./templ-main/cmd/templ/templ generate`

Say your filename was "users.templ", it would generate a "users_templ.go"

## Running

To run the program, you need to build the project, and then run: `go run *.go`
