package dto

import (
	"time"

	"github.com/tahmooress/weConnect-task/internal/entity"
)

type Statistics struct {
	ID              string    `json:"_id, omitempty"`
	SeriesReference string    `json:"series_reference,omitempty"`
	Period          string    `json:"period,omitempty"`
	DataValue       float64   `json:"data_value,omitempty"`
	Suppressed      string    `json:"suppressed,omitempty"`
	Status          string    `json:"status,omitempty"`
	Units           string    `json:"Units,omitempty"`
	Magnitude       int       `json:"magnitude,omitempty"`
	Subject         string    `json:"subject,omitempty"`
	Group           string    `json:"group,omitempty"`
	SeriesTitle1    string    `json:"Series_title_1,omitempty"`
	SeriesTitle2    string    `json:"Series_title_2,omitempty"`
	SeriesTitle3    string    `json:"Series_title_3,omitempty"`
	SeriesTitle4    string    `json:"Series_title_4,omitempty"`
	SeriesTitle5    string    `json:"Series_title_5,omitempty"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
}

func EntityToDto(s entity.Statistics) Statistics {
	return Statistics{
		ID:              s.ID.Hex(),
		SeriesReference: s.SeriesReference,
		Period:          s.Period,
		DataValue:       s.DataValue,
		Suppressed:      s.Suppressed,
		Status:          s.Status,
		Units:           s.Units,
		Magnitude:       s.Magnitude,
		Subject:         s.Subject,
		Group:           s.Group,
		SeriesTitle1:    s.SeriesTitle1,
		SeriesTitle2:    s.SeriesTitle2,
		SeriesTitle3:    s.SeriesTitle3,
		SeriesTitle4:    s.SeriesTitle4,
		SeriesTitle5:    s.SeriesTitle5,
		CreatedAt:       s.CreatedAt,
	}
}

func DtoToEntity(s Statistics) entity.Statistics {
	return entity.Statistics{
		SeriesReference: s.SeriesReference,
		Period:          s.Period,
		DataValue:       s.DataValue,
		Suppressed:      s.Suppressed,
		Status:          s.Status,
		Units:           s.Units,
		Magnitude:       s.Magnitude,
		Subject:         s.Subject,
		Group:           s.Group,
		SeriesTitle1:    s.SeriesTitle1,
		SeriesTitle2:    s.SeriesTitle2,
		SeriesTitle3:    s.SeriesTitle3,
		SeriesTitle4:    s.SeriesTitle4,
		SeriesTitle5:    s.SeriesTitle5,
		CreatedAt:       s.CreatedAt,
	}
}
