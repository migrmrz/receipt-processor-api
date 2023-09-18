package service

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/sirupsen/logrus"

	"fetch-rewards.com/mx/receipt-processor/pkg/handlers/rest"
)

func (s *Service) GetPoints(id string) (int, []string) {
	var currReceipt rest.ProcessReceiptRequest

	// look for received id in receipts list and store in currReceipt variable for further processing
	for _, receipt := range s.receipts {
		if id == receipt.ID {
			currReceipt.ID = receipt.ID
			currReceipt.Retailer = receipt.Retailer
			currReceipt.PurchaseDate = receipt.PurchaseDate
			currReceipt.PurchaseTime = receipt.PurchaseTime
			currReceipt.Total = receipt.Total
			currReceipt.Items = receipt.Items
		}
	}

	// calculate points according to each defined rule
	return runCalculations(currReceipt)
}

func runCalculations(receipt rest.ProcessReceiptRequest) (points int, errors []string) {
	var validations = []Validation{
		retailerValidation,
		totalAmountValidation,
		itemsValidation,
		itemDescriptionsValidation,
		dateTimeValidations,
	}

	for _, validation := range validations {
		valPoints, valErrors := validation(receipt)
		points += valPoints
		errors = append(errors, valErrors...)
	}

	return points, errors
}

// Validation type that will help run validations iteratively

type Validation func(rest.ProcessReceiptRequest) (int, []string)

// validation functions
func retailerValidation(receipt rest.ProcessReceiptRequest) (points int, errors []string) {
	// --- 1 point for each alphanumeric character in retailer name
	// check for letters and numbers only (no spaces or special characters)
	for _, r := range receipt.Retailer {
		if unicode.IsDigit(r) || unicode.IsLetter(r) {
			points += 1
		}
	}

	return points, nil
}

func totalAmountValidation(receipt rest.ProcessReceiptRequest) (points int, errors []string) {
	// --- 50 points if the total is round amount with no cents
	floatTotal, totalConvErr := strconv.ParseFloat(receipt.Total, 32) // check for valid float number
	if totalConvErr != nil {
		errors = append(
			errors,
			fmt.Sprintf(
				"receipt total is not in a valid format value and won't be considered: %s",
				totalConvErr.Error(),
			),
		)
	} else {
		// split total and check for "00" or "0" after "."
		centsStr := strings.SplitAfter(receipt.Total, ".")[1]
		if centsStr == "00" || centsStr == "0" {
			points += 50
		}
	}

	// --- 25 points if the total is multiple of 0.25
	if totalConvErr == nil && math.Mod(math.Floor(floatTotal*100)/100, 0.25) == 0.0 {
		points += 25
	}

	return points, errors
}

func itemsValidation(receipt rest.ProcessReceiptRequest) (points int, errors []string) {
	// --- 5 points for every 2 items
	if len(receipt.Items) > 1 {
		points += (len(receipt.Items) / 2) * 5
	}

	return points, nil
}

func itemDescriptionsValidation(receipt rest.ProcessReceiptRequest) (points int, errors []string) {
	// --- points are: if the trimmed lenght of item description is
	// multiple of 3, multiply the price by 0.2 and round to nearest int
	for _, item := range receipt.Items {
		if len(strings.Trim(item.ShortDescription, " "))%3 == 0 {
			price, err := strconv.ParseFloat(item.Price, 32)
			if err != nil {
				errors = append(
					errors,
					fmt.Sprintf(
						"item price is not in a valid format value and won't be considered: %s",
						err.Error(),
					),
				)
			}

			points += int(math.Ceil((math.Floor(price*100) / 100) * 0.2))
		}
	}

	return points, errors
}

func dateTimeValidations(receipt rest.ProcessReceiptRequest) (points int, errors []string) {
	// get and convert dateTime values
	purDateTime, dateParseErr := time.Parse(
		"2006-01-02 15:04",
		fmt.Sprint(receipt.PurchaseDate, " ", receipt.PurchaseTime),
	)
	if dateParseErr != nil {
		errors = append(errors, fmt.Sprintf(
			"date and/or time values are not in a valid format value and won't be considered: %s",
			dateParseErr.Error(),
		))
		logrus.Errorf(
			"date and/or time values are not in a valid format value and won't be considered: %s",
			dateParseErr.Error(),
		)
	}

	// --- 6 points if the date of purchase is odd
	if dateParseErr == nil && purDateTime.Day()%2 == 1 {
		points += 6
	}

	// 10 points if the time of purchase is after 2 pm and before 4 pm
	if dateParseErr == nil && (purDateTime.Hour() >= 14 && purDateTime.Hour() < 16) {
		points += 10
	}

	return points, errors
}
