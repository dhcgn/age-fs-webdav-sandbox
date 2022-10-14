package fileinfo

import (
	"io/fs"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func NewFileInfo(logger *logrus.Entry, fsinternal os.FileInfo) fs.FileInfo {
	return fileinfo{
		Logger:     logger,
		fsinternal: fsinternal,
	}
}

type fileinfo struct {
	Logger     *logrus.Entry
	fsinternal os.FileInfo
}

// IsDir implements fs.FileInfo
func (fi fileinfo) IsDir() bool {
	fi.Logger.WithField("result", fi.fsinternal.IsDir()).Debugln("IsDir")
	return fi.fsinternal.IsDir()
}

// ModTime implements fs.FileInfo
func (fi fileinfo) ModTime() time.Time {
	fi.Logger.WithField("result", fi.fsinternal.ModTime()).Debugln("ModTime")
	return fi.fsinternal.ModTime()
}

// Mode implements fs.FileInfo
func (fi fileinfo) Mode() fs.FileMode {
	fi.Logger.WithField("result", fi.fsinternal.Mode()).Debugln("Mode")
	return fi.fsinternal.Mode()
}

// Name implements fs.FileInfo
func (fi fileinfo) Name() string {
	fi.Logger.WithField("result", fi.fsinternal.Name()).Debugln("Name")
	return fi.fsinternal.Name()
}

// Size implements fs.FileInfo
func (fi fileinfo) Size() int64 {
	fi.Logger.WithField("result", fi.fsinternal.Size()).Debugln("Size")
	return fi.fsinternal.Size()
}

// Sys implements fs.FileInfo
func (fi fileinfo) Sys() any {
	fi.Logger.Debugln("Sys")
	return nil
}
