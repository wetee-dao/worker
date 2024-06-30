package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.42

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/vektah/gqlparser/v2/gqlerror"
	chain "github.com/wetee-dao/go-sdk"
	gtypes "github.com/wetee-dao/go-sdk/gen/types"
	"wetee.app/worker/graph/model"
	"wetee.app/worker/mint"
	"wetee.app/worker/mint/proof"
	"wetee.app/worker/store"
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
	worker := &chain.Worker{
		Client: client,
		Signer: mint.Signer,
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
	worker := &chain.Worker{
		Client: client,
		Signer: mint.Signer,
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
	worker := &chain.Worker{
		Client: client,
		Signer: mint.Signer,
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
	worker := &chain.Worker{
		Client: client,
		Signer: mint.Signer,
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
	worker := &chain.Worker{
		Client: client,
		Signer: mint.Signer,
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

// WorkerInfo is the resolver for the workerInfo field.
func (r *queryResolver) WorkerInfo(ctx context.Context) (*model.WorkerInfo, error) {
	root, err := store.GetRootUser()
	if err != nil {
		root = ""
	}
	var maddress = ""
	minter, err := mint.GetMintKey()
	if err == nil {
		maddress = minter.Address
	}

	_, _, report, err := proof.GetRemoteReport("")
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
	worker := &chain.Worker{
		Client: client,
		Signer: mint.Signer,
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
