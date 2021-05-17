package main

import (
	"embed"
	"log"
	"path"

	"github.com/alecthomas/kingpin"
	"github.com/go-openapi/analysis"
	"github.com/go-swagger/go-swagger/generator"
)

//go:embed templates/*
var templates embed.FS

//go:generate go-bindata -pkg main -o bindata.go templates

const (
	defaultSwaggerPath = "./swagger.json"
)

var (
	config struct {
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

func run(application, specPath string) error {
	assets, err := templates.ReadDir("templates")
	if err != nil {
		return err
	}

	for _, asset := range assets {
		data, err := templates.ReadFile(path.Join("templates", asset.Name()))
		if err != nil {
			return err
		}

		err = generator.AddFile(path.Join("templates", asset.Name()), string(data))
		if err != nil {
			return err
		}
	}

	opts := &generator.GenOpts{
		Spec:              specPath,
		Target:            "./",
		APIPackage:        "operations",
		ModelPackage:      "models",
		ServerPackage:     "restapi",
		ClientPackage:     "client",
		Principal:         "",
		DefaultScheme:     "http",
		IncludeModel:      true,
		IncludeValidator:  true,
		IncludeHandler:    true,
		IncludeParameters: true,
		IncludeResponses:  true,
		IncludeURLBuilder: true,
		IncludeMain:       true,
		IncludeSupport:    true,
		ValidateSpec:      true,
		FlattenOpts: &analysis.FlattenOpts{
			Minimal:      true,
			Verbose:      true,
			RemoveUnused: false,
			Expand:       false,
		},
		ExcludeSpec:       false,
		TemplateDir:       "",
		DumpData:          false,
		Models:            nil,
		Operations:        nil,
		Tags:              nil,
		Name:              application,
		FlagStrategy:      "go-flags",
		CompatibilityMode: "modern",
		ExistingModels:    "",
		Copyright:         "",
		Sections: generator.SectionOpts{
			Application: []generator.TemplateOpts{
				{
					Name:       "configure",
					Source:     "templates/config.gotmpl",
					Target:     "{{ joinFilePath .Target .ServerPackage }}",
					FileName:   "config.go",
					SkipExists: false,
					SkipFormat: false,
				},
				{
					Name:       "embedded_spec",
					Source:     "asset:swaggerJsonEmbed",
					Target:     "{{ joinFilePath .Target .ServerPackage }}",
					FileName:   "embedded_spec.go",
					SkipExists: false,
					SkipFormat: false,
				},
				{
					Name:       "server",
					Source:     "templates/api.gotmpl",
					Target:     "{{ joinFilePath .Target .ServerPackage }}",
					FileName:   "api.go",
					SkipExists: false,
					SkipFormat: false,
				},
			},
			Operations: []generator.TemplateOpts{
				{
					Name:       "parameters",
					Source:     "templates/parameter.gotmpl",
					Target:     "{{ if gt (len .Tags) 0 }}{{ joinFilePath .Target .ServerPackage .APIPackage .Package  }}{{ else }}{{ joinFilePath .Target .ServerPackage .Package  }}{{ end }}",
					FileName:   "{{ (snakize (pascalize .Name)) }}_parameters.go",
					SkipExists: false,
					SkipFormat: false,
				},
			},
			Models: []generator.TemplateOpts{
				{
					Name:       "definition",
					Source:     "asset:model",
					Target:     "{{ joinFilePath .Target .ModelPackage }}",
					FileName:   "{{ (snakize (pascalize .Name)) }}.go",
					SkipExists: false,
					SkipFormat: false,
				},
			},
		},
	}

	err = opts.EnsureDefaults()
	if err != nil {
		return err
	}

	err = generator.GenerateServer(application, nil, nil, opts)
	if err != nil {
		return err
	}

	return nil
}
