# omp-cli

A command line interface tool for open.mp projects.

## Features

- Build/compile open.mp projects with `ompcli build`
- Run open.mp projects with `ompcli run`
- Automatic detection of project structure
- Support for project configuration via `project.json`

## Installation

### From Source

1. Clone the repository:
```
git clone https://github.com/weltschmerzie/omp-cli.git
```

2. Build the project:
```
cd omp-cli
go build -o ompcli
```

3. Add the binary to your PATH or move it to a directory in your PATH:
```
# Linux/macOS
mv ompcli /usr/local/bin/

# Windows
# Move to a directory in your PATH
```

### Using Go Install

```
go install github.com/weltschmerzie/omp-cli@latest
```

## Usage

### Building a Project

```
ompcli build
```

Options:
- `-v, --verbose`: Enable verbose output

### Running a Project

```
ompcli run
```

Options:
- `-d, --debug`: Enable debug mode
- `-p, --port`: Port to run the server on (default: 7777)

## Project Configuration

You can configure your open.mp project using a `project.json` file:

```json
{
  "name": "my-gamemode",
  "version": "1.0.0",
  "main_file": "gamemode.pwn",
  "resources": [
    "maps/map1.json",
    "textures/logo.png"
  ],
  "plugins": [
    "plugins/streamer.dll",
    "plugins/mysql.dll"
  ],
  "server_cfg": "server.cfg",
  "author": "Your Name",
  "repository": "https://github.com/weltschmerzie/my-gamemode"
}
```

If no `project.json` file is found, the CLI will try to infer the configuration from the project structure.

## Requirements

- Go 1.16 or higher
- Pawn compiler (pawncc) for building projects
- open.mp server for running projects

## License

MIT 