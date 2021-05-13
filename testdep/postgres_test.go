package testdep_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/jackc/pgconn"
	"github.com/kiselev-nikolay/go-test-docker-dependencies/testdep"
)

const DockerTimeout = 20

func AssertNoError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
		t.FailNow()
		return
	}
}

func TestPostgres(t *testing.T) {
	port, err := testdep.FindFreePort()
	AssertNoError(t, err)
	pg := &testdep.Postgres{
		port,
		"test",
		"test",
		"test",
	}
	stop, err := pg.Run(DockerTimeout)
	AssertNoError(t, err)
	pgConn, err := pgconn.Connect(context.Background(), fmt.Sprintf("host=localhost port=%d user=test password=test dbname=test", port))
	AssertNoError(t, err)
	if pgConn.PID() == 0 {
		t.Fail()
	}
	err = stop()
	AssertNoError(t, err)
}

func TestPostgresPortUsed(t *testing.T) {
	port, err := testdep.FindFreePort()
	AssertNoError(t, err)
	pg := &testdep.Postgres{
		port,
		"test",
		"test",
		"test",
	}
	stop, _ := pg.Run(DockerTimeout)
	_, err = pg.Run(DockerTimeout)
	if err == nil {
		t.Errorf("expected error, because other container supposted to bind this port")
		t.Fail()
	}
	if !strings.Contains(err.Error(), "failed: port is already allocated") {
		t.Error(err.Error())
		t.Fail()
	}
	stop()
}

func TestPostgresStartTimeout(t *testing.T) {
	port, err := testdep.FindFreePort()
	AssertNoError(t, err)
	pg := &testdep.Postgres{
		port,
		"test",
		"test",
		"test",
	}
	stop, err := pg.Run(0)
	if err != testdep.HealthCheckTimeoutErr {
		t.Error(err)
		t.Fail()
	}
	stop()
}
