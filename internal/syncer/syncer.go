package syncer

import (
	"context"
	"encoding/json"
	"github.com/bufbuild/protocompile"
	"github.com/rs/zerolog/log"
	"github.com/tensved/snet-matrix-framework/pkg/blockchain"
	"github.com/tensved/snet-matrix-framework/pkg/db"
	ipfs "github.com/tensved/snet-matrix-framework/pkg/ipfs"
	"google.golang.org/protobuf/reflect/protoreflect"
	"os"
	"strings"
	"time"
)

type SnetSyncer struct {
	Ethereum        blockchain.Ethereum
	IPFSClient      ipfs.IPFSClient
	DB              db.Service
	FileDescriptors map[string][]protoreflect.FileDescriptor
	cancelFunc      context.CancelFunc
}

func New(eth blockchain.Ethereum, ipfs ipfs.IPFSClient, db db.Service) SnetSyncer {
	return SnetSyncer{
		Ethereum:        eth,
		IPFSClient:      ipfs,
		DB:              db,
		FileDescriptors: make(map[string][]protoreflect.FileDescriptor),
	}
}

func (s *SnetSyncer) SyncOnce() {
	orgs, _ := s.Ethereum.GetOrgs()
	for _, orgIDBytes := range orgs {
		borg, err := s.Ethereum.GetOrg(orgIDBytes)
		if err != nil {
			log.Error().Err(err).Msg("failed to get org")
			continue
		}
		var org blockchain.OrganizationMetaData

		metadataJSON, err := s.IPFSClient.GetIpfsFile(string(borg.OrgMetadataURI))
		if err != nil {
			log.Error().Err(err).Msg("failed to get ipfs file")
			continue
		}

		err = json.Unmarshal(metadataJSON, &org)
		if err != nil {
			log.Error().Err(err).Any("content", string(metadataJSON)).Msg("cannot unmarshal org metadata from ipfs")
			continue
		}

		org.Owner = borg.Owner.Hex()
		org.SnetID = strings.ReplaceAll(string(borg.Id[:]), "\u0000", "")
		dbOrg, dbGroups := org.DB()
		orgID, err := s.DB.CreateSnetOrg(dbOrg)
		if err != nil {
			log.Error().Err(err).Msg("failed to create org")
		}
		org.ID = orgID
		err = s.DB.CreateSnetOrgGroups(orgID, dbGroups)
		if err != nil {
			log.Error().Err(err).Msg("failed to create org group")
		}
		var service blockchain.Service
		for _, serviceIDBytes := range borg.ServiceIds {
			service, err = s.Ethereum.GetService(borg.Id, serviceIDBytes)
			if err != nil {
				log.Error().Err(err)
				continue
			}

			metadataJSON, err = s.IPFSClient.GetIpfsFile(string(service.MetadataURI))
			if err != nil {
				log.Error().Err(err).Msg("failed to get file from ipfs")
				return
			}

			var srvMeta blockchain.ServiceMetadata
			err = json.Unmarshal(metadataJSON, &srvMeta)
			if err != nil {
				log.Error().Err(err).Any("content", string(metadataJSON)).Msg("failed to unmarshal metadata from ipfs")
				return
			}

			log.Debug().Msgf("metadata of service: %+v", srvMeta)

			srvMeta.OrgID = orgID
			srvMeta.SnetID = strings.ReplaceAll(string(serviceIDBytes[:]), "\u0000", "")
			srvMeta.SnetOrgID = org.SnetID
			dbSrvMeta, err := srvMeta.DB()
			if err != nil {
				log.Error().Err(err)
			}
			srvMeta.ID, err = s.DB.CreateSnetService(dbSrvMeta)
			if err != nil {
				log.Error().Err(err).Int("id", srvMeta.ID).Str("snet-id", srvMeta.SnetID).Msg("failed to add snet_service")
			}

			content, err := s.IPFSClient.GetIpfsFile(srvMeta.ModelIpfsHash)
			if err != nil {
				log.Error().Err(err)
			}
			protoFiles, err := ipfs.ReadFilesCompressed(string(content))
			if err != nil {
				log.Error().Err(err)
			}

			for fileName, fileContent := range protoFiles {
				fd := getFileDescriptor(string(fileContent), fileName)
				if fd != nil {
					s.FileDescriptors[srvMeta.SnetID] = append(s.FileDescriptors[srvMeta.SnetID], fd)
					err := os.WriteFile(fileName, fileContent, 0600)
					if err != nil {
						return
					}
				}
			}
		}
	}
}

