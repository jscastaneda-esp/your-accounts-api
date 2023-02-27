package main

import (
	"errors"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestMainSuccess(t *testing.T) {
	osStat = func(name string) (os.FileInfo, error) {
		return nil, nil
	}
	dotenvLoad = func(filenames ...string) (err error) {
		return nil
	}
	newDB = func() *gorm.DB {
		return &gorm.DB{}
	}
	newServer = func(db *gorm.DB) {
		log.Printf("Starting server with database %v", db)
	}

	main()
}

func TestMainErrorLoadFileEnv(t *testing.T) {
	osStat = func(_ string) (os.FileInfo, error) {
		return nil, nil
	}
	dotenvLoad = func(_ ...string) (err error) {
		err = errors.New("Error loading file")
		return
	}

	var fatal = false
	logFatal = func(_ ...any) {
		fatal = true
	}

	require := require.New(t)
	main()
	require.True(fatal)
}
