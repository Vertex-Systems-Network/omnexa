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
	terminalErr    error
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
	if reader.terminalErr != nil {
		err := reader.terminalErr
		reader.terminalErr = nil
		return 0, err
	}
	if reader.verified {
		return 0, io.EOF
	}

	n, err := reader.reader.Read(buffer)
	if n > 0 {
		reader.read += int64(n)
		if reader.read > reader.expectedLength {
			return n, safeFailure(
				codeLengthInvalid,
				failure.CategoryValidation,
				"storage object length is invalid",
				failure.WithDetail("upload body exceeds the declared content length"),
			)
		}
		_, _ = reader.hasher.Write(buffer[:n])
		if reader.read == reader.expectedLength {
			reader.verified = true
			if actual := hex.EncodeToString(reader.hasher.Sum(nil)); actual != reader.expectedHash {
				integrityErr := safeFailure(
					codeIntegrityFailed,
					failure.CategoryValidation,
					"storage object integrity check failed",
				)
				if n > 0 {
					reader.terminalErr = integrityErr
					return n, nil
				}
				return 0, integrityErr
			}
		}
	}

	if err == io.EOF {
		if reader.read != reader.expectedLength {
			return n, safeFailure(
				codeLengthInvalid,
				failure.CategoryValidation,
				"storage object length is invalid",
				failure.WithDetail("upload body ended before the declared content length"),
			)
		}
		if reader.verified {
			return n, io.EOF
		}
	}
	return n, err
}

type verifiedReadCloser struct {
	body           io.ReadCloser
	hasher         hash.Hash
	expectedHash   string
	expectedLength int64
	read           int64
	terminalErr    error
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
	if reader.terminalErr != nil {
		err := reader.terminalErr
		reader.terminalErr = nil
		return 0, err
	}
	if reader.verified {
		return 0, io.EOF
	}

	n, err := reader.body.Read(buffer)
	if n > 0 {
		reader.read += int64(n)
		if reader.read > reader.expectedLength {
			return n, reader.integrityError()
		}
		_, _ = reader.hasher.Write(buffer[:n])
	}

	if err == io.EOF {
		verifyErr := reader.verify()
		if verifyErr != nil {
			if n > 0 {
				reader.terminalErr = verifyErr
				return n, nil
			}
			return 0, verifyErr
		}
		reader.verified = true
	}
	return n, err
}

func (reader *verifiedReadCloser) Close() error {
	return reader.body.Close()
}

func (reader *verifiedReadCloser) verify() error {
	if reader.read != reader.expectedLength {
		return reader.integrityError()
	}
	actual := hex.EncodeToString(reader.hasher.Sum(nil))
	if actual != reader.expectedHash {
		return reader.integrityError()
	}
	return nil
}

func (reader *verifiedReadCloser) integrityError() error {
	return safeFailure(
		codeIntegrityFailed,
		failure.CategoryInvariant,
		"storage object integrity check failed",
	)
}
