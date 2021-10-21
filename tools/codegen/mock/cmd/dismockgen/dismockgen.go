// Package main provides the mockgen command that automatically generates
// mocked api actions by analyzing arikawa's api source files.
//
// Caveats
//
// The code generated should not be blindly trusted and all new methods
// generated should be carefully examined, since this tool is by no means able
// to automatically generate mocks for all endpoints.
//
// Particularly, methods that manually manipulate data passed into it, cannot
// be automatically generated, without further configuration.
// Furthermore, paginating methods will also not be successfully generated.
// Such methods should be added to the config with their exclude option set to
// true.
// Some methods may be able to be mocked after adding the necessary information
// to the config.
// Refer to the documentation inside the config for more information.
//
// Despite these caveats this tool is by no means useless.
// In fact most mocks CAN be automatically generated.
// Some may require additional information in config, but still the process of
// mocking an api action is much simplified and less error-prone.
//
// Furthermore, this tool is a great aid in spotting recently added methods,
// for which no mocks exist yet.
// A simple diff between an old and a new version of generated_mocks is
// sufficient, since the order of mocks will remain the same for subsequent
// uses of mockgen.
package main

import (
	"context"
	"log"
	"strings"

	"github.com/mavolin/dismock/tools/codegen/mock/pkg/gocmd"
	"github.com/mavolin/dismock/tools/codegen/mock/pkg/mockgen"
	"github.com/mavolin/dismock/tools/codegen/mock/pkg/parser"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	apiPath, err := apiPath()
	if err != nil {
		return err
	}

	p, err := parser.New(apiPath)
	if err != nil {
		return err
	}

	log.Println("parsing api files located in", apiPath)

	pfiles, err := p.Parse()
	if err != nil {
		return err
	}

	files, err := mockgen.FilesFromParserFiles(pfiles...)
	if err != nil {
		return err
	}

	return mockgen.Generate(files)
}

func apiPath() (string, error) {
	log.Println("downloading dependencies to module cache")

	if err := gocmd.DownloadModules(context.Background(), ""); err != nil {
		return "", err
	}

	log.Println("reading go.mod file to determine the current arikawa version")

	modFile, err := gocmd.ModFile("")
	if err != nil {
		return "", err
	}

	var arikawaPath string
	var arikawaVersion string

	for _, dep := range modFile.Require {
		if strings.HasPrefix(dep.Path, "github.com/diamondburned/arikawa/") {
			log.Println("found entry; arikawa is at version", dep.Version)

			arikawaPath, arikawaVersion = dep.Path, dep.Version
			break
		}
	}

	modCachePath, err := gocmd.Env(gocmd.ModCacheEnv)
	if err != nil {
		return "", err
	}

	apiPath := modCachePath + "/" + arikawaPath + "@" + arikawaVersion + "/api"
	return apiPath, nil
}
