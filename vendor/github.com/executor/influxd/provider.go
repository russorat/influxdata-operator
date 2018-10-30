package influxd 

import (
	"compress/gzip"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/pkg/errors"

	utilexec "k8s.io/utils/exec"

)

const (
	sqldumpCmd = "influxd"
	sqlCmd     = "influx"
)

// Executor creates backups using influxd.
type Executor struct {
	config *Config
}

// NewExecutor creates a provider capable of creating and restoring backups with the mysqldump
// tool.
func NewExecutor(executor *v1alpha1.DumpBackupExecutor, creds map[string]string) (*Executor, error) {
	cfg := NewConfig(executor, creds)
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	return &Executor{config: cfg}, nil
}

// Backup performs a full cluster backup using the influxd tool.
func (ex *Executor) Backup(backupDir string, clusterName string) (io.ReadCloser, string, error) {
	exec := utilexec.New()
	mysqldumpPath, err := exec.LookPath(sqldumpCmd)
	if err != nil {
		return nil, "", fmt.Errorf("influxd path: %v", err)
	}

	args := []string{
		"-u" + ex.config.username,
		"-p" + ex.config.password,
		"--single-transaction",
		"--skip-lock-tables",
		"--flush-privileges",
		"--set-gtid-purged=OFF",
		"--databases",
	}

	dbNames := make([]string, len(ex.config.databases))
	for i, database := range ex.config.databases {
		dbNames[i] = database.Name
	}

	cmd := exec.Command(mysqldumpPath, append(args, dbNames...)...)

	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	backupName := fmt.Sprintf("%s.%s.sql.gz", clusterName, time.Now().UTC().Format("20060102150405"))

	pr, pw := io.Pipe()
	zw := gzip.NewWriter(pw)
	cmd.SetStdout(zw)

	go func() {
		glog.V(4).Infof("running cmd: '%s %s'", mysqldumpPath, SanitizeArgs(append(args, dbNames...), ex.config.password))
		err = cmd.Run()
		zw.Close()
		if err != nil {
			pw.CloseWithError(errors.Wrap(err, "executing backup"))
		} else {
			pw.Close()
		}
	}()

	return pr, backupName, nil
}
