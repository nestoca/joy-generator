package generator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nestoca/joy-generator/internal/generator"
	"github.com/nestoca/joy/api/v1alpha1"
)

func TestRenderValuesWithValidInput(t *testing.T) {
	release := &v1alpha1.Release{
		ReleaseMetadata: v1alpha1.ReleaseMetadata{
			Name: "test-release",
		},
		Spec: v1alpha1.ReleaseSpec{
			Values: map[string]interface{}{
				"foo":       "bar",
				"templated": "{{ .Release.Name }}",
			},
		},
	}

	expected := `foo: bar
templated: 'test-release'
`

	renderedValues, err := generator.RenderValues(release)
	assert.Nil(t, err, "should not return an error")
	assert.Equal(t, expected, renderedValues, "should render the values")
}

func TestRenderValuesReturnsErrorIfInvalidTemplate(t *testing.T) {
	release := &v1alpha1.Release{
		ReleaseMetadata: v1alpha1.ReleaseMetadata{
			Name: "test-release",
		},
		Spec: v1alpha1.ReleaseSpec{
			Values: map[string]interface{}{
				"foo":       "bar",
				"templated": "{{ .Release.Name }",
			},
		},
	}

	renderedValues, err := generator.RenderValues(release)
	assert.NotNil(t, err, "should return an error")
	assert.Equal(t, "", renderedValues, "should not render the values")
}

func TestRenderValuesReturnsErrorIfMissingKey(t *testing.T) {
	release := &v1alpha1.Release{
		ReleaseMetadata: v1alpha1.ReleaseMetadata{
			Name: "test-release",
		},
		Spec: v1alpha1.ReleaseSpec{
			Values: map[string]interface{}{
				"foo":       "bar",
				"templated": "{{ .Foo.Name }}",
			},
		},
	}

	renderedValues, err := generator.RenderValues(release)
	assert.NotNil(t, err, "should return an error")
	assert.Equal(t, "", renderedValues, "should not render the values")
}
