package appserver

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/casualjim/go-appserver/middleware"
)

var (
	Version   string
	BuildDate string
	GitCommit string
	GitState  string
)

func NewVersionInfo() VersionInfo {
	ver := VersionInfo{
		Version:   "dev",
		BuildDate: BuildDate,
		GitCommit: GitCommit,
		GitState:  "",
	}
	if Version != "" {
		ver.Version = Version
		ver.GitState = "clean"
	}
	if GitState != "" {
		ver.GitState = GitState
	}
	return ver
}

type VersionInfo struct {
	Version   string `json:"version,omitempty"`
	BuildDate string `json:"buildDate,omitempty"`
	GitCommit string `json:"gitCommit,omitempty"`
	GitState  string `json:"gitState,omitempty"`
}

func (v VersionInfo) String() string {
	var buf bytes.Buffer
	buf.WriteString("Version: ")
	buf.WriteString(v.Version)
	buf.WriteString("\n")
	buf.WriteString("Build date: ")
	buf.WriteString(v.BuildDate)
	buf.WriteString("\n")
	buf.WriteString("Commit: ")
	buf.WriteString(v.GitCommit)
	buf.WriteString("\n")
	buf.WriteString("Working tree: ")
	buf.WriteString(v.GitState)
	buf.WriteString("\n")
	return buf.String()
}

func VersionHandler(log middleware.Logger, info VersionInfo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		enc := json.NewEncoder(w)
		if err := enc.Encode(info); err != nil {
			log.Printf("failed to write version response: %v", err)
		}
	}
}
