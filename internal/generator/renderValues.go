package generator

import (
	"bytes"
	"fmt"
	"github.com/nestoca/joy/api/v1alpha1"
	"gopkg.in/yaml.v3"
	"text/template"
)

// RenderValues renders the values of the given release and produces a yaml string, processing any go template
// directives found in the values.
func RenderValues(release *v1alpha1.Release) (string, error) {
	valuesStr, err := yaml.Marshal(release.Spec.Values)
	if err != nil {
		return "", fmt.Errorf("marshalling release values: %w", err)
	}

	tpl, err := template.New("values").Parse(string(valuesStr))
	if err != nil {
		return "", fmt.Errorf("parsing values template: %w", err)
	}

	tpl.Option("missingkey=error")

	var result bytes.Buffer
	type TemplateData struct {
		Release *v1alpha1.Release
	}

	err = tpl.Execute(&result, TemplateData{
		Release: release,
	})
	if err != nil {
		return "", fmt.Errorf("executing values template: %w", err)
	}

	return result.String(), nil
}
