package loader

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/drrzmr/tf-inspect/filter"
	"github.com/drrzmr/tf-inspect/loader/section"
	"github.com/drrzmr/tf-inspect/util"
	"github.com/kr/pretty"

	"github.com/hashicorp/hcl"
)

// LoadResources read all resources on .tf files from rootDir recursively
func LoadResources(rootDir string, subFilter []string) ([]section.Section, error) {
	return load(rootDir, ".tf", "resource", subFilter)
}

func getSectionData(v interface{}) map[string]interface{} {

	switch sectionData := v.(type) {
	case []map[string]interface{}:
		ret := map[string]interface{}{}

		for _, elem := range sectionData {
			for key, value := range elem {
				ret[key] = value
			}
		}
		pretty.Println(ret)
		return ret

	default:
		panic("invalid type")
	}

}

func LoadDir(rootDir string) Files {

	files := Files{}
	filter.ForEachFile(rootDir, ".tf", func(fn string, data map[string]interface{}) {
		files[fn] = map[string]map[string][]interface{}{}
		filter.ForEachSection2(data, func(sectionName string, sectionData interface{}) {
			if _, ok := files[fn][sectionName]; !ok {
				files[fn][sectionName] = map[string][]interface{}{}
			}
			filter.ForEachSectionData2(sectionData, func(key string, value interface{}) {
				files[fn][sectionName][key] = append(files[fn][sectionName][key], value)
			})
		})
	})

	return files
}

func load(rootDir string, extension string, sectionName string, filters []string) ([]section.Section, error) {

	fileNames, err := listFileNames(rootDir, ".tf")
	if err != nil {
		return nil, err
	}

	sections, err := loadSection(rootDir, fileNames, sectionName)
	if err != nil {
		return nil, err
	}

	if filters == nil {
		return sections, nil
	}

	filteredSections := sectionFilter(sections, filters)
	return filteredSections, nil
}

func listFileNames(baseDir string, extension string) ([]string, error) {

	list := []string{}

	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if ext := filepath.Ext(path); ext != extension {
			return nil
		}
		list = append(list, path)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return list, nil
}

func loadSection(rootDir string, fileNames []string, sectionName string) ([]section.Section, error) {

	sections := []section.Section{}

	for _, fn := range fileNames {
		var err error
		var buf []byte

		if buf, err = ioutil.ReadFile(fn); err != nil {
			return nil, err
		}

		var v interface{}
		if err = hcl.Unmarshal(buf, &v); err != nil {
			return nil, err
		}

		switch dict := v.(type) {
		case map[string]interface{}:
			ivalue, ok := dict[sectionName]
			if !ok {
				continue
			}

			switch data := ivalue.(type) {
			case []map[string]interface{}:
				sections = append(sections, section.Section{
					Namespace: path.Dir(fn[len(rootDir)+1:]),
					Filename:  fn,
					Data:      data,
					Name:      sectionName,
				})
			default:
				return nil, fmt.Errorf("could not parse section: %s", sectionName)
			}

		default:
			return nil, fmt.Errorf("could not load file: %s", fn)
		}
	}
	return sections, nil
}

func sectionFilter(sections []section.Section, filters []string) []section.Section {

	filteredSections := []section.Section{}

	filter.ForEachSection(sections, func(s section.Section) {

		data := []map[string]interface{}{}
		filter.ForEachSectionData(s, func(dataItem map[string]interface{}) {
			for key := range dataItem {
				if !util.InList(key, filters) {
					continue
				}
				data = append(data, dataItem)
			}
		})
		if len(data) > 0 {
			filteredSections = append(filteredSections, section.Section{
				Data:      data,
				Filename:  s.Filename,
				Name:      s.Name,
				Namespace: s.Namespace,
			})
		}
	})

	return filteredSections
}
