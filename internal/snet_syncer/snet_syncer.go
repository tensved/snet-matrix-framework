package snet_syncer

import (
	"context"
	"encoding/json"
	"github.com/bufbuild/protocompile"
	"github.com/rs/zerolog/log"
	"github.com/tensved/snet-matrix-framework/pkg/blockchain"
	"github.com/tensved/snet-matrix-framework/pkg/db"
	ipfs "github.com/tensved/snet-matrix-framework/pkg/ipfs"
	"google.golang.org/protobuf/reflect/protoreflect"
	"strings"
	"time"
)

type SnetSyncer struct {
	Ethereum        blockchain.Ethereum
	IPFSClient      ipfs.IPFSClient
	DB              db.Service
	FileDescriptors map[string][]protoreflect.FileDescriptor
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
	log.Info().Msg("SnetSyncer now working...")

	orgs, _ := s.Ethereum.GetOrgs()
	for _, orgIDBytes := range orgs {
		borg, err := s.Ethereum.GetOrg(orgIDBytes)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get org")
			continue
		}
		var org blockchain.OrganizationMetaData

		metadataJson, err := s.IPFSClient.GetIpfsFile(string(borg.OrgMetadataURI))
		if err != nil {
			log.Error().Err(err).Msg("Failed to get ipfs file")
			continue
		}

		err = json.Unmarshal(metadataJson, &org)
		if err != nil {
			log.Error().Err(err).Any("content", string(metadataJson)).Msg("Can't unmarshal org metadata from ipfs")
			continue
		}

		org.Owner = borg.Owner.Hex()
		org.SnetID = strings.ReplaceAll(string(borg.Id[:]), "\u0000", "")
		dbOrg, dbGroups := org.DB()
		orgID, err := s.DB.CreateSnetOrg(dbOrg)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create org")
		}
		org.ID = orgID
		err = s.DB.CreateSnetOrgGroups(orgID, dbGroups)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create org group")
		}

		for _, serviceIDBytes := range borg.ServiceIds {
			service, err := s.Ethereum.GetService(borg.Id, serviceIDBytes)
			if err != nil {
				log.Error().Err(err)
				continue
			}

			metadataJson, err = s.IPFSClient.GetIpfsFile(string(service.MetadataURI))
			if err != nil {
				log.Error().Err(err).Msg("Failed to get file from ipfs")
				return
			}

			var srvMeta blockchain.ServiceMetadata
			err = json.Unmarshal(metadataJson, &srvMeta)
			if err != nil {
				log.Error().Err(err).Any("content", string(metadataJson)).Msg("Failed to unmarshal metadata from ipfs")
				return
			}

			log.Debug().Msgf("Metadata of service: %+v", srvMeta)

			srvMeta.OrgID = orgID
			srvMeta.SnetID = strings.ReplaceAll(string(serviceIDBytes[:]), "\u0000", "")
			srvMeta.SnetOrgID = org.SnetID
			srvMeta.ID, err = s.DB.CreateSnetService(srvMeta.DB())
			if err != nil {
				log.Error().Err(err).Int("id", srvMeta.ID).Str("snet-id", srvMeta.SnetID).Msg("Failed to add snet_service")
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
				s.FileDescriptors[srvMeta.SnetID] = append(s.FileDescriptors[srvMeta.SnetID], fd)
			}
		}
	}
}

func (s *SnetSyncer) Start() {
	log.Info().Msg("SnetSyncer ticker started")
	ticker := time.NewTicker(100 * time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.SyncOnce()
		}
	}
}

func getFileDescriptor(protoContent, name string) (ds protoreflect.FileDescriptor) {
	accessor := protocompile.SourceAccessorFromMap(map[string]string{
		name: protoContent,
	})
	compiler := protocompile.Compiler{
		Resolver:       &protocompile.SourceResolver{Accessor: accessor},
		SourceInfoMode: protocompile.SourceInfoStandard,
	}
	fds, err := compiler.Compile(context.Background(), name)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create file descriptor")
		return
	}
	ds = fds.FindFileByPath(name)
	return
}

