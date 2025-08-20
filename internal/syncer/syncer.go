package syncer

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/bufbuild/protocompile"
	"github.com/rs/zerolog/log"
	"github.com/tensved/snet-matrix-framework/pkg/blockchain"
	"github.com/tensved/snet-matrix-framework/pkg/db"
	ipfs "github.com/tensved/snet-matrix-framework/pkg/ipfs"
	"google.golang.org/protobuf/reflect/protoreflect"
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
	logger := log.With().Logger()

	logger.Info().Msg("starting SNET synchronization")

	startTime := time.Now()

	// Clear file descriptors to prevent duplicates
	s.FileDescriptors = make(map[string][]protoreflect.FileDescriptor)

	orgs, err := s.Ethereum.GetOrgs()
	if err != nil {
		logger.Error().Err(err).Msg("failed to get organizations from blockchain")
		return
	}

	logger.Info().
		Int("organizations_count", len(orgs)).
		Msg("found organizations in blockchain")

	processedOrgs := 0
	processedServices := 0

	for i, orgIDBytes := range orgs {
		orgIDStr := strings.ReplaceAll(string(orgIDBytes[:]), "\u0000", "")

		logger.Info().
			Int("org_index", i).
			Str("org_id", orgIDStr).
			Msg("processing organization")

		borg, err := s.Ethereum.GetOrg(orgIDBytes)
		if err != nil {
			logger.Error().
				Err(err).
				Str("org_id", orgIDStr).
				Int("org_index", i).
				Msg("failed to get organization from blockchain")
			continue
		}

		if !borg.Found {
			logger.Warn().
				Str("org_id", orgIDStr).
				Int("org_index", i).
				Msg("organization not found in blockchain")
			continue
		}

		var org blockchain.OrganizationMetaData

		if len(borg.OrgMetadataURI) == 0 {
			logger.Warn().
				Str("org_id", orgIDStr).
				Int("org_index", i).
				Msg("organization has no metadata URI")
			continue
		}

		if s.IPFSClient.HttpApi == nil {
			logger.Error().Msg("IPFSClient.HttpApi is nil")
			return
		}

		logger.Debug().
			Str("org_id", orgIDStr).
			Str("metadata_uri", string(borg.OrgMetadataURI)).
			Msg("fetching organization metadata from IPFS")

		metadataJSON, err := s.IPFSClient.GetIpfsFile(string(borg.OrgMetadataURI))
		if err != nil {
			logger.Error().
				Err(err).
				Str("org_id", orgIDStr).
				Int("org_index", i).
				Str("metadata_uri", string(borg.OrgMetadataURI)).
				Msg("failed to get organization metadata from IPFS")
			continue
		}

		err = json.Unmarshal(metadataJSON, &org)
		if err != nil {
			logger.Error().
				Err(err).
				Str("org_id", orgIDStr).
				Int("org_index", i).
				Msg("failed to unmarshal organization metadata from IPFS")
			continue
		}

		org.Owner = borg.Owner.Hex()
		org.SnetID = strings.ReplaceAll(string(borg.Id[:]), "\u0000", "")

		if s.DB == nil {
			logger.Error().Msg("DB is nil")
			return
		}

		logger.Debug().
			Str("org_id", orgIDStr).
			Str("org_name", org.OrgName).
			Msg("saving organization to database")

		dbOrg, dbGroups := org.DB()
		orgID, err := s.DB.CreateSnetOrg(dbOrg)
		if err != nil {
			logger.Error().
				Err(err).
				Str("org_id", orgIDStr).
				Int("org_index", i).
				Msg("failed to create organization in database")
		}
		org.ID = orgID
		err = s.DB.CreateSnetOrgGroups(orgID, dbGroups)
		if err != nil {
			logger.Error().
				Err(err).
				Str("org_id", orgIDStr).
				Int("org_index", i).
				Msg("failed to create organization groups in database")
		}

		logger.Info().
			Str("org_id", orgIDStr).
			Int("org_index", i).
			Int("services_count", len(borg.ServiceIds)).
			Msg("processing organization services")

		processedOrgs++

		var service blockchain.Service
		for j, serviceIDBytes := range borg.ServiceIds {
			serviceIDStr := strings.ReplaceAll(string(serviceIDBytes[:]), "\u0000", "")

			logger.Info().
				Int("org_index", i).
				Int("service_index", j).
				Str("org_id", orgIDStr).
				Str("service_id", serviceIDStr).
				Msg("processing service")

			service, err = s.Ethereum.GetService(borg.Id, serviceIDBytes)
			if err != nil {
				logger.Error().
					Err(err).
					Str("org_id", orgIDStr).
					Str("service_id", serviceIDStr).
					Int("org_index", i).
					Int("service_index", j).
					Msg("failed to get service from blockchain")
				continue
			}

			logger.Debug().
				Str("org_id", orgIDStr).
				Str("service_id", serviceIDStr).
				Str("metadata_uri", string(service.MetadataURI)).
				Msg("fetching service metadata from IPFS")

			metadataJSON, err = s.IPFSClient.GetIpfsFile(string(service.MetadataURI))
			if err != nil {
				logger.Error().
					Err(err).
					Str("org_id", orgIDStr).
					Str("service_id", serviceIDStr).
					Int("org_index", i).
					Int("service_index", j).
					Str("metadata_uri", string(service.MetadataURI)).
					Msg("failed to get service metadata from IPFS")
				continue
			}

			var srvMeta blockchain.ServiceMetadata
			err = json.Unmarshal(metadataJSON, &srvMeta)
			if err != nil {
				logger.Error().
					Err(err).
					Str("org_id", orgIDStr).
					Str("service_id", serviceIDStr).
					Int("org_index", i).
					Int("service_index", j).
					Msg("failed to unmarshal service metadata from IPFS")
				continue
			}

			logger.Info().
				Str("org_id", orgIDStr).
				Str("service_id", serviceIDStr).
				Int("org_index", i).
				Int("service_index", j).
				Interface("metadata", srvMeta).
				Msg("retrieved service metadata")
			if len(srvMeta.Groups) > 0 {
				logger.Info().
					Str("org_id", orgIDStr).
					Str("service_id", serviceIDStr).
					Interface("endpoints", srvMeta.Groups[0].Endpoints).
					Msg("service endpoints")
				logger.Info().
					Str("org_id", orgIDStr).
					Str("service_id", serviceIDStr).
					Interface("daemon_addresses", srvMeta.Groups[0].DaemonAddresses).
					Msg("service daemon addresses")
			}

			srvMeta.OrgID = orgID
			srvMeta.SnetID = strings.ReplaceAll(string(serviceIDBytes[:]), "\u0000", "")
			srvMeta.SnetOrgID = org.SnetID
			dbSrvMeta, err := srvMeta.DB()
			if err != nil {
				logger.Error().Err(err)
			}
			srvMeta.ID, err = s.DB.CreateSnetService(dbSrvMeta)
			if err != nil {
				logger.Error().
					Err(err).
					Int("id", srvMeta.ID).
					Str("snet-id", srvMeta.SnetID).
					Msg("failed to add snet_service")
			}

			// Try ServiceApiSource first, then ModelIpfsHash, then both if available
			var protoHashes []string
			if srvMeta.ServiceApiSource != "" {
				protoHashes = append(protoHashes, srvMeta.ServiceApiSource)
			}
			if srvMeta.ModelIpfsHash != "" {
				protoHashes = append(protoHashes, srvMeta.ModelIpfsHash)
			}

			if len(protoHashes) == 0 {
				logger.Error().
					Str("snet_id", srvMeta.SnetID).
					Msg("both ModelIpfsHash and ServiceApiSource are empty")
				continue
			}

			logger.Info().
				Str("snet_id", srvMeta.SnetID).
				Str("model_ipfs_hash", srvMeta.ModelIpfsHash).
				Msg("model IPFS hash")
			logger.Info().
				Str("snet_id", srvMeta.SnetID).
				Str("service_api_source", srvMeta.ServiceApiSource).
				Msg("service API source")

			// Try each hash until one works
			var content []byte
			var successfulHash string
			var protoErr error

			for _, protoHash := range protoHashes {
				logger.Info().
					Str("snet_id", srvMeta.SnetID).
					Str("hash", protoHash).
					Msg("trying to get proto files from IPFS")
				content, protoErr = s.IPFSClient.GetIpfsFile(protoHash)
				if protoErr == nil {
					successfulHash = protoHash
					logger.Info().
						Str("snet_id", srvMeta.SnetID).
						Str("successful_hash", successfulHash).
						Msg("successfully got proto files from IPFS")
					break
				}
				logger.Error().
					Err(protoErr).
					Str("hash", protoHash).
					Str("snet_id", srvMeta.SnetID).
					Msg("failed to get proto files from IPFS, trying next hash")
			}

			if protoErr != nil {
				logger.Error().
					Err(protoErr).
					Strs("hashes", protoHashes).
					Str("snet_id", srvMeta.SnetID).
					Msg("failed to get proto files from all IPFS hashes")
				continue
			}

			logger.Info().
				Str("snet_id", srvMeta.SnetID).
				Int("content_size", len(content)).
				Msg("received content from IPFS")
			logger.Info().
				Str("snet_id", srvMeta.SnetID).
				Str("full_content", string(content)).
				Msg("FULL IPFS CONTENT")

			protoFiles, err := ipfs.ReadFilesCompressed(string(content))
			if err != nil {
				logger.Error().
					Err(err).
					Str("hash", successfulHash).
					Str("snet_id", srvMeta.SnetID).
					Msg("failed to read compressed proto files")
				continue
			}

			logger.Info().
				Str("snet_id", srvMeta.SnetID).
				Int("files_count", len(protoFiles)).
				Msg("extracted proto files count")
			for fileName, fileContent := range protoFiles {
				logger.Info().
					Str("snet_id", srvMeta.SnetID).
					Str("file_name", fileName).
					Int("file_size", len(fileContent)).
					Msg("proto file details")
			}

			protoFilesMap := make(map[string]string)
			for fileName, fileContent := range protoFiles {
				protoFilesMap[fileName] = string(fileContent)
			}

			// Add training.proto without importing google/protobuf/descriptor.proto
			trainingProtoContent := `syntax = "proto3";
package training;
option go_package = "github.com/singnet/snet-daemon/v5/training;training";
import "google/protobuf/descriptor.proto";

// Methods that the service provider must implement
service Model {

  // Free
  // Can pass the address of the model creator
  rpc create_model(NewModel) returns (ModelID) {}

  // Free
  rpc validate_model_price(ValidateRequest) returns (PriceInBaseUnit) {}

  // Paid
  rpc upload_and_validate(stream UploadInput) returns (StatusResponse) {}

  // Paid
  rpc validate_model(ValidateRequest) returns (StatusResponse) {}

  // Free, one signature for both train_model_price & train_model methods
  rpc train_model_price(ModelID) returns (PriceInBaseUnit) {}

  // Paid
  rpc train_model(ModelID) returns (StatusResponse) {}

  // Free
  rpc delete_model(ModelID) returns (StatusResponse) {
    // After model deletion, the status becomes DELETED in etcd
  }

  // Free
  rpc get_model_status(ModelID) returns (StatusResponse) {}
}

message ModelResponse {
  string model_id = 1;
  Status status = 2;
  string created_date = 3;
  string updated_date = 4;
  string name = 5;
  string description = 6;
  string grpc_method_name = 7;
  string grpc_service_name = 8;

  // List of all addresses that will have access to this model
  repeated string address_list = 9;

  // Access to the model is granted only for use and viewing
  bool is_public = 10;

  string training_data_link = 11;

  string created_by_address = 12;
  string updated_by_address = 13;
}

// Used as input for new_model requests
// The service provider decides whether to use these fields; returning model_id is mandatory
message NewModel {
  string name = 1;
  string description = 2;
  string grpc_method_name = 3;
  string grpc_service_name = 4;

  // List of all addresses that will have access to this model
  repeated string address_list = 5;

  // Set this to true if you want your model to be accessible by other AI consumers
  bool is_public = 6;

  // These parameters will be passed by the daemon
  string organization_id = 7;
  string service_id = 8;
  string group_id = 9;
}

// This structure must be used by the service provider
message ModelID {
  string model_id = 1;
}

// This structure must be used by the service provider
// Used in the train_model_price method to get the training/validation price
message PriceInBaseUnit {
  uint64 price = 1; // cogs, weis, afet, aasi, etc.
}

enum Status {
  CREATED = 0;
  VALIDATING = 1;
  VALIDATED = 2;
  TRAINING = 3;
  READY_TO_USE = 4; // After training is completed
  ERRORED = 5;
  DELETED = 6;
}

message StatusResponse {
  Status status = 1;
}

message UploadInput {
  string model_id = 1;
  bytes data = 2;
  string file_name = 3;
  uint64 file_size = 4; // in bytes
  uint64 batch_size = 5;
  uint64 batch_number = 6;
  uint64 batch_count = 7;
}

message ValidateRequest {
  string model_id = 2;
  string training_data_link = 3;
}

// Temporarily removed extensions to test compilation`

			protoFilesMap["training.proto"] = trainingProtoContent

			// Add full google/protobuf/descriptor.proto
			googleDescriptorProto := `syntax = "proto2";

package google.protobuf;

option go_package = "google.golang.org/protobuf/types/descriptorpb;descriptorpb";
option java_package = "com.google.protobuf";
option java_outer_classname = "DescriptorProtos";
option java_multiple_files = true;
option cc_enable_arenas = true;
option objc_class_prefix = "GPB";

message FileDescriptorSet {
  repeated FileDescriptorProto file = 1;
}

message FileDescriptorProto {
  optional string name = 1;
  optional string package = 2;
  repeated string dependency = 3;
  repeated int32 public_dependency = 10;
  repeated int32 weak_dependency = 11;
  repeated DescriptorProto message_type = 4;
  repeated EnumDescriptorProto enum_type = 5;
  repeated ServiceDescriptorProto service = 6;
  repeated FieldDescriptorProto extension = 7;
  optional FileOptions options = 8;
  optional SourceCodeInfo source_code_info = 9;
  optional string syntax = 12;
}

message DescriptorProto {
  optional string name = 1;
  repeated FieldDescriptorProto field = 2;
  repeated FieldDescriptorProto extension = 6;
  repeated DescriptorProto nested_type = 3;
  repeated EnumDescriptorProto enum_type = 4;
  repeated DescriptorProto.ExtensionRange extension_range = 5;
  repeated OneofDescriptorProto oneof_decl = 8;
  optional MessageOptions options = 7;
  repeated ReservedRange reserved_range = 9;
  repeated string reserved_name = 10;

  message ExtensionRange {
    optional int32 start = 1;
    optional int32 end = 2;
    optional ExtensionRangeOptions options = 3;
  }

  message ReservedRange {
    optional int32 start = 1;
    optional int32 end = 2;
  }
}

message ExtensionRangeOptions {
  // The parser stores options it doesn't recognize here. See above.
  repeated UninterpretedOption uninterpreted_option = 999;
}

message FieldDescriptorProto {
  enum Type {
    TYPE_DOUBLE = 1;
    TYPE_FLOAT = 2;
    TYPE_INT64 = 3;
    TYPE_UINT64 = 4;
    TYPE_INT32 = 5;
    TYPE_FIXED64 = 6;
    TYPE_FIXED32 = 7;
    TYPE_BOOL = 8;
    TYPE_STRING = 9;
    TYPE_GROUP = 10;
    TYPE_MESSAGE = 11;
    TYPE_BYTES = 12;
    TYPE_UINT32 = 13;
    TYPE_ENUM = 14;
    TYPE_SFIXED32 = 15;
    TYPE_SFIXED64 = 16;
    TYPE_SINT32 = 17;
    TYPE_SINT64 = 18;
  }

  enum Label {
    LABEL_OPTIONAL = 1;
    LABEL_REQUIRED = 2;
    LABEL_REPEATED = 3;
  }

  optional string name = 1;
  optional int32 number = 3;
  optional Label label = 4;
  optional Type type = 5;
  optional string type_name = 6;
  optional string extendee = 2;
  optional string default_value = 7;
  optional FieldOptions options = 8;
  optional string json_name = 10;
  optional int32 oneof_index = 9;
  optional bool proto3_optional = 17;
}

message OneofDescriptorProto {
  optional string name = 1;
  optional OneofOptions options = 2;
}

message EnumDescriptorProto {
  optional string name = 1;
  repeated EnumValueDescriptorProto value = 2;
  optional EnumOptions options = 3;
}

message EnumValueDescriptorProto {
  optional string name = 1;
  optional int32 number = 2;
  optional EnumValueOptions options = 3;
}

message ServiceDescriptorProto {
  optional string name = 1;
  repeated MethodDescriptorProto method = 2;
  optional ServiceOptions options = 3;
}

message MethodDescriptorProto {
  optional string name = 1;
  optional string input_type = 2;
  optional string output_type = 3;
  optional MethodOptions options = 4;
  optional bool client_streaming = 5;
  optional bool server_streaming = 6;
}

message FileOptions {
  optional string java_package = 1;
  optional string java_outer_classname = 8;
  optional bool java_multiple_files = 10;
  optional bool java_generate_equals_and_hash = 20;
  optional bool java_string_check_utf8 = 27;
  optional OptimizeMode optimize_for = 9;
  optional string go_package = 11;
  optional bool cc_generic_services = 16;
  optional bool java_generic_services = 17;
  optional bool py_generic_services = 18;
  optional bool php_generic_services = 42;
  optional bool deprecated = 23;
  optional bool cc_enable_arenas = 31;
  optional string objc_class_prefix = 36;
  repeated UninterpretedOption uninterpreted_option = 999;

  enum OptimizeMode {
    SPEED = 1;
    CODE_SIZE = 2;
    LITE_RUNTIME = 3;
  }
}

message MessageOptions {
  optional bool message_set_wire_format = 1;
  optional bool no_standard_descriptor_accessor = 2;
  optional bool deprecated = 3;
  optional bool map_entry = 7;
  repeated UninterpretedOption uninterpreted_option = 999;
}

message FieldOptions {
  optional CType ctype = 1;
  optional bool packed = 2;
  optional JSType jstype = 6;
  optional bool lazy = 5;
  optional bool deprecated = 3;
  optional bool weak = 10;
  repeated UninterpretedOption uninterpreted_option = 999;

  enum CType {
    STRING = 0;
    CORD = 1;
    STRING_PIECE = 2;
  }

  enum JSType {
    JS_NORMAL = 0;
    JS_STRING = 1;
    JS_NUMBER = 2;
  }
}

message OneofOptions {
  repeated UninterpretedOption uninterpreted_option = 999;
}

message EnumOptions {
  optional bool allow_alias = 2;
  optional bool deprecated = 3;
  repeated UninterpretedOption uninterpreted_option = 999;
}

message EnumValueOptions {
  optional bool deprecated = 1;
  repeated UninterpretedOption uninterpreted_option = 999;
}

message ServiceOptions {
  optional bool deprecated = 33;
  repeated UninterpretedOption uninterpreted_option = 999;
}

message MethodOptions {
  optional bool deprecated = 33;
  optional IdempotencyLevel idempotency_level = 34;
  repeated UninterpretedOption uninterpreted_option = 999;

  enum IdempotencyLevel {
    IDEMPOTENCY_UNKNOWN = 0;
    NO_SIDE_EFFECTS = 1;
    IDEMPOTENT = 2;
  }
}

message UninterpretedOption {
  repeated NamePart name = 2;
  optional string identifier_value = 3;
  optional uint64 positive_int_value = 4;
  optional int64 negative_int_value = 5;
  optional double double_value = 6;
  optional bytes string_value = 7;
  optional string aggregate_value = 8;
}

message NamePart {
  optional string name_part = 1;
  optional bool is_extension = 2;
}

message SourceCodeInfo {
  repeated Location location = 1;
}

message Location {
  repeated int32 path = 1;
  repeated int32 span = 2;
  optional string leading_comments = 3;
  optional string trailing_comments = 4;
  repeated string leading_detached_comments = 6;
}

message GeneratedCodeInfo {
  repeated Annotation annotation = 1;
}

message Annotation {
  repeated int32 path = 1;
  optional string source_file = 2;
  optional int32 begin = 3;
  optional int32 end = 4;
}`

			protoFilesMap["google/protobuf/descriptor.proto"] = googleDescriptorProto

			hasValidFileDescriptor := false
			for fileName, fileContent := range protoFiles {
				modifiedContent := string(fileContent)
				if fileName == "main.proto" {
					lines := strings.Split(modifiedContent, "\n")
					var filteredLines []string
					for _, line := range lines {
						if !strings.Contains(line, "option (training.") {
							filteredLines = append(filteredLines, line)
						}
					}
					modifiedContent = strings.Join(filteredLines, "\n")
				}

				tempProtoFilesMap := make(map[string]string)
				for k, v := range protoFilesMap {
					tempProtoFilesMap[k] = v
				}
				tempProtoFilesMap[fileName] = modifiedContent

				fd := getFileDescriptorWithDependencies(tempProtoFilesMap, fileName)
				if fd != nil {
					s.FileDescriptors[srvMeta.SnetID] = append(s.FileDescriptors[srvMeta.SnetID], fd)
					hasValidFileDescriptor = true
					err := os.WriteFile(fileName, fileContent, 0600)
					if err != nil {
						return
					}
				}
			}

			if hasValidFileDescriptor {
				logger.Info().
					Str("snet_id", srvMeta.SnetID).
					Msg("successfully created file descriptors")
			} else {
				logger.Warn().
					Str("snet_id", srvMeta.SnetID).
					Msg("failed to create file descriptors, but will still be saved to database")
			}
		}
	}

	logger.Info().
		Int("processed_organizations", processedOrgs).
		Int("processed_services", processedServices).
		Dur("duration", time.Since(startTime)).
		Msg("snet syncer successfully")
}

