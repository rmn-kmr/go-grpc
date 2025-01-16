package callback

import (
	"context"

	"github.com/rmnkmr/lsp/internal/app"
	"github.com/rmnkmr/protonium"
)

type Callback struct {
	App *app.App
}

func ListenerServer(ctx context.Context, c *Callback) []protonium.Option {
	//client, err := pubsub.NewClient(ctx, c.App.Config.Gcp.ProjectID)
	//if err != nil {
	//	panic(err)
	//}
	//var options []protonium.Option
	//options = []protonium.Option{
	//	func() protonium.Option {
	//		l := &pubsub.Listener{
	//			Environment:    c.App.Environment,
	//			Client:         client,
	//			Message:        nil,
	//			SubscriptionID: fmt.Sprintf("%s__lsp__listener__rmnkmr.lsp.v1.Callback", c.App.Environment),
	//		}
	//		return pubsub.ListenerOption(l, c.OnHTTPCallback, func(data []byte) (*HTTPCallback, error) {
	//			var msg HTTPCallback
	//			if err := json.Unmarshal(data, &msg); err != nil {
	//				return nil, err
	//			}
	//			return &msg, nil
	//		})
	//	}()}
	//
	//// filter options for nil values
	//var targetOptions []protonium.Option
	//for _, o := range options {
	//	if o != nil {
	//		targetOptions = append(targetOptions, o)
	//	}
	//}
	//return targetOptions
	return nil
}
