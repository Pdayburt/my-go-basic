package mongo

import (
	"testing"
)

func TestMongo(t *testing.T) {
	/*	//初始化操作时间
		ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelFunc()
		monitor := &event.CommandMonitor{
			Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
				fmt.Println(startedEvent.Command)
			},
			Succeeded: func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {

			},
			Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {

			},
		}
		opts := options.Client().ApplyURI("mongodb://localhost:27017").
			SetMonitor(monitor)
		client, err := mongo.Connect(ctx, opts)
		assert.NoError(t, err)

		mdb := client.Database("webook")
		col := mdb.Collection("articles")

		//col.InsertOne()*/
}
