package metadata

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/evanoberholster/imagemeta"
)

// ImageType: 	image/jpeg
// Make:
// Model:
// LensMake:
// LensModel:
// CameraSerial:
// LensSerial:
// Image Size: 	0x0
// Orientation: 	Unknown
// ShutterSpeed:
// Aperture: 	0.00
// ISO: 		0
// Flash: 		No Flash
// Focal Length: 	0.00mm
// Fl 35mm Eqv: 	0.00mm
// Exposure Prgm: 	Not Defined
// Metering Mode: 	Unknown
// Exposure Mode: 	Auto
// Date Modified: 	0001-01-01 00:00:00 +0000 UTC
// Date Created: 	0001-01-01 00:00:00 +0000 UTC
// Date Original: 	0001-01-01 00:00:00 +0000 UTC
// Date GPS: 	0001-01-01 00:00:00 +0000 UTC
// Artist:
// Copyright:
// Software:
// Image Desc:
// GPS Altitude: 	0.00
// GPS Latitude: 	0.000000
// GPS Longitude: 	0.000000
func getMetadata0(file string) {

	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	e, err := imagemeta.Decode(f)
	if err != nil {
		panic(err)
	}
	fmt.Println(e)
}

//	{
//		"place_id":49016033,
//		"licence":"Data © OpenStreetMap contributors, ODbL 1.0. http://osm.org/copyright",
//		"osm_type":"way",
//		"osm_id":359637646,
//		"lat":"36.0911488",
//		"lon":"28.0881070",
//		"class":"historic",
//		"type":"archaeological_site",
//		"place_rank":30,
//		"importance":5.18844892041965e-05,
//		"addresstype":"historic",
//		"name":"Lindos",
//		"display_name":"Lindos, Ακροπόλεως, Lindos, Municipality of Rhodes, Rhodes Regional Unit, South Aegean, Aegean, 851 07, Greece",
//		"address":
//		{
//			"historic":"Lindos","road":"Ακροπόλεως",
//			"town":"Lindos",
//			"municipality":"Municipality of Rhodes",
//			"county":"Rhodes Regional Unit","state_district":"South Aegean",
//			"ISO3166-2-lvl5":"GR-L","state":"Aegean","postcode":"851 07",
//			"country":"Greece",
//			"country_code":"gr"
//		},
//			"boundingbox":["36.0894062","36.0929831","28.0868152","28.0902161"]
//		}
func getNominatim0(lat, lon float64) error {
	url := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?lat=%f&lon=%f&format=json&accept-language=en", lat, lon)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creando request: %v", err)
	}
	req.Header.Set("User-Agent", "MyApp/1.0") // Requerido

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error haciendo request: %v", err)
	}
	defer resp.Body.Close()

	// Verificar el status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error leyendo respuesta: %v", err)
	}

	fmt.Println(string(body))
	return nil
}

func getInfo() {
	// getMetadata("images/eaf06d9a-b644-40bf-ac09-0216d5468545.jpeg")
	//
	//	if err := getNominatim(36.091436, 28.088505); err != nil {
	//		fmt.Printf("Error: %v\n", err)
	//	}
}
