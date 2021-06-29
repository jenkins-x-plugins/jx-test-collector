package resources

import (
	"github.com/jenkins-x/jx-helpers/v3/pkg/options"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/jenkins-x/jx-helpers/v3/pkg/files"
	"github.com/jenkins-x/jx-helpers/v3/pkg/kube"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"sigs.k8s.io/yaml"
)

// Options the options for the resources operation
type Options struct {
	options.BaseOptions

	// Dir the directory to dump resources to
	Dir string

	// Namespace the namespace to query resources from
	Namespace string

	// DynamicClient the client to access kubernetes resources
	DynamicClient dynamic.Interface
}

var (
	// ResourceGVRs the resources
	ResourceGVRs = []schema.GroupVersionResource{
		// core resources
		{Group: "", Version: "v1", Resource: "pods"},

		// jx resources
		{Group: "jenkins.io", Version: "v1", Resource: "pipelineactivities"},

		// tekton resources
		{Group: "tekton.dev", Version: "v1alpha1", Resource: "pipelines"},
		{Group: "tekton.dev", Version: "v1alpha1", Resource: "pipelineruns"},
		{Group: "tekton.dev", Version: "v1alpha1", Resource: "taskruns"},
		{Group: "tekton.dev", Version: "v1alpha1", Resource: "tasks"},
	}

	// ResourceMap the map for fake clients
	ResourceMap = map[schema.GroupVersionResource]string{}
)

func init() {
	for _, r := range ResourceGVRs {
		ResourceMap[r] = strings.Title(r.Resource + "List")
	}
}

// Validate validates the options
func (o *Options) Validate(dir string) error {
	o.Dir = dir
	var err error
	o.DynamicClient, err = kube.LazyCreateDynamicClient(o.DynamicClient)
	if err != nil {
		return errors.Wrapf(err, "failed to create kubernetes dynamic client")
	}
	return nil
}

// Run will implement this command
func (o *Options) Run() error {
	err := o.Validate(o.Dir)
	if err != nil {
		return errors.Wrap(err, "invalid options")
	}

	dynClient := o.DynamicClient
	ns := o.Namespace

	ctx := o.GetContext()
	for _, r := range ResourceGVRs {
		log := logrus.WithFields(map[string]interface{}{
			"Namespace": ns,
			"Group":     r.Group,
			"Resource":  r.Resource,
		})
		resources, err := dynClient.Resource(r).Namespace(ns).List(ctx, metav1.ListOptions{})
		if err != nil && !apierrors.IsNotFound(err) {
			// probably RBAC related
			log.WithError(err).Error("cannot list resources")
			continue
		}
		if resources == nil {
			continue
		}
		for i := range resources.Items {
			resource := &resources.Items[i]
			if r.Group == "" {
				r.Group = "core"
			}
			if r.Version == "" {
				r.Version = "v1"
			}
			dir := filepath.Join(o.Dir, r.Group, r.Version, r.Resource)
			ns := resource.GetNamespace()
			name := resource.GetName()

			if ns != "" {
				dir = filepath.Join(dir, ns)
			}
			err := os.MkdirAll(dir, files.DefaultDirWritePermissions)
			if err != nil {
				return errors.Wrapf(err, "failed to create directory %s", dir)
			}

			fileName := filepath.Join(dir, name+".yaml")
			data, err := yaml.Marshal(resource)
			if err != nil {
				return errors.Wrapf(err, "failed to marshal resource to YAML for file %s", fileName)
			}

			err = ioutil.WriteFile(fileName, data, files.DefaultFileWritePermissions)
			if err != nil {
				return errors.Wrapf(err, "failed to save file %s", fileName)
			}
		}
	}
	return nil
}
