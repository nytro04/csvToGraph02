// func parseStatementsCSV(r io.Reader) ([]Statement, error) {
	
// 	csvFile := csv.NewReader(r)

// 	records, err := csvFile.ReadAll()
// 	if err != nil {
// 		return nil, err
// 	}

// 	if len(records) == 0 {
// 		return nil, fmt.Errorf("could not parse csv: missing header")
// 	}

// 	headers := records[0]

// 	if len(headers) != 4 {
// 		return nil, fmt.Errorf("could not parse csv: invalid headers")
// 	}

// 	if headers[0] != "DATE" {
// 		return nil, fmt.Errorf("could not parse csv: missing DATE")
// 	}

// 	if headers[1] != "DESCRIPTION" {
// 		return nil, fmt.Errorf("could not parse csv: missing DESCRIPTION")
// 	}

// 	if headers[2] != "TYPE" {
// 		return nil, fmt.Errorf("could not parse csv: missing TYPE")
// 	}

// 	if headers[3] != "AMOUNT" {
// 		return nil, fmt.Errorf("could not parse csv: missing AMOUNT")
// 	}

// 	records = records[1:] // Skip's the firt row which are the label.
// 	credit := 0
// 	debit := 0
// 	ss := make([]Statement, 0)

// 	for _, record := range records {
// 		amount, err := strconv.Atoi(record[3])
// 		if err != nil {
// 			return nil, fmt.Errorf("could not parse amount as integer %v", err)
// 		}

// 		if record[2] == "C" {
// 			credit += amount
// 		} else if record[2] == "D" {
// 			debit += amount
// 		}

// 		ss = append(ss, Statement{credit, debit})
// 	}

// 	return ss, nil
// }