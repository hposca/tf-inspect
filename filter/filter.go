package filter

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/drrzmr/tf-inspect/loader/section"
	"github.com/hashicorp/hcl"
)

type forEachResourceCallback func(section section.Section, resourceName string, resourceValue map[string]interface{})
type forEachSectionCallback func(section section.Section)
type forEachSectionDataCallback func(data map[string]interface{})

type forEachSectionCallback2 func(sectionName string, sectionData interface{})

func ForEachSection(sections []section.Section, callback forEachSectionCallback) {
	for _, section := range sections {
		callback(section)
	}
}

func ForEachSection2(data interface{}, callback forEachSectionCallback2) {

	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			callback(key, value)
		}
	default:
		panic(fmt.Sprintf("unsupported type: %T", v))
	}

}

func ForEachSectionData(section section.Section, callback forEachSectionDataCallback) {
	for _, dataItem := range section.Data {
		callback(dataItem)
	}
}

type forEachSectionDataCallback2 func(key string, value interface{})

func ForEachSectionData2(sectionData interface{}, callback forEachSectionDataCallback2) {

	switch data := sectionData.(type) {
	case []map[string]interface{}:
		for _, d := range data {
			for key, value := range d {
				callback(key, value)
			}
		}
	}
}

func ForEachResource(sections []section.Section, resourceType string, callback forEachResourceCallback) {

	for _, section := range sections {
		for _, data := range section.Data {
			for key, value := range data {
				if key != resourceType {
					continue
				}
				for _, resource := range value.([]map[string]interface{}) {
					for resourceName, resourceValue := range resource {
						for _, value := range resourceValue.([]map[string]interface{}) {
							callback(section, resourceName, value)
						}
					}
				}
			}
		}
	}
}

type forEachFileCallback func(filename string, data map[string]interface{})

func ForEachFile(rootDir string, extension string, callback forEachFileCallback) {

	err := filepath.Walk(rootDir, func(fn string, info os.FileInfo, err error) error {
		if ext := filepath.Ext(fn); ext != extension {
			return nil
		}

		var buf []byte
		if buf, err = ioutil.ReadFile(fn); err != nil {
			return err
		}

		var v interface{}
		if err = hcl.Unmarshal(buf, &v); err != nil {
			return err
		}

		switch data := v.(type) {
		case map[string]interface{}:
			callback(fn, data)

		default:
			panic(fmt.Sprintf("unsupported type: %T", data))
		}

		return nil

	})

	if err != nil {
		panic(err)
	}
}
