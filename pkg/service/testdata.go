package service

import "fetch-rewards.com/mx/receipt-processor/pkg/handlers/rest"

var testReceiptsOk = []rest.ProcessReceiptRequest{
	{
		ID:           "77b341b0-24db-4188-9954-ff353073f09a",
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Total:        "35.35",
		Items: []rest.Item{
			{
				ShortDescription: "Mountain Dew 12PK",
				Price:            "6.49",
			},
			{
				ShortDescription: "Emils Cheese Pizza",
				Price:            "12.25",
			},
			{
				ShortDescription: "Knorr Creamy Chicken",
				Price:            "1.26",
			},
			{
				ShortDescription: "Doritos Nacho Cheese",
				Price:            "3.35",
			},
			{
				ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ",
				Price:            "12.00",
			},
		},
	},
	{
		ID:           "6dd1a55a-48bc-460d-8643-4d616358acd1",
		Retailer:     "M&M Corner Market",
		PurchaseDate: "2022-03-20",
		PurchaseTime: "14:33",
		Total:        "9.00",
		Items: []rest.Item{
			{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
			{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
			{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
			{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
		},
	},
}

var testReceiptsError = []rest.ProcessReceiptRequest{
	{
		ID:           "801786e4-9515-40e7-8c31-8ce973526f5e",
		Retailer:     "Walgreens",
		PurchaseDate: "Sunday 2 January 2022",
		PurchaseTime: "08:13",
		Total:        "$2.65",
		Items: []rest.Item{
			{
				ShortDescription: "Pepsi - 12-oz",
				Price:            "1.25",
			},
			{
				ShortDescription: "Dasani",
				Price:            "1.40",
			},
		},
	},
	{
		ID:           "cf912fd8-9727-4183-b6fa-c108b5af6c4c",
		Retailer:     "Target",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "01:13 PM",
		Total:        "1.25",
		Items: []rest.Item{
			{
				ShortDescription: "Pepsi - 12-oz",
				Price:            "1.25",
			},
		},
	},
	{
		ID:           "afe8bad1-3538-49b3-8847-7dc45e6563b9",
		Retailer:     "Target",
		PurchaseDate: "2023-01-02",
		PurchaseTime: "14:13",
		Total:        "1.25",
		Items: []rest.Item{
			{
				ShortDescription: "Pepsi - 12oz",
				Price:            "$1.25",
			},
		},
	},
}
