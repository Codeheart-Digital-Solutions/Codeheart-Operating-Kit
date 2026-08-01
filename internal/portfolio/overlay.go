package portfolio

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"
)

type OverlayEvidence struct {
	Path     string
	Present  bool
	Bytes    []byte
	identity os.FileInfo
	file     *os.File
}

func ReadOverlay(root, homeID string) (OverlayEvidence, error) {
	path := filepath.Join(root, filepath.FromSlash(OverlayPath))
	absoluteRoot, err := filepath.Abs(root)
	if err != nil {
		return OverlayEvidence{}, err
	}
	repositoryRoot, err := os.OpenRoot(absoluteRoot)
	if err != nil {
		return OverlayEvidence{}, err
	}
	defer repositoryRoot.Close()
	info, err := repositoryRoot.Lstat(filepath.FromSlash(OverlayPath))
	if os.IsNotExist(err) {
		return OverlayEvidence{Path: path}, nil
	}
	if err != nil {
		return OverlayEvidence{}, fmt.Errorf("portfolio_overlay_unavailable: %w", err)
	}
	if !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
		return OverlayEvidence{}, fmt.Errorf("portfolio_overlay_unsafe: overlay is not a contained regular file")
	}
	file, err := repositoryRoot.Open(filepath.FromSlash(OverlayPath))
	if err != nil {
		return OverlayEvidence{}, fmt.Errorf("portfolio_overlay_unavailable: %w", err)
	}
	opened, err := file.Stat()
	if err != nil || !os.SameFile(info, opened) {
		_ = file.Close()
		return OverlayEvidence{}, fmt.Errorf("portfolio_overlay_changed: overlay identity changed while opening")
	}
	data, err := io.ReadAll(file)
	if err != nil {
		_ = file.Close()
		return OverlayEvidence{}, fmt.Errorf("portfolio_overlay_unavailable: %w", err)
	}
	value, err := state.DecodeAndValidateYAML(state.PortfolioOverlaySchema, data)
	if err != nil {
		_ = file.Close()
		return OverlayEvidence{}, fmt.Errorf("portfolio_overlay_invalid: %w", err)
	}
	if state.AsString(value["coordination_home_id"]) != homeID {
		_ = file.Close()
		return OverlayEvidence{}, fmt.Errorf("portfolio_overlay_home_mismatch: overlay belongs to a different coordination home")
	}
	return OverlayEvidence{Path: path, Present: true, Bytes: append([]byte{}, data...), identity: opened, file: file}, nil
}

func (evidence OverlayEvidence) VerifyUnchanged() error {
	if !evidence.Present {
		if _, err := os.Lstat(evidence.Path); os.IsNotExist(err) {
			return nil
		} else if err != nil {
			return err
		}
		return fmt.Errorf("portfolio_overlay_changed: overlay appeared during scan")
	}
	info, err := os.Lstat(evidence.Path)
	if err != nil {
		return fmt.Errorf("portfolio_overlay_changed: %w", err)
	}
	if !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 || evidence.identity == nil || !os.SameFile(evidence.identity, info) {
		return fmt.Errorf("portfolio_overlay_changed: strategic overlay identity changed during scan")
	}
	if evidence.file == nil {
		return fmt.Errorf("portfolio_overlay_changed: strategic overlay descriptor is unavailable")
	}
	opened, err := evidence.file.Stat()
	if err != nil || !os.SameFile(info, opened) {
		return fmt.Errorf("portfolio_overlay_changed: strategic overlay descriptor identity changed")
	}
	if _, err := evidence.file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("portfolio_overlay_changed: %w", err)
	}
	data, err := io.ReadAll(evidence.file)
	if err != nil {
		return fmt.Errorf("portfolio_overlay_changed: %w", err)
	}
	if !bytes.Equal(data, evidence.Bytes) {
		return fmt.Errorf("portfolio_overlay_changed: strategic overlay bytes changed during scan")
	}
	return nil
}

func (evidence OverlayEvidence) Close() {
	if evidence.file != nil {
		_ = evidence.file.Close()
	}
}
