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
)

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
	f, err := os.OpenFile(fmt.Sprintf("%s/ansible.log", e.cfg.LogDirectory), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	ansiblePlaybookConnectionOptions := &options.AnsibleConnectionOptions{
		SSHCommonArgs: e.cfg.AnsibleSSHExtraArgs,
	}

	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
		Inventory: e.cfg.HostsFile,
		//ExtraVarsFile: []string{e.cfg.VarsFile},
	}

	executorTimeMeasurement := measure.NewExecutorTimeMeasurement(
		execute.NewDefaultExecute(
			execute.WithEnvVar("ANSIBLE_FORCE_COLOR", "true"),
			execute.WithWrite(f),
		),
	)

	pb := &playbook.AnsiblePlaybookCmd{
		Playbooks:         []string{e.cfg.PlaybookFile},
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
