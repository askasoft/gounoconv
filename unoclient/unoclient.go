package unoclient

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/askasoft/pango/fsu"
	"github.com/askasoft/pango/gog"
	"github.com/askasoft/pango/mag"
	"github.com/askasoft/pango/net/xmlrpc"
)

type Option struct {
	Local bool // is the unoserver runs on local machine

	ConvertTo     string   // convert option: The file type/extension of the output file (ex: pdf)
	FilterName    string   // convert option: The export filter to use when converting
	FilterOptions []string // convert option: The Export filter options, in name=value format
	UpdateIndex   bool     // convert option: Updates the indexes before conversion. Can be time consuming.
	InFilterName  string   // convert option: The LibreOffice input filter to use (ex: 'writer8').

	FileType string // compare option: The file type/extension of the result file (ex: pdf).
}

func (o *Option) param(v string) any {
	return gog.If(v == "", nil, any(v))
}

func (o *Option) convert_to() any {
	return o.param(o.ConvertTo)
}

func (o *Option) filter_name() any {
	return o.param(o.FilterName)
}

func (o *Option) infilter_name() any {
	return o.param(o.InFilterName)
}

func (o *Option) file_type() any {
	return o.param(o.FileType)
}

type OptionBuilder func(*Option)

func WithLocal(local bool) OptionBuilder {
	return func(co *Option) {
		co.Local = local
	}
}

func WithConvertTo(to string) OptionBuilder {
	return func(co *Option) {
		co.ConvertTo = to
	}
}

func WithFilterName(fn string) OptionBuilder {
	return func(co *Option) {
		co.FilterName = fn
	}
}

func WithFilterOptions(ops ...string) OptionBuilder {
	return func(co *Option) {
		co.FilterOptions = append(co.FilterOptions, ops...)
	}
}

func WithUpdateIndex(ui bool) OptionBuilder {
	return func(co *Option) {
		co.UpdateIndex = ui
	}
}

func WithInFilterName(fn string) OptionBuilder {
	return func(co *Option) {
		co.InFilterName = fn
	}
}

func WithFileType(ft string) OptionBuilder {
	return func(co *Option) {
		co.FileType = ft
	}
}

func buildOption(obs ...OptionBuilder) (op Option) {
	for _, ob := range obs {
		ob(&op)
	}
	return
}

type UnoClient xmlrpc.Client

func (uc *UnoClient) call(ctx context.Context, method string, result any, params ...any) error {
	return (*xmlrpc.Client)(uc).Call(ctx, method, result, params...)
}

type UnoInfo struct {
	API           string            `xmlrpc:"api"`
	UnoServer     string            `xmlrpc:"unoserver"`
	ImportFilters map[string]string `xmlrpc:"import_filters"`
	ExportFilters map[string]string `xmlrpc:"export_filters"`
}

func (ui *UnoInfo) String() string {
	sb := &strings.Builder{}

	fmt.Fprintf(sb, "api           : %s\n", ui.API)
	fmt.Fprintf(sb, "unoserver     : %s\n", ui.UnoServer)

	fmt.Fprint(sb, "import_filters:\n")
	ifs := mag.Keys(ui.ImportFilters)
	sort.Strings(ifs)
	for _, f := range ifs {
		fmt.Fprintf(sb, "    %-34s: %s\n", f, ui.ImportFilters[f])
	}

	fmt.Fprint(sb, "export_filters:\n")
	efs := mag.Keys(ui.ExportFilters)
	sort.Strings(efs)
	for _, f := range efs {
		fmt.Fprintf(sb, "    %-34s: %s\n", f, ui.ExportFilters[f])
	}

	return sb.String()
}

func (uc *UnoClient) Info(ctx context.Context) (info UnoInfo, err error) {
	err = uc.call(ctx, "info", &info)
	return
}

func (uc *UnoClient) Convert(ctx context.Context, inData []byte, obs ...OptionBuilder) (outData []byte, err error) {
	op := buildOption(obs...)

	err = uc.call(ctx, "convert", &outData, nil, inData, nil,
		op.convert_to(), op.filter_name(), op.FilterOptions, op.UpdateIndex, op.infilter_name())

	return
}

func (uc *UnoClient) ConvertFile(ctx context.Context, inFile, outFile string, obs ...OptionBuilder) (err error) {
	var (
		inPath  any
		outPath any
		inData  []byte
		outData []byte
	)

	op := buildOption(obs...)

	if op.Local {
		inPath, outPath = inFile, outFile
	} else {
		inData, err = fsu.ReadFile(inFile)
		if err != nil {
			return err
		}
	}

	if op.ConvertTo == "" {
		op.ConvertTo = strings.TrimLeft(filepath.Ext(outFile), ".")
	}

	err = uc.call(ctx, "convert", &outData, inPath, inData, outPath,
		op.convert_to(), op.filter_name(), op.FilterOptions, op.UpdateIndex, op.infilter_name())

	if err == nil && !op.Local {
		err = fsu.WriteFile(outFile, outData, 0660)
	}
	return
}

func (uc *UnoClient) Compare(ctx context.Context, oldData, newData []byte, obs ...OptionBuilder) (outData []byte, err error) {
	op := buildOption(obs...)

	err = uc.call(ctx, "compare", &outData, nil, oldData, nil, newData, nil, op.file_type())

	return
}

func (uc *UnoClient) CompareFile(ctx context.Context, oldFile, newFile, outFile string, obs ...OptionBuilder) (err error) {
	var (
		oldPath any
		newPath any
		outPath any
		oldData []byte
		newData []byte
		outData []byte
	)

	op := buildOption(obs...)

	if op.Local {
		oldPath, newPath, outPath = oldFile, newFile, outFile
	} else {
		oldData, err = fsu.ReadFile(oldFile)
		if err != nil {
			return err
		}

		newData, err = fsu.ReadFile(newFile)
		if err != nil {
			return err
		}
	}

	if op.FileType == "" {
		op.FileType = strings.TrimLeft(filepath.Ext(outFile), ".")
	}

	err = uc.call(ctx, "compare", &outData, oldPath, oldData, newPath, newData, outPath, op.file_type())

	if err == nil && !op.Local {
		err = fsu.WriteFile(outFile, outData, 0660)
	}
	return
}
