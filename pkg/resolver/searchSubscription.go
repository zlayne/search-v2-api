// Copyright Contributors to the Open Cluster Management project
package resolver

import (
	"context"
	"time"

	klog "k8s.io/klog/v2"

	"github.com/stolostron/search-v2-api/graph/model"
	"github.com/stolostron/search-v2-api/pkg/config"
)

func SearchSubscription(ctx context.Context, input []*model.SearchInput) (<-chan []*SearchResult, error) {
	ch := make(chan []*SearchResult)

	// You can (and probably should) handle your channels in a central place outside of `schema.resolvers.go`.
	// For this example we'll simply use a Goroutine with a simple loop.
	go func() {
		// Handle deregistration of the channel here. Note the `defer`
		defer close(ch)

		for {
			searchResult, err := Search(ctx, input)

			if err != nil {
				klog.Errorf("Error occurred durign the search subscription request: %s", err)
			}

			// The subscription may have got closed due to the client disconnecting.
			// Hence we do send in a select block with a check for context cancellation.
			// This avoids goroutine getting blocked forever or panicking,
			select {
			case <-ctx.Done(): // This runs when context gets cancelled. Subscription closes.
				klog.Infof("Search subscription Closed")
				// Handle deregistration of the channel here. `close(ch)`
				return // Remember to return to end the routine.
			
			case ch <- searchResult: // This is the actual send.
				// Our message went through, do nothing	
			}

			// Wait 10 seconds for next search reuslt send.
			time.Sleep(time.Duration(config.Cfg.SubscriptionPollInterval) * time.Second)
			klog.Infof("Search subscription new poll interval")
		}
	}()

	// We return the channel and no error.
	return ch, nil
}