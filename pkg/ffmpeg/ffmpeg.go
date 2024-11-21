package ffmpeg

import (
	"io"
	"os/exec"
)

type AudioOptions struct {
	Filename   string
	Bitrate    string // e.g., "128k"
	SampleRate string // e.g., "44100"
	Channels   string // e.g., "2"
	Codec      string // e.g., "mp3", "aac"
	Quality    string // e.g., "2" (for -q:a)
}

func StreamAudio(opts AudioOptions, w io.Writer) error {
	// Set defaults if options are missing
	if opts.Bitrate == "" {
		opts.Bitrate = "128k"
	}
	if opts.SampleRate == "" {
		opts.SampleRate = "44100"
	}
	if opts.Channels == "" {
		opts.Channels = "2"
	}
	if opts.Codec == "" {
		opts.Codec = "mp3"
	}
	if opts.Quality == "" {
		opts.Quality = "2"
	}

	// Build the ffmpeg command arguments
	args := []string{
		"-i", opts.Filename,
		"-vn", // No video
		"-c:a", opts.Codec,
		"-b:a", opts.Bitrate,
		"-ar", opts.SampleRate,
		"-ac", opts.Channels,
		"-q:a", opts.Quality,
		"-f", opts.Codec,
		"pipe:1",
	}

	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = w
	cmd.Stderr = nil // Optionally handle stderr for logging

	err := cmd.Run()
	return err
}
