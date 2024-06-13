package render

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"
)

type Read struct {
	ctx  context.Context //nolint:containedctx
	fs   afero.Fs
	base string
}

func NewRead(ctx context.Context, fs afero.Fs, base string) *Read {
	return &Read{
		ctx:  ctx,
		fs:   fs,
		base: base,
	}
}

func (r *Read) path(fileName string) string {
	if filepath.IsAbs(fileName) {
		return fileName
	}
	return filepath.Join(r.base, fileName)
}

func (r *Read) Run(fileName string) (string, error) {
	s, err := afero.ReadFile(r.fs, r.path(fileName))
	if err != nil {
		return "", fmt.Errorf("read a file: %w", err)
	}
	return string(s), nil
}
