package details

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"

	log "github.com/sirupsen/logrus"
)

var (
	ErrSectionTypeNotFound   = errors.New("section type not found")
	ErrSectionContainsNoType = errors.New("section contains no type")
)

type ResponseNormalizerInterface interface {
	Normalize(response Response) (NormalizedResponse, error)
}

type ResponseNormalizer struct {
	logger *log.Logger
}

func NewResponseNormalizer(logger *log.Logger) ResponseNormalizer {
	return ResponseNormalizer{
		logger: logger,
	}
}

func (n ResponseNormalizer) Normalize(response Response) (NormalizedResponse, error) {
	var err error

	resp := NormalizedResponse{
		ID: response.ID,
	}

	resp.Header, err = n.SectionTypeHeader(response)
	if err != nil {
		// header is crucial for all transactions, therefore we cannot simply raise a warning but have to fail.
		return resp, fmt.Errorf("could not deserialize section type header: %w", err)
	}

	tableSections, err := n.SectionsTypeTable(response)
	if err != nil {
		n.logger.Warnf("could not deserialize table sections: %v", err)
	}

	for _, tableSection := range tableSections {
		switch tableSection.Title {
		case SectionTitleOverview:
			resp.Overview = &NormalizedResponseOverviewSection{tableSection}
		case SectionTitlePerformance:
			resp.Performance = &NormalizedResponsePerformanceSection{tableSection}
		case SectionTitleTransaction:
			resp.Transaction = &NormalizedResponseTransactionSection{tableSection}
		case SectionTitleBusiness:
			resp.Security = &NormalizedResponseSecuritySection{tableSection}
		case SectionTitleSavingPlan:
		default:
			n.logger.Warnf("unknown section type: %v", tableSection.Title)
		}
	}

	resp.Documents, err = n.SectionTypeDocuments(response)
	if err != nil {
		n.logger.Warnf("could not deserialize documents: %v", err)
	}

	return resp, nil
}

func (n ResponseNormalizer) SectionTypeHeader(response Response) (*NormalizedResponseHeaderSection, error) {
	var sections []NormalizedResponseHeaderSection

	if err := n.deserializeSections(response, &sections, ResponseSectionTypeValueHeader); err != nil {
		return nil, err
	}

	if len(sections) == 0 {
		return nil, ErrSectionTypeNotFound
	}

	return &sections[0], nil
}

func (n ResponseNormalizer) SectionsTypeTable(response Response) ([]NormalizedResponseTableSection, error) {
	var sections []NormalizedResponseTableSection

	if err := n.deserializeSections(
		response,
		&sections,
		ResponseSectionTypeValueTable,
		ResponseSectionTypeValueHorizontalTable,
	); err != nil {
		return sections, err
	}

	return sections, nil
}

func (n ResponseNormalizer) SectionTypeDocuments(response Response) (*NormalizedResponseDocumentsSection, error) {
	var sections []NormalizedResponseDocumentsSection

	if err := n.deserializeSections(response, &sections, ResponseSectionTypeValueDocuments); err != nil {
		return nil, err
	}

	if len(sections) == 0 {
		return nil, ErrSectionTypeNotFound
	}

	return &sections[0], nil
}

func (n ResponseNormalizer) deserializeSections(response Response, v any, needles ...string) error {
	sections := make([]map[string]any, 0)

	for _, section := range response.Sections {
		sectionType, found := section["type"]
		if !found {
			return ErrSectionContainsNoType
		}

		sectionTypeString, valid := sectionType.(string)
		if !valid {
			return fmt.Errorf("section type is not a string: %v", sectionType)
		}

		if !slices.Contains(needles, sectionTypeString) {
			continue
		}

		sections = append(sections, section)
	}

	if len(sections) == 0 {
		return ErrSectionTypeNotFound
	}

	sectionsBytes, err := json.Marshal(sections)
	if err != nil {
		return fmt.Errorf("could not marshal %s section: %w", needles, err)
	}

	if err = json.Unmarshal(sectionsBytes, v); err != nil {
		return fmt.Errorf("could not unmarshal %s section: %w", needles, err)
	}

	return nil
}
