package module

import (
	"path"
	"strings"

	"github.com/awalterschulze/gographviz"
	"github.com/drrzmr/tf-inspect/loader"
	"github.com/kr/pretty"
)

// Do exec module command
func Do(rootDir string, files loader.Files) int {

	graph := gographviz.NewGraph()
	//graph := map[string][]string{}

	files.ForEachModules(func(fn string, name string, module map[string]interface{}) {

		dir := path.Dir(fn)
		join := path.Join(dir, module["source"].(string))

		moduleName := strings.Replace(dir[len(rootDir):], "-", "_", -1)
		moduleName = strings.Replace(moduleName, "/", "_", -1)
		usedModuleName := strings.Replace(join[len(rootDir):], "-", "_", -1)
		usedModuleName = strings.Replace(usedModuleName, "/", "_", -1)

		//fmt.Printf("module: %s -> %s\n", moduleName, usedModuleName)

		/*
			if _, ok := graph[moduleName]; !ok {
				graph[moduleName] = []string{}
			}
		*/
		//graph[moduleName] = append(graph[moduleName], usedModuleName)

		//graph.AddNode(moduleName, usedModuleName, attrs)
		if !graph.IsNode(moduleName) {
			graph.AddNode("tf", moduleName, nil)
		}

		if !graph.IsNode(usedModuleName) {
			graph.AddNode("tf", usedModuleName, nil)
		}

		graph.AddEdge(moduleName, usedModuleName, true, nil)

	})

	//pretty.Println(graph)

	ast, err := graph.WriteAst()
	if err != nil {
		panic(err)
	}

	//s := generateGraph(graph)

	pretty.Println(ast.String())
	return 0
}

/*
func generateGraph(in map[string][]string) string {

	out := gographviz.NewGraph()

	for k := range in {
		out.AddNode("tf", k, nil)
	}

	for k, v := range in {
		for _, e := range v {
			pretty.Println(k, "->", e)
			out.AddEdge(k, e, true, nil)
		}
	}

	return out.String()
}
*/
