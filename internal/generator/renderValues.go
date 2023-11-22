package generator

import (
	"bytes"
	"fmt"
	"text/template"

	"gopkg.in/yaml.v3"

	"github.com/nestoca/joy/api/v1alpha1"
)

// RenderValues renders the values of the given release and produces a yaml string, processing any go template
// directives found in the values.
func RenderValues(release *v1alpha1.Release) (string, error) {
	buf := &bytes.Buffer{}
	encoder := yaml.NewEncoder(buf)
	encoder.SetIndent(2)
	err := encoder.Encode(release.Spec.Values)
	if err != nil {
		return "", fmt.Errorf("marshalling release values: %w", err)
	}

	tpl, err := template.New("values").Parse(string(buf.Bytes()))
	if err != nil {
		return "", fmt.Errorf("parsing values template: %w", err)
	}

	tpl.Option("missingkey=error")

	var result bytes.Buffer
	type TemplateData struct {
		Release     *v1alpha1.Release
		Environment *v1alpha1.Environment
	}

	err = tpl.Execute(&result, TemplateData{
		Release:     release,
		Environment: release.Environment,
	})
	if err != nil {
		return "", fmt.Errorf("executing values template: %w", err)
	}

	return result.String(), nil
}
