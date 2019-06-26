package mongo

import (
	"context"
	"io"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	docker "github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/mee6aas/zeep/api"
)

// DeployOptions holds options to deploy the mongodb
type DeployOptions struct {
	port    string
	Port    uint16 // Port number of the mongodb to deploy
	Pull    bool   // Pull the mongodb image if specified version is not exists
	Version string // Version of the mongodb to deploy
	Network string // Networt which container is connected
}

// DeployOption is a option to deploy a mongodb
type DeployOption func(*DeployOptions)

// DeployDeployedBody OK response to ContainerCreate operation
type DeployDeployedBody struct {
	ContainerID string
}

// Deploy create a mongodb container
func Deploy(ctx context.Context, setters ...DeployOption) (res DeployDeployedBody, err error) {
	var (
		img       string
		dbURI     string
		createRes container.ContainerCreateCreatedBody
		dbClient  *mongo.Client
	)

	args := &DeployOptions{
		Port:    api.DefaultMongoPort,
		Pull:    false,
		Version: api.DefaultMongoVersion,
	}

	for _, setter := range setters {
		setter(args)
	}

	img = "mongo:" + args.Version
	dbURI = "mongodb://localhost:" + strconv.Itoa(int(args.Port))

	create := func() (e error) {
		createRes, e = engineClient.ContainerCreate(ctx, &container.Config{
			Image: img,
			ExposedPorts: nat.PortSet{
				nat.Port(strconv.Itoa(int(args.Port)) + "/tcp"): struct{}{},
			},
			Cmd: []string{
				"--port",
				strconv.Itoa(int(args.Port)),
			},
		}, &container.HostConfig{
			NetworkMode: container.NetworkMode(args.Network),
		}, nil, "")
		return e
	}
	checkAlive := func() (e error) {
		c, cc := context.WithTimeout(context.Background(), time.Second*3)
		defer cc()

		// TODO: Is it right implementation?
		go func() {
			select {
			case <-ctx.Done():
			case <-c.Done():
			}

			cc()
		}()

		if dbClient, e = mongo.Connect(c, options.Client().ApplyURI(dbURI)); err != nil {
			e = errors.Wrapf(e, "Failed to create mongodb client %s", dbURI)
			return
		}

		if e = dbClient.Ping(c, readpref.Primary()); err != nil {
			e = errors.Wrapf(e, "Failed to find mongodb server at %s", dbURI)
			return
		}

		return
	}

	defer func() {
		if err == nil {
			return
		}

		if len(createRes.ID) == 0 {
			return
		}

		engineClient.ContainerRemove(context.Background(), createRes.ID, types.ContainerRemoveOptions{
			Force: true,
		})
	}()

	if err = checkAlive(); err == nil {
		// Already exists
		// TODO: return error with exists information
		err = errors.Errorf("Mongo server already exists at %s", dbURI)
		return
	}

	if err = create(); err == nil {
		// no error
	} else if !docker.IsErrNotFound(err) {
		// failed to create with unknown error
		err = errors.Wrapf(err, "Failed to create container for %s", img)
		return
	} else if !args.Pull {
		// image not found on local and no pull allowed
		err = errors.Wrapf(err, "Image %s not exists on local", img)
		return
	} else {
		// image not found but pull allowed
		if out, e := engineClient.ImagePull(ctx, img, types.ImagePullOptions{}); e == nil {
			io.Copy(ioutil.Discard, out)
			out.Close()
		} else {
			// image not found on hub
			err = errors.Wrapf(e, "Failed to pull image %s", img)
			return
		}

		// image found on hub
		// try creating container again
		if err = create(); err != nil {
			err = errors.Wrapf(err, "Failed to create container for %s", img)
			return
		}
	}

	if err = engineClient.ContainerStart(ctx, createRes.ID, types.ContainerStartOptions{}); err != nil {
		err = errors.Wrapf(err, "Failed to start container %s", createRes.ID)
		return
	}

	if err = checkAlive(); err != nil {
		err = errors.Wrap(err, "Failed to check alive")
		return
	}

	res.ContainerID = createRes.ID

	return
}

// WithVersion overrides the version of the mongodb to deploy.
func WithVersion(version string) DeployOption {
	return func(args *DeployOptions) {
		args.Version = version
	}
}

// WithPort overrides the port number of the mongodb to expose.
func WithPort(port uint16) DeployOption {
	return func(args *DeployOptions) {
		args.Port = port
	}
}

// WithPull overrides the option to pull mongodb image.
func WithPull() DeployOption {
	return func(args *DeployOptions) {
		args.Pull = true
	}
}

// WithoutPull overrides the option not to pull mongodb image.
func WithoutPull() DeployOption {
	return func(args *DeployOptions) {
		args.Pull = false
	}
}

// WithNetwork overrides the network name which container is connected.
func WithNetwork(network string) DeployOption {
	return func(args *DeployOptions) {
		args.Network = network
	}
}
