package portfolio

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/plancatalog"
)

type remoteBatchReader struct {
	repository *GitRepository
	stream     *BoundCommandStream
	stdout     *bufio.Reader
}

func newRemoteBatchReader(ctx context.Context, repository *GitRepository) (*remoteBatchReader, error) {
	runner, ok := repository.runner.(BoundStreamingRunner)
	if !ok {
		return nil, nil
	}
	if err := repository.verifyAuthority(); err != nil {
		return nil, err
	}
	stream, err := runner.StartBound(ctx, repository.mirrorFile, repository.Path, "git", safeGitArgs("cat-file", "--batch")...)
	if err != nil {
		return nil, err
	}
	return &remoteBatchReader{repository: repository, stream: stream, stdout: bufio.NewReader(stream.Stdout)}, nil
}

func (reader *remoteBatchReader) Read(blob plancatalog.GitBlob) ([]byte, error) {
	if blob.Mode != plancatalog.GitModeRegular && blob.Mode != plancatalog.GitModeExecutable {
		return nil, fmt.Errorf("remote_blob_mode_unreadable: %s has mode %s", blob.Path, blob.Mode)
	}
	if _, err := io.WriteString(reader.stream.Stdin, blob.ObjectID+"\n"); err != nil {
		return nil, err
	}
	header, err := reader.stdout.ReadString('\n')
	if err != nil {
		return nil, err
	}
	fields := strings.Fields(header)
	if len(fields) != 3 || fields[0] != blob.ObjectID || fields[1] != "blob" {
		return nil, fmt.Errorf("remote_blob_unavailable: %s", strings.TrimSpace(header))
	}
	size, err := strconv.ParseInt(fields[2], 10, 64)
	if err != nil || size < 0 {
		return nil, fmt.Errorf("remote_blob_invalid_size: %s", blob.Path)
	}
	if size > plancatalog.MaxPlanSourceBytes {
		if _, discardErr := io.CopyN(io.Discard, reader.stdout, size+1); discardErr != nil {
			return nil, discardErr
		}
		return nil, &plancatalog.CodedError{Code: "source_too_large", Err: fmt.Errorf("remote source exceeds %d-byte plan-catalog limit: %s", plancatalog.MaxPlanSourceBytes, blob.Path)}
	}
	data := make([]byte, size)
	if _, err := io.ReadFull(reader.stdout, data); err != nil {
		return nil, err
	}
	if terminator, err := reader.stdout.ReadByte(); err != nil || terminator != '\n' {
		return nil, fmt.Errorf("remote_blob_invalid_terminator: %v", err)
	}
	return data, nil
}

func (reader *remoteBatchReader) Close() error {
	_ = reader.stream.Stdin.Close()
	result, waitErr := reader.stream.Wait()
	_ = reader.stream.Stdout.Close()
	if waitErr != nil {
		return fmt.Errorf("remote_blob_batch_failed: %w: %s", waitErr, strings.TrimSpace(string(result.Stderr)))
	}
	return reader.repository.verifyAuthority()
}
