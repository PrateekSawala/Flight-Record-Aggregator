package main

import (
	"fmt"
	"time"

	"github.com/PrateekSawala/flight-ticket-aggregator/domain"
	"github.com/nyaruka/phonenumbers"
)

func IsPNR(PNRs map[string]bool, value string) error {
	if len(value) < domain.FlightRecordPNREntryLength || !domain.AlphanumericRegex.MatchString(value) {
		return domain.ErrInvalidPNR
	}
	if PNRs[value] {
		return domain.ErrPNRRepeated
	}
	return nil
}

func IsMailValid(email string) error {
	if len(email) < domain.FlightRecordEmailEntryMinimumLength || len(email) > domain.FlightRecordEmailEntryMaximumLength || !domain.EmailRegex.MatchString(email) {
		return domain.ErrInvalidMail
	}
	return nil
}

func IsPhoneValid(phonenumber string) error {
	phonenumber = fmt.Sprintf("%s%s", domain.PhoneDialCode, phonenumber)
	phonenumberResponse, phoneParseError := phonenumbers.Parse(phonenumber, domain.Empty)
	if phoneParseError != nil || !phonenumbers.IsValidNumber(phonenumberResponse) {
		return domain.ErrInvalidPhone
	}
	return nil
}

func IsTicketingDateBeforeTravelDate(ticketingDate string, travelDate string) error {
	bookingDate, err := time.Parse(domain.FlightRecordTimeFormat, ticketingDate)
	if err != nil {
		return domain.ErrTicketingDate
	}
	tripDate, err := time.Parse(domain.FlightRecordTimeFormat, travelDate)
	if err != nil {
		return domain.ErrTravelDate
	}
	if !bookingDate.Before(tripDate) {
		return domain.ErrInvalidBooking
	}
	return nil
}

func IsCabinValid(cabin string) error {
	if !domain.ValidFlightCabins[cabin] {
		return domain.ErrInvalidCabin
	}
	return nil
}

func IsRecordsColumnValid(flightRecord []string) error {
	if len(flightRecord) != domain.FlightRecordEntriesLength {
		return domain.ErrInvalidFile
	}
	for index, record := range flightRecord {
		if domain.ValidRecordIndexes[index] != record {
			return domain.ErrInvalidFile
		}
	}
	return nil
}

func IsRecordValid(recordedPNRs map[string]bool, flightRecord domain.FlightRecord) error {
	if err := IsMailValid(flightRecord.Email); err != nil {
		return err
	}
	if err := IsPhoneValid(flightRecord.MobilePhone); err != nil {
		return err
	}
	if err := IsTicketingDateBeforeTravelDate(flightRecord.TicketingDate, flightRecord.TravelDate); err != nil {
		return err
	}
	if err := IsPNR(recordedPNRs, flightRecord.PNR); err != nil {
		return err
	}
	if err := IsCabinValid(flightRecord.BookedCabin); err != nil {
		return err
	}
	return nil
}
