package transform

import (
	"path/filepath"

	"github.com/fusor/cpma/pkg/env"
	"github.com/fusor/cpma/pkg/io"
	"github.com/fusor/cpma/pkg/transform/oauth"
	"github.com/fusor/cpma/pkg/transform/secrets"
	"github.com/sirupsen/logrus"
)

// OCP4InstallMsg message about using generated manifests
const OCP4InstallMsg = `To install OCP4 run the installer as follow in order to add CRDs:
' /openshift-install --dir $INSTALL_DIR create install-config'
'./openshift-install --dir $INSTALL_DIR create manifests'
# Copy generated CRD manifest files  to '$INSTALL_DIR/openshift/'
# Edit them if needed, then run installation:
'./openshift-install --dir $INSTALL_DIR  create cluster'`

// Cluster contains a cluster
type Cluster struct {
	Master Master
}

// Master is a cluster Master
type Master struct {
	OAuth   oauth.CRD
	Secrets []secrets.Secret
}

// Manifest to be exported for use with OCP 4
type Manifest struct {
	Name string
	CRD  []byte
}

// Config contains CPMA configuration information
type Config struct {
	OutputDir string
	Hostname  string
}

// Runner a generic transform runner
type Runner struct {
	Config string
}

// Extraction is a generic data extraction
type Extraction interface {
	Transform() (Output, error)
	Validate() error
}

// Transform is a generic transform
type Transform interface {
	Extract() (Extraction, error)
	Type() string
}

// Output is a generic output type
type Output interface {
	Flush() error
}

// GetFile allows to mock file retrieval
var GetFile = io.GetFile

//Start generating manifests to be used with Openshift 4
func Start() {
	config := LoadConfig()
	runner := NewRunner(config)

	runner.Transform([]Transform{
		OAuthTransform{
			Config: &config,
		},
		SDNTransform{
			Config: &config,
		},
		RegistriesTransform{
			Config: &config,
		},
	})
}

// LoadConfig collects and stores configuration for CPMA
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

// Transform is the process run to complete a transform
func (r Runner) Transform(transforms []Transform) {
	logrus.Info("TransformRunner::Transform")

	// For each transform, extract the data, validate it, and run the transform.
	// Handle any errors, and finally flush the output to it's desired destination
	// NOTE: This should be parallelized with channels unless the transforms have
	// some dependency on the outputs of others
	for _, transform := range transforms {
		extraction, err := transform.Extract()
		if err != nil {
			HandleError(err, transform.Type())
			continue
		}

		if err := extraction.Validate(); err != nil {
			HandleError(err, transform.Type())
			continue
		}

		output, err := extraction.Transform()
		if err != nil {
			HandleError(err, transform.Type())
			continue
		}

		if err := output.Flush(); err != nil {
			HandleError(err, transform.Type())
			continue
		}
	}
}

// NewRunner creates a new Runner
func NewRunner(config Config) *Runner {
	return &Runner{}
}

// HandleError handles errors
func HandleError(err error, transformType string) error {
	logrus.Warnf("Skipping %s, see error below\n", transformType)
	logrus.Warnf("%s\n", err)
	return err
}