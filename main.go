package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/askasoft/gounoconv/unoclient"
	"github.com/askasoft/pango/asg"
	"github.com/askasoft/pango/fsu"
	"github.com/askasoft/pango/gog"
	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/str"
)

func usage() {
	help := `
Usage: %s <command> [options]
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
`
	fmt.Printf(help, filepath.Base(os.Args[0]))
}

func options(uo *unoclient.Option) []unoclient.OptionBuilder {
	return []unoclient.OptionBuilder{
		unoclient.WithLocal(uo.Local),
		unoclient.WithConvertTo(uo.ConvertTo),
		unoclient.WithInFilterName(uo.InFilterName),
		unoclient.WithFilterName(uo.FilterName),
		unoclient.WithFilterOptions(uo.FilterOptions...),
		unoclient.WithUpdateIndex(uo.UpdateIndex),
		unoclient.WithFileType(uo.FileType),
	}
}

func readFile(file string) []byte {
	if file == "-" {
		bb := &bytes.Buffer{}
		if _, err := iox.Copy(bb, os.Stdin); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return bb.Bytes()
	}

	data, err := fsu.ReadFile(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return data
}

func writeFile(file string, data []byte) {
	if file == "-" {
		if _, err := iox.Copy(os.Stdout, bytes.NewReader(data)); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	} else {
		if err := fsu.WriteFile(file, data, 0660); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

var quiet bool

func printf(format string, args ...any) {
	if !quiet {
		fmt.Printf(format, args...)
	}
}

func main() {
	var (
		debug     bool
		quiet     bool
		host      string
		port      int
		protocol  string
		location  string
		filteropt string
		uo        unoclient.Option
	)

	flag.BoolVar(&debug, "debug", false, "")
	flag.BoolVar(&quiet, "quiet", false, "")
	flag.StringVar(&host, "host", "localhost", "")
	flag.IntVar(&port, "port", 2003, "")
	flag.StringVar(&protocol, "protocol", "http", "")
	flag.StringVar(&location, "location", "auto", "")
	flag.StringVar(&location, "host-location", "auto", "")

	flag.StringVar(&uo.ConvertTo, "convert-to", "", "")
	flag.StringVar(&uo.InFilterName, "input-filter", "", "")
	flag.StringVar(&uo.FilterName, "output-filter", "", "")
	flag.StringVar(&uo.FilterName, "filter", "", "")
	flag.StringVar(&filteropt, "filter-options", "", "")
	flag.StringVar(&filteropt, "filter-option", "", "")
	flag.BoolVar(&uo.UpdateIndex, "update-index", false, "")

	flag.StringVar(&uo.FileType, "file-type", "", "")

	flag.CommandLine.Usage = usage
	flag.Parse()

	cw := log.NewConsoleWriter()
	cw.SetFormat("%t [%p] - %m%n%T")

	log.SetWriter(cw)
	log.SetLevel(gog.If(debug, log.LevelTrace, log.LevelInfo))

	switch location {
	case "local":
		uo.Local = true
	case "remote":
		uo.Local = false
	default:
		uo.Local = asg.Contains([]string{"127.0.0.1", "localhost"}, host)
	}
	uo.FilterOptions = str.FieldsByte(filteropt, ',')

	uc := unoclient.UnoClient{
		Endpoint: fmt.Sprintf("%s://%s:%d", protocol, host, port),
		Logger:   log.GetLogger("UNO"),
	}

	ops := options(&uo)

	arg := flag.Arg(0)
	switch arg {
	case "info":
		info, err := uc.Info(context.TODO())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(info.String())
	case "convert":
		inFile, outFile := flag.Arg(1), flag.Arg(2)
		if inFile == "" {
			fmt.Fprintf(os.Stderr, "Missing argument <infile> !\n")
			os.Exit(1)
		}
		if outFile == "" {
			fmt.Fprintf(os.Stderr, "Missing argument <outfile> !\n")
			os.Exit(1)
		}
		if outFile == "-" && uo.ConvertTo == "" {
			fmt.Fprintf(os.Stderr, "Missing -convert-to option !\n")
			os.Exit(1)
		}

		printf("Convert %s --> %s ... ", inFile, outFile)
		if inFile == "-" || outFile == "-" {
			inData := readFile(inFile)

			outData, err := uc.Convert(context.TODO(), inData, ops...)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			writeFile(outFile, outData)
		} else {
			err := uc.ConvertFile(context.TODO(), inFile, outFile, ops...)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}
		printf("OK.\n")
	case "compare":
		oldFile, newFile, outFile := flag.Arg(1), flag.Arg(2), flag.Arg(3)
		if oldFile == "" {
			fmt.Fprintf(os.Stderr, "Missing argument <oldfile> !\n")
			os.Exit(1)
		}
		if newFile == "" {
			fmt.Fprintf(os.Stderr, "Missing argument <newfile> !\n")
			os.Exit(1)
		}
		if outFile == "" {
			fmt.Fprintf(os.Stderr, "Missing argument <outfile> !\n")
			os.Exit(1)
		}

		printf("Compare '%s' '%s' --> %s ... ", oldFile, newFile, outFile)
		if oldFile == "-" || newFile == "-" || outFile == "-" {
			oldData := readFile(oldFile)
			newData := readFile(newFile)

			outData, err := uc.Compare(context.TODO(), oldData, newData, ops...)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			writeFile(outFile, outData)
		} else {
			err := uc.CompareFile(context.TODO(), oldFile, newFile, outFile, ops...)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}
		printf("OK.\n")
	default:
		fmt.Fprintf(os.Stderr, "Invalid command %q\n\n", arg)
		usage()
		os.Exit(1)
	}
}
