package agent

import "time"

// Event represents a scheduled event for an outlet
type Event struct {
	Date        time.Time
	Description string
	OutletID    string
}

// GlobalEvents is a hardcoded list of events for demonstration
// In a real app, this might come from a DB
var GlobalEvents = []Event{
	// January Events
	{
		Date:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Description: "New Year's Day - High demand for party supplies and groceries",
		OutletID:    "Outlet-1",
	},
	{
		Date:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Description: "New Year's Day - Closed for holiday, prepare for reopening",
		OutletID:    "Outlet-2",
	},
	{
		Date:        time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		Description: "First Weekend - Restock required for weekend shoppers",
		OutletID:    "Outlet-3",
	},
	{
		Date:        time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
		Description: "Monday Rush - Office supplies and electronics low",
		OutletID:    "Outlet-4",
	},
	{
		Date:        time.Date(2024, 1, 12, 0, 0, 0, 0, time.UTC),
		Description: "Mid-month Sale Preparation - Stock up for promotional event",
		OutletID:    "Outlet-1",
	},
	// February Events
	{
		Date:        time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
		Description: "Month Start - Inventory audit and restocking",
		OutletID:    "Outlet-2",
	},
	{
		Date:        time.Date(2024, 2, 14, 0, 0, 0, 0, time.UTC),
		Description: "Valentine's Day - High demand for special items",
		OutletID:    "Outlet-3",
	},
	// March Events
	{
		Date:        time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
		Description: "Spring Season Start - Seasonal inventory adjustment",
		OutletID:    "Outlet-4",
	},
	{
		Date:        time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC),
		Description: "Mid-month Promotion - Electronics clearance sale",
		OutletID:    "Outlet-1",
	},
	// April Events
	{
		Date:        time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC),
		Description: "Spring Sale - High demand for seasonal products",
		OutletID:    "Outlet-2",
	},
	{
		Date:        time.Date(2024, 4, 15, 0, 0, 0, 0, time.UTC),
		Description: "Tax Season End - Office supplies rush",
		OutletID:    "Outlet-3",
	},
	// May Events
	{
		Date:        time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
		Description: "May Day - Holiday preparation",
		OutletID:    "Outlet-4",
	},
	{
		Date:        time.Date(2024, 5, 15, 0, 0, 0, 0, time.UTC),
		Description: "Mid-May Promotion - Summer stock arrival",
		OutletID:    "Outlet-1",
	},
	// June Events
	{
		Date:        time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
		Description: "Summer Start - Seasonal inventory shift",
		OutletID:    "Outlet-2",
	},
	{
		Date:        time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC),
		Description: "Mid-June Sale - Electronics promotion",
		OutletID:    "Outlet-3",
	},
	// July Events
	{
		Date:        time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC),
		Description: "Independence Day - Holiday shopping surge",
		OutletID:    "Outlet-4",
	},
	{
		Date:        time.Date(2024, 7, 20, 0, 0, 0, 0, time.UTC),
		Description: "Summer Peak - High demand for groceries",
		OutletID:    "Outlet-1",
	},
	// August Events
	{
		Date:        time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC),
		Description: "Back to School - Supplies and electronics needed",
		OutletID:    "Outlet-2",
	},
	{
		Date:        time.Date(2024, 8, 15, 0, 0, 0, 0, time.UTC),
		Description: "Mid-August Clearance - Summer stock reduction",
		OutletID:    "Outlet-3",
	},
	// September Events
	{
		Date:        time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC),
		Description: "Fall Season Start - Autumn inventory preparation",
		OutletID:    "Outlet-4",
	},
	{
		Date:        time.Date(2024, 9, 15, 0, 0, 0, 0, time.UTC),
		Description: "Mid-September Sale - Electronics upgrade",
		OutletID:    "Outlet-1",
	},
	// October Events
	{
		Date:        time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
		Description: "October Fest - Special event preparation",
		OutletID:    "Outlet-2",
	},
	{
		Date:        time.Date(2024, 10, 31, 0, 0, 0, 0, time.UTC),
		Description: "Halloween - Party supplies and groceries surge",
		OutletID:    "Outlet-3",
	},
	// November Events
	{
		Date:        time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
		Description: "Holiday Season Prep - Early stocking",
		OutletID:    "Outlet-4",
	},
	{
		Date:        time.Date(2024, 11, 28, 0, 0, 0, 0, time.UTC),
		Description: "Black Friday - Massive electronics and grocery demand",
		OutletID:    "Outlet-1",
	},
	// December Events
	{
		Date:        time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
		Description: "Holiday Shopping - Peak season begins",
		OutletID:    "Outlet-2",
	},
	{
		Date:        time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC),
		Description: "Christmas Day - Post-holiday restocking needed",
		OutletID:    "Outlet-3",
	},
	{
		Date:        time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
		Description: "New Year's Eve - Party supplies final rush",
		OutletID:    "Outlet-4",
	},
}
