package unoclient

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/askasoft/pango/fsu"
	"github.com/askasoft/pango/log"
)

func testClient(t *testing.T) *UnoClient {
	endpoint := os.Getenv("UNOSERVER")

	if endpoint == "" {
		t.Skip("UNOSERVER not set")
	}

	logger := log.GetLogger("UNO")

	return &UnoClient{Endpoint: endpoint, Logger: logger}
}

func TestInfo(t *testing.T) {
	client := testClient(t)

	info, err := client.Info(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(info.String())
}

func TestConvertRemote(t *testing.T) {
	client := testClient(t)

	inFile, outFile := "testdata/hello.txt", "testdata/hello.pdf"

	defer os.Remove(outFile)

	if err := client.ConvertFile(context.TODO(), inFile, outFile); err != nil {
		t.Fatal(err)
	}

	if err := fsu.FileExists(outFile); err != nil {
		t.Fatal(err)
	}
}

func TestConvertLocal(t *testing.T) {
	client := testClient(t)

	inFile, outFile := "/tmp/hello.txt", "/tmp/hello.pdf"

	if err := client.ConvertFile(context.TODO(), inFile, outFile, WithLocal(true)); err != nil {
		t.Fatal(err)
	}
}

func TestCompareRemote(t *testing.T) {
	client := testClient(t)

	oldFile, newFile, outFile := "testdata/hello.txt", "testdata/greet.txt", "testdata/compare.txt"

	defer os.Remove(outFile)

	if err := client.CompareFile(context.TODO(), oldFile, newFile, outFile); err != nil {
		t.Fatal(err)
	}

	if err := fsu.FileExists(outFile); err != nil {
		t.Fatal(err)
	}
}

func TestCompareLocal(t *testing.T) {
	client := testClient(t)

	oldFile, newFile, outFile := "/tmp/hello.txt", "/tmp/hello2.txt", "/tmp/compare.txt"

	if err := client.CompareFile(context.TODO(), oldFile, newFile, outFile, WithLocal(true)); err != nil {
		t.Fatal(err)
	}
}
