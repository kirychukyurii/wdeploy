package ansible

import (
	"context"
	"fmt"
	"github.com/apenella/go-ansible/pkg/execute"
	"github.com/apenella/go-ansible/pkg/execute/measure"
	"github.com/apenella/go-ansible/pkg/options"
	"github.com/apenella/go-ansible/pkg/playbook"
	"github.com/apenella/go-ansible/pkg/stdoutcallback/results"
	"github.com/kirychukyurii/wdeploy/internal/config"
	"github.com/kirychukyurii/wdeploy/internal/constants"
	"github.com/kirychukyurii/wdeploy/internal/lib/logger"
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
	ansiblePlaybookConnectionOptions := &options.AnsibleConnectionOptions{
		SSHCommonArgs: e.cfg.AnsibleSSHExtraArgs,
	}

	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
		Inventory:     e.cfg.HostsFile,
		ExtraVarsFile: []string{e.cfg.VarsFile},
	}

	executorTimeMeasurement := measure.NewExecutorTimeMeasurement(
		execute.NewDefaultExecute(
			execute.WithEnvVar("ANSIBLE_FORCE_COLOR", "true"),
			execute.WithTransformers(
				results.Prepend(fmt.Sprintf("Webitel v%s", e.cfg.WebitelVersion)),
				results.LogFormat(constants.TimeFormat, results.Now),
			),
		),
		measure.WithShowDuration(),
	)

	pb := &playbook.AnsiblePlaybookCmd{
		Playbooks:         []string{e.cfg.PlaybookFile},
		ConnectionOptions: ansiblePlaybookConnectionOptions,
		Options:           ansiblePlaybookOptions,
		Exec:              executorTimeMeasurement,
		//StdoutCallback:    "json",
	}

	err := pb.Run(context.TODO())
	if err != nil {
		e.logger.Zap.Error(err)
	}

	e.logger.Zap.Info(executorTimeMeasurement.Duration())

	return nil
}
