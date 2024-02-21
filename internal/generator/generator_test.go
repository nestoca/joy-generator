package generator

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nestoca/joy/api/v1alpha1"
	joy "github.com/nestoca/joy/pkg"
	"github.com/nestoca/joy/pkg/catalog"
)

func TestGenerator(t *testing.T) {
	cases := []struct {
		Name            string
		Release         *v1alpha1.Release
		DefaultChart    string
		ValueMapping    *joy.ValueMapping
		ExpectedRelease *v1alpha1.Release
		ExpectedValues  string
	}{
		{
			Name: "release chart is not overridden by default chart",
			Release: &v1alpha1.Release{
				ApiVersion:      "joy.nesto.ca/v1alpha1",
				Kind:            "Release",
				ReleaseMetadata: v1alpha1.ReleaseMetadata{Name: "app"},
				Spec: v1alpha1.ReleaseSpec{
					Version: "v1",
					Chart: v1alpha1.ReleaseChart{
						RepoUrl: "test",
						Name:    "release-chart",
						Version: "0.0.0",
					},
					Values: map[string]any{},
				},
				Environment: &v1alpha1.Environment{
					ApiVersion:          "joy.nesto.ca/v1alpha1",
					Kind:                "Environment",
					EnvironmentMetadata: v1alpha1.EnvironmentMetadata{Name: "test"},
				},
				Project: &v1alpha1.Project{
					ApiVersion:      "joy.nesto.ca/v1alpha1",
					Kind:            "Project",
					ProjectMetadata: v1alpha1.ProjectMetadata{Name: "test"},
				},
			},
			DefaultChart: "default/default",
			ExpectedRelease: &v1alpha1.Release{
				ApiVersion:      "joy.nesto.ca/v1alpha1",
				Kind:            "Release",
				ReleaseMetadata: v1alpha1.ReleaseMetadata{Name: "app"},
				Spec: v1alpha1.ReleaseSpec{
					Project: "",
					Version: "v1",
					Chart: v1alpha1.ReleaseChart{
						Name:    "release-chart",
						RepoUrl: "test",
						Version: "0.0.0",
					},
					Values: map[string]interface{}{},
				},
				Environment: &v1alpha1.Environment{
					ApiVersion:          "joy.nesto.ca/v1alpha1",
					Kind:                "Environment",
					EnvironmentMetadata: v1alpha1.EnvironmentMetadata{Name: "test"},
				},
				Project: &v1alpha1.Project{
					ApiVersion:      "joy.nesto.ca/v1alpha1",
					Kind:            "Project",
					ProjectMetadata: v1alpha1.ProjectMetadata{Name: "test"},
				},
			},
			ExpectedValues: "{}\n",
		},
		{
			Name: "renders default chart and value mappings",
			Release: &v1alpha1.Release{
				ApiVersion:      "joy.nesto.ca/v1alpha1",
				Kind:            "Release",
				ReleaseMetadata: v1alpha1.ReleaseMetadata{Name: "app"},
				Spec: v1alpha1.ReleaseSpec{
					Version: "v1",
					Chart:   v1alpha1.ReleaseChart{Version: "test-version"},
					Values:  map[string]any{},
				},
				Environment: &v1alpha1.Environment{
					ApiVersion:          "joy.nesto.ca/v1alpha1",
					Kind:                "Environment",
					EnvironmentMetadata: v1alpha1.EnvironmentMetadata{Name: "test"},
				},
				Project: &v1alpha1.Project{
					ApiVersion:      "joy.nesto.ca/v1alpha1",
					Kind:            "Project",
					ProjectMetadata: v1alpha1.ProjectMetadata{Name: "test"},
				},
			},
			DefaultChart: "nesto.repo/test-chart",
			ValueMapping: &joy.ValueMapping{Mappings: map[string]any{
				"annotations.test": true,
				"image":            "image@{{ .Release.Spec.Version }}",
			}},
			ExpectedRelease: &v1alpha1.Release{
				ApiVersion:      "joy.nesto.ca/v1alpha1",
				Kind:            "Release",
				ReleaseMetadata: v1alpha1.ReleaseMetadata{Name: "app"},
				Spec: v1alpha1.ReleaseSpec{
					Project: "",
					Version: "v1",
					Chart: v1alpha1.ReleaseChart{
						Name:    "test-chart",
						RepoUrl: "nesto.repo",
						Version: "test-version",
					},
					Values: map[string]interface{}{},
				},
				Environment: &v1alpha1.Environment{
					ApiVersion:          "joy.nesto.ca/v1alpha1",
					Kind:                "Environment",
					EnvironmentMetadata: v1alpha1.EnvironmentMetadata{Name: "test"},
				},
				Project: &v1alpha1.Project{
					ApiVersion:      "joy.nesto.ca/v1alpha1",
					Kind:            "Project",
					ProjectMetadata: v1alpha1.ProjectMetadata{Name: "test"},
				},
			},
			ExpectedValues: "annotations:\n    test: true\nimage: image@v1\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			generator := New(func() (*JoyContext, error) {
				return &JoyContext{
					Catalog: BuildCatalogFromRelease(tc.Release),
					Config:  &joy.Config{DefaultChart: tc.DefaultChart, ValueMapping: tc.ValueMapping},
				}, nil
			})

			results, err := generator.Run()
			require.NoError(t, err)

			require.Len(t, results, 1)
			require.EqualValues(t, tc.ExpectedRelease, results[0].Release)
			require.Equal(t, tc.ExpectedValues, results[0].Values)
		})
	}
}

func BuildCatalogFromRelease(release *v1alpha1.Release) *catalog.Catalog {
	return &catalog.Catalog{
		Releases: &catalog.ReleaseList{
			Items: []*catalog.Release{
				{
					Name:     release.ReleaseMetadata.Name,
					Releases: []*v1alpha1.Release{release},
				},
			},
		},
	}
}
