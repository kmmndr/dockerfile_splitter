package main

import (
	"flag"
	"fmt"

	"github.com/kmmndr/dockerfile_splitter"
)

func main() {
	dockerfilePtr := flag.String("dockerfile", "Dockerfile", "Source Dockerfile")
	baseImagePtr := flag.String("base-image", "localhost/application", "Resulting base image")

	flag.Parse()

	fmt.Printf("Dockerfile: %s\n", *dockerfilePtr)
	fmt.Printf("Base image: %s\n", *baseImagePtr)

	dockerfile := dockerfile_splitter.NewDockerfile(*dockerfilePtr, *baseImagePtr)

	dockerfile.WriteLayers()
}
