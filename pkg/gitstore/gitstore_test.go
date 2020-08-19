package gitstore_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/jenkins-x/jx-helpers/pkg/cmdrunner"
	"github.com/jenkins-x/jx-helpers/pkg/cmdrunner/fakerunner"
	"github.com/jenkins-x/jx-helpers/pkg/files"
	"github.com/jenkins-x/jx-test-collector/pkg/gitstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGitStore(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test-jx-test-collector-")
	require.NoError(t, err, "failed to create temp dir")
	t.Logf("running in dir %s", tmpDir)

	o := &gitstore.Options{}
	ns := "jx"
	kubeClient := fake.NewSimpleClientset(
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "jx-boot",
				Namespace: ns,
			},
			Data: map[string][]byte{
				"url":      []byte("https://github.com/jenkins-x/jenkins-x-versions-test.git"),
				"username": []byte("myuser"),
				"password": []byte("mypwd"),
			},
		},
	)
	runner := &fakerunner.FakeRunner{
		CommandRunner: func(c *cmdrunner.Command) (string, error) {
			if c.Name == "git" && len(c.Args) > 0 && c.Args[0] == "push" {
				// lets disable git pushing
				t.Logf("disabling command: %s\n", c.CLI())
				return "fake git pushing", nil
			}
			return cmdrunner.DefaultCommandRunner(c)
		},
	}
	o.CommandRunner = runner.Run

	err = o.Validate(kubeClient, tmpDir)
	require.NoError(t, err, "failed to run Validate()")

	cloneURL, err := o.GitCloneURL()
	require.NoError(t, err, "failed to create git clone URL")
	assert.Equal(t, "https://myuser:mypwd@github.com/jenkins-x/jenkins-x-versions-test.git", cloneURL, "git clone URL")

	err = o.Setup()
	require.NoError(t, err, "failed to run Setup()")

	// lets add a file into git...
	outDir := filepath.Join(tmpDir, "logs", "jx", "mypod")
	err = os.MkdirAll(outDir, files.DefaultDirWritePermissions)
	require.NoError(t, err, "failed to create dir %s", outDir)

	outFile := filepath.Join(outDir, "container.log")
	err = ioutil.WriteFile(outFile, []byte("Hello\nWorld!\n"), files.DefaultFileWritePermissions)
	require.NoError(t, err, "failed to save file %s", outFile)

	text, err := o.Sync()
	require.NoError(t, err, "failed to run Sync()")

	t.Logf("Sync returned: %s\n", text)
}
