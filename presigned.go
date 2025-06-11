package main

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"
)

func generatePresignedURL(s3Client *s3.Client, bucket, key string, expireTime time.Duration) (string, error) {
	// Create a presign client
	presignClient := s3.NewPresignClient(s3Client)

	// Create the presigned request
	presignedReq, err := presignClient.PresignGetObject(context.Background(), &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}, s3.WithPresignExpires(expireTime))
	if err != nil {
		return "", err
	}

	return presignedReq.URL, nil
}

func (cfg *apiConfig) dbVideoToSignedVideo(video database.Video) (database.Video, error) {
	// If VideoURL is nil or empty, return as is
	if video.VideoURL == nil || *video.VideoURL == "" {
		return video, nil
	}

	// Split the video URL on comma to get bucket and key
	parts := strings.Split(*video.VideoURL, ",")
	if len(parts) != 2 {
		// If it's not in the expected format, return as is
		return video, nil
	}

	bucket := parts[0]
	key := parts[1]

	// Generate presigned URL (expires in 1 hour)
	presignedURL, err := generatePresignedURL(cfg.s3Client, bucket, key, time.Hour)
	if err != nil {
		return video, err
	}

	// Create a copy of the video with the presigned URL
	signedVideo := video
	signedVideo.VideoURL = &presignedURL

	return signedVideo, nil
}
