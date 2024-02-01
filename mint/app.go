package mint

import (
	"context"
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
	chain "github.com/wetee-dao/go-sdk"
	gtype "github.com/wetee-dao/go-sdk/gen/types"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"wetee.app/worker/dao"
	"wetee.app/worker/util"
)

func (m *Minter) DoWithAppState(ctx *context.Context, c ContractStateWrap, stage uint32, head types.Header) error {
	if c.App == nil || c.WorkState == nil {
		return errors.New("app is nil")
	}

	app := c.App
	state := c.WorkState

	// 状态为停止状态，停止Pod
	if uint64(app.Status) == 2 {
		m.StopApp(c.ContractState.WorkId)
		return nil
	}

	_, err := m.CheckAppStatus(ctx, c)
	if err != nil {
		util.LogWithRed("checkPodStatus", err)
		return err
	}

	// 判断是否上传工作证明
	// Check if work proof needs to be uploaded
	if uint64(app.Status) == 1 && uint64(head.Number)-state.BlockNumber >= uint64(stage) {
		util.LogWithRed("=========================================== WorkProofUpload APP")

		workId := c.ContractState.WorkId
		name := util.GetWorkTypeStr(workId) + "-" + fmt.Sprint(workId.Id)
		nameSpace := AccountToAddress(c.ContractState.User[:])

		// 获取log和硬件资源使用量
		// Get log and hardware resource usage
		logs, crs, err := m.getMetricInfo(*ctx, workId, nameSpace, name, uint64(head.Number)-state.BlockNumber)
		if err != nil {
			util.LogWithRed("getMetricInfo", err)
			return err
		}

		// 获取log hash
		// Get log hash
		logHash, err := getWorkLogHash(name, logs, state.BlockNumber)
		if err != nil {
			util.LogWithRed("getWorkLogHash", err)
			return err
		}

		// 获取计算资源hash
		// Get Computing resource hash
		crHash, cr, err := getWorkCrHash(name, crs, state.BlockNumber)
		if err != nil {
			util.LogWithRed("getWorkCrHash", err)
			return err
		}

		// 初始化worker对象
		worker := chain.Worker{
			Client: m.ChainClient,
			Signer: Signer,
		}

		// 上传工作证明
		// Upload work proof
		err = worker.WorkProofUpload(c.ContractState.WorkId, logHash, crHash, gtype.Cr{
			Cpu:  cr[0],
			Mem:  cr[1],
			Disk: 0,
		}, []byte(""), false)
		if err != nil {
			util.LogWithRed("WorkProofUpload", err)
			return err
		}
	}

	return nil
}

// checkAppStatus check app status
// 校对应用状态
func (m *Minter) CheckAppStatus(ctx *context.Context, state ContractStateWrap) (*v1.Pod, error) {
	address := AccountToAddress(state.ContractState.User[:])
	nameSpace := m.K8sClient.CoreV1().Pods(address)
	workID := state.ContractState.WorkId
	name := util.GetWorkTypeStr(workID) + "-" + fmt.Sprint(workID.Id)

	app := state.App
	if uint8(app.Status) == 2 {
		m.StopApp(workID)
		return nil, errors.New("app stop")
	}

	pod, err := nameSpace.Get(*ctx, name, metav1.GetOptions{})
	if err != nil && err.Error() != "pods \""+name+"\" not found" {
		return nil, err
	}

	version := state.Version
	if pod.ObjectMeta.Annotations["version"] != fmt.Sprint(version) {
		err = m.CreateApp(ctx, state.ContractState.User[:], workID, app, version)
		if err != nil {
			return nil, err
		}
		pod, err = nameSpace.Get(*ctx, name, metav1.GetOptions{})
	}

	return pod, err
}

