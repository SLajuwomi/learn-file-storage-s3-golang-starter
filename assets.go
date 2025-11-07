package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func (cfg apiConfig) ensureAssetsDir() error {
	if _, err := os.Stat(cfg.assetsRoot); os.IsNotExist(err) {
		return os.Mkdir(cfg.assetsRoot, 0755)
	}
	return nil
}

func (cfg *apiConfig) getObjectURL(key string) string {
	objectURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", cfg.s3Bucket, cfg.s3Region, key)
	return objectURL
}

func getAssetPath(videoID string, mediaType string) string {
	ext := mediaTypeToExt(mediaType)
	return fmt.Sprintf("%s%s", videoID, ext)
}

func (cfg apiConfig) getAssetDiskPath(assetPath string) string {
	return filepath.Join(cfg.assetsRoot, assetPath)
}

func (cfg apiConfig) getAssetURL(assetPath string) string {
	return fmt.Sprintf("http://localhost:%s/assets/%s", cfg.port, assetPath)
}

func create32ByteHex() string {
	randomVideoID := make([]byte, 32)
	rand.Read(randomVideoID)
	randomVideoIDString := base64.RawURLEncoding.EncodeToString(randomVideoID)
	return randomVideoIDString
}

func mediaTypeToExt(mediaType string) string {
	parts := strings.Split(mediaType, "/")
	if len(parts) != 2 {
		return ".bin"
	}
	return "." + parts[1]
}

func getVideoAspectRatio(filePath string) (string, error) {
	command := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)
	var buffer bytes.Buffer
	var stdErr bytes.Buffer
	command.Stdout = &buffer
	command.Stderr = &stdErr
	err := command.Run()
	if err != nil {
		log.Print(stdErr.String())
		return "", err
	}
	var vidInfo VideoInfo
	err = json.Unmarshal(buffer.Bytes(), &vidInfo)
	if err != nil {
		return "", err
	}

	ratio := float64(vidInfo.Streams[0].Width) / float64(vidInfo.Streams[0].Height)
	rounded := math.Round(ratio*100) / 100
	if rounded == 1.78 {
		return "16:9", nil
	}
	if rounded == 0.56 {
		return "9:16", nil
	}
	return "other", nil
}
