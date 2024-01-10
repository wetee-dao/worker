package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.42

import (
	"context"
	"fmt"
)

// ClusterRegister is the resolver for the cluster_register field.
func (r *mutationResolver) ClusterRegister(ctx context.Context, input string) (string, error) {
	panic(fmt.Errorf("not implemented: ClusterRegister - cluster_register"))
}

// Worker is the resolver for the worker field.
func (r *queryResolver) Worker(ctx context.Context) ([]string, error) {
	return []string{"127.0.0.1:8080"}, nil
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }