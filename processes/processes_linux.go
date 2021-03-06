// +build linux

// Originally inspired by the MIT-licensed
// https://github.com/mitchellh/go-ps/blob/master/process_unix.go

package processes

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/elastic/beats/libbeat/logp"
)

func (ps *Processes) Refresh() error {
	procs, err := processes(ps.exposeCmdline, ps.exposeEnviron, ps.procFSDirectory)
	if err != nil {
		return err
	}
	ps.byInode = make(map[uint64]*UnixProcess)

	for _, p := range procs {
		err := p.Refresh(ps.exposeCmdline, ps.exposeEnviron, ps.procFSDirectory)
		if err != nil {
			return err
		}
		for _, inode := range p.inodes {
			ps.byInode[inode] = p
		}
	}

	return nil
}

func findSocketsOfPid(procFSDirectory string, pid int) (inodes []uint64, err error) {

	dirname := filepath.Join(procFSDirectory, strconv.Itoa(pid), "fd")
	procfs, err := os.Open(dirname)
	if err != nil {
		return []uint64{}, fmt.Errorf("Open: %s", err)
	}
	defer procfs.Close()

	names, err := procfs.Readdirnames(0)
	if err != nil {
		return []uint64{}, fmt.Errorf("Readdirnames: %s", err)
	}

	for _, name := range names {
		link, err := os.Readlink(filepath.Join(dirname, name))
		if err != nil {
			logp.Debug("procs", "Readlink %s: %s", name, err)
			continue
		}

		if strings.HasPrefix(link, "socket:[") {
			inode, err := strconv.ParseInt(link[8:len(link)-1], 10, 64)
			if err != nil {
				logp.Debug("procs", "ParseInt: %s:", err)
				continue
			}

			inodes = append(inodes, uint64(inode))
		}
	}

	return inodes, nil
}

// Refresh reloads all the data associated with this process.
func (p *UnixProcess) Refresh(exposeCmdline, exposeEnviron bool, procFSDirectory string) error {

	inodes, err := findSocketsOfPid(procFSDirectory, p.pid)
	p.inodes = inodes

	if exposeCmdline {
		p.Cmdline, err = readFile(procFSDirectory, p.pid, "cmdline")
	}
	if exposeEnviron {
		p.Environ, err = readFile(procFSDirectory, p.pid, "environ")
	}

	return err
}

func readFile(procFSDirectory string, pid int, filename string) (string, error) {
	path := fmt.Sprintf("%s/%d/%s", procFSDirectory, pid, filename)
	bytes, err := ioutil.ReadFile(path)
	return string(bytes), err
}

func processes(exposeCmdline, exposeEnviron bool, procFSDirectory string) ([]*UnixProcess, error) {
	d, err := os.Open(procFSDirectory)
	if err != nil {
		return nil, err
	}
	defer d.Close()

	results := make([]*UnixProcess, 0, 50)
	for {
		fis, err := d.Readdir(10)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		for _, fi := range fis {
			// We only care about directories, since all pids are dirs
			if !fi.IsDir() {
				continue
			}

			// We only care if the name starts with a numeric
			name := fi.Name()
			if name[0] < '0' || name[0] > '9' {
				continue
			}

			// From this point forward, any errors we just ignore, because
			// it might simply be that the process doesn't exist anymore.
			pid, err := strconv.ParseInt(name, 10, 0)
			if err != nil {
				continue
			}

			p, err := newUnixProcess(int(pid), exposeCmdline, exposeEnviron, procFSDirectory)
			if err != nil {
				continue
			}

			results = append(results, p)
		}
	}

	return results, nil
}

func newUnixProcess(pid int, exposeCmdline, exposeEnviron bool, procFSDirectory string) (*UnixProcess, error) {
	p := &UnixProcess{pid: pid}
	return p, p.Refresh(exposeCmdline, exposeEnviron, procFSDirectory)
}

func (ps *Processes) FindProcessByInode(inode uint64) *UnixProcess {
	if inode == 0 {
		return nil
	}

	proc := ps.byInode[inode]
	if proc == nil {
		// Refesh and try again
		ps.Refresh()

		return ps.byInode[inode]
	}
	return proc
}
