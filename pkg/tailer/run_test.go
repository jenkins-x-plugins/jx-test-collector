package tailer_test

import (
	"io/ioutil"
	"testing"

	"github.com/jenkins-x/jx-test-collector/pkg/tailer"
	"github.com/stretchr/testify/require"
)

func TestTailerValidate(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test-jx-test-collector-")
	require.NoError(t, err, "failed to create temp dir")

	t.Logf("running in dir %s", tmpDir)
	p := &tailer.Options{}

	p.Dir = tmpDir

	err = p.ValidateOptions()
	require.NoError(t, err, "failed to ValidateOptions()")
}
