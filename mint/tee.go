package mint

import (
	gtypes "github.com/wetee-dao/go-sdk/gen/types"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// 为 Deployment 节点添加机密设置
func (m *Minter) DeploymentTEEWrap(deployment *appsv1.Deployment, version *gtypes.TEEVersion) {
	if version.IsSGX {
		deployment.Spec.Template.Spec.Containers[0].Resources.Limits["alibabacloud.com/sgx_epc_MiB"] = *resource.NewQuantity(int64(10), resource.DecimalExponent)
		deployment.Spec.Template.Spec.Containers[0].Resources.Requests["alibabacloud.com/sgx_epc_MiB"] = *resource.NewQuantity(int64(10), resource.DecimalExponent)
	} else if version.IsCVM {
		// TODO add TDX
		deployment.Spec.Template.Spec.NodeSelector = map[string]string{"TEE": "CVM-SEV"}
		deployment.Spec.Template.ObjectMeta.Annotations = map[string]string{
			"io.containerd.cri.runtime-handler":                "kata-qemu-sev",
			"io.katacontainers.config.pre_attestation.enabled": "true",
			"io.katacontainers.config.pre_attestation.uri":     "192.168.111.121:30005",
		}
		var KATAQUEMUSEV = "kata-qemu-sev"
		deployment.Spec.Template.Spec.RuntimeClassName = &KATAQUEMUSEV
	}
}

// 为 Pod 节点添加机密设置
func (m *Minter) PodTEEWrap(pod *v1.Pod, version *gtypes.TEEVersion) {
	if version.IsSGX {
		pod.Spec.Containers[0].Resources.Limits["alibabacloud.com/sgx_epc_MiB"] = *resource.NewQuantity(int64(10), resource.DecimalExponent)
		pod.Spec.Containers[0].Resources.Requests["alibabacloud.com/sgx_epc_MiB"] = *resource.NewQuantity(int64(10), resource.DecimalExponent)
	} else if version.IsCVM {
		// TODO add TDX
		pod.Spec.NodeSelector = map[string]string{"TEE": "CVM-SEV"}
	}
}