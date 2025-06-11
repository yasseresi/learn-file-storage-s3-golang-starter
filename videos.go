package main

import (
	"bytes"
	"encoding/json"
	"math"
	"os/exec"
)

type FFProbeOutput struct {
	Streams []FFProbeStream `json:"streams"`
}

type FFProbeStream struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

func getVideoAspectRatio(filePath string) (string, error) {
	// Create the ffprobe command
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)

	// Set up stdout buffer
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	// Run the command
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	// Unmarshal the JSON output
	var probeOutput FFProbeOutput
	err = json.Unmarshal(stdout.Bytes(), &probeOutput)
	if err != nil {
		return "", err
	}

	// Find the first video stream
	for _, stream := range probeOutput.Streams {
		if stream.Width > 0 && stream.Height > 0 {
			// Calculate aspect ratio
			ratio := float64(stream.Width) / float64(stream.Height)

			// Determine the aspect ratio category
			// 16:9 = 1.777...
			// 9:16 = 0.5625
			if math.Abs(ratio-16.0/9.0) < 0.1 {
				return "16:9", nil
			} else if math.Abs(ratio-9.0/16.0) < 0.1 {
				return "9:16", nil
			} else {
				return "other", nil
			}
		}
	}

	return "other", nil
}

func aspectRatioToPrefix(aspectRatio string) string {
	switch aspectRatio {
	case "16:9":
		return "landscape"
	case "9:16":
		return "portrait"
	default:
		return "other"
	}
}

func processVideoForFastStart(filePath string) (string, error) {
	// Create output file path by appending .processing
	outputPath := filePath + ".processing"

	// Create ffmpeg command for fast start processing
	cmd := exec.Command("ffmpeg", "-i", filePath, "-c", "copy", "-movflags", "faststart", "-f", "mp4", outputPath)

	// Run the command
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return outputPath, nil
}
