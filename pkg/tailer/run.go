package tailer

import (
	"context"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"text/template"
	"time"

	"github.com/jenkins-x/jx-helpers/pkg/kube"
	"github.com/jenkins-x/jx-test-collector/pkg/gitstore"
	"github.com/jenkins-x/jx-test-collector/pkg/masker"
	"github.com/jenkins-x/jx-test-collector/pkg/resources"
	"github.com/jenkins-x/jx-test-collector/pkg/web"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

// Options the configuration options for the tailer
type Options struct {
	// Web REST API
	Web web.Options

	// Resources for dumping kubernetes resources
	Resources resources.Options

	// GitStore takes care of storing files in git
	GitStore gitstore.Options

	// Dir is the work directory. If not specified a temporary directory is created on startup.
	Dir string `env:"WORK_DIR"`

	// LogPath the path within Dir where we store pod logs
	LogPath string `env:"LOG_PATH,default=logs"`

	// ResourcePath the path within Dir where we store resources
	ResourcePath string `env:"RESOURCE_PATH,default=resources"`

	// Namespace the namespace polled. Defaults to all of them
	Namespace string `env:"NAMESPACE"`

	// PollDuration duration between polls
	PollDuration time.Duration `env:"POLL_DURATION"`

	// NoLoop disable the polling loop so that a single poll is performed only
	NoLoop bool `env:"NO_LOOP"`

	// NoResourceApply disable the applying of resources in a git repository at `.jx/git-operator/resources/*.yaml`
	NoResourceApply bool `env:"NO_RESOURCE_APPLY"`

	// KubeClient is used to lazy create the repo client and launcher
	KubeClient kubernetes.Interface

	// Masker for masking secrets in logs
	Masker *masker.Client

	Timestamps    bool
	Exclude       []*regexp.Regexp
	Include       []*regexp.Regexp
	Since         time.Duration
	AllNamespaces bool
	LabelSelector labels.Selector
	TailLines     *int64
	Template      *template.Template
}

// Run polls for git changes
func (o *Options) Run() error {
	err := o.ValidateOptions()
	if err != nil {
		return errors.Wrap(err, "invalid options")
	}

	err = o.GitStore.Setup()
	if err != nil {
		return errors.Wrapf(err, "failed to setup git store")
	}

	go func() {
		err := o.Web.Run()
		if err != nil {
			logrus.WithError(err).Fatal("failed to serve http")
		}
	}()

	namespace := o.Namespace

	logrus.Infof("tailing logs of pods in namespace :%s", namespace)

	ctx := context.Background()
	kubeClient := o.KubeClient

	o.Masker, err = masker.NewMasker(kubeClient, namespace)
	if err != nil {
		return errors.Wrapf(err, "failed to create masker")
	}

	added, removed, err := o.Watch(ctx, kubeClient.CoreV1().Pods(namespace), o.LabelSelector)
	if err != nil {
		return errors.Wrap(err, "failed to set up watch")
	}

	tails := make(map[string]*Tail)

	podLogDir := filepath.Join(o.Dir, o.LogPath)

	go func() {
		for p := range added {
			id := p.GetID()
			if tails[id] != nil {
				continue
			}

			tail := NewTail(o.Masker, podLogDir, p.Namespace, p.Pod, p.Container, p.App, o.Template, &TailOptions{
				Timestamps:   o.Timestamps,
				SinceSeconds: int64(o.Since.Seconds()),
				Exclude:      o.Exclude,
				Include:      o.Include,
				Namespace:    o.AllNamespaces,
				TailLines:    o.TailLines,
			})
			tails[id] = tail

			tail.Start(ctx, kubeClient.CoreV1().Pods(p.Namespace))
		}
	}()

	go func() {
		for p := range removed {
			id := p.GetID()
			if tails[id] == nil {
				continue
			}
			tails[id].Close()
			delete(tails, id)
		}
	}()

	<-ctx.Done()

	return nil
}

// ValidateOptions validates the options and lazily creates any resources required
func (o *Options) ValidateOptions() error {
	o.Web.Sync = func() (string, error) {
		err := o.Resources.Run()
		if err != nil {
			return "", errors.Wrapf(err, "failed to get kubernetes resources")
		}
		return o.GitStore.Sync()
	}

	var err error
	o.KubeClient, err = kube.LazyCreateKubeClient(o.KubeClient)
	if err != nil {
		return errors.Wrapf(err, "failed to create kube client")
	}
	if o.LabelSelector == nil {
		o.LabelSelector = labels.NewSelector()
	}
	if o.PollDuration.Milliseconds() == int64(0) {
		o.PollDuration = time.Second * 30
	}
	if o.Dir == "" {
		o.Dir, err = ioutil.TempDir("", "jx-test-collector-")
		if err != nil {
			return errors.Wrapf(err, "failed to create temp dir")
		}
	}
	logrus.Infof("writing files to dir: %s", o.Dir)

	err = o.GitStore.Validate(o.KubeClient, o.Dir)
	if err != nil {
		return errors.Wrapf(err, "failed to validate GitStore")
	}

	err = o.Resources.Validate(filepath.Join(o.Dir, o.ResourcePath))
	if err != nil {
		return errors.Wrapf(err, "failed to setup resource fetcher")
	}
	return nil
}

// MatchPod for filtering on the pod
func (o *Options) MatchPod(_ *corev1.Pod) bool {
	return true
}

// MatchesContainerStatus matches a container
func (o *Options) MatchesContainerStatus(pod *corev1.Pod, c corev1.ContainerStatus) bool {
	return true
}

// MatchesContainer returns true if we should match this container
func (o *Options) MatchesContainer(pod *corev1.Pod, c corev1.Container) bool {
	return true

}
