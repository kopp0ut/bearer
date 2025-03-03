package flag

var (
	ParentProcessIDFlag = Flag{
		Name:       "parent-process-id",
		ConfigName: "process.parent-process-id",
		Value:      -1,
	}

	PortFlag = Flag{
		Name:       "port",
		ConfigName: "process.port",
		Shorthand:  "p",
		Value:      "",
		Usage:      "Set the server's listening port.",
	}

	WorkerIDFlag = Flag{
		Name:       "worker-id",
		ConfigName: "process.worker-id",
		Value:      "",
		Usage:      "Set the worker's identifier.",
	}
)

type ProcessFlagGroup struct {
	ParentProcessIDFlag *Flag
	PortFlag            *Flag
	WorkerIDFlag        *Flag
}

type ProcessOptions struct {
	ParentProcessID int
	WorkerID        string `mapstructure:"worker-id" json:"worker-id" yaml:"worker-id"`
	Port            string `mapstructure:"port" json:"port" yaml:"port"`
}

func NewProcessGroup() *ProcessFlagGroup {
	return &ProcessFlagGroup{
		ParentProcessIDFlag: &ParentProcessIDFlag,
		PortFlag:            &PortFlag,
		WorkerIDFlag:        &WorkerIDFlag,
	}
}

func (f *ProcessFlagGroup) Name() string {
	return "process"
}

func (f *ProcessFlagGroup) Flags() []*Flag {
	return []*Flag{
		f.ParentProcessIDFlag,
		f.PortFlag,
		f.WorkerIDFlag,
	}
}

func (f *ProcessFlagGroup) ToOptions() (ProcessOptions, error) {
	port := getString(f.PortFlag)
	workerID := getString(f.WorkerIDFlag)

	return ProcessOptions{
		ParentProcessID: getInteger(f.ParentProcessIDFlag),
		Port:            port,
		WorkerID:        workerID,
	}, nil
}
