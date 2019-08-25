package utils

import (
	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()
var LogSummary = logrus.New().WithField("tag", "summary")
