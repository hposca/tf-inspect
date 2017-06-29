package cmd

import (
	"fmt"
	"os"

	docopt "github.com/docopt/docopt-go"
	"github.com/drrzmr/tf-inspect/cmd/module"
	"github.com/drrzmr/tf-inspect/cmd/sg"

	"github.com/drrzmr/tf-inspect/loader"
	"github.com/kr/pretty"
)

// Parse the command line
func Parse() int {

	doc := `Tf Inspector.

Usage:
	hcl-inspector sg --list-names <path>
	hcl-inspector module <path>
	hcl-inspector test <path>
	hcl-inspector -h | --help
	hcl-inspector --version

Options:
	-h --help     Show this screen.
	--version     Show version.`

	arguments, err := docopt.Parse(doc, nil, true, "Tf Inspector 0.0.1", false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could initialize command parser\n")
		return 1
	}

	var rootDir string
	if rootDir, err = getRootDir(arguments); err != nil {
		fmt.Fprintf(os.Stdout, err.Error())
		return -1
	}

	if arguments["sg"].(bool) {
		resourceSection, err := loader.LoadResources(rootDir, []string{
			"aws_security_group",
			"aws_security_group_rule",
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "count not load resources: %v", err)
			return 4
		}

		return sg.Do(resourceSection, arguments)
	}

	if arguments["module"].(bool) {
		return module.Do(rootDir, loader.LoadDir(rootDir))
	}

	if arguments["test"].(bool) {
		files := loader.LoadDir(rootDir)
		count := 0

		for k, v := range files {

			if module, ok := v["module"]; !ok {
				continue
			} else {
				pretty.Println(k)
				for _, m := range module {
					pretty.Println(m)
					count++
				}

			}
		}
		pretty.Println("count: ", count)
		return 0
	}

	fmt.Fprintf(os.Stderr, "could undestand what do you want\n")
	return 5
}

func getRootDir(arguments map[string]interface{}) (string, error) {

	root, ok := arguments["<path>"].(string)
	if !ok {
		return "", fmt.Errorf("path must be a string")
	}

	if root == "." {
		var err error
		if root, err = os.Getwd(); err != nil {
			return "", err
		}
	}

	return root, nil
}
