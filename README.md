# xpdf-wasm

## Xpdf & Wazero

This project uses the Xpdf C++ library by Glyph & Cog, LLC (https://www.xpdfreader.com/about.html) to process the PDF
documents.

We use a Webassembly version of Xpdf that is compiled with [Emscripten](https://emscripten.org/) and runs in the [Wazero Go](https://github.com/tetratelabs/wazero) runtime.

The compiled WebAssembly files are in the wasm folder and are embedded in the Go binary.

The instructions on how to compile the WebAssembly files can be found in the build folder.

## Usage
## Getting started

### From binary

Download the binary from the latest release for your platform and save it as `xpdf-wasm`.

You can also use the `install` tool for this:

```bash
sudo install xpdf-wasm-linux-x64 /usr/local/bin/xpdf-wasm
```

### From source

Make sure you have a working Go development environment.

Clone the repository:

```bash
git clone https://github.com/jerbob92/xpdf-wasm.git
```

Move into the directory:

```bash
cd xpdf-wasm
```

Run the command:

```bash
go run main.go [pdfdetach, pdffonts, pdfimages, pdfinfo, pdftohtml, pdftopng, pdftoppm, pdftops, pdftotext] [arguments]
```

Or to compile and run xpdf-wasm:

```bash
go build -o xpdf-wasm main.go
./xpdf-wasm [pdfdetach, pdffonts, pdfimages, pdfinfo, pdftohtml, pdftopng, pdftoppm, pdftops, pdftotext] [arguments]
```

Output:

`./xpdf-wasm pdfinfo`

```text
pdfinfo version 4.04 [www.xpdfreader.com]
Copyright 1996-2022 Glyph & Cog, LLC
Usage: pdfinfo [options] <PDF-file>
  -f <int>          : first page to convert
  -l <int>          : last page to convert
  -box              : print the page bounding boxes
  -meta             : print the document metadata (XML)
  -rawdates         : print the undecoded date strings directly from the PDF file
  -enc <string>     : output text encoding name
  -opw <string>     : owner password (for encrypted files)
  -upw <string>     : user password (for encrypted files)
  -cfg <string>     : configuration file to use in place of .xpdfrc
  -v                : print copyright and version info
  -h                : print usage information
  -help             : print usage information
  --help            : print usage information
  -?                : print usage information
```

## File paths

Because you can tell Wazero which folders have to be mounted in WebAssembly, you have full control over the filesystem.

By default, xpdf-wasm will mount the full root disk in Wazero on non-Windows environments.
On Windows environments, xpdf-wasm will get the volume of the current working directory and mount that as the root.

All paths given to xpdf-wasm in WebAssembly mode have to be in POSIX style and have to be absolute, so for
example: `/home/user/Downloads/file.pdf`. If you have mounted `/home/user/`on the root, then the path you would have to
give is `/Downloads/file.pdf`, this is the same on Windows, so no backward slashes or volume names in paths.
