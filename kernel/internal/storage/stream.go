package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

type integrityReader struct {
	reader         io.Reader
	hasher         hash.Hash
	expectedHash   string
	expectedLength int64
	read           int64
	verified       bool
}

func newIntegrityReader(reader io.Reader, expectedLength int64, expectedHash string) io.Reader {
	return &integrityReader{
		reader:         reader,
		hasher:         sha256.New(),
		expectedHash:   expectedHash,
		expectedLength: expectedLength,
	}
}

func (reader *integrityReader) Read(buffer []byte) (int, error) {
	if reader.verified {
		return 0, io.EOF
	}
	if len(buffer) == 0 {
		return 0, nil
	}

	remaining := reader.expectedLength - reader.read
	if remaining < 0 {
		return 0, uploadLengthFailure("upload body exceeds the declared content length")
	}
	if remaining == 0 {
		if err := reader.finish(); err != nil {
			return 0, err
		}
		return 0, io.EOF
	}

	limit := len(buffer)
	if int64(limit) > remaining {
		limit = int(remaining)
	}
	n, err := reader.reader.Read(buffer[:limit])
	if n > 0 {
		reader.read += int64(n)
		_, _ = reader.hasher.Write(buffer[:n])
	}

	if reader.read == reader.expectedLength {
		finishErr := reader.finish()
		if finishErr != nil {
			return n, finishErr
		}
		return n, io.EOF
	}

	if err == io.EOF {
		return n, uploadLengthFailure("upload body ended before the declared content length")
	}
	return n, err
}

func (reader *integrityReader) finish() error {
	var probe [1]byte
	n, err := reader.reader.Read(probe[:])
	if n > 0 {
		return uploadLengthFailure("upload body exceeds the declared content length")
	}
	if err == nil {
		return uploadLengthFailure("upload body did not terminate at the declared content length")
	}
	if err != io.EOF {
		return err
	}

	actual := hex.EncodeToString(reader.hasher.Sum(nil))
	if actual != reader.expectedHash {
		return safeFailure(
			codeIntegrityFailed,
			failure.CategoryValidation,
			"storage object integrity check failed",
		)
	}
	reader.verified = true
	return nil
}

func uploadLengthFailure(detail string) error {
	return safeFailure(
		codeLengthInvalid,
		failure.CategoryValidation,
		"storage object length is invalid",
		failure.WithDetail(detail),
	)
}

type verifiedReadCloser struct {
	body           io.ReadCloser
	hasher         hash.Hash
	expectedHash   string
	expectedLength int64
	read           int64
	verified       bool
}

func newVerifiedReadCloser(body io.ReadCloser, expectedLength int64, expectedHash string) io.ReadCloser {
	return &verifiedReadCloser{
		body:           body,
		hasher:         sha256.New(),
		expectedHash:   expectedHash,
		expectedLength: expectedLength,
	}
}

func (reader *verifiedReadCloser) Read(buffer []byte) (int, error) {
	if reader.verified {
		return 0, io.EOF
	}
	if len(buffer) == 0 {
		return 0, nil
	}

	remaining := reader.expectedLength - reader.read
	if remaining < 0 {
		return 0, reader.integrityError()
	}
	if remaining == 0 {
		if err := reader.finish(); err != nil {
			return 0, err
		}
		return 0, io.EOF
	}

	limit := len(buffer)
	if int64(limit) > remaining {
		limit = int(remaining)
	}
	n, err := reader.body.Read(buffer[:limit])
	if n > 0 {
		reader.read += int64(n)
		_, _ = reader.hasher.Write(buffer[:n])
	}

	if reader.read == reader.expectedLength {
		finishErr := reader.finish()
		if finishErr != nil {
			return n, finishErr
		}
		return n, io.EOF
	}

	if err == io.EOF {
		return n, reader.integrityError()
	}
	return n, err
}

func (reader *verifiedReadCloser) Close() error {
	return reader.body.Close()
}

func (reader *verifiedReadCloser) finish() error {
	var probe [1]byte
	n, err := reader.body.Read(probe[:])
	if n > 0 || err == nil {
		return reader.integrityError()
	}
	if err != io.EOF {
		return err
	}
	if reader.read != reader.expectedLength {
		return reader.integrityError()
	}
	actual := hex.EncodeToString(reader.hasher.Sum(nil))
	if actual != reader.expectedHash {
		return reader.integrityError()
	}
	reader.verified = true
	return nil
}

func (reader *verifiedReadCloser) integrityError() error {
	return safeFailure(
		codeIntegrityFailed,
		failure.CategoryInvariant,
		"storage object integrity check failed",
	)
}
