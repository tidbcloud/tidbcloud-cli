//go:build !windows

// Copyright 2026 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fuse

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/tidbcloud/tidbcloud-cli/internal/service/fs"

	gofuse "github.com/hanwen/go-fuse/v2/fuse"
)

// ReadOnlyFS is a read-only FUSE filesystem for TiDB Cloud FS.
type ReadOnlyFS struct {
	gofuse.RawFileSystem

	client *fs.Client
	opts   *MountOptions
	inodes *inodeManager
	dirs   *dirCache
	uid    uint32
	gid    uint32
}

// inodeManager manages inode to path mappings.
type inodeManager struct {
	mu      sync.RWMutex
	byInode map[uint64]*inodeEntry
	byPath  map[string]uint64
	nextIno uint64
}

type inodeEntry struct {
	Path     string
	IsDir    bool
	RefCnt   uint64
	Attr     *fs.StatResult
	AttrTime time.Time
}

// dirCache caches directory listings.
type dirCache struct {
	mu      sync.RWMutex
	entries map[string]*dirCacheEntry
	ttl     time.Duration
}

type dirCacheEntry struct {
	Items     []fs.FileInfo
	Timestamp time.Time
}

// NewReadOnlyFS creates a new read-only FUSE filesystem.
func NewReadOnlyFS(client *fs.Client, opts *MountOptions) *ReadOnlyFS {
	fsys := &ReadOnlyFS{
		RawFileSystem: gofuse.NewDefaultRawFileSystem(),
		client:        client,
		opts:          opts,
		inodes: &inodeManager{
			byInode: make(map[uint64]*inodeEntry),
			byPath:  make(map[string]uint64),
			nextIno: gofuse.FUSE_ROOT_ID,
		},
		dirs: &dirCache{
			entries: make(map[string]*dirCacheEntry),
			ttl:     5 * time.Second,
		},
		uid: uint32(os.Getuid()),
		gid: uint32(os.Getgid()),
	}

	// Initialize root inode
	fsys.inodes.byPath["/"] = gofuse.FUSE_ROOT_ID
	fsys.inodes.byInode[gofuse.FUSE_ROOT_ID] = &inodeEntry{
		Path:   "/",
		IsDir:  true,
		RefCnt: 1,
	}

	return fsys
}

// GetAttr implements GetAttr.
func (f *ReadOnlyFS) GetAttr(cancel <-chan struct{}, in *gofuse.GetAttrIn, out *gofuse.AttrOut) gofuse.Status {
	f.inodes.mu.RLock()
	entry, ok := f.inodes.byInode[in.NodeId]
	f.inodes.mu.RUnlock()

	if !ok {
		return gofuse.ENOENT
	}

	// Check cache
	if entry.Attr != nil && time.Since(entry.AttrTime) < time.Second {
		f.fillAttr(entry, out)
		return gofuse.OK
	}

	// Special-case root directory: backend may not support Stat("/")
	if entry.Path == "/" {
		entry.IsDir = true
		entry.Attr = &fs.StatResult{IsDir: true, Size: 0}
		entry.AttrTime = time.Now()
		f.fillAttr(entry, out)
		return gofuse.OK
	}

	// Fetch from server
	stat, err := f.client.Stat(entry.Path)
	if err != nil {
		return gofuse.ENOENT
	}

	entry.Attr = stat
	entry.AttrTime = time.Now()
	entry.IsDir = stat.IsDir

	f.fillAttr(entry, out)
	return gofuse.OK
}

func (f *ReadOnlyFS) fillAttr(entry *inodeEntry, out *gofuse.AttrOut) {
	out.Attr.Ino = f.inodes.getIno(entry.Path)
	if entry.Attr != nil {
		out.Attr.Size = uint64(entry.Attr.Size)
	} else {
		out.Attr.Size = 0
	}
	out.Attr.Mode = 0o755
	if entry.IsDir {
		out.Attr.Mode |= syscall.S_IFDIR
		out.Attr.Nlink = 2
	} else {
		out.Attr.Mode |= syscall.S_IFREG
		out.Attr.Nlink = 1
	}
	out.Attr.Uid = f.uid
	out.Attr.Gid = f.gid
	out.SetTimeout(1 * time.Second)
}

