package tailer

import (
	"context"
	"io/ioutil"
	"regexp"
	"text/template"
	"time"

	"github.com/jenkins-x/jx-helpers/pkg/cmdrunner"
	"github.com/jenkins-x/jx-helpers/pkg/gitclient"
	"github.com/jenkins-x/jx-helpers/pkg/gitclient/cli"
	"github.com/jenkins-x/jx-helpers/pkg/kube"
	"github.com/jenkins-x/jx-logging/pkg/log"
	"github.com/jenkins-x/jx-test-collector/pkg/constants"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

// Options the configuration options for the tailer
type Options struct {
	GitClient gitclient.Interface

	// CommandRunner used to run git commands if no GitClient provided
	CommandRunner cmdrunner.CommandRunner

	// KubeClient is used to lazy create the repo client and launcher
	KubeClient kubernetes.Interface

	// Dir is the work directory. If not specified a temporary directory is created on startup.
	Dir string `env:"WORK_DIR"`

	// Namespace the namespace polled. Defaults to all of them
	Namespace string `env:"NAMESPACE"`

	// GitBinary name of the git binary; defaults to `git`
	GitBinary string `env:"GIT_BINARY"`

	// PollDuration duration between polls
	PollDuration time.Duration `env:"POLL_DURATION"`

	// NoLoop disable the polling loop so that a single poll is performed only
	NoLoop bool `env:"NO_LOOP"`

	// NoResourceApply disable the applying of resources in a git repository at `.jx/git-operator/resources/*.yaml`
	NoResourceApply bool `env:"NO_RESOURCE_APPLY"`

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

	if o.Namespace != "" {
		log.Logger().Infof("looking in namespace %s for Secret resources with selector %s", o.Namespace, constants.DefaultSelector)
	}

	err = o.Poll()
	if err != nil {
		return err
	}
	return nil
}

// Poll polls the available repositories
func (o *Options) Poll() error {
	err := o.ValidateOptions()
	if err != nil {
		return errors.Wrap(err, "invalid options")
	}

	namespace := o.Namespace

	logrus.Infof("tailing logs of pods in namespace :%s", namespace)

	ctx := context.Background()
	kubeClient := o.KubeClient

	added, removed, err := o.Watch(ctx, kubeClient.CoreV1().Pods(namespace), o.LabelSelector)
	if err != nil {
		return errors.Wrap(err, "failed to set up watch")
	}

	tails := make(map[string]*Tail)

	go func() {
		for p := range added {
			id := p.GetID()
			if tails[id] != nil {
				continue
			}

			tail := NewTail(o.Dir, p.Namespace, p.Pod, p.Container, o.Template, &TailOptions{
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
	if o.LabelSelector == nil {
		o.LabelSelector = labels.NewSelector()
	}
	if o.PollDuration.Milliseconds() == int64(0) {
		o.PollDuration = time.Second * 30
	}
	if o.GitClient == nil {
		o.GitClient = cli.NewCLIClient(o.GitBinary, o.CommandRunner)
	}
	var err error
	o.KubeClient, err = kube.LazyCreateKubeClient(o.KubeClient)
	if err != nil {
		return errors.Wrapf(err, "failed to create kube client")
	}
	if o.Dir == "" {
		o.Dir, err = ioutil.TempDir("", "jx-test-collector-")
		if err != nil {
			return errors.Wrapf(err, "failed to create temp dir")
		}
	}
	logrus.Infof("writing files to dir: %s", o.Dir)
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
