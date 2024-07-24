package ipfsutils

import (
	"archive/tar"
	"context"
	"fmt"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/kubo/client/rpc"
	"github.com/rs/zerolog/log"
	"github.com/tensved/snet-matrix-framework/internal/config"
	"io"
	"net/http"
	"regexp"
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
	httpClient := http.Client{
		Timeout: 5 * time.Second,
	}
	ifpsClient, err := rpc.NewURLApiWithClient(config.IPFS.IPFSProviderURL, &httpClient)
	if err != nil {
		log.Fatal().Err(err).Msg("Connection failed to IPFS")
	}
	return IPFSClient{ifpsClient}
}

// ReadFilesCompressed reads all files from a compressed tar archive and returns them as a map.
//
// Parameters:
//   - compressedFile: A string representing the compressed file content.
//
// Returns:
//   - protofiles: A map where keys are file names and values are file contents.
//   - err: An error if the operation fails.
func ReadFilesCompressed(compressedFile string) (protofiles map[string][]byte, err error) {
	f := strings.NewReader(compressedFile)
	tarReader := tar.NewReader(f)
	protofiles = map[string][]byte{}
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Error().Err(err).Msg("Failed to get organizations")
			return nil, err
		}
		name := header.Name
		switch header.Typeflag {
		case tar.TypeDir:
			log.Debug().Any("dir_name", name).Msg("Directory Name")
		case tar.TypeReg:
			data := make([]byte, header.Size)
			_, err := tarReader.Read(data)
			if err != nil && err != io.EOF {
				log.Error().Err(err)
				return nil, err
			}
			protofiles[name] = data
		default:
			err = fmt.Errorf("%s : %c %s %s\n", "Unknown file Type", header.Typeflag, "in file", name)
			log.Error().Err(err)
			return nil, err
		}
	}
	return protofiles, nil
}

// RemoveSpecialCharacters removes all non-alphanumeric characters from the provided hash string.
//
// Parameters:
//   - hash: A string representing the hash to be cleaned.
//
// Returns:
//   - string: The cleaned hash string.
func RemoveSpecialCharacters(hash string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9=]")
	if err != nil {
		log.Error().Err(err).Msg("Failed to remove special characters from hash")
	}
	return reg.ReplaceAllString(hash, "")
}

// GetIpfsFile retrieves a file from IPFS using the provided hash.
//
// Parameters:
//   - hash: A string representing the IPFS hash of the file.
//
// Returns:
//   - content: A byte slice containing the file content.
//   - err: An error if the operation fails.
func (ipfsClient IPFSClient) GetIpfsFile(hash string) (content []byte, err error) {
	hash = strings.TrimPrefix(hash, "ipfs://")
	hash = RemoveSpecialCharacters(hash)

	cID, err := cid.Parse(hash)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse CID")
		return
	}

	req := ipfsClient.Request("cat", cID.String())
	resp, err := req.Send(context.Background())
	defer func(resp *rpc.Response) {
		err = resp.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close resp, got error")
		}
	}(resp)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get content, got error")
		return
	}
	if resp == nil {
		log.Error().Msg("Failed to get content, response is nil")
		return
	}
	if resp.Error != nil {
		log.Error().Err(resp.Error).Msg("Failed to get content, got error")
		return
	}
	fileContent, err := io.ReadAll(resp.Output)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read content, got error")
		return
	}

	// Create a CID manually to check CID.
	_, c, err := cid.CidFromBytes(append(cID.Bytes(), fileContent...))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse CID")
		return
	}

	// To test if two CIDs are equivalent, be sure to use the 'Equals' method.
	if !c.Equals(cID) {
		log.Error().Msg("CIDs not equal!")
	}

	return fileContent, err
}
