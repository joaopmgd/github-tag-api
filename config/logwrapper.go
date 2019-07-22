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
	databaseConnectionError           = Event{10, "Error while trying to connect to the PostgreSQL database: %s"}
	couldNotParseRequestBody          = Event{11, "Could not parse request body : %s"}
	repoNotFound                      = Event{11, "Repository with id %s was not found"}
	repoAlreadyTagged                 = Event{12, "Repository already has the tag %s"}
	stringToInt64Error                = Event{13, "Error while converting the string %s to int64"}
	pageIsBiggerThanRequestValues     = Event{14, "Requested page is bigger than requested value limit %s, offset %s"}
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
		"GITHUB_HEALTH_STATUS":       os.Getenv("GITHUB_HEALTH_STATUS"),
		"HOST":                       os.Getenv("HOST"),
	}
	if os.Getenv("GITHUB_PROPERTIES_ENDPOINT") == "" ||
		os.Getenv("GITHUB_USER_STARRED") == "" ||
		os.Getenv("GITHUB_HEALTH_STATUS") == "" ||
		os.Getenv("HOST") == "" {
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

// DatabaseConnectionError details the error while connectiong to the database
func (l *StandardLogger) DatabaseConnectionError(reason string) {
	l.Errorf(databaseConnectionError.message, reason)
}

// CouldNotParseRequestBody logs if the body request could no be parsed
func (l *StandardLogger) CouldNotParseRequestBody(err string) {
	l.Errorf(couldNotParseRequestBody.message, err)
}

// RepoNotFound logs if the repo is not found
func (l *StandardLogger) RepoNotFound(id string) {
	l.Errorf(repoNotFound.message, id)
}

// RepoAlreadyTagged logs if the repo already has the tag
func (l *StandardLogger) RepoAlreadyTagged(tag string) {
	l.Errorf(repoAlreadyTagged.message, tag)
}

// StringToInt64Error details the error while trying to convert a string to a int number
func (l *StandardLogger) StringToInt64Error(number string) {
	l.Errorf(stringToInt64Error.message, number)
}

// PageIsBiggerThanRequestValues details a warning while the requested page is is bigger than the requested value
func (l *StandardLogger) PageIsBiggerThanRequestValues(limit, offset string) {
	l.Errorf(pageIsBiggerThanRequestValues.message, limit, offset)
}
