# Enhanced github CLI
ghe grabs the github CLI and enhances the experience, no more selecting 700 stupid questions! 

> [!IMPORTANT]
> Very WIP early stages


## Install
- Need to have go installed

#### Go
You can use the `setup.sh` with the commands below
````bash
curl -O https://raw.githubusercontent.com/GMkonan/ghe/main/setup.sh

chmod +x setup.sh && ./setup.sh
``````

#### Manual
Clone the repo and run `go build` and `go install` 

> [!NOTE]
> Both methods will install the exec in your `GOBIN` so you need to have it in `PATH`

### ideas
- default to normal gh commands when user tries them
- ask for committing and stuff when creating repo
- add command for `gh auth` (awful amount of questions)

For now way to run it is `go build` and then `go install`

### REF
- https://cli.github.com/manual/gh_repo_create
