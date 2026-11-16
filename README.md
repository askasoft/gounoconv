gounoconv
=====================================================================

[![Build Status](https://github.com/askasoft/gounoconv/actions/workflows/build.yml/badge.svg)](https://github.com/askasoft/gounoconv/actions?query=branch%3Amaster) 
[![codecov](https://codecov.io/gh/askasoft/gounoconv/branch/master/graph/badge.svg)](https://codecov.io/gh/askasoft/gounoconv) 
[![MIT](https://img.shields.io/badge/license-MIT-green)](https://opensource.org/licenses/MIT)


gounoconv is a golang XML-RPC client tool for unoserver https://github.com/unoconv/unoserver .


### Prepare unoserver
```sh
sudo apt install -y python3-pip fonts-noto-cjk libreoffice
#sudo pip3 install unoserver
sudo apt install -y pipx libreoffice python3-uno

sudo useradd -m -s /bin/bash unoserver
sudo su - unoserver
pipx ensurepath
source ~/.bashrc
pipx install --system-site-packages unoserver
which unoserver

echo '
[Unit]
Description = Unoserver for document conversions
Documentation = https://github.com/unoconv/unoserver
After = syslog.target network.target local-fs.target

[Service]
Type = simple
ExecStart = /home/unoserver/.local/bin/unoserver --interface 0.0.0.0 --uno-interface 127.0.0.1
TimeoutStopSec = 5
User = unoserver
Group = unoserver

[Install]
WantedBy = multi-user.target
' | sudo tee /usr/lib/systemd/system/unoserver.service

sudo systemctl daemon-reload
sudo systemctl enable unoserver
sudo systemctl start unoserver
```

## Install

```sh
go install github.com/askasoft/gounoconv@latest
```

## Usage
```
Usage: unoclient <command> [options]
  <command>:
    info        Print unoserver information.
    convert <infile> <outfile>
                Convert document <infile> to <outfile>.
      [options]:
        -convert-to CONVERT_TO
                The file type/extension of the output file (ex: pdf).
        -input-filter INPUT_FILTER
                The LibreOffice input filter to use (ex: writer8).
        -output-filter OUTPUT_FILTER, -filter OUTPUT_FILTER
                The export filter to use when converting.
        -filter-options FILTER_OPTIONS, -filter-option FILTER_OPTIONS
                Pass options for the export filter, in name[=value] format.
                Comma separated list for multiple options.
        -update-index
                Updates the indexes before conversion. Can be time consuming.
    compare <oldfile> <newfile> <outfile>
                Compare documents <oldfile> <outfile> to <outfile>.
      [options]:
        -file-type FILE_TYPE
                The file type/extension of the result file (ex: pdf).
  <general options>:
    -h | -help  Print this help message.
    -host HOST  The host the server runs on.
    -port PORT  The port used by the server.
    -protocol {http,https}
                The protocol used by the server.
    -location {auto,remote,local}, -host-location {auto,remote,local}
                The host location determines the handling of files.
                If you run the client on the same machine as the server,
                it can be set to local, and the files are sent as paths.
                If they are different machines, it is remote and the files
                are sent as binary data. Default is auto, and it will send
                the file as a path if the host is 127.0.0.1 or localhost.
    -debug      Print the debug log.
    -quiet      Do not print information message.
  <notes>:
    * Use - for stdin or stdout.
```

## Example
```sh
./unoclient -host 127.0.0.1  info

./unoclient -host 127.0.0.1  convert  hello.docx  hello.pdf

./unoclient -host 127.0.0.1  compare  old.docx  new.docx  out.pdf
```
