## Notes

* We're moving connection backups into the main package using Marshal / Unmarshal

## Current

* Refactor HTTP tracker

## Future

* Refactor HTTP tracker tests
* Refactor tracker main()
* Refactor controller / command line interface
* Cover all TODOs
* Implement new features

## Features

### Controller

* Look at caddy
* Arguments for providing config & pidfile
* Default paths for config and pid resepect XDG
* Detect and use systemd when possible
  * What did I mean by this lol

### Config

* Calculate logical defaults by default
  * Ie spawn 1-1.5x nproc worker threads by default
* Refactor setting names
* Determine how much memory to pre allocate for maps based on the last N runs maximum size
  * Create a routine that wakes up every minuite to check peak usage, write changes to a file that just has last N runs

### Backups

* Throw out pg backups
  * Make backup config setting just an optional path
