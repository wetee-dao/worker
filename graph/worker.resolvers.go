package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.42

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	subkey "github.com/vedhavyas/go-subkey/v2"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"github.com/wetee-dao/go-sdk/core"
	"github.com/wetee-dao/go-sdk/module"
	"github.com/wetee-dao/go-sdk/pallet/balances"
	gtypes "github.com/wetee-dao/go-sdk/pallet/types"
	"github.com/wetee-dao/go-sdk/pallet/weteedsecret"
	"github.com/wetee-dao/go-sdk/pallet/weteeworker"
	"wetee.app/worker/graph/model"
	"wetee.app/worker/internal/store"
	"wetee.app/worker/mint"
	"wetee.app/worker/mint/proof"
	"wetee.app/worker/util"
)

// ClusterRegister is the resolver for the cluster_register field.
func (r *mutationResolver) ClusterRegister(ctx context.Context, name string, ip string, domain string, port int, level int) (string, error) {
	if mint.MinterIns.ChainClient == nil {
		return "", gqlerror.Errorf("Invalid chain client")
	}
	client := mint.MinterIns.ChainClient
	if client == nil {
		return "", gqlerror.Errorf("Cant connect to chain")
	}
	worker := &module.Worker{
		Client: client,
		Signer: mint.MinterIns.Signer,
	}

	ipstrs := strings.Split(ip, ".")
	if len(ipstrs) != 4 {
		return "", gqlerror.Errorf("Ip address format error")
	}

	iparr := []uint8{}
	for _, ipstr := range ipstrs {
		i, err := strconv.Atoi(ipstr)
		if err != nil {
			return "", gqlerror.Errorf("Ip address int transfer error")
		}
		iparr = append(iparr, uint8(i))
	}

	// iparr
	err := worker.ClusterRegister(name, []gtypes.Ip{
		{
			Ipv4: gtypes.OptionTUint32{
				IsNone: true,
			},
			Ipv6: gtypes.OptionTU128{
				IsNone: true,
			},
			Domain: gtypes.OptionTByteSlice{
				IsSome:       true,
				AsSomeField0: []byte(domain),
			},
		},
	}, uint32(port), uint8(level), false)

	if err != nil {
		return "", gqlerror.Errorf("Chain call error:" + err.Error())
	}
	return "ok", nil
}

// ClusterMortgage is the resolver for the cluster_mortgage field.
func (r *mutationResolver) ClusterMortgage(ctx context.Context, cpu int, mem int, cvmCPU int, cvmMem int, disk int, gpu int, deposit int64) (string, error) {
	if mint.MinterIns.ChainClient == nil {
		return "", gqlerror.Errorf("Invalid chain client")
	}
	client := mint.MinterIns.ChainClient
	if client == nil {
		return "", gqlerror.Errorf("Cant connect to chain")
	}
	worker := &module.Worker{
		Client: client,
		Signer: mint.MinterIns.Signer,
	}

	id, err := store.GetClusterId()
	if err != nil {
		return "", gqlerror.Errorf("Cant get cluster id:" + err.Error())
	}
	err = worker.ClusterMortgage(id, uint32(cpu), uint32(mem), uint32(cvmCPU), uint32(cvmMem), uint32(disk), uint32(gpu), uint64(deposit), false)
	if err != nil {
		return "", gqlerror.Errorf("Chain call error:" + err.Error())
	}
	return "ok", nil
}

// ClusterUnmortgage is the resolver for the cluster_unmortgage field.
func (r *mutationResolver) ClusterUnmortgage(ctx context.Context, id int64) (string, error) {
	if mint.MinterIns.ChainClient == nil {
		return "", gqlerror.Errorf("Invalid chain client")
	}
	client := mint.MinterIns.ChainClient
	if client == nil {
		return "", gqlerror.Errorf("Cant connect to chain")
	}
	worker := &module.Worker{
		Client: client,
		Signer: mint.MinterIns.Signer,
	}

	clusterID, err := store.GetClusterId()
	if err != nil {
		return "", gqlerror.Errorf("Cant get cluster id:" + err.Error())
	}

	err = worker.ClusterUnmortgage(clusterID, uint64(id), false)
	if err != nil {
		return "", gqlerror.Errorf("Chain call error:" + err.Error())
	}
	return "ok", nil
}

