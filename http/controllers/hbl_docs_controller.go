package controllers

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type archiveDownloadFile struct {
	Filename string `json:"filename"`
	URL      string `json:"url"`
}

type archiveDownloadResponse struct {
	Files []archiveDownloadFile `json:"files"`
}

type downloadedArchiveFile struct {
	Name string
	Data []byte
}

// DownloadHBLDocsArchive handles POST /api/v1/hbl-docs/download-archive.
// It accepts either:
// - {"files":[{"filename":"...","url":"..."}]}
// - [{"filename":"...","url":"..."}]
func DownloadHBLDocsArchive(ctx *gin.Context) {
	requestBody, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
		return
	}

	files, err := parseArchiveDownloadFiles(requestBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(files) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "files must contain at least one downloadable item"})
		return
	}

	downloadedFiles, err := downloadArchiveFiles(ctx.Request.Context(), files)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	zipData, err := buildArchiveZip(downloadedFiles)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to build zip archive"})
		return
	}

	ctx.Header("Content-Type", "application/zip")
	ctx.Header("Content-Disposition", `attachment; filename="hbl-documents.zip"`)
	ctx.Data(http.StatusOK, "application/zip", zipData)
}

func parseArchiveDownloadFiles(body []byte) ([]archiveDownloadFile, error) {
	var wrapped archiveDownloadResponse
	if err := json.Unmarshal(body, &wrapped); err == nil {
		return normalizeArchiveFiles(wrapped.Files), nil
	}

	var direct []archiveDownloadFile
	if err := json.Unmarshal(body, &direct); err == nil {
		return normalizeArchiveFiles(direct), nil
	}

	return nil, fmt.Errorf("invalid request body: expected {\"files\": [...]} or a JSON array of files")
}

func normalizeArchiveFiles(files []archiveDownloadFile) []archiveDownloadFile {
	normalized := make([]archiveDownloadFile, 0, len(files))
	for _, file := range files {
		file.Filename = strings.TrimSpace(file.Filename)
		file.URL = strings.TrimSpace(file.URL)
		if file.URL == "" {
			continue
		}
		normalized = append(normalized, file)
	}
	return normalized
}

func downloadArchiveFiles(ctx context.Context, files []archiveDownloadFile) ([]downloadedArchiveFile, error) {
	client := &http.Client{Timeout: 120 * time.Second}
	downloaded := make([]downloadedArchiveFile, 0, len(files))

	for i, file := range files {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, file.URL, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request for %q: %w", file.URL, err)
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to download %q: %w", file.URL, err)
		}

		body, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			return nil, fmt.Errorf("failed to read response for %q: %w", file.URL, readErr)
		}

		if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
			return nil, fmt.Errorf("upstream returned %d for %q", resp.StatusCode, file.URL)
		}

		name := sanitizeArchiveEntryName(file.Filename, file.URL, i)
		downloaded = append(downloaded, downloadedArchiveFile{
			Name: name,
			Data: body,
		})
	}

	return downloaded, nil
}

func sanitizeArchiveEntryName(filename, rawURL string, index int) string {
	name := strings.TrimSpace(filename)
	if name == "" {
		if parsed, err := url.Parse(rawURL); err == nil {
			name = path.Base(parsed.Path)
		}
	}
	if name == "" || name == "." || name == "/" {
		name = fmt.Sprintf("document-%d.pdf", index+1)
	}

	name = strings.ReplaceAll(name, "\\", "/")
	name = path.Base(name)
	if name == "." || name == "/" || name == "" {
		return fmt.Sprintf("document-%d.pdf", index+1)
	}
	return name
}

func buildArchiveZip(files []downloadedArchiveFile) ([]byte, error) {
	var buffer bytes.Buffer
	zw := zip.NewWriter(&buffer)

	for _, file := range files {
		entry, err := zw.Create(file.Name)
		if err != nil {
			zw.Close()
			return nil, fmt.Errorf("failed to create zip entry %q: %w", file.Name, err)
		}
		if _, err := entry.Write(file.Data); err != nil {
			zw.Close()
			return nil, fmt.Errorf("failed to write zip entry %q: %w", file.Name, err)
		}
	}

	if err := zw.Close(); err != nil {
		return nil, fmt.Errorf("failed to finalize zip archive: %w", err)
	}

	return buffer.Bytes(), nil
}
