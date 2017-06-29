package sg

import (
	"github.com/drrzmr/tf-inspect/filter"
	"github.com/drrzmr/tf-inspect/loader/section"
	"github.com/kr/pretty"
)

type securityGroup struct {
	name       string
	vpcID      string
	namePrefix string
	tags       map[string]string
}

// Do exec sg command
func Do(sections []section.Section, arguments map[string]interface{}) int {

	sg := map[string]securityGroup{}

	filter.ForEachResource(sections, "aws_security_group", func(s section.Section, name string, value map[string]interface{}) {

		pretty.Println(s.Namespace)
		pretty.Println(s.Filename)
		pretty.Println(s.Name)
		pretty.Println(s)
		pretty.Println(name)
		pretty.Println(value)
	})

	pretty.Println(sg)

	if arguments["--list-names"].(bool) {
		//	return securityGroupsListNames(tags)
	}

	return 0
}
