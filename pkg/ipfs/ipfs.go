package ipfsutils

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/kubo/client/rpc"
	"github.com/rs/zerolog/log"
	"github.com/tensved/snet-matrix-framework/internal/config"
)

// IPFSClient wraps the IPFS HTTP API client.
type IPFSClient struct {
	*rpc.HttpApi
}

// Init initializes a new IPFS client with a specified timeout from the configuration.
//
// Returns:
//   - IPFSClient: An instance of IPFSClient with the configured HTTP client.
func Init() IPFSClient {
	duration := 5 * time.Second
	httpClient := http.Client{
		Timeout: duration,
	}
	ipfsClient, err := rpc.NewURLApiWithClient(config.IPFS.IPFSProviderURL, &httpClient)
	if err != nil {
		log.Fatal().Err(err).Msg("connection failed to IPFS")
	}
	return IPFSClient{ipfsClient}
}

// ReadFilesCompressed reads all files from a compressed tar archive and returns them as a map.
// Supports both regular tar archives and gzip compressed tar archives.
//
// Parameters:
//   - compressedFile: A string representing the compressed file content.
//
// Returns:
//   - protofiles: A map where keys are file names and values are file contents.
//   - err: An error if the operation fails.
func ReadFilesCompressed(compressedFile string) (map[string][]byte, error) {
	// Check that content is not an HTML error page
	if strings.Contains(compressedFile, "<html>") || strings.Contains(compressedFile, "403 Forbidden") || strings.Contains(compressedFile, "404 Not Found") {
		return nil, fmt.Errorf("received HTML error page instead of tar archive")
	}

	// Check minimum size for tar archive
	if len(compressedFile) < 100 {
		return nil, fmt.Errorf("content too small to be a valid archive (size: %d)", len(compressedFile))
	}

	var reader io.Reader

	// Check if this is a gzip archive (first 2 bytes: 0x1f 0x8b)
	if len(compressedFile) >= 2 && compressedFile[0] == 0x1f && compressedFile[1] == 0x8b {
		log.Debug().Msg("detected gzip archive, attempting to decompress")
		gzReader, err := gzip.NewReader(strings.NewReader(compressedFile))
		if err != nil {
			return nil, fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer gzReader.Close()
		reader = gzReader
	} else {
		reader = strings.NewReader(compressedFile)
	}

	tarReader := tar.NewReader(reader)
	protoFiles := map[string][]byte{}

	filesRead := 0

	for {
		header, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Error().Err(err).Msg("failed to read tar header")
			return nil, fmt.Errorf("failed to read tar header: %w", err)
		}
		name := header.Name
		switch header.Typeflag {
		case tar.TypeDir:
			log.Debug().Any("dir_name", name).Msg("directory name")
		case tar.TypeReg:
			data := make([]byte, header.Size)
			_, err = tarReader.Read(data)
			if err != nil && !errors.Is(err, io.EOF) {
				log.Error().Err(err).Str("file_name", name).Msg("failed to read file content")
				return nil, fmt.Errorf("failed to read file %s: %w", name, err)
			}
			protoFiles[name] = data
			filesRead++
			log.Debug().Str("file_name", name).Int("file_size", len(data)).Msg("successfully read file")
		default:
			err = fmt.Errorf("unknown file type: %c %s %s", header.Typeflag, "in file", name)
			log.Error().Err(err)
			return nil, err
		}
	}

	if filesRead == 0 {
		return nil, fmt.Errorf("no files found in tar archive")
	}

	log.Info().Int("files_read", filesRead).Msg("successfully read files from tar archive")
	return protoFiles, nil
}

// Deprecated: RemoveSpecialCharacters removes all non-alphanumeric characters from the provided hash string.
//
// Parameters:
//   - hash: A string representing the hash to be cleaned.
//
// Returns:
//   - string: The cleaned hash string.
func RemoveSpecialCharacters(hash string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9=]")
	return reg.ReplaceAllString(hash, "")
}

// formatHash improved hash formatting function
func formatHash(hash string) string {
	hash = strings.Replace(hash, "ipfs://", "", -1)
	hash = strings.Replace(hash, "filecoin://", "", -1)
	hash = RemoveSpecialCharacters(hash)
	return hash
}

// GetIpfsFile retrieves a file from IPFS using the provided hash.
// Improved version with better CID error handling
//
// Parameters:
//   - hash: A string representing the IPFS hash of the file.
//
// Returns:
//   - content: A byte slice containing the file content.
//   - err: An error if the operation fails.
func (ipfsClient IPFSClient) GetIpfsFile(hash string) ([]byte, error) {
	hash = formatHash(hash)

	cID, err := cid.Parse(hash)
	if err != nil {
		log.Error().Err(err).Str("hash", hash).Msg("failed to parse CID, trying to continue")
		return nil, fmt.Errorf("invalid cid: %w", err)
	}

	req := ipfsClient.Request("cat", cID.String())
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	resp, err := req.Send(ctx)
	if err != nil {
		log.Error().Err(err).Str("hash", hash).Msg("failed to get content from IPFS")
		return nil, err
	}

	defer func(resp *rpc.Response) {
		if resp != nil {
			err = resp.Close()
			if err != nil {
				log.Error().Err(err).Msg("failed to close IPFS response")
			}
		}
	}(resp)

	if resp == nil {
		log.Error().Str("hash", hash).Msg("IPFS response is nil")
		return nil, fmt.Errorf("IPFS response is nil")
	}

	if resp.Error != nil {
		log.Error().Err(resp.Error).Str("hash", hash).Msg("IPFS response contains error")
		return nil, resp.Error
	}

	fileContent, err := io.ReadAll(resp.Output)
	if err != nil {
		log.Error().Err(err).Str("hash", hash).Msg("failed to read IPFS content")
		return nil, err
	}

	if len(fileContent) == 0 {
		log.Error().Str("hash", hash).Msg("IPFS returned empty content")
		return nil, fmt.Errorf("IPFS returned empty content")
	}

	contentStr := string(fileContent)
	if strings.Contains(contentStr, "<html>") || strings.Contains(contentStr, "403 Forbidden") || strings.Contains(contentStr, "404 Not Found") {
		log.Error().Str("hash", hash).Str("content_preview", contentStr[:min(100, len(contentStr))]).Msg("IPFS returned HTML error page")
		return nil, fmt.Errorf("IPFS returned HTML error page instead of file content")
	}

	if len(fileContent) > 0 {
		_, c, err := cid.CidFromBytes(append(cID.Bytes(), fileContent...))
		if err != nil {
			log.Warn().Err(err).Str("hash", hash).Msg("failed to verify CID, but content retrieved successfully")
		} else if !c.Equals(cID) {
			log.Warn().Str("expected_hash", hash).Str("actual_hash", c.String()).Msg("CID verification failed, but content retrieved")
		}
	}

	log.Debug().Str("hash", hash).Int("content_size", len(fileContent)).Msg("successfully retrieved content from IPFS")
	return fileContent, nil
}

// min returns the minimum of two numbers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
