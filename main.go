package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/dlopes7/go-mssql-connector/connector"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io/ioutil"
	"os"
	"path"
	"time"
)

type Query struct {
	Name        string `json:"name"`
	QueryString string `json:"query"`
}

type Config struct {
	Queries               []Query `json:"queries"`
	Host                  string  `json:"host"`
	User                  string  `json:"username"`
	Password              string  `json:"password"`
	Port                  int     `json:"port"`
	Database              string  `json:"database"`
	LogLevel              string  `json:"logLevel"`
	WindowsAuthentication bool    `json:"windowsAuthentication"`
}

var log *logrus.Logger

func setupLog(tempFolder string, endpointID string) {
	log = logrus.New()
	logDirPath := path.Join(tempFolder, "log")

	if _, err := os.Stat(logDirPath); os.IsNotExist(err) {
		fmt.Printf("Creating log folder: %s\n", logDirPath)
		_ = os.Mkdir(logDirPath, os.ModePerm)
	}

	logFilePath := path.Join(logDirPath, fmt.Sprintf("%s.log", endpointID))
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    5,
		MaxBackups: 5,
	}
	log.SetOutput(lumberjackLogger)
	log.SetLevel(logrus.InfoLevel)
}

func loadConfig(endpointID string, tempFolder string) *Config {
	fileName := fmt.Sprintf("%s.json", endpointID)
	configFilePath := path.Join(tempFolder, "config", fileName)

	log.WithFields(logrus.Fields{"configFilePath": configFilePath}).Info("Reading configuration")
	configFile, err := os.Open(configFilePath)
	if err != nil {
		log.Fatalf("Could not read %s: %s", fileName, err.Error())
	}
	defer configFile.Close()

	c := new(Config)
	byteValue, _ := ioutil.ReadAll(configFile)
	err = json.Unmarshal(byteValue, c)
	if err != nil {
		log.Fatalf("Could not parse the configuration file: %s", err.Error())
	}

	level, err := logrus.ParseLevel(c.LogLevel)
	if err != nil {
		log.SetLevel(logrus.InfoLevel)
	} else {
		log.WithFields(logrus.Fields{"level": level}).Info(fmt.Sprintf("Setting log level from %s", fileName))
		log.SetLevel(level)
	}
	if os.Getenv("DATABASE_PASSWORD") != "" {
		c.Password = os.Getenv("DATABASE_PASSWORD")
	}

	return c
}

func writeResponse(endpointID string, tempFolder string, response *connector.Response) {
	jsonResponse, err := json.MarshalIndent(response, "", "  ")
	if err == nil {

		resultPath := path.Join(tempFolder, "results")
		if _, err := os.Stat(resultPath); os.IsNotExist(err) {
			_ = os.Mkdir(resultPath, os.ModePerm)
		}

		resultFilePath := path.Join(resultPath, fmt.Sprintf("%s.json", endpointID))
		err := ioutil.WriteFile(resultFilePath, jsonResponse, os.ModePerm)
		if err != nil {
			log.Errorf("Could not write results: %s", err)
		}
	}
}

func main() {

	response := connector.NewResponse()

	endpointID := flag.String("endpoint", "", "Endpoint ID")
	tempFolder := flag.String("tempfolder", "", "Temp Folder")
	flag.Parse()

	if endpointID == nil || *endpointID == "" {
		response.Error = true
		response.ErrorMessage = "The parameter endpoint must be a valid endpoint ID"
		fmt.Printf("%+v", response)
	} else if tempFolder == nil || *tempFolder == "" {
		response.Error = true
		response.ErrorMessage = "The parameter tempfolder must be a path"
		fmt.Printf("%+v", response)
	} else {

		setupLog(*tempFolder, *endpointID)
		c := loadConfig(*endpointID, *tempFolder)
		dbConnection := &connector.MSSQLConnector{
			Log: log,
		}

		if len(c.Queries) > 0 {
			db, err := dbConnection.GetDB(c.Host, c.Port, c.User, c.Password, c.Database, c.WindowsAuthentication)
			if err != nil {
				response.Error = true
				response.ErrorMessage = err.Error()
				log.WithFields(logrus.Fields{"Error": err.Error()}).Error("Error obtaining DB")
			} else {
				defer db.Close()
				for _, query := range c.Queries {
					log.WithFields(logrus.Fields{"Query": query.Name, "QueryString": query.QueryString}).Info("Running query")
					start := time.Now()
					qr := connector.Query(query.QueryString, db)
					qr.Duration = time.Now().Sub(start).Nanoseconds()
					qr.Name = query.Name
					response.Queries = append(response.Queries, qr)
				}
			}
		}
		writeResponse(*endpointID, *tempFolder, response)
	}

}
