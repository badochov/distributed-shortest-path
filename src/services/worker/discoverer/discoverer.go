package discoverer

import (
	"context"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type WorkerInstanceStatus struct {
	Status v1.PodPhase
	Ip     string
	Name   string
}

type Deps struct {
	Client        *kubernetes.Clientset
	Namespace     string
	LabelSelector string
}

type WorkerInstance struct {
	Ip string
}

type RegionData struct {
	RegionId  int
	Instances []WorkerInstance
}

type Discoverer interface {
	// InstancesChan returns channel streaming data regarding updates of set of instances.
	InstancesChan() <-chan RegionData
	// InstanceStatuses returns channel streaming data regarding updates of statuses on instances.
	// Useful for checking health of the worker.
	InstanceStatuses() <-chan WorkerInstanceStatus

	Run(ctx context.Context) error
}

type discoverer struct {
	client        *kubernetes.Clientset
	ch            <-chan RegionData
	statuses      <-chan WorkerInstanceStatus
	labelSelector string
	namespace     string
}

func (d *discoverer) InstanceStatuses() <-chan WorkerInstanceStatus {
	return d.statuses
}

func (d *discoverer) InstancesChan() <-chan RegionData {
	return d.ch
}

func (d *discoverer) Run(ctx context.Context) error {
	endpointWatcher, err := d.client.CoreV1().Endpoints(d.namespace).Watch(ctx, metav1.ListOptions{
		LabelSelector: d.labelSelector,
	})
	if err != nil {
		return err
	}
	podWatcher, err := d.client.CoreV1().Pods(d.namespace).Watch(ctx, metav1.ListOptions{
		LabelSelector: d.labelSelector,
	})
	if err != nil {
		return err
	}
	//d.client.AppsV1().Deployments().ApplyScale()

	// Endpoints
	ch := make(chan RegionData)
	d.ch = ch

	go func() {
		for ev := range endpointWatcher.ResultChan() {
			endpoints := ev.Object.(*v1.Endpoints)

			var res RegionData
			var err error

			regStr := endpoints.Labels["region"]
			res.RegionId, err = strconv.Atoi(regStr)
			if err != nil {
				panic(err)
			}

			for _, subset := range endpoints.Subsets {
				for _, address := range subset.Addresses {
					res.Instances = append(res.Instances, WorkerInstance{
						Ip: address.IP,
					})
				}
			}
			ch <- res
		}
	}()

	// Pods
	podCh := make(chan WorkerInstanceStatus)
	d.statuses = podCh
	go func() {
		for ev := range podWatcher.ResultChan() {
			pod := ev.Object.(*v1.Pod)
			podCh <- WorkerInstanceStatus{
				Status: pod.Status.Phase,
				Ip:     pod.Status.PodIP,
				Name:   pod.Name,
			}
		}
	}()

	return nil
}

func New(deps Deps) Discoverer {
	return &discoverer{
		client:        deps.Client,
		namespace:     deps.Namespace,
		labelSelector: deps.LabelSelector,
	}
}
