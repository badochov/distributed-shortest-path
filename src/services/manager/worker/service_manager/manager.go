package service_manager

import (
	"context"
	"fmt"

	"github.com/badochov/distributed-shortest-path/src/libs/db"
	"github.com/hashicorp/go-multierror"
	v1 "k8s.io/api/autoscaling/v1"
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
	GetReplicas(ctx context.Context, region db.RegionId) (replicas int32, err error)
}

type manager struct {
	client                   *kubernetes.Clientset
	namespace                string
	numRegions               int
	workerDeploymentTemplate string
}

func (m *manager) GetReplicas(ctx context.Context, region db.RegionId) (int32, error) {
	name := fmt.Sprintf(m.workerDeploymentTemplate, region)
	s, err := m.getScale(ctx, name)
	if err != nil {
		return 0, err
	}
	return s.Spec.Replicas, nil
}

func (m *manager) Rescale(ctx context.Context, replicas int32) error {
	var err error

	// TOOD [wprzytula] Run the rescales in parallel.
	for i := 0; i < m.numRegions; i++ {
		name := fmt.Sprintf(m.workerDeploymentTemplate, i)
		// TODO [wprzytula] Think if retires should be here on on client sied. Add backoff to retries.
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

func (m *manager) getScale(ctx context.Context, deploymentName string) (*v1.Scale, error) {
	return m.client.AppsV1().Deployments(m.namespace).GetScale(ctx, deploymentName, metav1.GetOptions{})
}

func (m *manager) rescaleDeployment(ctx context.Context, deploymentName string, replicas int32) error {
	s, err := m.getScale(ctx, deploymentName)
	if err != nil {
		return err
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
