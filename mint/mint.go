package mint

import (
	"fmt"
	"sync"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"k8s.io/client-go/kubernetes"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"wetee.app/worker/db"
	"wetee.app/worker/mint/chain"
	"wetee.app/worker/mint/chain/gen/system"
	"wetee.app/worker/mint/chain/gen/types"
	"wetee.app/worker/util"
)

// Minter
// 矿工
type Minter struct {
	K8sClient     *kubernetes.Clientset
	MetricsClient *versioned.Clientset
	ChainClient   *chain.ChainClient
}

var (
	MinterIns *Minter
	lock      sync.Mutex
	Signer    *signature.KeyringPair
)

// InitMint
// 初始化矿工
func InitMint(mgr manager.Manager) error {
	// 创建K8s Client
	clientset, err := kubernetes.NewForConfig(mgr.GetConfig())
	if err != nil {
		return err
	}

	// 创建Metrics Client
	metricsClient, err := versioned.NewForConfig(mgr.GetConfig())
	if err != nil {
		return err
	}

	client, err := chain.ClientInit()
	if err != nil {
		return err
	}

	// 初始化minter
	lock.Lock()
	MinterIns = &Minter{
		K8sClient:     clientset,
		MetricsClient: metricsClient,
		ChainClient:   client,
	}
	lock.Unlock()

	// 获取签名账户
	Signer, err = chain.GetMintKey()
	return err
}

// start mint
// 开始挖矿
func StartMint() {
	fmt.Println("MintKey => ", Signer.Address)
	client := MinterIns.ChainClient
	chainAPI := client.Api

	// 获取挖矿状态
	worker := chain.Worker{
		Client: client,
		Signer: Signer,
	}

	// 挖矿开始
mintStart:

	// 等待集群开启
	for {
		clusterId, err := worker.Getk8sClusterAccounts(Signer.PublicKey)
		if err != nil {
			fmt.Println("ClusterId => ", err)
			time.Sleep(time.Second * 10)
			continue
		}
		db.SetClusterId(clusterId)

		break
	}

	clusterId, _ := db.GetClusterId()
	fmt.Println("ClusterId => ", clusterId)

	// 触发区块监听
	sub, err := chainAPI.RPC.Chain.SubscribeFinalizedHeads()

	// 失败后等待10秒重新尝试
	if err != nil {
		util.LogWithRed("SubscribeNewHeads", err)
		time.Sleep(time.Second * 10)
		goto mintStart
	}
	defer sub.Unsubscribe()

	for {
		head := <-sub.Chan()
		fmt.Printf("Chain is at block: #%v\n", head.Number)
		blockHash, _ := chainAPI.RPC.Chain.GetBlockHash(uint64(head.Number))

		// 读取/处理新的区块信息
		events, err := system.GetEvents(chainAPI.RPC.State, blockHash)
		if err != nil {
			util.LogWithRed("GetEventsLatest", err)
			continue
		}

		for _, event := range events {
			e := event.Event
			if e.IsWeteeWorker {
				startEvent := e.AsWeteeWorkerField0
				if startEvent.IsWorkRuning {
					workId := startEvent.AsWorkRuningWorkId1
					user := startEvent.AsWorkRuningUser0
					fmt.Println("===========================================WorkRuning ID: ", workId)
					err = CreateOrUpdatePod(user[:], workId, blockHash.Hex())
					fmt.Println("===========================================CreateOrUpdatePod error: ", err)
				}
			}
			if e.IsWeteeApp {
				appEvent := e.AsWeteeAppField0
				if appEvent.IsWorkStopped {
					workId := appEvent.AsWorkStoppedWorkId1
					fmt.Println("===========================================WorkStopped", workId)
					err := StopPod(workId)
					fmt.Println("===========================================StopPod error: ", err)
				}
				if appEvent.IsWorkUpdated {
					workId := appEvent.AsWorkUpdatedWorkId1
					user := appEvent.AsWorkUpdatedUser0
					fmt.Println("===========================================WorkUpdated ID: ", workId)
					err = CreateOrUpdatePod(user[:], workId, blockHash.Hex())
					fmt.Println("===========================================CreateOrUpdatePod error: ", err)
				}
				// e.AsWeteeAppField0.IsClusterCreated
			}
			// if e.IsWeteeTask {
			// fmt.Println("e.AsWeteeTaskField0")
			// }
		}

		// 获取合约列表
		cs, err := worker.GetClusterContracts(clusterId, &blockHash)
		fmt.Println("GetClusterContracts", cs)
		if err != nil {
			fmt.Println("GetClusterContracts", err)
			continue
		}
		stage, err := worker.GetStage()
		if err != nil {
			fmt.Println("GetStage", err)
			continue
		}

		// 校对合约状态
		for _, c := range cs {
			err := checkPodStatus(c, blockHash.Hex())
			fmt.Println("checkPodStatus", err)

			state, err := worker.GetWorkContract(c.ContractState.WorkId, clusterId)

			// 上传工作证明
			if err != nil && (uint64(head.Number)-state.BlockNumber) >= uint64(stage) {
				fmt.Println("===========================================WorkProofUpload")
				worker.WorkProofUpload(c.ContractState.WorkId, []string{}, []string{}, []byte(""))
			}
		}
	}
}

func GetWorkTypeStr(work types.WorkId) string {
	if work.Wtype.IsAPP {
		return "app"
	}

	if work.Wtype.IsTASK {
		return "task"
	}

	return "unknown"
}

// func getPodInfo() {
// podLogOpts := &corev1.PodLogOptions{
// 	Container: "worker",
// 	SinceTime: &metav1.Time{
// 		Time: time.Now().Add(-1 * time.Minute),
// 	},
// }
// req := clientset.CoreV1().Pods("default").GetLogs("testnginx", podLogOpts)
// podLogs, err := req.Stream(ctx)
// if err != nil {
// 	fmt.Println("error in opening stream " + err.Error())
// 	return
// }
// defer podLogs.Close()

// // Read the logs line by line
// logs := ""
// scanner := bufio.NewScanner(podLogs)
// for scanner.Scan() {
// 	logs += scanner.Text() + "\n"
// }
// if err := scanner.Err(); err != nil {
// 	fmt.Printf("failed to read log line: %v", err)
// }

// fmt.Println("logs: ================================================")
// fmt.Println(logs)
// fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")

// 获取Pod的内存使用情况
// podMetrics, err := metricsClient.MetricsV1beta1().PodMetricses("default").List(ctx, metav1.ListOptions{})
// if err != nil {
// 	fmt.Println(err)
// 	return
// }

// for i := 0; i < len(podMetrics.Items); i++ {
// 	pod := podMetrics.Items[i]
// 	// 遍历Pod的容器，获取内存使用情况
// 	for _, container := range pod.Containers {
// 		fmt.Printf("Pod %s CPU使用情况: %s \n", pod.Name, container.Usage.Cpu())
// 		fmt.Printf("Pod %s 内存使用情况: %s \n", pod.Name, container.Usage.Memory())
// 	}
// }
// }