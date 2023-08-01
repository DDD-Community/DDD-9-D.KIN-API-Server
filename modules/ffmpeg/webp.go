package ffmpeg

import (
	"context"
	"d.kin-app/internal/awsx/lambdax"
	"d.kin-app/internal/typex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var ffmpegPath = typex.ByLazy(func() string {
	if lambdax.IsLambdaRuntime() {
		return "/opt/ffmpeg/ffmpeg"
	}

	return "ffmpeg"
})

func EncodeWebP(ctx context.Context, r io.Reader) (*os.File, error) {
	file, err := os.CreateTemp(os.TempDir(), "image-")
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(file, r)
	if err != nil {
		file.Close()
		return nil, err
	}
	if err = file.Close(); err != nil {
		return nil, err
	}

	inFilename := file.Name()
	outFilename := filepath.Join(os.TempDir(), fmt.Sprintf("image-%d.webp", time.Now().UnixNano()))
	err = exec.CommandContext(
		ctx,
		ffmpegPath.Value(),
		"-v", "quiet",
		"-i", inFilename,
		"-vcodec", "libwebp",
		outFilename,
	).Run()
	if err != nil {
		return nil, err
	}

	return os.Open(outFilename)
}
