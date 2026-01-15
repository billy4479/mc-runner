package main

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
)

const (
	githubLatestReleaseAPI = "https://api.github.com/repos/cloudflare/cloudflared/releases/latest"
	cacheSubdir            = "cloudflared-cache"
)

type ghRelease struct {
	TagName string    `json:"tag_name"`
	Assets  []ghAsset `json:"assets"`
}
type ghAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

var flags string

func main() {
	fmt.Printf("flags: %s\n", flags)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	cloudflaredPath, err := ensureCloudflared(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to ensure cloudflared:", err)
		os.Exit(1)
	}

	if err := execInto(cloudflaredPath, strings.Fields(flags)); err != nil {
		fmt.Fprintln(os.Stderr, "failed to exec cloudflared:", err)
		os.Exit(1)
	}
}

func ensureCloudflared(ctx context.Context) (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(cacheDir, cacheSubdir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	binName := "cloudflared"
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}
	binPath := filepath.Join(dir, binName)

	if fileExists(binPath) {
		// Try "cloudflared update" best-effort.
		// Note: Cloudflare docs mention some update flows work best when installed as a service.
		// We still attempt it, and fall back to re-download if it fails.

		err := tryUpdate(ctx, binPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while updating cloudflared: %v\n", err)
		}
		return binPath, nil
	}

	// Download latest.
	tmpPath, err := downloadLatestCloudflared(ctx, dir)
	if err != nil {
		return "", err
	}

	// Move into place atomically when possible.
	if runtime.GOOS == "windows" {
		// Windows can't rename over an existing file reliably; ensure removed already.
		_ = os.Remove(binPath)
	}
	if err := os.Rename(tmpPath, binPath); err != nil {
		// Cross-device rename fallback.
		if err2 := copyFile(tmpPath, binPath); err2 != nil {
			return "", fmt.Errorf("rename failed (%v) and copy failed: %w", err, err2)
		}
		_ = os.Remove(tmpPath)
	}

	// Ensure executable bit on unix.
	if runtime.GOOS != "windows" {
		_ = os.Chmod(binPath, 0o755)
	}

	return binPath, nil
}

func tryUpdate(ctx context.Context, binPath string) error {
	cctx, cancel := context.WithTimeout(ctx, 45*time.Second)
	defer cancel()

	cmd := exec.CommandContext(cctx, binPath, "update")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func downloadLatestCloudflared(ctx context.Context, outDir string) (string, error) {
	release, err := fetchLatestRelease(ctx)
	if err != nil {
		return "", err
	}

	assetName, kind, err := desiredAssetName()
	if err != nil {
		return "", err
	}

	var url string
	for _, a := range release.Assets {
		if a.Name == assetName {
			url = a.BrowserDownloadURL
			break
		}
	}
	if url == "" {
		// Helpful error: list available assets for this release.
		var names []string
		for _, a := range release.Assets {
			names = append(names, a.Name)
		}
		return "", fmt.Errorf("asset %q not found in latest release %q. available: %s",
			assetName, release.TagName, strings.Join(names, ", "))
	}

	// Download to a temp file in outDir.
	tmpDownload := filepath.Join(outDir, assetName+".download")
	if err := httpDownload(ctx, url, tmpDownload); err != nil {
		return "", err
	}

	// If tgz, extract cloudflared out of it to a temp binary file.
	if kind == "tgz" {
		tmpBin := filepath.Join(outDir, "cloudflared.extracted.tmp")
		if err := extractCloudflaredFromTGZ(tmpDownload, tmpBin); err != nil {
			return "", err
		}
		_ = os.Remove(tmpDownload)
		if runtime.GOOS != "windows" {
			_ = os.Chmod(tmpBin, 0o755)
		}
		return tmpBin, nil
	}

	// Otherwise it's already a binary.
	if runtime.GOOS != "windows" {
		_ = os.Chmod(tmpDownload, 0o755)
	}
	return tmpDownload, nil
}

func fetchLatestRelease(ctx context.Context) (*ghRelease, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, githubLatestReleaseAPI, nil)
	if err != nil {
		return nil, err
	}

	// GitHub API likes a UA.
	req.Header.Set("User-Agent", "cloudflared-bootstrapper/1.0")
	// If you hit rate limits, you can set GITHUB_TOKEN env and uncomment:
	// if tok := os.Getenv("GITHUB_TOKEN"); tok != "" {
	// 	req.Header.Set("Authorization", "Bearer "+tok)
	// }

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<10))
		return nil, fmt.Errorf("github api status %s: %s", resp.Status, strings.TrimSpace(string(b)))
	}

	var r ghRelease
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

func desiredAssetName() (name string, kind string, err error) {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	switch goos {
	case "linux":
		switch goarch {
		case "amd64":
			return "cloudflared-linux-amd64", "bin", nil
		case "arm64":
			return "cloudflared-linux-arm64", "bin", nil
		default:
			return "", "", fmt.Errorf("unsupported linux arch: %s", goarch)
		}
	case "darwin":
		switch goarch {
		case "amd64":
			return "cloudflared-darwin-amd64.tgz", "tgz", nil
		case "arm64":
			return "cloudflared-darwin-arm64.tgz", "tgz", nil
		default:
			return "", "", fmt.Errorf("unsupported darwin arch: %s", goarch)
		}
	case "windows":
		switch goarch {
		case "amd64":
			return "cloudflared-windows-amd64.exe", "bin", nil
		default:
			return "", "", fmt.Errorf("unsupported windows arch: %s", goarch)
		}
	default:
		return "", "", fmt.Errorf("unsupported OS: %s", goos)
	}
}

func httpDownload(ctx context.Context, url, outPath string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "cloudflared-bootstrapper/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<10))
		return fmt.Errorf("download status %s: %s", resp.Status, strings.TrimSpace(string(b)))
	}

	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	_, err = io.Copy(f, resp.Body)
	return err
}

func extractCloudflaredFromTGZ(tgzPath, outBinPath string) error {
	f, err := os.Open(tgzPath)
	if err != nil {
		return err
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}

		// Usually the tar contains a single "cloudflared" file.
		base := filepath.Base(hdr.Name)
		if base != "cloudflared" {
			continue
		}
		if hdr.Typeflag != tar.TypeReg {
			continue
		}

		out, err := os.Create(outBinPath)
		if err != nil {
			return err
		}
		if _, err := io.Copy(out, tr); err != nil {
			_ = out.Close()
			return err
		}
		return out.Close()
	}

	return fmt.Errorf("cloudflared binary not found inside tgz: %s", tgzPath)
}

func execInto(binPath string, args []string) error {
	argv := append([]string{binPath}, args...)

	if runtime.GOOS == "windows" {
		// No syscall.Exec; spawn and forward stdio.
		cmd := exec.Command(binPath, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		return cmd.Run()
	}

	// Replace current process.
	return syscall.Exec(binPath, argv, os.Environ())
}

func fileExists(path string) bool {
	st, err := os.Stat(path)
	return err == nil && !st.IsDir()
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}
