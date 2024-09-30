package app

import (
	"github.com/tensved/snet-matrix-framework/internal/config"
	"github.com/tensved/snet-matrix-framework/internal/grpcmanager"
	"github.com/tensved/snet-matrix-framework/internal/logger"
	"github.com/tensved/snet-matrix-framework/internal/matrix"
	"github.com/tensved/snet-matrix-framework/internal/server"
	"github.com/tensved/snet-matrix-framework/internal/syncer"
	"github.com/tensved/snet-matrix-framework/pkg/blockchain"
	"github.com/tensved/snet-matrix-framework/pkg/db"
	ipfs "github.com/tensved/snet-matrix-framework/pkg/ipfs"
	"regexp"
)

type App struct {
	DB           db.Service
	Fiber        *server.FiberServer
	Ethereum     blockchain.Ethereum
	MatrixClient matrix.Service
	IPFSClient   ipfs.IPFSClient
	Syncer       syncer.SnetSyncer
	GRPCManager  *grpcmanager.GRPCClientManager
}

func New() App {
	logger.Setup()
	config.Init()
	config.IPFS.HashCutterRegexp = regexp.MustCompile("[^a-zA-Z0-9=]")
	database := db.New()
	eth := blockchain.Init()
	ipfsClient := ipfs.Init()
	snetSyncer := syncer.New(eth, ipfsClient, database)
	grpcManager := grpcmanager.NewGRPCClientManager()
	app := App{DB: database, Fiber: server.New(database), MatrixClient: matrix.New(database, snetSyncer, grpcManager, eth), IPFSClient: ipfsClient, Ethereum: eth, Syncer: snetSyncer, GRPCManager: grpcManager}

	app.Syncer.DB = app.DB
	app.Syncer.Ethereum = app.Ethereum
	app.Syncer.IPFSClient = app.IPFSClient

	return app
}
