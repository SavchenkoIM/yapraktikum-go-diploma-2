package cli

import (
	"fmt"
	"golang.org/x/exp/maps"
	proto "passwordvault/internal/proto/gen"
	"slices"
	"strings"
)

// RecordType
type RecordType string

var acceptedDataTypes = map[string]proto.DataType{
	"any":  proto.DataType_UNSPECIFIED,
	"cred": proto.DataType_CREDENTIALS,
	"card": proto.DataType_CREDIT_CARD,
	"file": proto.DataType_BLOB,
	"note": proto.DataType_TEXT_NOTE,
}

func (d RecordType) GetType() (proto.DataType, error) {
	dd := strings.ToLower(string(d))
	knownKeys := maps.Keys(acceptedDataTypes)
	if slices.Contains(knownKeys, dd) {
		return acceptedDataTypes[dd], nil
	} else {
		errorMsg := strings.Builder{}
		errorMsg.WriteString(fmt.Sprintf("unknown data type: %s\n", dd))
		errorMsg.WriteString("Accepted data types are:\n")
		for k, _ := range acceptedDataTypes {
			errorMsg.WriteString(fmt.Sprintf("%s\n", k))
		}
		return proto.DataType_UNSPECIFIED, fmt.Errorf(errorMsg.String())
	}
}

// MetadataFilters
type MetadataFilters []string

func (m MetadataFilters) GetFilters() ([]*proto.MetaDataKV, error) {
	md := make([]*proto.MetaDataKV, 0)
	for _, filter := range m {
		filterParts := strings.Split(filter, "=")
		if len(filterParts) != 2 {
			return nil, fmt.Errorf(`filter must match "key=value" pattern`)
		}
		md = append(md, &proto.MetaDataKV{
			Name:  filterParts[0],
			Value: filterParts[1],
		})
	}

	return md, nil
}