// CreateOrUpdateApp create or update app
// 校对应用链上状态后创建或更新应用
func (m *Minter) CreateApp(ctx *context.Context, user []byte, workID gtype.WorkId, app *gtype.TeeApp, version uint64) error {
	saddress := AccountToAddress(user)
	errc := m.checkNameSpace(*ctx, saddress)
	if errc != nil {
		return errc
	}

	nameSpace := m.K8sClient.CoreV1().Pods(saddress)
	name := util.GetWorkTypeStr(workID) + "-" + fmt.Sprint(workID.Id)

	err := dao.SetSecrets(workID, &dao.Secrets{
		Env: map[string]string{
			"": "",
		},
	})
	if err != nil {
		return err
	}

	existingPod, err := nameSpace.Get(*ctx, name, metav1.GetOptions{})
	if err == nil {
		if uint8(app.Status) == 2 {
			m.StopApp(workID)
			return nil
		}
		existingPod.ObjectMeta.Annotations = map[string]string{
			"version": fmt.Sprint(version),
		}
		existingPod.Spec.Containers[0].Image = string(app.Image)
		existingPod.Spec.Containers[0].Ports[0].ContainerPort = int32(app.Port[0])
		_, err = nameSpace.Update(*ctx, existingPod, metav1.UpdateOptions{})
		fmt.Println("================================================= Update", err)
	} else {
		// 用于应用联系控制面板的凭证
		wid, err := dao.SealAppID(workID)
		if err != nil {
			return err
		}
		pod := &v1.Pod{
			TypeMeta: metav1.TypeMeta{
				Kind:       "App",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
				Annotations: map[string]string{
					"version": fmt.Sprint(version),
				},
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name:  "c1",
						Image: string(app.Image),
						Ports: []v1.ContainerPort{
							{
								Name:          string(app.Name) + "0",
								ContainerPort: int32(app.Port[0]),
								Protocol:      "TCP",
							},
						},
						Env: []v1.EnvVar{
							{
								Name:  "APPID",
								Value: wid,
							},
							{
								Name:  "IN_TEE",
								Value: string("1"),
							},
						},
						Resources: v1.ResourceRequirements{
							Limits: v1.ResourceList{
								"alibabacloud.com/sgx_epc_MiB": *resource.NewQuantity(int64(20), resource.DecimalExponent),
							},
							Requests: v1.ResourceList{
								"alibabacloud.com/sgx_epc_MiB": *resource.NewQuantity(int64(20), resource.DecimalExponent),
							},
						},
					},
				},
			},
		}
		_, err = nameSpace.Create(*ctx, pod, metav1.CreateOptions{})
		fmt.Println("================================================= Create", err)
	}

	return err
}

func (m *Minter) UpdateApp(ctx *context.Context, user []byte, workID gtype.WorkId, app *gtype.TeeApp, version uint64) error {
	saddress := AccountToAddress(user)
	nameSpace := m.K8sClient.CoreV1().Pods(saddress)
	name := util.GetWorkTypeStr(workID) + "-" + fmt.Sprint(workID.Id)

	existingPod, err := nameSpace.Get(*ctx, name, metav1.GetOptions{})
	if err == nil {
		if uint8(app.Status) == 2 {
			m.StopApp(workID)
			return nil
		}
		existingPod.ObjectMeta.Annotations = map[string]string{
			"version": fmt.Sprint(version),
		}
		existingPod.Spec.Containers[0].Image = string(app.Image)
		existingPod.Spec.Containers[0].Ports[0].ContainerPort = int32(app.Port[0])
		_, err = nameSpace.Update(*ctx, existingPod, metav1.UpdateOptions{})
		fmt.Println("================================================= Update", err)
	}

	return err
}

// StopApp
// 停止应用
func (m *Minter) StopApp(workID gtype.WorkId) error {
	ctx := context.Background()
	user, err := chain.GetAccount(m.ChainClient, workID)
	if err != nil {
		return err
	}

	saddress := AccountToAddress(user[:])

	nameSpace := m.K8sClient.CoreV1().Pods(saddress)
	name := util.GetWorkTypeStr(workID) + "-" + fmt.Sprint(workID.Id)
	return nameSpace.Delete(ctx, name, metav1.DeleteOptions{})
}
