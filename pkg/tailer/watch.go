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
// https://github.com/wercker/stern/blob/master/stern/watch.go

package tailer

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/typed/core/v1"
)

// Target is a target to watch
type Target struct {
	Namespace string
	Pod       string
	Container string
	App       string
}

// GetID returns the ID of the object
func (t *Target) GetID() string {
	return fmt.Sprintf("%s-%s-%s", t.Namespace, t.Pod, t.Container)
}

// Watch starts listening to Kubernetes events and emits modified
// containers/pods. The first result is targets added, the second is targets
// removed
func (o *Options) Watch(ctx context.Context, i v1.PodInterface, labelSelector labels.Selector) (chan *Target, chan *Target, error) {
	watcher, err := i.Watch(metav1.ListOptions{Watch: true, LabelSelector: labelSelector.String()})
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to set up watch")
	}

	added := make(chan *Target)
	removed := make(chan *Target)

	go func() {
		for {
			select {
			case e := <-watcher.ResultChan():
				if e.Object == nil {
					// Closed because of error
					return
				}

				pod := e.Object.(*corev1.Pod)
				if pod == nil {
					continue
				}

				if !o.MatchPod(pod) {
					continue
				}

				switch e.Type {
				case watch.Added, watch.Modified:
					var statuses []corev1.ContainerStatus
					statuses = append(statuses, pod.Status.InitContainerStatuses...)
					statuses = append(statuses, pod.Status.ContainerStatuses...)

					for _, c := range statuses {
						if !o.MatchesContainerStatus(pod, c) {
							continue
						}

						app := ""
						if pod.Labels != nil {
							app = pod.Labels["app"]
							if app == "" {
								app = pod.Labels["app.kubernetes.io/managed-by"]
							}
							if app == "tekton-pipelines" {
								owner := pod.Labels["owner"]
								repository := pod.Labels["repository"]
								branch := pod.Labels["branch"]
								if owner != "" {
									app = filepath.Join(app, owner)
								}
								if repository != "" {
									app = filepath.Join(app, repository)
								}
								if branch != "" {
									app = filepath.Join(app, branch)
								}
							}
						}
						added <- &Target{
							Namespace: pod.Namespace,
							Pod:       pod.Name,
							Container: c.Name,
							App:       app,
						}
					}
				case watch.Deleted:
					var containers []corev1.Container
					containers = append(containers, pod.Spec.Containers...)
					containers = append(containers, pod.Spec.InitContainers...)

					for _, c := range containers {
						if !o.MatchesContainer(pod, c) {
							continue
						}

						removed <- &Target{
							Namespace: pod.Namespace,
							Pod:       pod.Name,
							Container: c.Name,
						}
					}
				}
			case <-ctx.Done():
				watcher.Stop()
				close(added)
				close(removed)
				return
			}
		}
	}()

	return added, removed, nil
}
