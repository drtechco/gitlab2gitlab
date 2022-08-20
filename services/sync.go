package services

import (
	"context"
	"drtech.co/gl2gl/orm"
	"drtech.co/gl2gl/orm/model"
	"sync"
	"time"
)

type SyncPipeStatus int

const (
	SyncPipeStatusNotInit             SyncPipeStatus = 101
	SyncPipeStatusIniting             SyncPipeStatus = 102
	SyncPipeStatusFromClientInitError SyncPipeStatus = 103
	SyncPipeStatusToClientInitError   SyncPipeStatus = 104
	SyncPipeStatusInitOk              SyncPipeStatus = 105

	SyncPipeStatusGetFromGroups           SyncPipeStatus = 206
	SyncPipeStatusGetFromGroupsError      SyncPipeStatus = 207
	SyncPipeStatusGetToGroups             SyncPipeStatus = 208
	SyncPipeStatusGetToGroupsError        SyncPipeStatus = 209
	SyncPipeStatusCreateToGroup           SyncPipeStatus = 210
	SyncPipeStatusCreateToGroupErr        SyncPipeStatus = 210
	SyncPipeStatusGetFromGroupProjects    SyncPipeStatus = 211
	SyncPipeStatusGetFromGroupProjectsErr SyncPipeStatus = 212
	SyncPipeStatusCreateToProject         SyncPipeStatus = 213
	SyncPipeStatusCreateToProjectErr      SyncPipeStatus = 214
	SyncPipeStatusGetFromBranches         SyncPipeStatus = 215
	SyncPipeStatusGetFromBranchesError    SyncPipeStatus = 216
	SyncPipeStatusGetToBranches           SyncPipeStatus = 217
	SyncPipeStatusGetToBranchesError      SyncPipeStatus = 218

	SyncPipeStatusGetToGroupProjects    SyncPipeStatus = 211
	SyncPipeStatusGetToGroupProjectsErr SyncPipeStatus = 212

	SyncPipeStatusGetFromProjects      SyncPipeStatus = 206
	SyncPipeStatusGetFromProjectsError SyncPipeStatus = 207

	SyncPipeStatusGetToProjectsError SyncPipeStatus = 207
	SyncPipeStatusDiffChecking       SyncPipeStatus = 206
)

var syncPipeMapLock sync.Mutex
var syncPipeMap map[int32]*SyncPipe

func Setup() error {

	makeSyncPipe()
	go run()
	return nil
}

func run() {
	for {
		err := loadConfig()
		if err != nil {
			//TODO log
			time.Sleep(10 * time.Second)
			break
		}
		syncPipeMapLock.Lock()
		for _, syncPipe := range syncPipeMap {
			err := syncPipe.Run()
			if err != nil {
				//TODO log
				time.Sleep(10 * time.Second)
			}
		}
	}

}

func loadConfig() error {
	fromToConfigM := orm.DbQuery().FromToConfig
	configList, err := fromToConfigM.WithContext(context.Background()).Find()
	if err != nil {
		return err
	}
	for _, config := range configList {
		err := initSyncPipe(config)
		if err != nil {
			//TODO log
		}
	}
	return nil
}

func makeSyncPipe() {
	syncPipeMapLock.Lock()
	defer syncPipeMapLock.Unlock()
	if syncPipeMap != nil {
		for _, syncPipe := range syncPipeMap {
			err := syncPipe.Stop()
			if err != nil {
				//TODO log
			}
		}
	}
	syncPipeMap = make(map[int32]*SyncPipe)
}

func initSyncPipe(config *model.FromToConfig) error {
	pipe := &SyncPipe{
		ConfigId:        config.ID,
		FromAddress:     config.FromAddress,
		FromAccessToken: config.FromAccessToken,
		Status:          SyncPipeStatusNotInit,
		ToAddress:       config.ToAddress,
		ToAccessToken:   config.ToAccessToken,
	}
	syncPipeMap[config.ID] = pipe
	return nil
}
