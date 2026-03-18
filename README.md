# mixcloud-go

## Move
The package is run using this command:
```
/path/to/binary -local <path/to/temp/local/folder> -credentials <path/to/cred.json>
```

The move pacakge is used to move files between google drive folders so that they are picked up and handled correctly further down the line. There are basically 3 transfers that take place:

### Auphonic preprocessing
1.1 Copy files from `3. Auphonic` on `VL Studio MacMini` to the `Auphonic preprocess` folder
1.2 Move files from `3. Auphonic` on `VL Studio MacMini` to the `1. Sent to mastering ` folder

### Auphonic postprocessing
1.1 Download files from `Auphonic postprocess` to the server
1.2 Move files from `3. Auphonic postprocess` to the `1. Sent to mastering` folder

### Standard upload 
1.1 Download files from `4. Upload folder` to the server
1.2 Move files from `4. Upload folder` to the `1. Sent to mastering` folder

## Metadata table
We use a SQLite3 DB for our metadata. This table can be initialized using the `init-db.sh` bash script. You can pass a path to the DB as argument or it defaults to `/var/lib/robot/metadata.db`.

The base_data table can be populated by csv (useful on very first deploy). This done by running the following set of commands:
```
sqlite3
.open <database_path>
.mode csv
.import <csv_path> <table_name>
```

## SoundCloud
SoundCloud uploads are part of the `upload` pipeline and are tracked in the same SQLite metadata table via the `soundcloud` and `soundcloud_urn` columns.

The uploader requires the following environment variables:
```
SOUNDCLOUD_CLIENT_ID=<app client id>
SOUNDCLOUD_CLIENT_SECRET=<app client secret>
SOUNDCLOUD_REDIRECT_URI=<registered redirect URI>
```

The program persists the SoundCloud token set in a separate JSON file, defaulting to `/var/lib/robot/soundcloud-token.json`. The token is refreshed on every run before uploads begin.

Initialize the token file once:
```
upload --soundcloud-init-auth --soundcloud-token /var/lib/robot/soundcloud-token.json
```

Open the printed URL in a browser, authorize the app, copy the returned `code` query parameter, then exchange it:
```
upload --soundcloud-auth-code <code> --soundcloud-token /var/lib/robot/soundcloud-token.json
```
