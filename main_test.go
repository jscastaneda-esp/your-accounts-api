package main

import (
	"errors"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	originalDotenvLoad func(filenames ...string) error
}

func (suite *TestSuite) SetupSuite() {
	suite.originalDotenvLoad = dotenvLoad

	osStat = func(name string) (os.FileInfo, error) {
		return nil, nil
	}
	newDB = func() {
		log.Println("Connect database")
	}
	newMongoClient = func() {
		log.Println("Connect mongo database")
	}
	newServer = func() {
		log.Println("Starting server")
	}
}

func (suite *TestSuite) SetupTest() {
	dotenvLoad = suite.originalDotenvLoad
}

func (suite *TestSuite) TestMainSuccess() {
	dotenvLoad = func(filenames ...string) error {
		return nil
	}

	main()
}

func (suite *TestSuite) TestMainErrorLoadFileEnv() {
	require := require.New(suite.T())

	dotenvLoad = func(_ ...string) (err error) {
		err = errors.New("Error loading file")
		return
	}

	var fatal = false
	logFatal = func(_ ...any) {
		fatal = true
	}

	main()
	require.True(fatal)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
