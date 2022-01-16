// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-drive/common"
	"go-drive/common/event"
	"go-drive/common/i18n"
	"go-drive/common/registry"
	"go-drive/common/task"
	"go-drive/common/utils"
	"go-drive/drive"
	"go-drive/server"
	"go-drive/server/search"
	"go-drive/server/thumbnail"
	"go-drive/storage"
)

// Injectors from wire.go:

func Initialize(ctx context.Context, ch *registry.ComponentsHolder) (*gin.Engine, error) {
	config, err := common.InitConfig(ch)
	if err != nil {
		return nil, err
	}
	bus := event.NewBus(ch)
	db, err := storage.NewDB(config, ch)
	if err != nil {
		return nil, err
	}
	driveDAO := storage.NewDriveDAO(db)
	pathMountDAO := storage.NewPathMountDAO(db)
	driveDataDAO := storage.NewDriveDataDAO(db)
	driveCacheDAO := storage.NewDriveCacheDAO(db, ch)
	rootDrive, err := drive.NewRootDrive(ctx, config, driveDAO, pathMountDAO, driveDataDAO, driveCacheDAO)
	if err != nil {
		return nil, err
	}
	pathPermissionDAO := storage.NewPathPermissionDAO(db)
	signer := utils.NewSigner()
	access, err := drive.NewAccess(ch, rootDrive, pathPermissionDAO, signer, bus)
	if err != nil {
		return nil, err
	}
	optionsDAO := storage.NewOptionsDAO(db)
	tunnyRunner := task.NewTunnyRunner(config, ch)
	service, err := search.NewService(ch, config, optionsDAO, rootDrive, tunnyRunner, bus)
	if err != nil {
		return nil, err
	}
	fileTokenStore, err := server.NewFileTokenStore(config, ch)
	if err != nil {
		return nil, err
	}
	maker, err := thumbnail.NewMaker(config, optionsDAO, ch)
	if err != nil {
		return nil, err
	}
	chunkUploader, err := server.NewChunkUploader(config)
	if err != nil {
		return nil, err
	}
	userDAO := storage.NewUserDAO(db)
	groupDAO := storage.NewGroupDAO(db)
	fileMessageSource, err := i18n.NewFileMessageSource(config)
	if err != nil {
		return nil, err
	}
	engine, err := server.InitServer(config, ch, bus, rootDrive, access, service, fileTokenStore, maker, signer, chunkUploader, tunnyRunner, optionsDAO, userDAO, groupDAO, driveDAO, driveCacheDAO, driveDataDAO, pathPermissionDAO, pathMountDAO, fileMessageSource)
	if err != nil {
		return nil, err
	}
	return engine, nil
}
