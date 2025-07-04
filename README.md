# omp-cli

A command line interface tool for open.mp projects.

## Features

- Initialize new open.mp projects with `ompcli init`
- Build/compile open.mp projects with `ompcli build`
- Run open.mp projects with `ompcli run`
- Automatic detection of project structure
- Support for project configuration via `project.json`
- Support for server configuration via `config.json`

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

### Initializing a Project

```
ompcli init
```

Options:
- `-n, --name`: Project name (default: current directory name)
- `-a, --author`: Project author
- `--pawncc-path`: Path to pawncc compiler (default: qawno)

This will create a basic project structure with:
- `project.json`: Project configuration file
- `config.json`: Server configuration file
- `gamemodes/` directory: Where your Pawn scripts will be stored

### Building a Project

```
ompcli build
```

Options:
- `-v, --verbose`: Enable verbose output

The build process will:
1. Compile the Pawn script specified in `main_file` using the pawncc compiler
2. Output the compiled AMX file to the path specified in `output_file`
3. Copy all necessary files to the build directory

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
  "main_file": "gamemodes/my-gamemode.pwn",
  "output_file": "gamemodes/my-gamemode.amx",
  "resources": [
    "maps/map1.json",
    "textures/logo.png"
  ],
  "plugins": [
    "plugins/streamer.dll",
    "plugins/mysql.dll"
  ],
  "server_cfg": "config.json",
  "author": "Your Name",
  "repository": "https://github.com/yourusername/my-gamemode",
  "pawncc_path": "qawno"
}
```

## Server Configuration

Open.MP uses `config.json` for server configuration:

```json
{
  "hostname": "My Open.MP Server",
  "port": 7777,
  "maxplayers": 100,
  "language": "English",
  "gamemode": "my-gamemode",
  "plugins": [
    "streamer",
    "mysql"
  ],
  "weburl": "open.mp",
  "rcon_password": "changeme",
  "password": ""
}
```

If no configuration files are found, the CLI will try to infer the configuration from the project structure.

## Requirements

- Go 1.16 or higher
- Pawn compiler (pawncc) for building projects
- open.mp server for running projects

## License

MIT 