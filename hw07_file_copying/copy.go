package main

import (
	"errors"
	"io"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3" //nolint
)

const blocksize int64 = 128 // Could be 4Kb, just like in most OSes, but this one to make copy slower

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func doCopy(inFp, outFp *os.File, limit int64) (err error) {
	var copied int64
	buffer := make([]byte, min(blocksize, limit))

	bar := pb.StartNew(int(limit))

	for copied < limit {
		toRead := min(blocksize, limit-copied)
		if toRead < int64(len(buffer)) {
			buffer = buffer[0:toRead]
		}
		nread, err := inFp.Read(buffer)
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return err
		}

		if nread < len(buffer) {
			buffer = buffer[0:nread]
		}

		nwritten, err := outFp.Write(buffer)
		if err != nil {
			return err
		}

		if nwritten != nread {
			return errors.New("read and write bytes mismatch")
		}

		// Make a pause to simulate slow copying device, so our progress bar will be visible
		time.Sleep(50 * time.Millisecond)

		copied += int64(nread)
		bar.SetCurrent(copied)
	}
	bar.Finish()
	return nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	fromFp, err := os.Open(fromPath)
	if err != nil {
		return err
	}

	defer fromFp.Close()

	if limit == 0 {
		limit, err = fromFp.Seek(0, io.SeekEnd)
		if err != nil {
			return err
		}
	}

	realOffset, err := fromFp.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	if realOffset != offset {
		return errors.New("invalid offset")
	}

	outFp, err := os.Create(toPath)
	if err != nil {
		return err
	}

	defer outFp.Close()

	return doCopy(fromFp, outFp, limit)
}
