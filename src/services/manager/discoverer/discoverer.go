package discoverer

import (
	"context"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type WorkerInstanceStatus struct {
	Status v1.PodPhase
	Ip     string
	Name   string
}

type Deps struct {
	Client *kubernetes.Clientset
}

type WorkerInstance struct {
	Ip   string
	Port int32
}

type Discoverer interface {
	InstancesChan() <-chan []WorkerInstance
	InstanceStatuses() <-chan WorkerInstanceStatus

	Run(ctx context.Context) error
}

type discoverer struct {
	client   *kubernetes.Clientset
	ch       <-chan []WorkerInstance
	statuses <-chan WorkerInstanceStatus
}

func (d *discoverer) InstanceStatuses() <-chan WorkerInstanceStatus {
	return d.statuses
}

func (d *discoverer) InstancesChan() <-chan []WorkerInstance {
	return d.ch
}

func (d *discoverer) Run(ctx context.Context) error {
	endpointWatcher, err := d.client.CoreV1().Endpoints("default").Watch(ctx, metav1.ListOptions{
		LabelSelector: "app = worker-mesh",
	})
	if err != nil {
		return err
	}
	podWatcher, err := d.client.CoreV1().Pods("default").Watch(ctx, metav1.ListOptions{
		LabelSelector: "app = worker",
	})
	if err != nil {
		return err
	}

	// Endpoints
	ch := make(chan []WorkerInstance)
	d.ch = ch

	go func() {
		for ev := range endpointWatcher.ResultChan() {
			endpoints := ev.Object.(*v1.Endpoints)

			var instances []WorkerInstance
			for _, subset := range endpoints.Subsets {
				var port int32
				for _, p := range subset.Ports {
					if p.Name == "worker-rpc" {
						port = p.Port
						break
					}
				}
				if port == 0 {
					panic("Port cannot be 0")
				}
				for _, address := range subset.Addresses {
					instances = append(instances, WorkerInstance{
						Ip:   address.IP,
						Port: port,
					})
				}
			}
			ch <- instances
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
	return &discoverer{client: deps.Client}
}
