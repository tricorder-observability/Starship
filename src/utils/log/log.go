// Package log initialize logurs
package log

import "github.com/sirupsen/logrus"

// init enables logging logging statement lines.
func init() {
	logrus.SetReportCaller(true)
}
