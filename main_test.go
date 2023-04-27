package main

import (
	"api-your-accounts/shared/infrastructure"
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
	newServer = func(_ bool) *infrastructure.Server {
		return infrastructure.NewServer(true)
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

	require.Panics(func() {
		main()
	})
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
