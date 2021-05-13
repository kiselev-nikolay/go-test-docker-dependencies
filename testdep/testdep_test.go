package testdep_test

import (
	"testing"

	"github.com/kiselev-nikolay/go-test-docker-dependencies/testdep"
)

func TestPostgres(t *testing.T) {
	stop, err := testdep.Postgres(5432, "test", "test", "test")
	if err != nil {
		t.Error(err)
		return
	}
	err = stop()
	if err != nil {
		t.Error(err)
		return
	}
}
