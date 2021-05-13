# Go Test Docker Dependencies

Go test docker dependencies is used to declare the dependencies required to run the test with something running in docker. For example this could be a database.

Currently available:
+ PostgreSQL

    ```go
    testdep.Postgres{
      Port:     5432,
      User:     "admin",
      Password: "Admin123",
      Database: "nicedb",
    }
    ```

_If you are missing something else add an [issue](https://github.com/kiselev-nikolay/go-test-docker-dependencies/issues/new) or fork this repo (PRs are welcome)!_


## Usage

Use go get to get this package

```shell
go get github.com/kiselev-nikolay/go-test-docker-dependencies/testdep
```

Import package as always

```go
import (
	"github.com/kiselev-nikolay/go-test-docker-dependencies/testdep"
)

func TestXXX(t *testing.T) {
	pg := &testdep.Postgres{} // As example
}
```


## Examples


### Create declared database

```go
timeoutSeconds := 5

pg := &testdep.Postgres{
	Port:     5432,
	User:     "admin",
	Password: "Admin123",
	Database: "nicedb",
}

stop, err := pg.Run(timeoutSeconds)
defer stop()

fmt.Printf("PostgreSQL connection string: %s", dockerPg.ConnString())
```

### Use with gorm

```go
func MustCreatePg() (string, func()) {
	dbPort, _ := testdep.FindFreePort()
	dockerPg := testdep.Postgres{
		Port:     dbPort,
		User:     "test",
		Password: "test",
		Database: "test",
	}
	stop, _ := dockerPg.Run(10)
	return dockerPg.ConnString(), func() {
		stop()
	}
}

func TestWithGorm(t *testing.T) {
	dsn, stop := MustCreatePg()
	defer stop()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// Do Gorm tests there...
}
```