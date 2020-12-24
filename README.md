## API
### use package
```go
import (
	"github.com/strapsi/go-docker"
)
```
### types
```go
type PsOptions struct {
	All bool
	FilterNames []string
}

type RunOptions struct {
	Image string
	Name string
	Force bool
	Env map[string]string
}
```
### docker
```go
// returns existing containers
docker.Ps(options *docker.PsOptions)

// create and start container
docker.Run(options *docker.RunOptions)
```

## Changelog
### v0.0.4
* added environment variables for docker run
* added force option to force an existing container to be removed
* added filter by name option to docker ps

### v0.0.3
* added docker run command

### v0.0.2
* added docker ps command
* removed test commands

### v0.0.1
* initialized package with test commands
