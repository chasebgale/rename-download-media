# rename-download-media

Add context command (windows) to format media folder and contained file

## To use:

1. Compile via `go build main.go`
2. Edit `install.reg` to point to the path of the executable you built
3. Double click `install.reg`
4. Profit (when right clicking folders you will now have the 'Rename Downloaded Media' command available)

## Notes:

1. On newer windows versions you _may_ have to click 'Show more options' to find the 'Rename Downloaded Media' command
2. This currently only works with folders/files using dot notation, e.g. `Night.of.the.Living.Dead.1968.1080p.REPACK.WEB-DL.AAC2.0.x264`
