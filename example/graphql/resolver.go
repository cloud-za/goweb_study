package graphql

import (
	"context"
)

type Resolver struct{}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) AddChannel(ctx context.Context, name string) (Channel, error) {
	panic("not implemented")
}
func (r *mutationResolver) UpdateChannel(ctx context.Context, id int, name *string) (Channel, error) {
	panic("not implemented")
}
func (r *mutationResolver) DeleteChannel(ctx context.Context, id int) (Channel, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Channels(ctx context.Context) ([]Channel, error) {
	panic("not implemented")
}

type subscriptionResolver struct{ *Resolver }

func (r *subscriptionResolver) SubscriptionChannelAdded(ctx context.Context) (<-chan Channel, error) {
	panic("not implemented")
}
func (r *subscriptionResolver) SubscriptionChannelDeleted(ctx context.Context) (<-chan Channel, error) {
	panic("not implemented")
}
func (r *subscriptionResolver) SubscriptionChannelUpdated(ctx context.Context) (<-chan Channel, error) {
	panic("not implemented")
}
