package service

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetpoints(t *testing.T) {
	cases := []struct {
		name              string
		ticketID          string
		expectedPoints    int
		expectedErrorsStr []string
	}{
		{
			name:              "get-points-target-ok",
			ticketID:          "77b341b0-24db-4188-9954-ff353073f09a",
			expectedPoints:    28,
			expectedErrorsStr: nil,
		},
		{
			name:              "get-points-m&m-corner-ok",
			ticketID:          "6dd1a55a-48bc-460d-8643-4d616358acd1",
			expectedPoints:    109,
			expectedErrorsStr: nil,
		},
		{
			name:           "get-points-walgreens-error",
			ticketID:       "801786e4-9515-40e7-8c31-8ce973526f5e",
			expectedPoints: 15,
			expectedErrorsStr: []string{
				"receipt total is not in a valid format value and won't be considered: " +
					"strconv.ParseFloat: parsing \"$2.65\": invalid syntax",
				"date and/or time values are not in a valid format value and won't be considered: " +
					"parsing time \"Sunday 2 January 2022 08:13\" as \"2006-01-02 15:04\": cannot " +
					"parse \"Sunday 2 January 2022 08:13\" as \"2006\"",
			},
		},
		{
			name:           "get-points-target-error",
			ticketID:       "cf912fd8-9727-4183-b6fa-c108b5af6c4c",
			expectedPoints: 31,
			expectedErrorsStr: []string{
				"date and/or time values are not in a valid format value and won't be considered: " +
					"parsing time \"2022-01-02 01:13 PM\": extra text: \" PM\"",
			},
		},
		{
			name:           "get-points-target-item-price-error",
			ticketID:       "afe8bad1-3538-49b3-8847-7dc45e6563b9",
			expectedPoints: 41,
			expectedErrorsStr: []string{
				"item price is not in a valid format value and won't be considered: " +
					"strconv.ParseFloat: parsing \"$1.25\": invalid syntax",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// init service
			srv := New()

			// process receipts from testdata first
			if strings.Contains(tc.name, "error") {
				for _, receipt := range testReceiptsError {
					_ = srv.ProcessReceipts(receipt)
				}
			} else {
				for _, receipt := range testReceiptsOk {
					_ = srv.ProcessReceipts(receipt)
				}
			}

			// call actual function to calculate points
			actualPoints, actualErrorsStr := srv.GetPoints(tc.ticketID)
			assert.Equal(t, tc.expectedPoints, actualPoints)
			assert.Equal(t, tc.expectedErrorsStr, actualErrorsStr)
		})
	}
}
