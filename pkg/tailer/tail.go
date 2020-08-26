//   Copyright 2016 Wercker Holding BV
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

// Based on source:
// https://github.com/wercker/stern/blob/master/stern/tail.go

package tailer

import (
	"bufio"
	"context"
	"hash/fnv"
	"os"
	"path/filepath"
	"regexp"
	"text/template"

	"github.com/fatih/color"
	"github.com/jenkins-x/jx-helpers/pkg/files"
	"github.com/jenkins-x/jx-test-collector/pkg/masker"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
)

type Tail struct {
	Dir            string
	Namespace      string
	PodName        string
	ContainerName  string
	Options        *TailOptions
	req            *rest.Request
	closed         chan struct{}
	podColor       *color.Color
	containerColor *color.Color
	tmpl           *template.Template
	log            *logrus.Entry
	masker         *masker.Client
}

type TailOptions struct {
	Timestamps   bool
	SinceSeconds int64
	Exclude      []*regexp.Regexp
	Include      []*regexp.Regexp
	Namespace    bool
	TailLines    *int64
}

// NewTail returns a new tail for a Kubernetes container inside a pod
func NewTail(masker *masker.Client, dir, namespace, podName, containerName, app string, tmpl *template.Template, options *TailOptions) *Tail {
	log := logrus.WithFields(
		map[string]interface{}{
			"Namespace": namespace,
			"Pod":       podName,
			"Container": containerName,
		})

	nsDir := filepath.Join(dir, namespace)
	if app != "" {
		nsDir = filepath.Join(nsDir, app)
	}
	podDir := filepath.Join(nsDir, podName)
	err := os.MkdirAll(podDir, files.DefaultDirWritePermissions)
	if err != nil {
		log.WithError(err).Errorf("failed to create dir: %s", podDir)
	}

	return &Tail{
		Dir:           podDir,
		log:           log,
		Namespace:     namespace,
		PodName:       podName,
		ContainerName: containerName,
		Options:       options,
		masker:        masker,
		closed:        make(chan struct{}),
		tmpl:          tmpl,
	}
}

var colorList = [][2]*color.Color{
	{color.New(color.FgHiCyan), color.New(color.FgCyan)},
	{color.New(color.FgHiGreen), color.New(color.FgGreen)},
	{color.New(color.FgHiMagenta), color.New(color.FgMagenta)},
	{color.New(color.FgHiYellow), color.New(color.FgYellow)},
	{color.New(color.FgHiBlue), color.New(color.FgBlue)},
	{color.New(color.FgHiRed), color.New(color.FgRed)},
}

func determineColor(podName string) (podColor, containerColor *color.Color) {
	hash := fnv.New32()
	hash.Write([]byte(podName))
	idx := hash.Sum32() % uint32(len(colorList))

	colors := colorList[idx]
	return colors[0], colors[1]
}

// Start starts tailing
func (t *Tail) Start(ctx context.Context, i v1.PodInterface) {
	t.podColor, t.containerColor = determineColor(t.PodName)

	go func() {
		req := i.GetLogs(t.PodName, &corev1.PodLogOptions{
			Follow:     true,
			Timestamps: t.Options.Timestamps,
			Container:  t.ContainerName,
			TailLines:  t.Options.TailLines,
			//SinceSeconds: &t.Options.SinceSeconds,
		})

		fileName := filepath.Join(t.Dir, t.ContainerName+".log")
		file, err := os.Create(fileName)
		if err != nil {
			t.log.WithError(err).Errorf("failed to create output")
			return
		}
		defer file.Close()

		writer := bufio.NewWriter(file)
		defer writer.Flush()

		stream, err := req.Stream()
		if err != nil {
			t.log.WithError(err).Warnf("Error opening stream to %s/%s: %s\n", t.Namespace, t.PodName, t.ContainerName)
			return
		}
		defer stream.Close()

		go func() {
			<-t.closed
			stream.Close()
			writer.Flush()
			file.Close()
		}()

		reader := bufio.NewReader(stream)

	OUTER:
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				return
			}

			str := string(line)

			for _, rex := range t.Options.Exclude {
				if rex.MatchString(str) {
					continue OUTER
				}
			}

			if len(t.Options.Include) != 0 {
				matches := false
				for _, rin := range t.Options.Include {
					if rin.MatchString(str) {
						matches = true
						break
					}
				}
				if !matches {
					continue OUTER
				}
			}

			t.Print(writer, str)
		}
	}()

	go func() {
		<-ctx.Done()
		close(t.closed)
	}()
}

// Close stops tailing
func (t *Tail) Close() {
	close(t.closed)
}

// Print prints a line to the file
func (t *Tail) Print(writer *bufio.Writer, msg string) {
	masked := t.masker.Mask(msg)
	writer.WriteString(masked)
	writer.Flush()
}
