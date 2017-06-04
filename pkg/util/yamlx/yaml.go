package yamlx

import "github.com/go-yaml/yaml"

func ToYaml(param interface{}) string {
	data, _ := yaml.Marshal(&param)
	return string(data)
}
