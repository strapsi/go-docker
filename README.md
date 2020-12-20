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
}

type RunOptions struct {
	Image string
	Name string
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
### v0.0.2
* added docker ps command
* added docker run command
* removed test commands

### v0.0.1
* initialized package with test commands