// Lookup implements Lookup.
func (f *ReadOnlyFS) Lookup(cancel <-chan struct{}, header *gofuse.InHeader, name string, out *gofuse.EntryOut) gofuse.Status {
	parentIno := header.NodeId

	f.inodes.mu.RLock()
	parentEntry, ok := f.inodes.byInode[parentIno]
	f.inodes.mu.RUnlock()

	if !ok {
		return gofuse.ENOENT
	}

	path := filepath.Join(parentEntry.Path, name)

	// Get or create inode
	ino := f.inodes.GetOrCreate(path)

	f.inodes.mu.RLock()
	entry := f.inodes.byInode[ino]
	f.inodes.mu.RUnlock()

	// Fetch attr
	stat, err := f.client.Stat(path)
	if err != nil {
		return gofuse.ENOENT
	}

	entry.Attr = stat
	entry.AttrTime = time.Now()
	entry.IsDir = stat.IsDir

	out.NodeId = ino
	out.Generation = 1
	out.Attr.Ino = ino
	if entry.Attr != nil {
		out.Attr.Size = uint64(entry.Attr.Size)
	} else {
		out.Attr.Size = 0
	}
	out.Attr.Mode = 0o755
	if entry.IsDir {
		out.Attr.Mode |= syscall.S_IFDIR
	} else {
		out.Attr.Mode |= syscall.S_IFREG
	}
	out.Attr.Nlink = 1
	out.Attr.Uid = f.uid
	out.Attr.Gid = f.gid

	return gofuse.OK
}

// Access implements Access.
func (f *ReadOnlyFS) Access(cancel <-chan struct{}, in *gofuse.AccessIn) gofuse.Status {
	return gofuse.OK
}

// OpenDir implements OpenDir.
func (f *ReadOnlyFS) OpenDir(cancel <-chan struct{}, in *gofuse.OpenIn, out *gofuse.OpenOut) gofuse.Status {
	out.Fh = in.NodeId
	return gofuse.OK
}

// ReadDir implements ReadDir.
func (f *ReadOnlyFS) ReadDir(cancel <-chan struct{}, in *gofuse.ReadIn, out *gofuse.DirEntryList) gofuse.Status {
	f.inodes.mu.RLock()
	entry, ok := f.inodes.byInode[in.NodeId]
	f.inodes.mu.RUnlock()

	if !ok {
		return gofuse.ENOENT
	}

	realEntries, err := f.getDirEntries(entry.Path)
	if err != nil {
		return gofuse.EIO
	}

	dotdotIno := in.NodeId
	if entry.Path != "/" {
		parentPath := filepath.Dir(entry.Path)
		dotdotIno = f.inodes.getIno(parentPath)
		if dotdotIno == 0 {
			dotdotIno = f.inodes.GetOrCreate(parentPath)
		}
	}

	type dirItem struct {
		name string
		ino  uint64
		mode uint32
	}

	items := make([]dirItem, 0, 2+len(realEntries))
	items = append(items, dirItem{".", in.NodeId, syscall.S_IFDIR | 0o555})
	items = append(items, dirItem{"..", dotdotIno, syscall.S_IFDIR | 0o555})
	for _, e := range realEntries {
		if e.Name == "" || e.Name == ":" {
			continue
		}
		childPath := filepath.Join(entry.Path, e.Name)
		childIno := f.inodes.GetOrCreate(childPath)
		mode := uint32(0o755)
		if e.IsDir {
			mode |= syscall.S_IFDIR
		} else {
			mode |= syscall.S_IFREG
		}
		items = append(items, dirItem{e.Name, childIno, mode})
	}

	offset := int(in.Offset)
	for i, item := range items {
		if i < offset {
			continue
		}
		de := gofuse.DirEntry{
			Name: item.name,
			Ino:  item.ino,
			Mode: item.mode,
			Off:  uint64(i + 1),
		}
		if !out.AddDirEntry(de) {
			break
		}
	}

	return gofuse.OK
}

// ReadDirPlus implements ReadDirPlus.
func (f *ReadOnlyFS) ReadDirPlus(cancel <-chan struct{}, in *gofuse.ReadIn, out *gofuse.DirEntryList) gofuse.Status {
	f.inodes.mu.RLock()
	entry, ok := f.inodes.byInode[in.NodeId]
	f.inodes.mu.RUnlock()

	if !ok {
		return gofuse.ENOENT
	}

	realEntries, err := f.getDirEntries(entry.Path)
	if err != nil {
		return gofuse.EIO
	}

	dotdotIno := in.NodeId
	if entry.Path != "/" {
		parentPath := filepath.Dir(entry.Path)
		dotdotIno = f.inodes.getIno(parentPath)
		if dotdotIno == 0 {
			dotdotIno = f.inodes.GetOrCreate(parentPath)
		}
	}

	type dirItem struct {
		name string
		ino  uint64
		mode uint32
	}

	items := make([]dirItem, 0, 2+len(realEntries))
	items = append(items, dirItem{".", in.NodeId, syscall.S_IFDIR | 0o555})
	items = append(items, dirItem{"..", dotdotIno, syscall.S_IFDIR | 0o555})
	for _, e := range realEntries {
		if e.Name == "" || e.Name == ":" {
			continue
		}
		childPath := filepath.Join(entry.Path, e.Name)
		childIno := f.inodes.GetOrCreate(childPath)
		mode := uint32(0o755)
		if e.IsDir {
			mode |= syscall.S_IFDIR
		} else {
			mode |= syscall.S_IFREG
		}
		items = append(items, dirItem{e.Name, childIno, mode})
	}

	offset := int(in.Offset)
	for i, item := range items {
		if i < offset {
			continue
		}
		de := gofuse.DirEntry{
			Name: item.name,
			Ino:  item.ino,
			Mode: item.mode,
			Off:  uint64(i + 1),
		}
		entryOut := out.AddDirLookupEntry(de)
		if entryOut == nil {
			break
		}
		entryOut.NodeId = item.ino
		entryOut.Generation = 1
		entryOut.Attr.Ino = item.ino
		entryOut.Attr.Size = 0
		if item.mode&syscall.S_IFDIR != 0 {
			entryOut.Attr.Mode = syscall.S_IFDIR | 0o755
		} else {
			entryOut.Attr.Mode = syscall.S_IFREG | 0o755
		}
		if item.mode&syscall.S_IFDIR != 0 {
			entryOut.Attr.Nlink = 2
		} else {
			entryOut.Attr.Nlink = 1
		}
		entryOut.Attr.Uid = f.uid
		entryOut.Attr.Gid = f.gid
		entryOut.SetEntryTimeout(1 * time.Second)
		entryOut.SetAttrTimeout(1 * time.Second)
	}

	return gofuse.OK
}

