# gok8slab

A CLI interface written in GO, to simulate k8s Problems in a ctf (Capture The Flag) style.

## Structure

``` tree
gok8slab/
│── cmd/                    # CLI commands
│   ├── root.go             # Root command
│   ├── list.go             # List courses
│   ├── start.go            # Start course
│   ├── stop.go             # Stop course
│   ├── checkflag.go        # Check user's flag
│── internal/               # Internal packages
│   ├── config/             # Configuration handling
│   ├── k8s/                # Kubernetes/OpenShift interactions
│   ├── git/                # Git operations
│   ├── course/             # Course handling
│   │   ├── course.go
│── courses/                # Course definitions
│── README.md               # Documentation
│── main.go                 # Entry point
│── go.mod                  # Go modules file

```

## Features

- Retrieve lessons/courses from a git repository.
- Deploy Kubernetes resources directly from YAML files.
- Commands to pull new courses, start and stop labs, and retrieve information.
- Logging and tracing capabilities.
- Dynamically updates bash/shell prompts with current lab information.

## Quickstart

`curl -sL https://raw.githubusercontent.com/ed1us3r/gok8slab/refs/heads/main/hack.sh | bash -s -- --install`

## Usage

```bash
# Pull courses
gok8slab pull

# Start a lab
gok8slab start

# Stop a lab
gok8slab stop

# Get info about the current lab
gok8slab info