func (s *SnetSyncer) Start(ctx context.Context) {
	// Store the cancel function for later use in Stop
	ctx, cancel := context.WithCancel(ctx)
	s.cancelFunc = cancel

	duration := 24 * time.Hour
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.SyncOnce()
		case <-ctx.Done():
			log.Info().Msg("snet syncer received shutdown signal, stopping sync process")
			return
		}
	}
}

func (s *SnetSyncer) Stop() {
	if s.cancelFunc != nil {
		s.cancelFunc()
	}
	log.Debug().Msg("snet syncer stopped successfully")
}

func getFileDescriptor(protoContent, name string) protoreflect.FileDescriptor {
	accessor := protocompile.SourceAccessorFromMap(map[string]string{
		name: protoContent,
	})
	compiler := protocompile.Compiler{
		Resolver:       &protocompile.SourceResolver{Accessor: accessor},
		SourceInfoMode: protocompile.SourceInfoStandard,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	fds, err := compiler.Compile(ctx, name)
	if err != nil {
		log.Error().Err(err).Msg("failed to create file descriptor")
		return nil
	}
	return fds.FindFileByPath(name)
}

func writeServiceSnetIDs(fileDescriptors map[string][]protoreflect.FileDescriptor) *strings.Builder {
	b := &strings.Builder{}
	b.WriteString("<div style=\"line-height: 0.8;\"><ol>")
	for snetID, descriptors := range fileDescriptors {
		for _, descriptor := range descriptors {
			b.WriteString("<li><strong>Snet ID: " + snetID + " Descriptor: " + string(descriptor.FullName().Name()) + "</strong></li>")
			services := descriptor.Services()
			if services != nil {
				b = writeServiceNames(b, services)
			}
		}
	}
	b.WriteString("</ol></div>")
	return b
}

func writeServiceNames(b *strings.Builder, services protoreflect.ServiceDescriptors) *strings.Builder {
	for i := range services.Len() {
		if services.Get(i) != nil {
			b.WriteString("<p><em>Service: " + string(services.Get(i).FullName().Name()) + "</em></p>")
			methods := services.Get(i).Methods()
			if methods != nil {
				b = writeMethodNames(b, methods)
			}
		}
	}
	return b
}

func writeMethodNames(b *strings.Builder, methods protoreflect.MethodDescriptors) *strings.Builder {
	b.WriteString("<p>🔁Methods: </p><ul>")
	for j := range methods.Len() {
		if methods.Get(j) != nil {
			b.WriteString("<li>" + string(methods.Get(j).FullName().Name()))
			b.WriteString("</li>")
		}
	}
	b.WriteString("</ul>")
	return b
}

func writeFields(b *strings.Builder, fields protoreflect.FieldDescriptors) *strings.Builder {
	b.WriteString("<pre><code>{")
	for n := range fields.Len() {
		if fields.Get(n).Message() != nil {
			messageFields := fields.Get(n).Message().Fields()
			if messageFields != nil {
				b.WriteString("\n    \"" + fields.Get(n).JSONName() + "\": {")
				for m := range messageFields.Len() {
					b.WriteString("\n        \"" + messageFields.Get(m).JSONName() + "\": " + messageFields.Get(m).Kind().String())
				}
				b.WriteString("\n    }")
			}
		} else {
			b.WriteString("\n    \"" + fields.Get(n).JSONName() + "\": " + fields.Get(n).Kind().String())
		}
	}
	b.WriteString("\n}</code></pre>")
	return b
}

func GetSnetServicesInfo(fileDescriptors map[string][]protoreflect.FileDescriptor) string {
	if fileDescriptors != nil {
		b := writeServiceSnetIDs(fileDescriptors)
		return b.String()
	}
	return ""
}
