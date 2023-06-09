package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/stolostron/search-v2-api/graph/generated"
	"github.com/stolostron/search-v2-api/graph/model"
	"github.com/stolostron/search-v2-api/pkg/resolver"
	klog "k8s.io/klog/v2"
)

// Search is the resolver for the search field.
func (r *queryResolver) Search(ctx context.Context, input []*model.SearchInput) ([]*resolver.SearchResult, error) {
	klog.V(3).Infof("--------- Received Search query with %d inputs ---------\n", len(input))
	return resolver.Search(ctx, input)
}

// SearchComplete is the resolver for the searchComplete field.
func (r *queryResolver) SearchComplete(ctx context.Context, property string, query *model.SearchInput, limit *int) ([]*string, error) {
	if limit != nil {
		klog.V(3).Infof("Received SearchComplete query with input property **%s** and limit %d", property, *limit)
	} else {
		klog.V(3).Infof("Received SearchComplete query with input property **%s**", property)
	}
	return resolver.SearchComplete(ctx, property, query, limit)
}

// SearchSchema is the resolver for the searchSchema field.
func (r *queryResolver) SearchSchema(ctx context.Context) (map[string]interface{}, error) {
	klog.V(3).Infoln("Received SearchSchema query")
	return resolver.SearchSchemaResolver(ctx)
}

// Messages is the resolver for the messages field.
func (r *queryResolver) Messages(ctx context.Context) ([]*model.Message, error) {
	klog.V(3).Infoln("Received Messages query")
	return resolver.Messages(ctx)
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
