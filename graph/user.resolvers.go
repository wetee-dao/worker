package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.42

import (
	"context"
	"encoding/json"

	subkey "github.com/vedhavyas/go-subkey/v2"
	"github.com/vedhavyas/go-subkey/v2/sr25519"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"wetee.app/worker/db"
	"wetee.app/worker/graph/model"
)

// LoginAndBindRoot is the resolver for the loginAndBindRoot field.
func (r *mutationResolver) LoginAndBindRoot(ctx context.Context, input model.LoginContent, signature string) (string, error) {
	account := &model.User{
		Address:   input.Address,
		Timestamp: input.Timestamp,
	}
	bt, _ := json.Marshal(account)
	str := subkey.EncodeHex(bt)

	rootUser, _ := db.GetRootUser()
	if rootUser != "" && rootUser != input.Address {
		return "", gqlerror.Errorf("Root user already exists")
	}

	// 解析地址
	_, pubkeyBytes, err := subkey.SS58Decode(input.Address)
	if err != nil {
		return "", gqlerror.Errorf("Bad address")
	}

	// 解析公钥
	pubkey, err := sr25519.Scheme{}.FromPublicKey(pubkeyBytes)
	if err != nil {
		return "", gqlerror.Errorf("Bad sr25519 address")
	}

	// 解析签名
	sig, chainerr := subkey.DecodeHex(signature)
	if !chainerr {
		return "", gqlerror.Errorf("Bad signature hex")
	}

	// 验证签名
	ok := pubkey.Verify(bt, sig)
	if !ok {
		// return "", gqlerror.Errorf("Bad signature")
	}

	// 设置根用户
	err = db.SetRootUser(input.Address)
	if err != nil {
		return "", gqlerror.Errorf("Set root user error: " + err.Error())
	}

	return str + "||" + signature, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
