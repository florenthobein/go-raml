package golang

import (
	"fmt"
	"path/filepath"

	log "github.com/Sirupsen/logrus"

	"github.com/Jumpscale/go-raml/codegen/commons"
	"github.com/Jumpscale/go-raml/codegen/resource"
	"github.com/Jumpscale/go-raml/raml"
)

// global variables
// it is needed for libraries support
var (
	// root import path
	globRootImportPath string

	// global value of API definition
	globAPIDef *raml.APIDefinition

	globLibRootURLs []string
)

// Server represents a Go server
type Server struct {
	apiDef           *raml.APIDefinition
	ResourcesDef     []resource.ResourceInterface
	PackageName      string // Name of the package this server resides in
	APIDocsDir       string // apidocs directory. apidocs won't be generated if it is empty
	withMain         bool   // true if we need to generate main file
	RootImportPath   string
	APIFilePerMethod bool     // true if we want to generate one API file per API method
	TargetDir        string   // root directory of the generated code
	libsRootURLs     []string // root URLs of the libraries
}

// NewServer creates a new Golang server
func NewServer(apiDef *raml.APIDefinition, packageName, apiDocsDir, rootImportPath string,
	withMain, apiFilePerMethod bool, targetDir string, libsRootURLs []string) Server {
	// global variables
	rootImportPath = setRootImportPath(rootImportPath, targetDir)
	globAPIDef = apiDef
	globRootImportPath = rootImportPath
	globLibRootURLs = libsRootURLs

	return Server{
		apiDef:           apiDef,
		PackageName:      packageName,
		APIDocsDir:       apiDocsDir,
		withMain:         withMain,
		RootImportPath:   rootImportPath,
		APIFilePerMethod: apiFilePerMethod,
		TargetDir:        targetDir,
		libsRootURLs:     libsRootURLs,
	}
}

// Generate generates all Go server files
func (gs Server) Generate() error {
	if gs.RootImportPath == "" {
		return fmt.Errorf("invalid import path = empty. please set --import-path or set target dir under gopath")
	}
	// helper package
	gh := goramlHelper{
		rootImportPath: gs.RootImportPath,
		packageName:    "goraml",
		packageDir:     "goraml",
	}
	if err := gh.generate(gs.TargetDir); err != nil {
		return err
	}

	if err := generateAllStructs(gs.apiDef, gs.TargetDir, gs.PackageName); err != nil {
		return err
	}

	// security scheme
	if err := generateSecurity(gs.apiDef.SecuritySchemes, gs.TargetDir, gs.PackageName); err != nil {
		log.Errorf("failed to generate security scheme:%v", err)
		return err
	}

	// genereate resources
	rds, err := gs.generateServerResources(gs.TargetDir)
	if err != nil {
		return err
	}
	gs.ResourcesDef = rds

	// libraries
	if err := generateLibraries(gs.apiDef.Libraries, gs.TargetDir, gs.libsRootURLs); err != nil {
		return err
	}

	// generate main
	if gs.withMain {
		// HTML front page
		if err := commons.GenerateFile(gs, "./templates/index.html.tmpl", "index.html", filepath.Join(gs.TargetDir, "index.html"), false); err != nil {
			return err
		}
		// main file
		return commons.GenerateFile(gs, "./templates/server_main_go.tmpl", "server_main_go", filepath.Join(gs.TargetDir, "main.go"), true)
	}

	return nil
}

// Title returns title of this server
func (gs Server) Title() string {
	return gs.apiDef.Title
}
