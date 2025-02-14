package binary

import (
	"testing"

	"github.com/portainer/portainer/pkg/libhelm/libhelmtest"
	"github.com/portainer/portainer/pkg/libhelm/options"
	"github.com/stretchr/testify/assert"
)

func Test_SearchRepo(t *testing.T) {
	libhelmtest.EnsureIntegrationTest(t)
	is := assert.New(t)

	hpm := NewHelmBinaryPackageManager("")

	type testCase struct {
		name    string
		url     string
		invalid bool
	}

	tests := []testCase{
		{"not a helm repo", "https://portainer.io", true},
		{"ingress helm repo", "https://kubernetes.github.io/ingress-nginx", false},
		{"portainer helm repo", "https://portainer.github.io/k8s/", false},
		{"gitlap helm repo with trailing slash", "https://charts.gitlab.io/", false},
		{"elastic helm repo with trailing slash", "https://helm.elastic.co/", false},
		{"fabric8.io helm repo with trailing slash", "https://fabric8.io/helm/", false},
		{"lensesio helm repo without trailing slash", "https://lensesio.github.io/kafka-helm-charts", false},
	}

	for _, test := range tests {
		func(tc testCase) {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				response, err := hpm.SearchRepo(options.SearchRepoOptions{Repo: tc.url})
				if tc.invalid {
					is.Errorf(err, "error expected: %s", tc.url)
				} else {
					is.NoError(err, "no error expected: %s", tc.url)
				}

				if err == nil {
					is.NotEmpty(response, "response expected")
				}
			})
		}(test)
	}
}
