# MCTravelerPresence
The MCTraveler Minecraft server's rich presence client written in Golang.

The program is of course entirely safe and very efficient, using only:
 - Around 3MBs of memory
 - Effectively 0% CPU usage on even the slowest of CPUs
 - Effectively 0mbps network
 - 0% disk activity

## Supported Platforms
- Windows

macOS and Linux support may be added in the future.

## Installing
Download and run the latest executable from https://github.com/Blazzike/MCTravelerPresence/releases

### Uninstalling
The rich presence can be disabled from starting using the "Startup" tab in Task Manager. It can be uninstalled by deleting the executable from `shell:startup`

## Contributing
1. Fork this repository
2. Clone the newly create repository down
3. `cd` into the clone directory and run `go get -t`
4. Make your changes!
5. Push the changes to your fork.
6. Open a pull request!

## Building
Run `go build -o build\MCTravelerPresence.exe -ldflags -H=windowsgui .`

The built executable will be in the build directory.