package service_manager

import (
	"context"

	"github.com/hashicorp/go-multierror"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Deps struct {
	Client                   *kubernetes.Clientset
	Namespace                string
	NumRegions               int
	WorkerDeploymentTemplate string
}

type WorkerServiceManager interface {
	// Rescale rescales worker service to desired amount of replicas. Rescale to 0 to shut down.
	Rescale(ctx context.Context, replicas int32) error
}

type manager struct {
	client                   *kubernetes.Clientset
	namespace                string
	numRegions               int
	workerDeploymentTemplate string
}

func (m *manager) Rescale(ctx context.Context, replicas int32) error {
	var err error

	for i := 0; i < m.numRegions; i++ {
		name := ""

		const retries = 3
		if scaleErr := m.rescaleDeploymentWithRetries(ctx, name, replicas, retries); scaleErr != nil {
			err = multierror.Append(scaleErr, err)
		}
	}

	return err
}

func (m *manager) rescaleDeploymentWithRetries(ctx context.Context, deploymentName string, replicas int32, retries int) error {
	var err error

	for i := 0; i < retries; i++ {
		if rescaleErr := m.rescaleDeployment(ctx, deploymentName, replicas); rescaleErr != nil {
			err = multierror.Append(err, rescaleErr)
		} else {
			return nil
		}
	}
	return err
}

func (m *manager) rescaleDeployment(ctx context.Context, deploymentName string, replicas int32) error {
	s, err := m.client.AppsV1().Deployments(m.namespace).GetScale(ctx, deploymentName, metav1.GetOptions{})
	if err != nil {
		return nil
	}

	sc := *s
	sc.Spec.Replicas = replicas

	_, err = m.client.AppsV1().Deployments(m.namespace).UpdateScale(ctx, deploymentName, &sc, metav1.UpdateOptions{})
	return err
}

func New(deps Deps) WorkerServiceManager {
	return &manager{
		client:                   deps.Client,
		namespace:                deps.Namespace,
		numRegions:               deps.NumRegions,
		workerDeploymentTemplate: deps.WorkerDeploymentTemplate,
	}
}
