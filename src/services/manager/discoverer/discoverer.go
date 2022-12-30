package discoverer

import (
	"context"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Deps struct {
	Client *kubernetes.Clientset
}

type WorkerInstance struct {
	Ip   string
	Port int32
}

type Discoverer interface {
	InstancesChan() <-chan []WorkerInstance

	Run(ctx context.Context) error
}

type discoverer struct {
	client *kubernetes.Clientset
	ch     <-chan []WorkerInstance
}

func (d *discoverer) InstancesChan() <-chan []WorkerInstance {
	return d.ch
}

func (d *discoverer) Run(ctx context.Context) error {
	watcher, err := d.client.CoreV1().Endpoints("default").Watch(ctx, metav1.ListOptions{
		LabelSelector: "app = worker",
	})
	if err != nil {
		return err
	}
	ch := make(chan []WorkerInstance)
	d.ch = ch

	go func() {
		for ev := range watcher.ResultChan() {
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

	return nil
}

func New(deps Deps) Discoverer {
	return &discoverer{client: deps.Client}
}
