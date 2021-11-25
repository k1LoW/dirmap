# dirmap

:file_folder: `dirmap` is a tool for generating a directory map. 

It extracts a part of the document from markdown or source code of each directory and uses it as overview of the directory.

``` console
$ dirmap generate
.
├── .github/
│   └── workflows/
├── cmd/ ... Commands.
├── config/ ... Configuration file.
├── matcher/ ... Implementation to find the string that will be the overview from the code or Markdown.
├── output/ ... Output format of the directory map.
├── scanner/ ... Implementation of scanning the target directory and its overview from the file system based on the configuration.
├── scripts/ ... scripts for Dockerfile.
└── version/ ... Version.
```
