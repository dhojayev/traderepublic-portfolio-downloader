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

type TransactionResponseNormalizer struct {
	logger *log.Logger
}

func NewTransactionResponseNormalizer(logger *log.Logger) TransactionResponseNormalizer {
	return TransactionResponseNormalizer{
		logger: logger,
	}
}

func (n TransactionResponseNormalizer) Normalize(response Response) (NormalizedResponse, error) {
	var err error

	resp := NormalizedResponse{
		ID: response.ID,
	}

	resp.Header, err = n.SectionTypeHeader(response)
	if err != nil {
		// header is crucial for all transactions, therefore we cannot simply raise a warning but have to fail.
		return resp, fmt.Errorf("could not deserialize header section: %w", err)
	}

	tableSections, err := n.SectionsTypeTable(response)
	if err != nil {
		n.logger.Warnf("could not deserialize table sections: %v", err)
	}

	for _, tableSection := range tableSections {
		switch tableSection.Title {
		case SectionTitleOverview:
			resp.Overview = NormalizedResponseOverviewSection{tableSection}
		case SectionTitlePerformance:
			resp.Performance = NormalizedResponsePerformanceSection{tableSection}
		case SectionTitleTransaction, SectionTitleTransactionAlt:
			resp.Transaction = NormalizedResponseTransactionSection{tableSection}
		case SectionTitleSavingPlan:
		default:
			n.logger.Warnf("unknown section title: %v", tableSection.Title)
		}
	}

	resp.Documents, err = n.SectionTypeDocuments(response)
	if err != nil {
		n.logger.Debugf("could not deserialize documents section: %v", err)
	}

	return resp, nil
}

func (n TransactionResponseNormalizer) SectionTypeHeader(response Response) (NormalizedResponseHeaderSection, error) {
	var sections []NormalizedResponseHeaderSection

	if err := n.deserializeSections(response, &sections, ResponseSectionTypeValueHeader); err != nil {
		return NormalizedResponseHeaderSection{}, err
	}

	if len(sections) == 0 {
		return NormalizedResponseHeaderSection{}, ErrSectionTypeNotFound
	}

	return sections[0], nil
}

func (n TransactionResponseNormalizer) SectionsTypeTable(response Response) ([]NormalizedResponseTableSection, error) {
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

func (n TransactionResponseNormalizer) SectionTypeDocuments(
	response Response,
) (NormalizedResponseDocumentsSection, error) {
	var sections []NormalizedResponseDocumentsSection

	if err := n.deserializeSections(response, &sections, ResponseSectionTypeValueDocuments); err != nil {
		return NormalizedResponseDocumentsSection{}, err
	}

	if len(sections) == 0 {
		return NormalizedResponseDocumentsSection{}, ErrSectionTypeNotFound
	}

	return sections[0], nil
}

func (n TransactionResponseNormalizer) deserializeSections(response Response, v any, sectionTypes ...string) error {
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

		if !slices.Contains(sectionTypes, sectionTypeString) {
			continue
		}

		sections = append(sections, section)
	}

	if len(sections) == 0 {
		return fmt.Errorf("section types '%v' were not found: %w", sectionTypes, ErrSectionTypeNotFound)
	}

	sectionsBytes, err := json.Marshal(sections)
	if err != nil {
		return fmt.Errorf("could not marshal %s section: %w", sectionTypes, err)
	}

	if err = json.Unmarshal(sectionsBytes, v); err != nil {
		return fmt.Errorf("could not unmarshal %s section: %w", sectionTypes, err)
	}

	return nil
}

type ActivityLogResponseNormalizer struct {
	TransactionResponseNormalizer
}

func NewActivityResponseNormalizer(logger *log.Logger) ActivityLogResponseNormalizer {
	return ActivityLogResponseNormalizer{TransactionResponseNormalizer: NewTransactionResponseNormalizer(logger)}
}

func (n ActivityLogResponseNormalizer) Normalize(response Response) (NormalizedResponse, error) {
	var err error

	resp := NormalizedResponse{
		ID: response.ID,
	}

	resp.Header, err = n.SectionTypeHeader(response)
	if err != nil {
		return resp, fmt.Errorf("could not deserialize header section: %w", err)
	}

	resp.Documents, err = n.SectionTypeDocuments(response)
	if err != nil {
		return resp, fmt.Errorf("could not deserialize documents section: %w", err)
	}

	return resp, nil
}
