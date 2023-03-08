# Dockerfile Splitter

Command line tool to split multi-layers Dockerfiles into multiple Dockerfiles.

## Installation

### Using Golang

```
go install github.com/kmmndr/dockerfile_splitter/cmd/dockerfile_splitter@latest
```

## Usage

Running the tool will try to split Dockerfile layers into multiple Dockerfiles.

You may override defaults by specifying command line arguments

```
$ dockerfile_splitter --help
Usage of dockerfile_splitter:
  -base-image string
        Resulting base image (default "localhost/application")
  -dockerfile string
        Source Dockerfile (default "Dockerfile")
```
