package ipfsutils

import (
	"archive/tar"
	"context"
	"fmt"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/kubo/client/rpc"
	"github.com/rs/zerolog/log"
	"io"
	"matrix-ai-framework/internal/config"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type IPFSClient struct {
	*rpc.HttpApi
}

func Init() IPFSClient {
	httpClient := http.Client{
		//Timeout: time.Duration(config.GetInt(config.IpfsTimeout)) * time.Second,
		Timeout: 5 * time.Second,
	}
	ifpsClient, err := rpc.NewURLApiWithClient(config.IPFS.IPFSProviderURL, &httpClient)
	if err != nil {
		log.Fatal().Err(err).Msg("Connection failed to IPFS")
	}
	return IPFSClient{ifpsClient}
}

// ReadFilesCompressed - read all files which have been compressed, there can be more than one file
// We need to start reading the proto files associated with the service.
// proto files are compressed and stored as modelipfsHash
func ReadFilesCompressed(compressedFile string) (protofiles map[string][]byte, err error) {
	f := strings.NewReader(compressedFile)
	tarReader := tar.NewReader(f)
	protofiles = map[string][]byte{}
	for true {
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
			//log.Debug().Any("file Name", name).Msg("file Name")
			data := make([]byte, header.Size)
			_, err := tarReader.Read(data)
			if err != nil && err != io.EOF {
				log.Error().Err(err)
				return nil, err
			}
			protofiles[name] = data
		default:
			err = fmt.Errorf(fmt.Sprintf("%s : %c %s %s\n",
				"Unknown file Type ",
				header.Typeflag,
				"in file",
				name,
			))
			log.Error().Err(err)
			return nil, err
		}
	}
	return protofiles, nil
}

func RemoveSpecialCharacters(hash string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9=]")
	if err != nil {
		log.Error().Err(err).Msg("Failed to remove special characters from hash")
	}
	return reg.ReplaceAllString(hash, "")
}

func (ipfsClient IPFSClient) GetIpfsFile(hash string) (content []byte, err error) {
	hash = strings.TrimPrefix(hash, "ipfs://")
	hash = RemoveSpecialCharacters(hash)

	cID, err := cid.Parse(hash)
	if err != nil {
		log.Error().Err(err)
		return
	}

	req := ipfsClient.Request("cat", cID.String())
	resp, err := req.Send(context.Background())
	defer func(resp *rpc.Response) {
		err := resp.Close()
		if err != nil {
			log.Error().Err(err)
		}
	}(resp)
	if err != nil {
		log.Error().Err(err)
		return
	}
	if resp == nil {
		log.Error().Msg("resp is nil!")
		return
	}
	if resp.Error != nil {
		log.Err(resp.Error)
		return
	}
	fileContent, err := io.ReadAll(resp.Output)
	if err != nil {
		log.Error().Err(err)
		return
	}

	// Create a cid manually to check cid
	_, c, err := cid.CidFromBytes(append(cID.Bytes(), fileContent...))
	if err != nil {
		log.Error().Err(err)
		return
	}

	// To test if two cid's are equivalent, be sure to use the 'Equals' method:
	if !c.Equals(cID) {
		log.Print("cIDs not equals!")
	}

	return fileContent, err
}
