# robot-go

`robot-go` automates the Vers Libre audio pipeline. It pulls source files from Google Drive, enriches them with metadata, uploads finished shows to external platforms, archives them, and tracks progress in SQLite.

The repository currently builds two Go binaries:

- `move`: pulls files from the configured Google Drive folders into local staging storage and moves the originals to the sent folder.
- `upload`: reads staged audio files, loads metadata from SQLite and Google Sheets, uploads to Mixcloud, SoundCloud, Radiocult, and Google Drive, then moves completed files into a local archive.

## Repository layout

- `packages/move`: Google Drive ingest step
- `packages/upload`: metadata, upload, and archive step
- `modules/verslibre.nix`: NixOS service and timer module
- `init-db.sql`: base SQLite schema
- `init-db.sh`: helper to initialize the SQLite database

## Build and development

### With Nix

Build the binaries from the flake:

```bash
nix build .#vl-move
nix build .#vl-upload
```

Enter the development shell:

```bash
nix develop
```

### With Go

The Go module lives in `packages/`:

```bash
cd packages
go test ./...
go build ./move
go build ./upload
```

## Runtime layout

Both binaries default to `/var/lib/robot` as their working directory.

The uploader expects these subdirectories under the base path:

- `to_upload`: staged audio files waiting to be uploaded
- `picture`: cover art downloaded from Google Drive

The archive path is separate and must be passed explicitly to `upload`.

## Google credentials

Both binaries require a Google service account credentials JSON file.

Default path:

```text
/etc/robot/cred.json
```

Override it with `--credentials`.

## `move`

`move` stages files locally for the uploader.

```bash
move --local /var/lib/robot --credentials /etc/robot/cred.json
```

Supported flags:

- `--local`: base path for local staging, default `/var/lib/robot`
- `--credentials`: Google credentials JSON, default `/etc/robot/cred.json`

Current behavior:

1. Copies files from the `auphonic_macmini` Drive folder into `auphonic_preprocess`
2. Moves those source files into the `sent` Drive folder
3. Downloads files from the `auphonic` Drive folder into `<local>/to_upload`
4. Moves those source files into the `sent` Drive folder
5. Downloads files from the `upload_source` Drive folder into `<local>/to_upload`
6. Moves those source files into the `sent` Drive folder

The Drive folder IDs are currently compiled into the binary.

## `upload`

`upload` processes files in `<local>/to_upload` and uploads each show to the configured destinations.

```bash
upload \
  --local /var/lib/robot \
  --archive /var/lib/robot/archive \
  --credentials /etc/robot/cred.json \
  --metadata /var/lib/robot/metadata.sql \
  --soundcloud-token /var/lib/robot/soundcloud-token.json
```

Required flags:

- `--archive`: local archive directory for completed uploads

Optional flags:

- `--local`: base path, default `/var/lib/robot`
- `--credentials`: Google credentials JSON, default `/etc/robot/cred.json`
- `--metadata`: SQLite database path, default `/var/lib/robot/metadata.sql`
- `--soundcloud-token`: persisted SoundCloud token JSON, default `/var/lib/robot/soundcloud-token.json`
- `--soundcloud-init-auth`: print a SoundCloud authorization URL and exit
- `--soundcloud-auth-code <code>`: exchange an authorization code for a persisted token and exit

### Input filename format

The uploader derives metadata from the audio filename. It expects underscore-separated names where:

- field 1 is the date in `YYYYMMDD` format
- field 3 is the show tag, before the file extension

Example:

```text
20260318_show_test-show.mp3
```

For that file, the uploader uses:

- date: `20260318`
- tag: `test-show`

### Upload flow

For each staged audio file, the uploader:

1. Loads base metadata from `base_data` in SQLite
2. Applies per-date overrides from the `meta` worksheet in the configured Google Sheet
3. Inserts or reuses a row in the `metadata` table
4. Downloads the matching cover image from Google Drive
5. Writes ID3 tags into the MP3
6. Uploads to Mixcloud
7. Uploads to SoundCloud and stores the returned track URN
8. Uploads live shows to Radiocult
9. Uploads the MP3 to the Google Drive archive
10. Moves fully uploaded files into the local archive

Uploads are tracked per destination in SQLite so repeated runs can continue where a previous run left off.

## Database

Initialize the SQLite database:

```bash
./init-db.sh /var/lib/robot/metadata.sql
```

If no path is provided, `init-db.sh` defaults to `/var/lib/robot/metadata.sql`.

The schema creates two tables:

- `base_data`: default metadata per show tag
- `metadata`: per-episode upload state

The uploader also performs a lightweight schema migration on startup to ensure the `soundcloud` and `soundcloud_urn` columns exist.

### Bootstrapping `base_data`

You can import initial show metadata from CSV:

```bash
sqlite3 /var/lib/robot/metadata.sql
.mode csv
.import /path/to/base_data.csv base_data
```

## Environment variables

The uploader reads credentials for external services from environment variables.

### Mixcloud

```bash
API_KEY=<mixcloud access token>
```

### Radiocult

```bash
RADIOCULT_API=<radiocult api key>
STATION_ID=<radiocult station id>
```

### SoundCloud

```bash
SOUNDCLOUD_CLIENT_ID=<app client id>
SOUNDCLOUD_CLIENT_SECRET=<app client secret>
SOUNDCLOUD_REDIRECT_URI=<registered redirect URI>
```

## SoundCloud authorization

The uploader persists a refreshable SoundCloud token in JSON.

Initialize auth:

```bash
upload --soundcloud-init-auth --soundcloud-token /var/lib/robot/soundcloud-token.json
```

Open the printed URL, authorize the app, then take the returned `code` query parameter and exchange it:

```bash
upload --soundcloud-auth-code <code> --soundcloud-token /var/lib/robot/soundcloud-token.json
```

On normal upload runs, the token is refreshed automatically before uploads start.

## NixOS module

The flake exposes a NixOS module at `nixosModules.verslibre`.

It defines:

- a `verslibre.service` systemd service
- a `verslibre.timer` timer, default `OnCalendar = "Hourly"`

Important module options:

- `services.verslibre.basePath`
- `services.verslibre.credPath`
- `services.verslibre.archivePath`
- `services.verslibre.dbPath`
- `services.verslibre.soundcloudTokenPath`
- `services.verslibre.timerConfig`

The service reads environment variables from:

```text
/etc/vl-upload.env
```

and runs `move` first, then `upload`.

## Tests

Current automated tests cover metadata schema migration and upload-state bookkeeping:

```bash
cd packages
go test ./...
```
