gounoconv
=====================================================================

[![Build Status](https://github.com/askasoft/gounoconv/actions/workflows/build.yml/badge.svg)](https://github.com/askasoft/gounoconv/actions?query=branch%3Amaster) 
[![codecov](https://codecov.io/gh/askasoft/gounoconv/branch/master/graph/badge.svg)](https://codecov.io/gh/askasoft/gounoconv) 
[![MIT](https://img.shields.io/badge/license-MIT-green)](https://opensource.org/licenses/MIT)


gounoconv is a golang XML-RPC client tool for unoserver https://github.com/unoconv/unoserver .


## Build

```sh
git clone https://github.com/askasoft/gounoconv.git

go build
```

## Usage

```sh
./gounoconv -endpoint http://127.0.0.1:2003  convert  hello.docx  hello.pdf

./gounoconv -endpoint http://127.0.0.1:2003  compare  old.docx  new.docx  out.pdf
```
