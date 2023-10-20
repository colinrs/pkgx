package os

import "syscall"

// DiskUsage ..
type DiskUsage struct {
	stat *syscall.Statfs_t
}

// NewDiskUsage Returns an object holding the disk usage of volumePath
// diskUsage function assumes volumePath is a valid path
func NewDiskUsage(volumePath string) *DiskUsage {

	var stat syscall.Statfs_t
	err := syscall.Statfs(volumePath, &stat)
	if err != nil {
		return nil
	}
	return &DiskUsage{&stat}
}

// Free Total free bytes on file system
func (diskUsage *DiskUsage) Free() uint64 {
	return diskUsage.stat.Bfree * uint64(diskUsage.stat.Bsize)
}

// Available Total available bytes on file system to an unpriveleged user
func (diskUsage *DiskUsage) Available() uint64 {
	return diskUsage.stat.Bavail * uint64(diskUsage.stat.Bsize)
}

// Size Total size of the file system
func (diskUsage *DiskUsage) Size() uint64 {
	return diskUsage.stat.Blocks * uint64(diskUsage.stat.Bsize)
}

// Used Total bytes used in file system
func (diskUsage *DiskUsage) Used() uint64 {
	return diskUsage.Size() - diskUsage.Free()
}

// Usage UsagePercentage of use on the file system
func (diskUsage *DiskUsage) Usage() float32 {
	return float32(diskUsage.Used()) / float32(diskUsage.Size())
}