// Open implements Open.
func (f *ReadOnlyFS) Open(cancel <-chan struct{}, in *gofuse.OpenIn, out *gofuse.OpenOut) gofuse.Status {
	f.inodes.mu.RLock()
	entry, ok := f.inodes.byInode[in.NodeId]
	f.inodes.mu.RUnlock()

	if !ok {
		return gofuse.ENOENT
	}

	if entry.IsDir {
		return gofuse.EISDIR
	}

	out.Fh = in.NodeId
	out.OpenFlags = gofuse.FOPEN_KEEP_CACHE

	return gofuse.OK
}

// Read implements Read.
func (f *ReadOnlyFS) Read(cancel <-chan struct{}, in *gofuse.ReadIn, buf []byte) (gofuse.ReadResult, gofuse.Status) {
	f.inodes.mu.RLock()
	entry, ok := f.inodes.byInode[in.NodeId]
	f.inodes.mu.RUnlock()

	if !ok {
		return nil, gofuse.ENOENT
	}

	reader, err := f.client.ReadStreamRange(context.Background(), entry.Path, int64(in.Offset), int64(len(buf)))
	if err != nil {
		return nil, gofuse.EIO
	}
	defer reader.Close()

	n, err := reader.Read(buf)
	if err != nil && err != io.EOF {
		return nil, gofuse.EIO
	}

	return gofuse.ReadResultData(buf[:n]), gofuse.OK
}

// Release implements Release.
func (f *ReadOnlyFS) Release(cancel <-chan struct{}, in *gofuse.ReleaseIn) {
	// Nothing to do for read-only
}

// ReleaseDir implements ReleaseDir.
func (f *ReadOnlyFS) ReleaseDir(in *gofuse.ReleaseIn) {
	// Nothing to do
}

// Forget implements Forget.
func (f *ReadOnlyFS) Forget(nodeID uint64, nlookup uint64) {
	f.inodes.Forget(nodeID, nlookup)
}

// getDirEntries gets directory entries with caching.
func (f *ReadOnlyFS) getDirEntries(path string) ([]fs.FileInfo, error) {
	f.dirs.mu.RLock()
	cache, ok := f.dirs.entries[path]
	f.dirs.mu.RUnlock()

	if ok && time.Since(cache.Timestamp) < f.dirs.ttl {
		return cache.Items, nil
	}

	entries, err := f.client.List(path)
	if err != nil {
		return nil, err
	}

	f.dirs.mu.Lock()
	f.dirs.entries[path] = &dirCacheEntry{
		Items:     entries,
		Timestamp: time.Now(),
	}
	f.dirs.mu.Unlock()

	return entries, nil
}

// inodeManager methods

func (im *inodeManager) GetOrCreate(path string) uint64 {
	im.mu.Lock()
	defer im.mu.Unlock()

	if ino, ok := im.byPath[path]; ok {
		return ino
	}

	ino := atomic.AddUint64(&im.nextIno, 1)
	im.byPath[path] = ino
	im.byInode[ino] = &inodeEntry{
		Path:   path,
		RefCnt: 1,
	}
	return ino
}

func (im *inodeManager) getIno(path string) uint64 {
	im.mu.RLock()
	defer im.mu.RUnlock()
	return im.byPath[path]
}

func (im *inodeManager) Forget(ino uint64, nlookup uint64) {
	im.mu.Lock()
	defer im.mu.Unlock()

	entry, ok := im.byInode[ino]
	if !ok {
		return
	}

	if entry.RefCnt <= nlookup {
		delete(im.byInode, ino)
		delete(im.byPath, entry.Path)
	} else {
		entry.RefCnt -= nlookup
	}
}
