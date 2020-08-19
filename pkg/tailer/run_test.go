package tailer_test

import (
	"io/ioutil"
	"testing"

	"github.com/jenkins-x/jx-helpers/pkg/cmdrunner/fakerunner"
	"github.com/jenkins-x/jx-test-collector/pkg/tailer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestTailerValidate(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test-jx-test-collector-")
	require.NoError(t, err, "failed to create temp dir")

	t.Logf("running in dir %s", tmpDir)
	o := &tailer.Options{}

	o.Dir = tmpDir

	ns := "jx"

	runner := &fakerunner.FakeRunner{}
	o.GitStore.CommandRunner = runner.Run

	o.Namespace = ns
	o.KubeClient = fake.NewSimpleClientset(
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "jx-boot",
				Namespace: ns,
			},
			Data: map[string][]byte{
				"url":      []byte("https://github.com/myorg/myrepo.git"),
				"username": []byte("myuser"),
				"password": []byte("mypwd"),
			},
		},
	)

	err = o.ValidateOptions()
	require.NoError(t, err, "failed to ValidateOptions()")

	cloneURL, err := o.GitStore.GitCloneURL()
	require.NoError(t, err, "failed to create git clone URL")
	assert.Equal(t, "https://myuser:mypwd@github.com/myorg/myrepo.git", cloneURL, "git clone URL")
}
