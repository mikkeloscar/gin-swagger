package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"text/template"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
)

//go:generate go-bindata -pkg main -o bindata.go config.yaml templates

const (
	swaggerTmplConfig  = "config.yaml"
	defaultSwaggerPath = "./swagger.json"
)

var (
	assetDirFmt = path.Join(os.TempDir(), "gin-swagger-%d")
	config      struct {
		Application string
		SwaggerPath string
	}
)

func main() {
	kingpin.Flag("application", "Name of the application (passed directly to swagger).").
		Required().Short('A').StringVar(&config.Application)
	kingpin.Flag("spec", "the spec file to use.").
		Short('f').Default(defaultSwaggerPath).StringVar(&config.SwaggerPath)
	kingpin.Parse()

	err := run(config.Application, config.SwaggerPath)
	if err != nil {
		log.Fatalf("failed to run swagger: %s", err)
	}
}

func run(application, swagger string) error {
	tmpDir, err := writeAssets()
	if err != nil {
		return err
	}
	defer func() error {
		return os.RemoveAll(tmpDir)
	}()

	cmd := exec.Command("swagger",
		"generate",
		"server",
		"-A",
		application,
		"-f",
		swagger,
		"-C",
		path.Join(tmpDir, swaggerTmplConfig),
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func writeAssets() (string, error) {
	tmp := fmt.Sprintf(assetDirFmt, time.Now().UTC().UnixNano())
	err := RestoreAssets(tmp, "templates")
	if err != nil {
		return "", err
	}

	d, err := Asset(swaggerTmplConfig)
	if err != nil {
		return "", err
	}

	t, err := template.New(swaggerTmplConfig).Parse(string(d))
	if err != nil {
		return "", err
	}

	fd, err := os.Create(path.Join(tmp, swaggerTmplConfig))
	if err != nil {
		return "", err
	}
	defer fd.Close()

	w := bufio.NewWriter(fd)

	err = t.Execute(w, map[string]string{"TmpDir": tmp})
	if err != nil {
		return "", err
	}

	err = w.Flush()
	if err != nil {
		return "", err
	}

	return tmp, nil
}
