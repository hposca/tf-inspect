package module

import (
	"fmt"
	"path"
	"strings"

	"github.com/awalterschulze/gographviz"
	"github.com/drrzmr/tf-inspect/loader"
)

// Do exec module command
func Do(rootDir string, files loader.Files) int {

	graph := gographviz.NewGraph()
	graph.SetDir(true)

	files.ForEachModules(func(fn string, name string, module map[string]interface{}) {

		dir := path.Dir(fn)
		join := path.Join(dir, module["source"].(string))

		rootDirLen := len(rootDir)
		if strings.HasSuffix(rootDir, "/") {
			fmt.Println(path.Dir(rootDir))
			rootDirLen--
		}

		moduleName := fmt.Sprintf("\"%s\"", dir[rootDirLen+1:])
		usedModuleName := fmt.Sprintf("\"%s\"", join[rootDirLen+1:])

		//fmt.Printf("rootDir: %s\ndir: %s\njoin: %s\n%s -> %s\n\n", rootDir, dir, join, moduleName, usedModuleName)

		//moduleName = strings.Replace(moduleName, "-", "_", -1)
		//moduleName = strings.Replace(moduleName, "/", "_", -1)
		//usedModuleName = strings.Replace(usedModuleName, "-", "_", -1)
		//usedModuleName = strings.Replace(usedModuleName, "/", "_", -1)

		if !graph.IsNode(moduleName) {
			graph.AddNode("tf", moduleName, nil)
		}

		if !graph.IsNode(usedModuleName) {
			graph.AddNode("tf", usedModuleName, nil)
		}

		graph.AddEdge(moduleName, usedModuleName, true, nil)

	})

	fmt.Print(graph.String())
	return 0
}