func writeServiceSnetIDs(builder strings.Builder, fileDescriptors map[string][]protoreflect.FileDescriptor) strings.Builder {
	builder.WriteString("<div style=\"line-height: 0.8;\"><ol>")
	for snetID, descriptors := range fileDescriptors {
		if descriptors != nil {
			for _, descriptor := range descriptors {
				if descriptor != nil {
					builder.WriteString("<li><strong>Path: " + descriptor.Path() + " Snet ID: " + snetID + " Descriptor: " + string(descriptor.FullName().Name()) + "</strong></li>")
					services := descriptor.Services()
					if services != nil {
						builder = writeServiceNames(builder, services)
					}
				}
			}
		}
	}
	builder.WriteString("</ol></div>")
	return builder
}

func writeServiceNames(builder strings.Builder, services protoreflect.ServiceDescriptors) strings.Builder {
	for i := 0; i < services.Len(); i++ {
		if services.Get(i) != nil {
			builder.WriteString("<p><em>Service: " + string(services.Get(i).FullName().Name()) + "</em></p>")
			methods := services.Get(i).Methods()
			if methods != nil {
				builder = writeMethodNames(builder, methods)
			}
		}
	}
	return builder
}

func writeMethodNames(builder strings.Builder, methods protoreflect.MethodDescriptors) strings.Builder {
	builder.WriteString("<p>🔁Methods: </p><ul>")
	for j := 0; j < methods.Len(); j++ {
		if methods.Get(j) != nil {
			builder.WriteString("<li>" + string(methods.Get(j).FullName().Name()) + "<br>")
			inputFields := methods.Get(j).Input().Fields()
			outputFields := methods.Get(j).Output().Fields()

			if inputFields != nil {
				builder.WriteString("<p>➡️Input:</p>")
				builder = writeFields(builder, inputFields)
			}
			if outputFields != nil {
				builder.WriteString("<p>➡️Output:</p>")
				builder = writeFields(builder, outputFields)
			}
			builder.WriteString("</li>")
		}
	}
	builder.WriteString("</ul>")
	return builder
}

func writeFields(builder strings.Builder, fields protoreflect.FieldDescriptors) strings.Builder {
	builder.WriteString("<pre><code>{")
	for n := 0; n < fields.Len(); n++ {
		if fields.Get(n).Message() != nil {
			messageFields := fields.Get(n).Message().Fields()
			if messageFields != nil {
				builder.WriteString("\n    \"" + fields.Get(n).JSONName() + "\": {")
				for m := 0; m < messageFields.Len(); m++ {
					builder.WriteString("\n        \"" + messageFields.Get(m).JSONName() + "\": " + messageFields.Get(m).Kind().String())
				}
				builder.WriteString("\n    }")
			}
		} else {
			builder.WriteString("\n    \"" + fields.Get(n).JSONName() + "\": " + fields.Get(n).Kind().String())
		}
	}
	builder.WriteString("\n}</code></pre>")
	return builder
}

func writeOutputs(builder strings.Builder, inputFields protoreflect.FieldDescriptors) strings.Builder {
	builder.WriteString("<p>➡️Input:</p>")
	builder.WriteString("<pre><code>{")
	for n := 0; n < inputFields.Len(); n++ {
		if inputFields.Get(n).Message() != nil {
			messageFields := inputFields.Get(n).Message().Fields()
			if messageFields != nil {
				builder.WriteString("\n    \"" + inputFields.Get(n).JSONName() + "\": {")
				for m := 0; m < messageFields.Len(); m++ {
					builder.WriteString("\n        \"" + messageFields.Get(m).JSONName() + "\": " + messageFields.Get(m).Kind().String())
				}
				builder.WriteString("\n    }")
			}
		} else {
			builder.WriteString("\n    \"" + inputFields.Get(n).JSONName() + "\": " + inputFields.Get(n).Kind().String())
		}
	}
	builder.WriteString("\n}</code></pre>")
	return builder
}

func (s *SnetSyncer) GetSnetServicesInfo() string {
	var builder strings.Builder
	if s.FileDescriptors != nil {
		builder = writeServiceSnetIDs(builder, s.FileDescriptors)
	}

	return builder.String()
}
