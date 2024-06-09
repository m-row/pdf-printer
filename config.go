package main

import (
	"os"
	"runtime/debug"
	"time"
)

type CommitInfo struct {
	FullSHA1 string
	Time     string
}

type Settings struct {
	Timezone string
	Port     string
	Domain   string
	Env      string
	AppVer   string
	AppCode  string
	AppDesc  string
	AppName  string

	CommitInfo CommitInfo
}

var (
	DOMAIN  = os.Getenv("DOMAIN")
	ENV     = os.Getenv("ENV")
	Port    = os.Getenv("PORT")
	AppVer  = os.Getenv("APP_VER")
	AppCode = os.Getenv("APP_CODE")
	AppDesc = os.Getenv("APP_DESC")
	AppName = os.Getenv("APP_NAME")
)

func getCommitInfo() CommitInfo {
	var commitInfo CommitInfo

	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				commitInfo.FullSHA1 = setting.Value
			}
			if setting.Key == "vcs.time" {
				commitInfo.Time = setting.Value
			}
		}
		return commitInfo
	}
	return commitInfo
}

func TimeNow() time.Time {
	return time.Now().UTC()
}

func GetSettings() *Settings {
	return &Settings{
		Env:     ENV,
		Port:    Port,
		Domain:  DOMAIN,
		AppVer:  AppVer,
		AppCode: AppCode,
		AppDesc: AppDesc,
		AppName: AppName,

		CommitInfo: getCommitInfo(),
	}
}
