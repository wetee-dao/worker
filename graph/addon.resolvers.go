package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.42

import (
	"context"
	"fmt"

	"github.com/vektah/gqlparser/v2/gqlerror"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"wetee.app/worker/mint"
)

// StartLocalWetee is the resolver for the StartLocalWetee field.
func (r *mutationResolver) StartLocalWetee(ctx context.Context, imageVersion string) (bool, error) {
	nameSpace := mint.MinterIns.K8sClient.AppsV1().Deployments("worker-addon")
	deployment := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "wetee-node",
			Annotations: map[string]string{"version": fmt.Sprint(imageVersion)},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "wetee-node"},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "wetee-node"},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "wetee",
							Image: "wetee/wetee-node:" + imageVersion,
							Ports: []v1.ContainerPort{
								{
									Name:          "wetee-node-9944",
									ContainerPort: 9944,
									Protocol:      "TCP",
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := nameSpace.Create(ctx, &deployment, metav1.CreateOptions{})
	if err != nil {
		return false, gqlerror.Errorf(err.Error())
	}

	// 创建集群内部服务
	ServiceSpace := mint.MinterIns.K8sClient.CoreV1().Services("worker-addon")
	service := v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "wetee-node",
		},
		Spec: v1.ServiceSpec{
			Selector:  map[string]string{"app": "wetee-node"},
			ClusterIP: "None",
			Ports: []v1.ServicePort{
				{
					Name:       "wetee-9944",
					Protocol:   "TCP",
					Port:       9944,
					TargetPort: intstr.FromInt(9944),
				},
			},
		},
	}
	_, err = ServiceSpace.Create(ctx, &service, metav1.CreateOptions{})
	if err != nil {
		return false, gqlerror.Errorf(err.Error())
	}

	return true, nil
}

// LinkWetee is the resolver for the LinkWetee field.
func (r *mutationResolver) LinkWetee(ctx context.Context, url string) (bool, error) {
	if url == "" || url == "local" {
		url = mint.DefaultChainUrl
	}

	err := mint.InitChainClient(url)
	if err != nil {
		return false, gqlerror.Errorf(err.Error())
	}
	return true, nil
}

// StartSgxPccs is the resolver for the start_sgx_pccs field.
func (r *mutationResolver) StartSgxPccs(ctx context.Context, imageVersion string, apiKey string) (bool, error) {
	nameSpace := mint.MinterIns.K8sClient.AppsV1().Deployments("worker-addon")
	deployment := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "sgx-pccs",
			Annotations: map[string]string{"version": fmt.Sprint(imageVersion)},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "sgx-pccs"},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "sgx-pccs"},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "pccs",
							Image: "wetee/pccs:" + imageVersion,
							Env: []v1.EnvVar{
								{
									Name:  "APIKEY",
									Value: apiKey,
								},
							},
							Ports: []v1.ContainerPort{
								{
									Name:          "pccs-8081",
									ContainerPort: 8081,
									Protocol:      "TCP",
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := nameSpace.Create(ctx, &deployment, metav1.CreateOptions{})
	if err != nil {
		return false, gqlerror.Errorf(err.Error())
	}

	// 创建集群内部服务
	ServiceSpace := mint.MinterIns.K8sClient.CoreV1().Services("worker-addon")
	service := v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "sgx-pccs",
		},
		Spec: v1.ServiceSpec{
			Selector:  map[string]string{"app": "sgx-pccs"},
			ClusterIP: "None",
			Ports: []v1.ServicePort{
				{
					Name:       "pccs-8081",
					Protocol:   "TCP",
					Port:       8081,
					TargetPort: intstr.FromInt(8081),
				},
			},
		},
	}
	_, err = ServiceSpace.Create(ctx, &service, metav1.CreateOptions{})
	if err != nil {
		return false, gqlerror.Errorf(err.Error())
	}
	return true, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
