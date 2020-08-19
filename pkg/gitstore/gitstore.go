package gitstore

import (
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"

	"github.com/jenkins-x/jx-helpers/pkg/cmdrunner"
	"github.com/jenkins-x/jx-helpers/pkg/gitclient"
	"github.com/jenkins-x/jx-helpers/pkg/gitclient/cli"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Options the interface to the storage
type Options struct {
	// Dir the directory of the git clone
	Dir string

	// GitClient the git client to perform git operations
	GitClient gitclient.Interface

	// CommandRunner used to run git commands if no GitClient provided
	CommandRunner cmdrunner.CommandRunner

	// GitBinary name of the git binary; defaults to `git`
	GitBinary string `env:"GIT_BINARY,default=git"`

	// URL the git URL to clone.
	//
	// If not specified it defaults to using the jx boot
	// secret installed via the Git Operator
	//
	// see: https://jenkins-x.io/docs/v3/guides/operator/
	URL string `env:"GIT_URL"`

	// Username the git username to clone and commit.
	//
	// If not specified it defaults to using the jx boot
	// secret installed via the Git Operator
	//
	// see: https://jenkins-x.io/docs/v3/guides/operator/
	Username string `env:"GIT_USERNAME"`

	// Token the git token to clone and commit.
	//
	// If not specified it defaults to using the jx boot
	// secret installed via the Git Operator
	//
	// see: https://jenkins-x.io/docs/v3/guides/operator/
	Token string `env:"GIT_TOKEN"`

	// Branch the git branch to use to store logs and resources
	Branch string `env:"GIT_BRANCH,default=gh-pages"`

	// JXNamespace the namespace Jenkins X is installed into.
	//
	// Used to find the jx-boot secret to get the URL, user and token for the git repository
	// if none is provided explicitly.
	//
	// see: https://jenkins-x.io/docs/v3/guides/operator/
	JXNamespace string `env:"JX_NAMESPACE,default=jx"`

	// SecretName the name of the Jenkins X boot secret to load the url/username/token from if not explicitly defined.
	//
	// see: https://jenkins-x.io/docs/v3/guides/operator/
	SecretName string `env:"SECRET_NAME,default=jx-boot"`
}

// Validate validates the options and lazily creates any resources required
func (o *Options) Validate(kubeClient kubernetes.Interface, dir string) error {
	o.Dir = dir
	if o.GitClient == nil {
		o.GitClient = cli.NewCLIClient(o.GitBinary, o.CommandRunner)
	}
	if o.Branch == "" {
		o.Branch = "gh-pages"
	}
	if o.URL == "" || o.Username == "" || o.Token == "" {
		name := o.SecretName
		ns := o.JXNamespace
		if ns == "" {
			ns = "jx"
		}

		// lets get the secret
		secret, err := kubeClient.CoreV1().Secrets(ns).Get(name, metav1.GetOptions{})
		if err != nil {
			if apierrors.IsNotFound(err) {
				return errors.Errorf("no Secret %s in namespace %s", name, ns)
			}
			return errors.Errorf("failed to load Secret %s in namespace %s", name, ns)
		}
		data := secret.Data
		if data != nil {
			gitURL := string(data["url"])
			username := string(data["username"])
			password := string(data["password"])

			if o.URL == "" {
				o.URL = gitURL
				if o.URL == "" {
					return errors.Errorf("secret %s in namespace %s does not have a url entry", name, ns)
				}
			}
			if o.Username == "" {
				o.Username = username
				if o.Username == "" {
					return errors.Errorf("secret %s in namespace %s does not have a username entry", name, ns)
				}
			}
			if o.Token == "" {
				o.Token = password
				if o.Token == "" {
					return errors.Errorf("secret %s in namespace %s does not have a password entry", name, ns)
				}
			}
		}
	}
	logrus.WithFields(map[string]interface{}{
		"URL":      o.URL,
		"Username": o.Username,
		"Git":      o.GitBinary,
	}).Infof("setup GitStore")
	return nil
}

// Setup sets up the storage in the given directory
func (o *Options) Setup() error {
	dir := o.Dir
	g := o.GitClient
	gitCloneURL, err := o.GitCloneURL()
	if err != nil {
		return errors.Wrapf(err, "failed to create the git clone URL")
	}

	parentDir := filepath.Dir(dir)

	text, err := g.Command(parentDir, "clone", gitCloneURL, "--branch", o.Branch, "--single-branch", dir)
	if err != nil {
		log := logrus.WithError(err).WithFields(map[string]interface{}{
			"URL":      o.URL,
			"Username": o.Username,
			"Branch":   o.Branch,
			"Output":   text,
		})
		log.Infof("assuming the remote branch does not exist so lets create it")

		_, err = gitclient.CloneToDir(g, gitCloneURL, dir)
		if err != nil {
			return errors.Wrapf(err, "failed to clone repository %s with user %s to directory: %s", o.URL, o.Username, dir)
		}

		// now lets create an empty orphan branch: see https://stackoverflow.com/a/13969482/2068211
		_, err = g.Command(dir, "checkout", "--orphan", o.Branch)
		if err != nil {
			return errors.Wrapf(err, "failed to checkout an orphan branch %s in dir %s", o.Branch, dir)
		}

		_, err = g.Command(dir, "rm", "--cached", "-r", ".")
		if err != nil {
			return errors.Wrapf(err, "failed to remove the cached git files in dir %s", dir)
		}

		// lets remove all the files other than .git
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			return errors.Wrapf(err, "failed to read files in dir %s", dir)
		}
		for _, f := range files {
			name := f.Name()
			if name == ".git" {
				continue
			}
			path := filepath.Join(dir, name)
			err = os.RemoveAll(path)
			if err != nil {
				return errors.Wrapf(err, "failed to remove path %s", path)
			}
		}
	}
	return nil
}

// Sync performs a synchronisation of any local files to the underlying storage engine
func (o *Options) Sync() (string, error) {
	dir := o.Dir
	g := o.GitClient
	answer := ""
	_, err := g.Command(dir, "add", "*")
	if err != nil {
		return answer, errors.Wrapf(err, "failed to add files to git")
	}

	changes, err := gitclient.HasChanges(g, dir)
	if err != nil {
		return answer, errors.Wrapf(err, "failed to check if there are changes in git")
	}

	if !changes {
		return "no changes", nil
	}

	_, err = g.Command(dir, "commit", "-a", "-m", "chore: latest logs")
	if err != nil {
		return "", errors.Wrapf(err, "failed to commit latest logs to dir %s", dir)
	}
	_, err = g.Command(dir, "push", "origin", o.Branch)
	if err != nil {
		return "", errors.Wrapf(err, "failed to push changes to git")
	}
	return "sync completed", nil
}

// GitCloneURL returns the git clone URL
func (o *Options) GitCloneURL() (string, error) {
	u, err := url.Parse(o.URL)
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse git URL %s", o.URL)
	}
	u.User = url.UserPassword(o.Username, o.Token)
	return u.String(), nil
}
