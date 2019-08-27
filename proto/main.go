package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
)

// Define enum of record types
const (
	Debit = iota
	Credit
	StartAutopay
	EndAutopay
)

// Header defines the metadata head of txnlog.dat
type Header struct {
	MagicString [4]byte // string == "MPS7"
	Version     byte
	NRecords    uint32
}

// Autopay defines the structure of the StartAutopay and EndAutopay records
type Autopay struct {
	Timestamp uint32 // ([4]byte) unix time
	UserID    uint64 // ([8]byte)
}

// CreditDebit defines the structure of the Credit and Debit records
type CreditDebit struct {
	Timestamp uint32  // [4]byte unix time
	UserID    uint64  // [8]byte
	Amount    float64 // [8]byte
}

// TestUser is the user ID for the last question
const TestUser uint64 = 2456938384156277127

func main() {
	var (
		totalDebit, totalCredit, testUserDebit, testUserCredit float64
		autopayStartCount, autopayEndCount                     int
	)

	content, err := ioutil.ReadFile("txnlog.dat")
	if err != nil {
		log.Fatal("Error reading txnlog.dat:", err)
	}

	fileBuf := bytes.NewBuffer(content)
	header := readHeader(fileBuf)

	if header.MagicString != [4]byte{'M', 'P', 'S', '7'} {
		log.Fatal("Error reading txnlog.dat: Incorrect format")
	}

	var i uint32
	for i < header.NRecords {
		nextRecord, err := fileBuf.ReadByte()
		if err != nil {
			break
		}

		if nextRecord <= Credit {
			record := readCreditDebit(fileBuf)
			if nextRecord == Debit {
				totalDebit += record.Amount
				if record.UserID == TestUser {
					testUserDebit += record.Amount
				}
			} else {
				totalCredit += record.Amount
				if record.UserID == TestUser {
					testUserCredit += record.Amount
				}
			}
		} else {
			readAutopay(fileBuf)
			if nextRecord == StartAutopay {
				autopayStartCount++
			} else {
				autopayEndCount++
			}
		}

		i++
	}

	fmt.Println("What is the total amount of dollars of debits?\n\t", totalDebit)
	fmt.Println("What is the total amount of dollars of credits?\n\t", totalCredit)
	fmt.Println("How many autopays were started?\n\t", autopayStartCount)
	fmt.Println("How many autopays were ended?\n\t", autopayEndCount)
	fmt.Printf("What is user %d's total debits? \n\t%f\n", TestUser, testUserDebit)
	fmt.Printf("What is user %d's total credits? \n\t%f\n", TestUser, testUserCredit)
	fmt.Println("What is balance of user ID 2456938384156277127?\n\t", testUserDebit-testUserCredit)

}

func readHeader(buf *bytes.Buffer) Header {
	var header Header
	hdrRdr := bytes.NewReader(buf.Next(9))
	if err := binary.Read(hdrRdr, binary.BigEndian, &header); err != nil {
		log.Fatal("readHeader failed:", err)
	}
	return header
}

func readCreditDebit(buf *bytes.Buffer) CreditDebit {
	var record CreditDebit
	recordRdr := bytes.NewReader(buf.Next(20))
	if err := binary.Read(recordRdr, binary.BigEndian, &record); err != nil {
		log.Fatal("readCreditDebit failed:", err)
	}
	return record
}

func readAutopay(buf *bytes.Buffer) Autopay {
	var record Autopay
	recordRdr := bytes.NewReader(buf.Next(12))
	if err := binary.Read(recordRdr, binary.BigEndian, &record); err != nil {
		log.Fatal("readAutopay failed:", err)
	}
	return record
}
