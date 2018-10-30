package executor

import (
	"io"
	"os"

	"github.com/dev9/prod/influxdata-operator/pkg/controller/backup/executor/influxd"
)

const (
	// MySQLDumpProvider denotes the mysqldump utility backup and restore provider.
	DumpProvider = "influxd"
)

// ExecutorProviders denotes the list of available ExecutorProviders.
var ExecutorProviders = [...]string{DumpProvider}

// Interface will execute backup operations via a tool such as mysqlbackup or
// mysqldump.
type Interface interface {
	// Backup runs a backup operation using the given credentials, returning the content.
	// TODO: default backupDir to allow streaming...
	Backup(backupDir string, clusterName string) (io.ReadCloser, string, error)
	// Restore restores the given content to the mysql node.
	Restore(content io.ReadCloser) error
}

// New builds a new backup executor.
func New(executor v1alpha1.BackupExecutor, creds map[string]string) (Interface, error) {
	return influxd.NewExecutor(executor.influxd, creds)
}

