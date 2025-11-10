package main

import (
	"strings"
	"time"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"
)

func (cfg *apiConfig) dbVideoToSignedVideo(video database.Video) (database.Video, error) {
	if video.VideoURL == nil {
		return video, nil
	}
	splitURL := strings.Split(*video.VideoURL, ",")
	bucket := splitURL[0]
	key := splitURL[1]
	presignedURL, err := generatePresignedURL(cfg.s3Client, bucket, key, 1*time.Hour)
	if err != nil {
		return video, err
	}
	video.VideoURL = &presignedURL
	return video, nil
}