// ClusterWithdrawal is the resolver for the cluster_withdrawal field.
func (r *mutationResolver) ClusterWithdrawal(ctx context.Context, id int64, ty model.WorkType, val int64) (string, error) {
	if mint.MinterIns.ChainClient == nil {
		return "", gqlerror.Errorf("Invalid chain client")
	}
	client := mint.MinterIns.ChainClient
	if client == nil {
		return "", gqlerror.Errorf("Cant connect to chain")
	}
	worker := &module.Worker{
		Client: client,
		Signer: mint.MinterIns.Signer,
	}

	err := worker.ClusterWithdrawal(gtypes.WorkId{
		Wtype: gtypes.WorkType{IsAPP: ty == model.WorkTypeApp, IsTASK: ty == model.WorkTypeTask},
		Id:    uint64(id),
	}, val, false)
	if err != nil {
		return "", gqlerror.Errorf("Chain call error:" + err.Error())
	}
	return "ok", nil
}

// ClusterStop is the resolver for the cluster_stop field.
func (r *mutationResolver) ClusterStop(ctx context.Context) (string, error) {
	if mint.MinterIns.ChainClient == nil {
		return "", gqlerror.Errorf("Invalid chain client")
	}
	client := mint.MinterIns.ChainClient
	if client == nil {
		return "", gqlerror.Errorf("Cant connect to chain")
	}
	worker := &module.Worker{
		Client: client,
		Signer: mint.MinterIns.Signer,
	}

	clusterID, err := store.GetClusterId()
	if err != nil {
		return "", gqlerror.Errorf("Cant get cluster id:" + err.Error())
	}

	err = worker.ClusterStop(clusterID, false)
	if err != nil {
		return "", gqlerror.Errorf("Chain call error:" + err.Error())
	}
	return "ok", nil
}

