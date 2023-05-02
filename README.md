# Magefile Extensions

![test](https://github.com/uwu-tools/magex/workflows/test/badge.svg)

This library provides helper methods to use with [mage](https://magefile.org).

Below is a sample of the type of helpers available. Full examples and
documentation is on [godoc](https://godoc.org/github.com/uwu-tools/magex).

```go
// +build mage

package main

import (
	"github.com/uwu-tools/magex/pkg"
	"github.com/uwu-tools/magex/shx"
)

// Check if packr2 is in the bin/ directory and is at least v2.
// If not, install packr@v2.8.0 into bin/
func EnsurePackr2() error {
	opts := pkg.EnsurePackageOptions{
		Name: "github.com/gobuffalo/packr/v2/packr2",
		DefaultVersion: "v2.8.0",
		VersionCommand: "version",
		Destination: "bin",
    }
   return pkg.EnsurePackageWith(opts)
}

// Install mage if it's not available, and ensure it's in PATH. We don't care which version
func Mage() error {
    return pkg.EnsureMage("")
}

// Run a docker registry in a container. Do not print stdout and only print
// stderr when the command fails even when -v is set.
//
// Useful for commands that you only care about when it fails, keeping unhelpful
// output out of your logs.
func StartRegistry() error {
    return shx.RunE("docker", "run", "-d", "-p", "5000:5000", "--name", "registry", "registry:2")
}

// Use go to download a tool, build and install it manually so 
// that it has version information embedded in the final binary.
func CustomInstallTool() error {
	err := shx.RunE("go", "get", "-u", "github.com/magefile/mage")
	if err != nil {
    		return err
	}
    
	src := filepath.Join(GOPATH(), "src/github.com/magefile/mage")
	return shx.Command("go", "run", "bootstrap.go").In(src).RunE()
}
```

## Attribution

This project is a fork of https://github.com/carolynvs/magex at
[0b6a1c6](https://github.com/carolynvs/magex/tree/0b6a1c6d5cba42cbad741c93546e99837b1c1fb9).

### In Memoriam

[Carolyn Van Slyck](https://github.com/carolynvs) was the original maintainer
of this project and is no longer with us.
She was a fan of [`mage`](https://magefile.org/) and I'd like to do my small
part by continuing her extensions project.

If you've had the opportunity to work with Carolyn and would like to leave a
message in her memory, please take a moment to do so on her CNCF
[memorial page](https://github.com/cncf/memorials/blob/main/carolyn-van-slyck.md).
