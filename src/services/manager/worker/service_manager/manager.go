package service_manager

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/badochov/distributed-shortest-path/src/libs/db"
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
	Rescale(ctx context.Context, region db.RegionId, replicas int32) error
	GetReplicas(ctx context.Context, region db.RegionId) (replicas int32, err error)
	WaitForRescale(ctx context.Context, region db.RegionId) error
}

type manager struct {
	client                   *kubernetes.Clientset
	namespace                string
	numRegions               int
	workerDeploymentTemplate string
}

func (m *manager) GetReplicas(ctx context.Context, region db.RegionId) (int32, error) {
	name := m.getDeploymentName(region)
	s, err := m.getScale(ctx, name)
	if err != nil {
		return 0, err
	}
	return s.Spec.Replicas, nil
}

func (m *manager) getDeploymentName(region db.RegionId) string {
	return fmt.Sprintf(m.workerDeploymentTemplate, region)
}

func (m *manager) getScale(ctx context.Context, deploymentName string) (*v1.Scale, error) {
	return m.client.AppsV1().Deployments(m.namespace).GetScale(ctx, deploymentName, metav1.GetOptions{})
}

func (m *manager) Rescale(ctx context.Context, regionId db.RegionId, replicas int32) error {
	deploymentName := m.getDeploymentName(regionId)
	s, err := m.getScale(ctx, deploymentName)
	if err != nil {
		return err
	}

	sc := *s
	sc.Spec.Replicas = replicas

	_, err = m.client.AppsV1().Deployments(m.namespace).UpdateScale(ctx, deploymentName, &sc, metav1.UpdateOptions{})
	return err
}

func (m *manager) WaitForRescale(ctx context.Context, regionId db.RegionId) error {
	deploymentName := m.getDeploymentName(regionId)

	for ctx.Err() == nil {
		s, err := m.client.AppsV1().Deployments(m.namespace).Get(ctx, deploymentName, metav1.GetOptions{})
		if err != nil {
			return err
		}
		if s.Spec.Replicas == nil {
			log.Printf("[STATUS][Region-%d] Nil spec replicas\n", regionId)
			continue
		}
		log.Printf("[STATUS][Region-%d] Spec replicas: %d, Status replicas %d, Ready replicas %d\n", regionId, *s.Spec.Replicas, s.Status.Replicas, s.Status.ReadyReplicas)
		if *s.Spec.Replicas == s.Status.Replicas && s.Status.ReadyReplicas == s.Status.Replicas {
			return nil
		}
		time.Sleep(time.Second)
	}
	return errors.New("timeout waiting for rescale")
}

func New(deps Deps) WorkerServiceManager {
	return &manager{
		client:                   deps.Client,
		namespace:                deps.Namespace,
		numRegions:               deps.NumRegions,
		workerDeploymentTemplate: deps.WorkerDeploymentTemplate,
	}
}
