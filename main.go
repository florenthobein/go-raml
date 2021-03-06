package main

import (
	"errors"
	"os"

	"github.com/Jumpscale/go-raml/commands"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

//Version define software version
var Version = "0.1-Dev"

//ApplicationName is the name of the application
var ApplicationName = "RAML code generation toolset"

var (
	serverCommand = &commands.ServerCommand{}
	clientCommand = &commands.ClientCommand{}
	capnpCommand  = &commands.CapnpCommand{}
	specCommand   = &commands.SpecCommand{}
	docsCommand   = &commands.DocsCommand{}
)

func main() {
	app := cli.NewApp()
	app.Name = ApplicationName
	app.Version = Version
	app.Usage = "Using a RAML specification, generate server and client code or a RAML specification from go code."

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	var debugLogging bool
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug, d",
			Usage:       "Enable debug logging",
			Destination: &debugLogging,
		},
	}
	app.Before = func(c *cli.Context) error {
		if debugLogging {
			log.SetLevel(log.DebugLevel)
			log.Debug("Debug logging enabled")
			log.Debug(ApplicationName, "-", Version)
		}
		return nil
	}
	app.Commands = []cli.Command{
		{
			Name:  "server",
			Usage: "Generate a server according to a RAML specification",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "language, l",
					Value:       "go",
					Usage:       "Language to construct a server for",
					Destination: &serverCommand.Language,
				},
				cli.StringFlag{
					Name:        "kind",
					Value:       "",
					Usage:       "Kind of server to generate (, sanic)",
					Destination: &serverCommand.Kind,
				},
				cli.StringFlag{
					Name:        "dir",
					Value:       ".",
					Usage:       "target directory",
					Destination: &serverCommand.Dir,
				},
				cli.StringFlag{
					Name:        "ramlfile",
					Value:       ".",
					Usage:       "Source raml file",
					Destination: &serverCommand.RamlFile,
				},
				cli.StringFlag{
					Name:        "package",
					Value:       "main",
					Usage:       "package name",
					Destination: &serverCommand.PackageName,
				},
				cli.BoolFlag{
					Name:        "no-main",
					Usage:       "Do not generate a main.go file",
					Destination: &serverCommand.NoMainGeneration,
				},
				cli.BoolFlag{
					Name:        "no-apidocs",
					Usage:       "Do not generate API Docs in /apidocs/ endpoint",
					Destination: &serverCommand.NoAPIDocs,
				},
				cli.StringFlag{
					Name:        "import-path",
					Value:       "",
					Usage:       "import path of the generated code. Set automatically if target dir under $GOPATH",
					Destination: &serverCommand.ImportPath,
				},
				cli.BoolFlag{
					Name:        "api-file-per-method",
					Usage:       "Generate one API implementation file per method (only for Go)",
					Destination: &serverCommand.APIFilePerMethod,
				},
				cli.StringFlag{
					Name:        "lib-root-urls",
					Usage:       "Array of libraries root URLs",
					Destination: &serverCommand.LibRootURLs,
				},
			},
			Action: func(c *cli.Context) {
				if err := serverCommand.Execute(); err != nil {
					log.Error(err)
				}
			},
		},
		{
			Name:  "client",
			Usage: "Create a client for a RAML specification",

			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "language, l",
					Value:       "go",
					Usage:       "Language to construct a client for",
					Destination: &clientCommand.Language,
				},
				cli.StringFlag{
					Name:        "dir",
					Value:       ".",
					Usage:       "target directory",
					Destination: &clientCommand.Dir,
				},
				cli.StringFlag{
					Name:        "ramlfile",
					Value:       ".",
					Usage:       "Source raml file",
					Destination: &clientCommand.RamlFile,
				},
				cli.StringFlag{
					Name:        "package",
					Value:       "client",
					Usage:       "package name",
					Destination: &clientCommand.PackageName,
				},
				cli.StringFlag{
					Name:        "import-path",
					Value:       "",
					Usage:       "import path of the generated code",
					Destination: &clientCommand.ImportPath,
				},
				cli.StringFlag{
					Name:        "kind",
					Value:       "requests",
					Usage:       "Kind of client to generate (requests,aiohttp)",
					Destination: &clientCommand.Kind,
				},
				cli.StringFlag{
					Name:        "lib-root-urls",
					Usage:       "Array of libraries root URLs",
					Destination: &clientCommand.LibRootURLs,
				},
			},
			Action: func(c *cli.Context) {
				if err := clientCommand.Execute(); err != nil {
					log.Error(err)
				}
			},
		},
		{
			Name:  "docs",
			Usage: "Generate API docs for a RAML specifications",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "format",
					Value:       "markdown",
					Usage:       "API documentation format, only markdown is supported at the moment.",
					Destination: &docsCommand.Format,
				},
				cli.StringFlag{
					Name:        "ramlfile",
					Value:       ".",
					Usage:       "Source raml file",
					Destination: &docsCommand.RamlFile,
				},
				cli.StringFlag{
					Name:        "output",
					Usage:       "Destination doc file",
					Destination: &docsCommand.OutputFile,
				},
			},
			Action: func(c *cli.Context) {
				if err := docsCommand.Execute(); err != nil {
					log.Error(err)
				}
			},
		},
		{
			Name:  "spec",
			Usage: "Generate a RAML specification from a go server",
			Action: func(c *cli.Context) {
				err := errors.New("Not implemented, check the roadmap")
				log.Error(err)
			},
		},
		{
			Name:  "capnp",
			Usage: "Create capnpn models",

			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "dir",
					Value:       ".",
					Usage:       "target directory",
					Destination: &capnpCommand.Dir,
				},
				cli.StringFlag{
					Name:        "ramlfile",
					Value:       ".",
					Usage:       "Source raml file",
					Destination: &capnpCommand.RAMLFile,
				},
				cli.StringFlag{
					Name:        "language, l",
					Value:       "plain",
					Usage:       "Language to construct capnpn models for",
					Destination: &capnpCommand.Language,
				},
				cli.StringFlag{
					Name:        "package",
					Value:       "main",
					Usage:       "package name - only for Go",
					Destination: &capnpCommand.Package,
				},
			},
			Action: func(c *cli.Context) {
				if err := capnpCommand.Execute(); err != nil {
					log.Error(err)
				}
			},
		}, {
			Name:  "spec",
			Usage: "Generate a RAML specification from a go server",
			Action: func(c *cli.Context) {
				err := errors.New("Not implemented, check the roadmap")
				log.Error(err)
			},
		},
	}

	app.Action = func(c *cli.Context) {
		cli.ShowAppHelp(c)
	}

	app.Run(os.Args)
}
