package discoverer

import (
	"context"
	"github.com/badochov/distributed-shortest-path/src/libs/db"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type WorkerInstanceStatus struct {
	Status   v1.PodPhase
	Ip       string
	RegionId int
}

func (w *WorkerInstanceStatus) IsRunning() bool {
	return w.Status == v1.PodRunning
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
	RegionId  db.RegionId
	Instances []WorkerInstance
}

type Discoverer interface {
	// RegionDataChan returns channel streaming data regarding updates of set of instances.
	RegionDataChan() <-chan RegionData
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

func (d *discoverer) RegionDataChan() <-chan RegionData {
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

	// Endpoints
	ch := make(chan RegionData)
	d.ch = ch

	go func() {
		for ev := range endpointWatcher.ResultChan() {
			endpoints := ev.Object.(*v1.Endpoints)

			regStr := endpoints.Labels["region"]
			regId, err := strconv.ParseUint(regStr, 10, 16)
			if err != nil {
				panic(err)
			}
			res := RegionData{RegionId: uint16(regId)}

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
			regStr := pod.Labels["region"]
			regId, err := strconv.Atoi(regStr)
			if err != nil {
				panic(err)
			}

			podCh <- WorkerInstanceStatus{
				Status:   pod.Status.Phase,
				Ip:       pod.Status.PodIP,
				RegionId: regId,
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
