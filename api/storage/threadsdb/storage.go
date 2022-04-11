package threadsdb

import (
	"context"

	"github.com/alecthomas/jsonschema"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"github.com/textileio/go-threads/api/client"
	"github.com/textileio/go-threads/core/thread"
	"github.com/textileio/go-threads/db"
	"google.golang.org/grpc"
)

type threadsdbConnectionParams struct {
	Host string `envconfig:"THREADS_DB_HOST"`
}

type threadsDbStorage struct {
	*client.Client
	*thread.ID
}

func NewThreadStorage() (*threadsDbStorage, error) {
	var params threadsdbConnectionParams
	envconfig.Process("", &params)

	logrus.Infof("Connecting to Threads database: %s", params.Host)

	cli, err := client.NewClient(params.Host, grpc.WithInsecure())

	if err != nil {
		return nil, err
	}
	threadId := thread.NewIDV1(thread.Raw, 32)
	_, err = cli.GetDBInfo(context.Background(), threadId)
	if err != nil {
		err = cli.NewDB(context.Background(), threadId)

		if err != nil {
			return nil, err
		}

		reflector := jsonschema.Reflector{}
		mySchema := reflector.Reflect(&UserSchema{})
		err = cli.NewCollection(context.Background(), threadId, db.CollectionConfig{
			Name:   "Users",
			Schema: mySchema,
			Indexes: []db.Index{
				{
					Path:   "did",
					Unique: true,
				},
				{
					Path:   "username",
					Unique: true,
				},
				{
					Path:   "name",
					Unique: false,
				}},
		})
		if err != nil {
			return nil, err
		}
		mySchema = reflector.Reflect(&CommunitySchema{})
		err = cli.NewCollection(context.Background(), threadId, db.CollectionConfig{
			Name:   "Communities",
			Schema: mySchema,
			Indexes: []db.Index{
				{
					Path:   "zid",
					Unique: true,
				},
				{
					Path:   "ownerDid",
					Unique: false,
				},
				{
					Path:   "ownerUsername",
					Unique: false,
				}},
		})
		if err != nil {
			return nil, err
		}
	}

	storage := threadsDbStorage{
		Client: cli,
		ID:     &threadId,
	}

	return &storage, nil
}
