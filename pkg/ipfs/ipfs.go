package ipfsutils

import (
	"archive/tar"
	"context"
	"errors"
	"fmt"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/kubo/client/rpc"
	"github.com/rs/zerolog/log"
	"github.com/tensved/snet-matrix-framework/internal/config"
	"io"
	"net/http"
	"strings"
	"time"
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
//
// Parameters:
//   - compressedFile: A string representing the compressed file content.
//
// Returns:
//   - protofiles: A map where keys are file names and values are file contents.
//   - err: An error if the operation fails.
func ReadFilesCompressed(compressedFile string) (map[string][]byte, error) {
	f := strings.NewReader(compressedFile)
	tarReader := tar.NewReader(f)
	protoFiles := map[string][]byte{}
	for {
		header, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Error().Err(err).Msg("failed to get organizations")
			return nil, err
		}
		name := header.Name
		switch header.Typeflag {
		case tar.TypeDir:
			log.Debug().Any("dir_name", name).Msg("directory name")
		case tar.TypeReg:
			data := make([]byte, header.Size)
			_, err = tarReader.Read(data)
			if err != nil && !errors.Is(err, io.EOF) {
				log.Error().Err(err)
				return nil, err
			}
			protoFiles[name] = data
		default:
			err = fmt.Errorf("unknown file type: %c %s %s", header.Typeflag, "in file", name)
			log.Error().Err(err)
			return nil, err
		}
	}
	return protoFiles, nil
}

// RemoveSpecialCharacters removes all non-alphanumeric characters from the provided hash string.
//
// Parameters:
//   - hash: A string representing the hash to be cleaned.
//
// Returns:
//   - string: The cleaned hash string.
func RemoveSpecialCharacters(hash string) string {
	return config.IPFS.HashCutterRegexp.ReplaceAllString(hash, "")
}

// GetIpfsFile retrieves a file from IPFS using the provided hash.
//
// Parameters:
//   - hash: A string representing the IPFS hash of the file.
//
// Returns:
//   - content: A byte slice containing the file content.
//   - err: An error if the operation fails.
func (ipfsClient IPFSClient) GetIpfsFile(hash string) ([]byte, error) {
	hash = strings.TrimPrefix(hash, "ipfs://")
	hash = RemoveSpecialCharacters(hash)

	cID, err := cid.Parse(hash)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse CID")
		return nil, err
	}

	req := ipfsClient.Request("cat", cID.String())
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	resp, err := req.Send(ctx)
	defer func(resp *rpc.Response) {
		err = resp.Close()
		if err != nil {
			log.Error().Err(err).Msg("failed to close resp, got error")
		}
	}(resp)
	if err != nil {
		log.Error().Err(err).Msg("failed to get content, got error")
		return nil, err
	}
	if resp == nil {
		log.Error().Msg("failed to get content, response is nil")
		return nil, err
	}
	if resp.Error != nil {
		log.Error().Err(resp.Error).Msg("failed to get content, got error")
		return nil, err
	}
	fileContent, err := io.ReadAll(resp.Output)
	if err != nil {
		log.Error().Err(err).Msg("failed to read content, got error")
		return nil, err
	}

	// Create a CID manually to check CID.
	_, c, err := cid.CidFromBytes(append(cID.Bytes(), fileContent...))
	if err != nil {
		log.Error().Err(err).Msg("failed to parse CID")
		return nil, err
	}

	// To test if two CIDs are equal, be sure to use the 'Equals' method.
	if !c.Equals(cID) {
		log.Error().Msg("CIDs are not equal")
	}

	return fileContent, err
}
