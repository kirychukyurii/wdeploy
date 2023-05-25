package ansible

import (
	"context"
	"fmt"
	"github.com/apenella/go-ansible/pkg/execute"
	"github.com/apenella/go-ansible/pkg/execute/measure"
	"github.com/apenella/go-ansible/pkg/options"
	"github.com/apenella/go-ansible/pkg/playbook"
	"github.com/kirychukyurii/wdeploy/internal/config"
	"github.com/kirychukyurii/wdeploy/internal/lib/logger"
	"io"
	"path/filepath"
)

const (
	playbookFileName    = "playbook.yml"
	unixyStdoutCallback = "unixy"
)

type Executor struct {
	cfg    config.Config
	logger logger.Logger
	writer io.Writer
}

func NewExecutor(cfg config.Config, logger logger.Logger, writer io.Writer) Executor {
	return Executor{
		cfg:    cfg,
		logger: logger,
		writer: writer,
	}
}

func (e Executor) RunPlaybook() error {
	ansiblePlaybookConnectionOptions := &options.AnsibleConnectionOptions{
		SSHCommonArgs: e.cfg.AnsibleSSHExtraArgs,
	}

	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
		Inventory:     e.cfg.ConfigFiles[config.InventoryConfig],
		ExtraVarsFile: []string{fmt.Sprintf("@%s", e.cfg.ConfigFiles[config.VarsConfig])},
	}

	executorTimeMeasurement := measure.NewExecutorTimeMeasurement(
		execute.NewDefaultExecute(
			execute.WithEnvVar("ANSIBLE_FORCE_COLOR", "true"),
			execute.WithEnvVar("ANSIBLE_STDOUT_CALLBACK", unixyStdoutCallback),
			execute.WithWrite(e.writer),
			execute.WithWriteError(e.writer),
		),
	)

	pb := &playbook.AnsiblePlaybookCmd{
		Playbooks:         []string{filepath.Join(e.cfg.PlaybookTempDir, playbookFileName)},
		ConnectionOptions: ansiblePlaybookConnectionOptions,
		Options:           ansiblePlaybookOptions,
		Exec:              executorTimeMeasurement,
	}

	err := pb.Run(context.TODO())
	if err != nil {
		e.logger.Zap.Error(err)
	}

	e.logger.Zap.Info(executorTimeMeasurement.Duration())

	return nil
}
