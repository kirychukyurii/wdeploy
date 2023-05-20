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
	"os"
	"path/filepath"
)

var playbookFileName = "playbook.yml"

type Executor struct {
	cfg    config.Config
	logger logger.Logger
}

func NewExecutor(cfg config.Config, logger logger.Logger) Executor {
	return Executor{
		cfg:    cfg,
		logger: logger,
	}
}

func (e Executor) RunPlaybook() error {
	f, err := os.OpenFile(e.cfg.GetAnsibleLogLocation(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		e.logger.Zap.Error(err)
	}

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
			execute.WithWrite(f),
			execute.WithWriteError(f),
		),
	)

	pb := &playbook.AnsiblePlaybookCmd{
		Playbooks:         []string{filepath.Join(e.cfg.PlaybookTempDir, playbookFileName)},
		ConnectionOptions: ansiblePlaybookConnectionOptions,
		Options:           ansiblePlaybookOptions,
		Exec:              executorTimeMeasurement,
		StdoutCallback:    "yaml",
	}

	err = pb.Run(context.TODO())
	if err != nil {
		e.logger.Zap.Error(err)
	}

	e.logger.Zap.Info(executorTimeMeasurement.Duration())

	return nil
}
