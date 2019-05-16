package report

import (
	"path/filepath"

	"github.com/fusor/cpma/pkg/env"
	"github.com/fusor/cpma/pkg/io"
	"github.com/sirupsen/logrus"
)

// Config contains CPMA configuration information
type Config struct {
	OutputDir string
	Hostname  string
}

// Runner a generic report runner
type Runner struct {
	Config string
}

// Extraction is a generic data extraction
type Extraction interface {
	Report() (Output, error)
	Validate() error
}

// Report is a generic report
type Report interface {
	Extract() (Extraction, error)
	Type() string
}

// Output is a generic output type
type Output interface {
	Flush() error
}

// Report to be exported for use with OCP 4
type ReportData struct {
	Name       string
	Type       string
	ReportInfo []byte
}

// GetFile allows to mock file retrieval
var GetFile = io.GetFile

//Start generating a report component
func Start() {
	config := LoadConfig()
	runner := NewRunner(config)
	runner.Report([]Report{
		SDNReport{
			Config: &config,
		},
	})
}

// LoadConfig collects and stores report for CPMA
func LoadConfig() Config {
	logrus.Info("Loaded config")

	config := Config{}
	config.OutputDir = env.Config().GetString("OutputDir")
	config.Hostname = env.Config().GetString("Source")

	return config
}

// Fetch files from the OCP3 cluster
func (config *Config) Fetch(path string) ([]byte, error) {
	dst := filepath.Join(config.OutputDir, config.Hostname, path)
	logrus.Infof("Fetching file: %s", dst)
	f, err := GetFile(config.Hostname, path, dst)
	if err != nil {
		return nil, err
	}
	logrus.Infof("File successfully loaded: %v", dst)

	return f, nil
}

// NewRunner creates a new Runner
func NewRunner(config Config) *Runner {
	return &Runner{}
}

// Report is the process run to complete a report
func (r Runner) Report(reports []Report) {
	logrus.Info("ReportRunner::Report")

	// For each report, extract the data, validate it, and run the report.
	// Handle any errors, and finally flush the output to it's desired destination
	// NOTE: This should be parallelized with channels unless the reports have
	// some dependency on the outputs of others
	for _, report := range reports {
		extraction, err := report.Extract()
		if err != nil {
			HandleError(err, report.Type())
			continue
		}

		if err := extraction.Validate(); err != nil {
			HandleError(err, report.Type())
			continue
		}

		output, err := extraction.Report()
		if err != nil {
			HandleError(err, report.Type())
			continue
		}

		if err := output.Flush(); err != nil {
			HandleError(err, report.Type())
			continue
		}
	}
}

// HandleError handles errors
func HandleError(err error, reportType string) error {
	logrus.Warnf("Skipping %s, see error below\n", reportType)
	logrus.Warnf("%s\n", err)
	return err
}
