package resources_test

import (
	"context"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/jenkins-x/jx-test-collector/pkg/resources"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic/fake"
)

func TestResources(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test-jx-test-collector-")
	require.NoError(t, err, "failed to create temp dir")
	t.Logf("running in dir %s", tmpDir)

	testResources := LoadTestResources(t, "test_data")

	o := &resources.Options{}
	o.Dir = tmpDir
	scheme := runtime.NewScheme()
	corev1.AddToScheme(scheme)

	o.Ctx = context.TODO()
	o.DynamicClient = fake.NewSimpleDynamicClientWithCustomListKinds(scheme, resources.ResourceMap, testResources...)

	err = o.Run()
	require.NoError(t, err, "failed to run Run()")
}

// LoadTestResources loads the test resources
func LoadTestResources(t *testing.T, dir string) []runtime.Object {
	files, err := ioutil.ReadDir(dir)
	require.NoError(t, err, "failed to read dir %s", dir)
	var dynObjects []runtime.Object
	for _, f := range files {
		name := f.Name()
		if f.IsDir() || !strings.HasSuffix(name, ".yaml") {
			continue

		}
		u := &unstructured.Unstructured{}
		path := filepath.Join(dir, name)

		data, err := ioutil.ReadFile(path)
		require.NoError(t, err, "failed to load file %s", path)

		err = yaml.Unmarshal(data, u)
		require.NoError(t, err, "failed to unmarshal YAML file %s", path)
		dynObjects = append(dynObjects, u)
	}
	return dynObjects
}
