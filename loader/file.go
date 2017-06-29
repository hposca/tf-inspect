package loader

type Files map[string]map[string]map[string][]interface{}

type ForEachModuleCallback func(fn string, name string, module map[string]interface{})

func (f Files) ForEachModules(callback ForEachModuleCallback) int {

	count := 0

	for fn, section := range f {
		if module, ok := section["module"]; !ok {
			continue
		} else {
			for name, m := range module {
				if len(m) > 1 {
					panic("unsupported multiple modules")
				}

				lmst := m[0].([]map[string]interface{})
				if len(lmst) > 1 {
					panic("buááá")
				}

				count++
				callback(fn, name, lmst[0])
			}
		}
	}

	return count
}