// StartForTest is the resolver for the start_for_test field.
func (r *mutationResolver) StartForTest(ctx context.Context) (bool, error) {
	client := mint.MinterIns.ChainClient
	if client == nil {
		return false, gqlerror.Errorf("Cant connect to chain")
	}

	gsigner := mint.MinterIns.Signer
	worker := &module.Worker{
		Client: client,
		Signer: mint.MinterIns.Signer,
	}

	// 1 unit of transfer
	bal, ok := new(big.Int).SetString("50000000000000000", 10)
	if !ok {
		panic(fmt.Errorf("failed to convert balance"))
	}

	minter, _ := types.NewMultiAddressFromAccountID(gsigner.PublicKey)
	minterWrap := gtypes.MultiAddress{
		IsId:       true,
		AsIdField0: minter.AsID,
	}
	c := balances.MakeTransferKeepAliveCall(minterWrap, types.NewUCompact(bal))
	signer, err := core.Sr25519PairFromSecret("//Alice", 42)
	if err != nil {
		return false, errors.New("Cant get signer:" + err.Error())
	}

	err = client.SignAndSubmit(&signer, c, false)
	if err != nil {
		return false, errors.New("Chain call error:" + err.Error())
	}
	time.Sleep(8 * time.Second)

	// 注册集群
	err = worker.ClusterRegister("baiL", []gtypes.Ip{
		{
			Ipv4: gtypes.OptionTUint32{
				IsNone: true,
			},
			Ipv6: gtypes.OptionTU128{
				IsNone: true,
			},
			Domain: gtypes.OptionTByteSlice{
				IsSome:       true,
				AsSomeField0: []byte("xiaobai.asyou.me"),
			},
		},
	}, uint32(30000), uint8(1), false)
	if err != nil {
		return false, errors.New("Chain ClusterRegister error:" + err.Error())
	}
	time.Sleep(10 * time.Second)

	// 获取集群id
	clusterId, err := worker.Getk8sClusterAccounts(gsigner.Public())
	if err != nil {
		return false, errors.New("Getk8sClusterAccounts:" + err.Error())
	}
	fmt.Println("ClusterId => ", clusterId)

	// 抵押集群
	err = worker.ClusterMortgage(
		clusterId, uint32(10000), uint32(10000),
		uint32(1000000), uint32(10000), uint32(1000000),
		uint32(10), uint64(1000000000000), false,
	)
	if err != nil {
		return false, errors.New("Chain ClusterMortgage error:" + err.Error())
	}

	// 创建 P2P 节点 p2p boot nodes
	port := util.GetEnvInt("P2P_PORT", 8881)
	call := weteeworker.MakeSetBootPeersCall([]gtypes.P2PAddr{
		{
			Ip: gtypes.Ip1{
				Ipv4: gtypes.OptionTUint32{
					IsNone: true,
				},
				Ipv6: gtypes.OptionTU128{
					IsNone: true,
				},
				Domain: gtypes.OptionTByteSlice{
					IsSome:       true,
					AsSomeField0: []byte("xiaobai.asyou.me"),
				},
			},
			Port: uint16(port),
			Id:   minter.AsID,
		},
	})

	err = client.SignAndSubmit(&signer, call, false)
	if err != nil {
		return false, errors.New("Chain call error:" + err.Error())
	}
	time.Sleep(10 * time.Second)

	var nodes []string = []string{
		"5GmiTJQfjKQnmoVQBFTFoBVTqGyJb8vJQRx5FTGxySMbytkt",
		"5Fk55vz7hXNKgize4yGoyqixiYbgj6djdDGHktw6ggKPePqE",
		"5FYJhLSqFRXRy17oHKgKjx7BvjMGXNu8ndAGt56ipwFSUNqu",
	}

	for _, node := range nodes {
		_, k, _ := subkey.SS58Decode(node)
		var a [32]byte
		copy(a[:], k)

		// 创建 dsecret 节点
		call = weteedsecret.MakeRegisterNodeCall(a)
		err = client.SignAndSubmit(&signer, call, false)
		if err != nil {
			return false, errors.New("Chain call error:" + err.Error())
		}
		time.Sleep(10 * time.Second)
	}

	return true, nil
}

// WorkerInfo is the resolver for the workerInfo field.
func (r *queryResolver) WorkerInfo(ctx context.Context) (*model.WorkerInfo, error) {
	root, err := store.GetRootUser()
	if err != nil {
		root = ""
	}
	var maddress = ""
	minter, _, err := mint.GetMintKey()
	if err == nil {
		maddress = minter.Address
	}

	report, _, err := proof.GetRemoteReport(minter, nil)
	if err != nil {
		report = nil
	}

	return &model.WorkerInfo{
		RootAddress: root,
		MintAddress: maddress,
		Report:      hex.EncodeToString(report),
	}, nil
}

// Worker is the resolver for the worker field.
func (r *queryResolver) Worker(ctx context.Context) ([]*model.Contract, error) {
	client := mint.MinterIns.ChainClient
	if client == nil {
		return nil, gqlerror.Errorf("Cant connect to chain")
	}
	worker := &module.Worker{
		Client: client,
		Signer: mint.MinterIns.Signer,
	}

	clusterID, err := store.GetClusterId()
	if err != nil {
		return nil, gqlerror.Errorf("Cant get cluster id:" + err.Error())
	}

	contracts, err := worker.GetClusterContracts(clusterID, nil)
	if err != nil {
		return nil, gqlerror.Errorf("GetClusterContracts:" + err.Error())
	}

	list := make([]*model.Contract, 0, len(contracts))
	for _, contract := range contracts {
		list = append(list, &model.Contract{
			StartNumber: fmt.Sprint(contract.ContractState.StartNumber),
			User:        hex.EncodeToString(contract.ContractState.User[:]),
			WorkID:      util.GetWorkTypeStr(contract.ContractState.WorkId) + "-" + fmt.Sprint(contract.ContractState.WorkId.Id),
		})
	}
	return list, nil
}
