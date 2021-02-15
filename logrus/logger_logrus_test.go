package logrus

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func TestInitLogger(t *testing.T) {
	t.Run("InitLogger() should return a default logger with two empty fields", func(t *testing.T) {
		bufferWriter := &bytes.Buffer{}
		InitLogger(bufferWriter)
		LLogger().Info("Test")
		fmt.Println(bufferWriter)

		logLine := bufferWriter.String()
		if !strings.Contains(logLine, "CID=") {
			t.Fatal("Not contain CID")
		}

		if !strings.Contains(logLine, "Usecase=") {
			t.Fatal("Not contain usecase")
		}
	})
}

func TestWithCid(t *testing.T) {
	t.Run("WithCid() should log the command id", func(t *testing.T) {
		bufferWriter := &bytes.Buffer{}
		InitLogger(bufferWriter)
		WithCid("123")
		LLogger().Info("Test")

		logLine := bufferWriter.String()
		if !strings.Contains(logLine, "CID=123") {
			t.Fatal("Not contain CID")
		}

		if !strings.Contains(logLine, "Usecase=") {
			t.Fatal("Not contain usecase")
		}
	})
}

func TestWithUsecase(t *testing.T) {
	t.Run("WithCid() should log the command id", func(t *testing.T) {
		bufferWriter := &bytes.Buffer{}
		InitLogger(bufferWriter)
		WithUsecase("test-usecase")
		LLogger().Info("Test")

		logLine := bufferWriter.String()
		if !strings.Contains(logLine, "CID=") {
			t.Fatal("Not contain CID")
		}

		if !strings.Contains(logLine, "Usecase=test-usecase") {
			t.Fatal("Not contain usecase")
		}
	})
}

func TestAttachStdout(t *testing.T) {
	t.Run("AttachStdout should attach another writer into stdout", func(t *testing.T) {

		// setup stdout to redirect to pipe
		oldStdout := os.Stdout
		rp, wp, _ := os.Pipe()
		os.Stdout = wp
		outC := make(chan string)

		bufferWriter := &bytes.Buffer{}
		InitLogger(bufferWriter)
		WithCid("123")
		WithUsecase("test-usecase")
		AttachStdout()

		LLogger().Info("Test Info")
		LLogger().Debug("Only in file")

		// now get content from pipe
		go func() {
			var buf bytes.Buffer
			io.Copy(&buf, rp)
			outC <- buf.String()
		}()
		wp.Close()
		os.Stdout = oldStdout
		stdoutLogLine := <-outC

		if !strings.Contains(stdoutLogLine, "CID=123") {
			t.Fatal("Not contain CID")
		}

		if !strings.Contains(stdoutLogLine, "Usecase=test-usecase") {
			t.Fatal("Not contain usecase")
		}

		if !strings.Contains(stdoutLogLine, "Test Info") {
			t.Fatal("Not contain log content")
		}

		fileLogLine := bufferWriter.String()
		fmt.Println(fileLogLine)
		if !strings.Contains(fileLogLine, "CID=123") {
			t.Fatal("Not contain CID")
		}

		if !strings.Contains(fileLogLine, "Usecase=test-usecase") {
			t.Fatal("Not contain usecase")
		}

		if !strings.Contains(fileLogLine, "Only in file") {
			t.Fatal("Not contain debug content")
		}
	})
}
