package dockerfile_splitter

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"

	"golang.org/x/exp/slices"
)

type DockerLayer struct {
	name       string
	dockerfile Dockerfile
}

type Dockerfile struct {
	filename       string
	layers         []DockerLayer
	lines          []string
	prependedLines []string
	baseImage      string
}

func NewDockerfile(filename string, baseImage string) *Dockerfile {
	dockerfile := Dockerfile{filename: filename, baseImage: baseImage}
	dockerfile.readContent()

	return &dockerfile
}

func (dd *Dockerfile) prependLine(line string) {
	d := &dd.currentLayer().dockerfile

	if slices.Contains(d.prependedLines, line) {
		return
	}

	d.prependedLines = append(d.prependedLines, line)

	d.lines = append(d.lines, "")
	copy(d.lines[1:], d.lines)
	d.lines[0] = line
}

func (dd *Dockerfile) appendLine(line string) {
	d := &dd.currentLayer().dockerfile
	d.lines = append(d.lines, line)
}

func (d *Dockerfile) addLayer(name string) {
	filename := fmt.Sprintf("%s-%s", d.filename, name)
	baseImage := fmt.Sprintf("%s-%s", d.baseImage, name)
	dockerLayer := DockerLayer{
		name:       name,
		dockerfile: Dockerfile{filename: filename, baseImage: baseImage},
	}
	d.layers = append(d.layers, dockerLayer)
}

func (d *Dockerfile) isLayer(name string) bool {
	for _, dl := range d.layers {
		if dl.name == name {
			return true
		}
	}
	return false
}

func (d *Dockerfile) currentLayer() *DockerLayer {
	if len(d.layers) == 0 {
		fmt.Fprintf(os.Stderr, "Source file has no layers\n")
		os.Exit(1)
	}

	return &d.layers[len(d.layers)-1]
}

func (d *Dockerfile) readContent() {
	file, err := os.Open(d.filename)
	if err != nil {
		log.Fatalf("failed to open")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()

		re := regexp.MustCompile(`^[[:space:]]*FROM ([^ ]*)( as (.*))*`)
		matches := re.FindStringSubmatch(line)
		if len(matches) > 0 {
			layer_source := matches[1]
			layer_name := matches[3]

			if len(layer_name) == 0 {
				layer_name = "final"
			}
			d.addLayer(layer_name)

			if d.isLayer(layer_source) {
				line = fmt.Sprintf("FROM %s-%s", d.baseImage, layer_source)
			} else {
				line = fmt.Sprintf("FROM %s", layer_source)
			}
		}

		re = regexp.MustCompile(`^[[:space:]]*COPY --from=([^ ]*) .*`)
		matches = re.FindStringSubmatch(line)
		if len(matches) > 0 {
			layer_source := matches[1]

			prepended_line := fmt.Sprintf("FROM %s-%s as %s", d.baseImage, layer_source, layer_source)

			d.prependLine(prepended_line)
		}

		d.appendLine(line)
	}
}

func (d *Dockerfile) PrintLayers() {
	for _, dl := range d.layers {
		dl.dockerfile.Print()
	}
}

func (d *Dockerfile) Print() {
	fmt.Printf("# %s\n", d.filename)
	fmt.Printf("# %s\n", d.baseImage)
	for _, line := range d.lines {
		fmt.Println(line)
	}
}

func (d *Dockerfile) WriteLayers() {
	for _, dl := range d.layers {
		fmt.Printf("Writting file %s\n", dl.dockerfile.filename)
		dl.dockerfile.Write()
	}
}

func (d *Dockerfile) Write() {
	file, err := os.Create(d.filename)
	if err != nil {
		log.Fatalf("failed to open")
	}
	defer file.Close()

	file.WriteString(fmt.Sprintf("# %s\n\n", d.baseImage))
	for _, line := range d.lines {
		file.WriteString(line)
		file.WriteString("\n")
	}
}