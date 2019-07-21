package config

import (
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

// Event stores messages to log later, from our standard interface
type Event struct {
	id      int
	message string
}

// StandardLogger enforces specific log message formats
type StandardLogger struct {
	*logrus.Logger
}

// NewLogger initializes the standard logger
func NewLogger() *StandardLogger {
	var baseLogger = logrus.New()
	baseLogger.Out = os.Stdout
	baseLogger.SetLevel(logrus.DebugLevel)

	// Create the log file if doesn't exist. And append to it if it already exists.
	// file, err := os.OpenFile(time.Now().Format(time.RFC3339)+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// if err == nil {
	// 	baseLogger.Out = file
	// } else {
	// 	baseLogger.Info("Failed to log to file, using default stderr")
	// }
	Formatter := new(logrus.TextFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true
	baseLogger.SetFormatter(Formatter)

	var standardLogger = &StandardLogger{baseLogger}
	standardLogger.Formatter = &logrus.JSONFormatter{}
	return standardLogger
}

// Declare variables to store log messages as new Events
var (
	initStatusArgMessage              = Event{1, "Initializing: %s"}
	environmentVariablesData          = Event{2, "Environment variables has been set"}
	listeningPort                     = Event{3, "Listening to the port %s"}
	settingUpRouters                  = Event{4, "Setting Routers..."}
	requestingData                    = Event{5, "Requesting data..."}
	noUserSet                         = Event{6, "The request must have an user set: %s"}
	unableToRequest                   = Event{7, "Unable to request data : %s"}
	restRequestTemplateCreationError  = Event{8, "Error while creating a REST template: %s"}
	restRequestTemplateExecutionError = Event{9, "Error while executing a REST template: %s"}
)

// InitFunction is a standard init function message
func (l *StandardLogger) InitFunction(argumentName string) {
	l.Infof(initStatusArgMessage.message, argumentName)
}

// EnvVariablesData logs the envrionment varibles, if there is a problem that one of them is not set it quits
func (l *StandardLogger) EnvVariablesData() {
	envVarsData := logrus.Fields{
		"GITHUB_PROPERTIES_ENDPOINT": os.Getenv("GITHUB_PROPERTIES_ENDPOINT"),
		"GITHUB_USER_STARRED":        os.Getenv("GITHUB_USER_STARRED"),
		"HOST":                       os.Getenv("HOST"),
	}
	if os.Getenv("GITHUB_PROPERTIES_ENDPOINT") == "" || os.Getenv("GITHUB_USER_STARRED") == "" || os.Getenv("HOST") == "" {
		l.WithFields(envVarsData).Error("Environment variables must be set.")
		os.Exit(0)
	}
	l.WithFields(envVarsData).Info(environmentVariablesData.message)
}

// ListeningPort logs the port exposed
func (l *StandardLogger) ListeningPort(message string) {
	l.Infof(listeningPort.message, message)
}

// SettingUpRouters the status of setting up routes
func (l *StandardLogger) SettingUpRouters() {
	l.Infof(settingUpRouters.message)
}

// MiddlewareRequest logs every request made
func (l *StandardLogger) MiddlewareRequest(r *http.Request) {
	requestData := logrus.Fields{
		"URL":    r.URL,
		"Header": r.Header,
		"Body":   r.Body,
	}
	l.WithFields(requestData).Infof(requestingData.message)
}

// NoUserSet logs that no user was sent with the request
func (l *StandardLogger) NoUserSet(user string) {
	l.Errorf(noUserSet.message, user)
}

// UnableToRequest logs an error encountered while request some data to and external server
func (l *StandardLogger) UnableToRequest(err string) {
	l.Errorf(unableToRequest.message, err)
}

// CreatingRestTemplateError logs the erro for creating a reat template
func (l *StandardLogger) CreatingRestTemplateError(err string) {
	l.Errorf(restRequestTemplateCreationError.message, err)
}

// ExecutinRestTemplateError logs the erro for creating a reat template
func (l *StandardLogger) ExecutinRestTemplateError(err string) {
	l.Errorf(restRequestTemplateExecutionError.message, err)
}
