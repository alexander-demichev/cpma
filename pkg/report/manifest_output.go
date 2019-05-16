package report

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/fusor/cpma/pkg/env"
	"github.com/sirupsen/logrus"
)

// ManifestOutput holds a collection of manifests to be written to fil
type ReportOutput struct {
	Reports []ReportData
}

// Flush calls DumpManifests to write file data
func (m ReportOutput) Flush() error {
	logrus.Info("Writing file data:")
	DumpManifests(m.Reports)
	return nil
}

// DumpManifests creates OCDs files
func DumpManifests(reports []ReportData) {
	for _, report := range reports {
		reportfile := filepath.Join(env.Config().GetString("OutputDir"), "reports", report.Type, report.Name)
		os.MkdirAll(path.Dir(reportfile), 0755)
		err := ioutil.WriteFile(reportfile, report.ReportInfo, 0644)
		logrus.Printf("Report:Added: %s", reportfile)
		if err != nil {
			logrus.Panic(err)
		}
	}
}
