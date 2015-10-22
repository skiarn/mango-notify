package models

import (
	"bufio"
	"os"
)

//Line represent a file line.
type Line struct {
	Number  int    `json:"number"`
	Content string `json:"content"`
	Sent    bool   `json:"sent"`
}

//Lines represent a list of line.
type Lines []Line

func (slice Lines) Len() int {
	return len(slice)
}

func (slice Lines) Less(i, j int) bool {
	return slice[i].Number < slice[j].Number
}

func (slice Lines) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

//UpdateChangedLines scans file and update lines in db that has changed.
func UpdateChangedLines(file string) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	r := bufio.NewReader(f)
	lineNr := 0
	s, err := readln(r)
	savedLine, err := GetLine(db, lineNr)
	if err != nil {
		return err
	}

	if savedLine == nil {
		//line dont exist yet.
		l := Line{Number: lineNr, Content: s, Sent: false}
		serr := SaveLine(db, &l)
		if serr != nil {
			return serr
		}
	} else if s != savedLine.Content {
		// content changed.
		l := Line{Number: lineNr, Content: s, Sent: false}
		serr := SaveLine(db, &l)
		if serr != nil {
			return serr
		}
	}

	lineNr++
	for err == nil {
		s, err = readln(r)
		savedLine, err := GetLine(db, lineNr)
		if err != nil {
			return err
		}

		if savedLine == nil {
			//line dont exist yet.
			l := Line{Number: lineNr, Content: s, Sent: false}
			serr := SaveLine(db, &l)
			if serr != nil {
				return serr
			}
		} else if s != savedLine.Content {
			// content changed.
			l := Line{Number: lineNr, Content: s, Sent: false}
			serr := SaveLine(db, &l)
			if serr != nil {
				return serr
			}
		}
		lineNr++
	}
	return nil
}

// Readln returns a single line (without the ending \n)
// from the input buffered reader.
// An error is returned iff there is an error with the
// buffered reader.
func readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool = true
		err      error
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}
