# rename-download-media

Add context command (windows) to rename media folder and contained file in a format Jellyfin expects: `Film Name (Year)`

The vast majority of public domain media you can acquire from the internet arrives in the format:

`Title.Year.Resolution.Bitrate.Encoding.Etc`, e.g. `Manos.the.Hands.of.Fate.1966.1080p.REPACK.WEB-DL.AAC2.0.x264`

This renaming tool uses the following formatting logic:

1. Search for a date string using a simple regex pattern match `\d{4}\.`, 4 digits followed by a period. Bail if not found.
2. Drop any characters after the date, wrap the date in parenthesis and replace periods with spaces.
3. Apply formatted name to the largest media file (".mp4", ".mkv" or ".avi") found within folder.
4. Apply formatted name to folder.

## To use:

1. Compile via `go build main.go`
   - If you have never used go, simply head to https://go.dev/doc/install and follow the instructions, then navigate to the code folder in file explorer, then right click the blank space and select `Open in terminal` - then type `go build main.go`
2. Edit `install.reg` to point to the path of the executable you built
   - Right-click executable, select `Copy as Path`, open `install.reg`, replace the file location on the last line
3. Double click `install.reg`
4. Profit (when right clicking folders you will now have the `Rename Downloaded Media` command available)

## Notes / Caveats:

1. On newer windows versions you _may_ have to click `Show more options` to find the `Rename Downloaded Media` command
2. This currently only works with folders/files using dot notation, e.g. `Manos.the.Hands.of.Fate.1966.1080p.REPACK.WEB-DL.AAC2.0.x264`
3. You can select multiple folders and run the command
4. This will drop periods from public domain films which contain periods, e.g. `D.O.A.`, and will require manual formatting
