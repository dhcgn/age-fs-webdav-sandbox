package mywebdav

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/net/webdav"
)

// A FS implements FileSystem using the native file system restricted to a
// specific directory tree.
//
// While the FileSystem.OpenFile method takes '/'-separated paths, a FS's
// string value is a filename on the native file system, not a URL, so it is
// separated by filepath.Separator, which isn't necessarily '/'.
//
// An empty FS is treated as ".".
type FS struct {
	Dir string
}

func NewFileSystem(dir string) FS {
	return FS{
		Dir: dir,
	}
}

func (f FS) resolve(name string) string {
	// This implementation is based on Dir.Open's code in the standard net/http package.
	if filepath.Separator != '/' && strings.IndexRune(name, filepath.Separator) >= 0 ||
		strings.Contains(name, "\x00") {
		return ""
	}
	dir := f.Dir
	if dir == "" {
		dir = "."
	}
	return filepath.Join(dir, filepath.FromSlash(slashClean(name)))
}

func (f FS) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	if name = f.resolve(name); name == "" {
		return os.ErrNotExist
	}
	return os.Mkdir(name, perm)
}

func (f FS) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	if name = f.resolve(name); name == "" {
		return nil, os.ErrNotExist
	}
	file, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (f FS) RemoveAll(ctx context.Context, name string) error {
	if name = f.resolve(name); name == "" {
		return os.ErrNotExist
	}
	if name == filepath.Clean(f.Dir) {
		// Prohibit removing the virtual root directory.
		return os.ErrInvalid
	}
	return os.RemoveAll(name)
}

func (f FS) Rename(ctx context.Context, oldName, newName string) error {
	if oldName = f.resolve(oldName); oldName == "" {
		return os.ErrNotExist
	}
	if newName = f.resolve(newName); newName == "" {
		return os.ErrNotExist
	}
	if root := filepath.Clean(f.Dir); root == oldName || root == newName {
		// Prohibit renaming from or to the virtual root directory.
		return os.ErrInvalid
	}
	return os.Rename(oldName, newName)
}

func (f FS) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	if name = f.resolve(name); name == "" {
		return nil, os.ErrNotExist
	}
	return os.Stat(name)
}

// slashClean is equivalent to but slightly more efficient than
// path.Clean("/" + name).
func slashClean(name string) string {
	if name == "" || name[0] != '/' {
		name = "/" + name
	}
	return path.Clean(name)
}