func (s *SnetSyncer) Start(ctx context.Context) {
	// Store the cancel function for later use in Stop
	ctx, cancel := context.WithCancel(ctx)
	s.cancelFunc = cancel

	duration := 24 * time.Hour
	// duration := 5 * time.Minute
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	log.Info().Msg("SNET syncer started successfully, waiting for sync interval")
	for {
		select {
		case <-ticker.C:
			log.Debug().Msg("sync interval triggered, starting sync process")
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

func getFileDescriptorWithDependencies(protoFiles map[string]string, name string) protoreflect.FileDescriptor {
	log.Info().Str("file_name", name).Msg("attempting to compile proto file")

	// Log available files for debugging
	for fileName := range protoFiles {
		log.Debug().Str("available_file", fileName).Msg("proto file available")
	}

	accessor := protocompile.SourceAccessorFromMap(protoFiles)
	compiler := protocompile.Compiler{
		Resolver:       &protocompile.SourceResolver{Accessor: accessor},
		SourceInfoMode: protocompile.SourceInfoStandard,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	fds, err := compiler.Compile(ctx, name)
	if err != nil {
		log.Error().Err(err).Str("file_name", name).Msg("failed to create file descriptor with dependencies")
		return nil
	}

	fd := fds.FindFileByPath(name)
	if fd != nil {
		log.Info().Str("file_name", name).Msg("successfully compiled proto file")
	} else {
		log.Warn().Str("file_name", name).Msg("compiled but could not find file descriptor")
	}
	return fd
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
