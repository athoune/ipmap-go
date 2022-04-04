package csv

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

// https://ftp.ripe.net/ripe/ipmap/geolocations_2022-04-04.csv.bz2

type Line struct {
	IP                  net.IP
	network             *net.IPNet
	Geolocation_id      string
	City_name           string
	State_name          string
	Country_name        string
	Country_code_alpha2 string
	Country_code_alpha3 string
	Location            *Location
	Score               float32
}

func (l Line) Network() net.IPNet {
	return *l.network
}

type Location struct {
	Latitude  float64
	Longitude float64
}

type CVS struct {
	scanner bufio.Scanner
}

func New(r io.Reader) *CVS {
	return &CVS{
		scanner: *bufio.NewScanner(r),
	}
}

func (c *CVS) Next() bool {
	return c.scanner.Scan()
}

func (c *CVS) Value() (Line, error) {
	line := c.scanner.Text()
	for _, patch := range []string{"Washington, D.C.", "WASHINGTON,_D.C", "Bonaire, Saint Eustatius and Saba"} {
		line = strings.ReplaceAll(line, patch, strings.Replace(patch, ",", "_", 1))
	}
	values := strings.Split(line, ",")
	if len(values) != 10 {
		fmt.Println(line)
		return Line{}, nil
	}
	l := Line{
		Geolocation_id:      values[1],
		City_name:           values[2],
		State_name:          values[3],
		Country_name:        values[4],
		Country_code_alpha2: values[5],
		Country_code_alpha3: values[6],
	}
	var err error
	l.IP, l.network, err = net.ParseCIDR(values[0])
	if err != nil {
		return l, err
	}
	if values[7] != "" && values[8] != "" {
		latitude, err := strconv.ParseFloat(values[7], 64)
		if err != nil {
			return l, fmt.Errorf("latitude parse error [%s] : %s", values[7], err)
		}
		longitude, err := strconv.ParseFloat(values[8], 64)
		if err != nil {
			return l, fmt.Errorf("longitude parse error [%s] : %s", values[8], err)
		}
		l.Location = &Location{
			Latitude:  latitude,
			Longitude: longitude,
		}
	}
	if values[9] != "" {
		score, err := strconv.ParseFloat(values[9], 32)
		if err != nil {
			return l, fmt.Errorf("score parse error [%s] : %s", values[9], err)
		}
		l.Score = float32(score)
	}
	return l, nil
}
