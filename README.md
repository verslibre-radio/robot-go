# mixcloud-go

## Move
The package is run using this command:
```
/path/to/binary <path/to/temp/local/folder>
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


