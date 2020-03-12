package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProfx(t *testing.T) {
	main()
}

func TestProfx_Crawl(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"profx", "-crawl"}

	err := os.Setenv("PROFX_RUN_COUNT", "1")
	assert.NoError(t, err)

	assert.Fail(t, "fix test run")
	// main()
}
