package stashfin

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"cloud.google.com/go/storage"
	"github.com/yeka/zip"
)

func unZipFile(byteData []byte, password string) ([]byte, error) {

	zipReader := bytes.NewReader(byteData)
	var resultBuffer []byte
	// Open the ZIP archive for reading.
	zipFile, err := zip.NewReader(zipReader, int64(len(byteData)))
	if err != nil {
		return resultBuffer, err
	}
	// defer r.Close()

	for _, f := range zipFile.File {
		if f.IsEncrypted() {
			f.SetPassword(password)
		}

		r, err := f.Open()
		if err != nil {
			return resultBuffer, err
		}

		buf, err := ioutil.ReadAll(r)
		if err != nil {
			return resultBuffer, err
		}
		r.Close()
		return buf, nil
	}
	return resultBuffer, nil
}

func unzipFile(file *storage.Reader, fileName string, password string) ([]byte, error) {

	zipDataBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	// Create a temporary file to write the zip data
	tmpZipFile := "/tmp/" + fileName // You can change the path as needed
	err = os.WriteFile(tmpZipFile, zipDataBytes, 0644)
	if err != nil {
		return nil, err
	}
	// Use the 'unzip' command-line tool to extract the zip file into memory as bytes.
	cmd := exec.Command("unzip", "-P", password, tmpZipFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	return output, nil
}

func parseAadhaarXML(xmlData []byte) (*OkAadhaarKycData, error) {
	// Unmarshal (parse) the XML data from the byte slice.
	okAadhaarKycData := OkAadhaarKycData{}

	err := xml.Unmarshal(xmlData, &okAadhaarKycData)
	if err != nil {
		return nil, err
	}
	return &okAadhaarKycData, nil
}

func createAadhaarXMLFileForStashfin(xmlData *OkAadhaarKycData) ([]byte, error) {
	data := SfAadhaarKycData{
		XMLName: xmlData.XMLName,
		CertificateData: SfCertificateData{
			KycRes: SfKycRes{
				Code:        "",
				Ret:         "",
				Transaction: "",
				UidData: SfUidData{
					Token: "",
					UID:   "",
					Poi: SfPoi{
						Dob:    xmlData.UidData.Poi.Dob,
						Gender: xmlData.UidData.Poi.Gender,
						Name:   xmlData.UidData.Poi.Name,
					},
					Poa: SfPoa{
						CareOf:  xmlData.UidData.Poa.CareOf,
						Country: xmlData.UidData.Poa.Country,
						Dist:    xmlData.UidData.Poa.Dist,
						Loc:     xmlData.UidData.Poa.Location,
						Pincode: xmlData.UidData.Poa.Pincode,
						State:   xmlData.UidData.Poa.State,
						Vtc:     xmlData.UidData.Poa.Vtc,
					},
					LData: SfLData{
						CareOf:   xmlData.UidData.Poa.CareOf,
						Country:  xmlData.UidData.Poa.Country,
						Dist:     xmlData.UidData.Poa.Dist,
						Language: "",
						Loc:      xmlData.UidData.Poa.Location,
						Name:     xmlData.UidData.Poi.Name,
						Pincode:  xmlData.UidData.Poa.Pincode,
						State:    xmlData.UidData.Poa.State,
						Vtc:      xmlData.UidData.Poa.Vtc,
					},
					Pht: xmlData.UidData.Pht,
				},
			},
		},
		Signature: SfSignature{
			XMLName: "http://www.w3.org/2000/09/xmldsig#",
			SignedInfo: SfSignedInfo{
				CanonicalizationMethod: SfCanonicalizationMethod(xmlData.Signature.SignedInfo.CanonicalizationMethod),
				SignatureMethod:        SfSignatureMethod(xmlData.Signature.SignedInfo.SignatureMethod),
				Reference: SfReference{
					URI: xmlData.Signature.SignedInfo.Reference.URI,
					Transforms: SfTransforms{
						Transform: SfTransform(xmlData.Signature.SignedInfo.Reference.Transforms.Transform),
					},
					DigestMethod: SfDigestMethod(xmlData.Signature.SignedInfo.Reference.DigestMethod),
					DigestValue:  xmlData.Signature.SignedInfo.Reference.DigestValue,
				},
			},
			SignatureValue: xmlData.Signature.SignatureValue,
			KeyInfo: SfKeyInfo{
				X509Data: SfX509Data{
					X509SubjectNames: xmlData.Signature.KeyInfo.X509Data.X509SubjectName,
					X509Certificates: xmlData.Signature.KeyInfo.X509Data.X509Certificate,
				},
			},
		},
	}

	// Marshal the struct into XML
	xmlFileData, err := xml.MarshalIndent(data, "", "    ")
	if err != nil {
		return nil, err
	}

	//// Print the XML as a string
	//fmt.Println(string(xmlFileData))
	return xmlFileData, err
}

type SfAadhaarKycData struct {
	XMLName         xml.Name          `xml:"Certificate"`
	CertificateData SfCertificateData `xml:"CertificateData"`
	Signature       SfSignature       `xml:"Signature"`
}

type SfCertificateData struct {
	KycRes SfKycRes `xml:"KycRes"`
}

type SfKycRes struct {
	Code        string    `xml:"code,attr"`
	Ret         string    `xml:"ret,attr"`
	Timestamp   time.Time `xml:"ts,attr"`
	TTL         time.Time `xml:"ttl,attr"`
	Transaction string    `xml:"txn,attr"`
	UidData     SfUidData `xml:"UidData"`
}

type SfUidData struct {
	Token string  `xml:"tkn,attr"`
	UID   string  `xml:"uid,attr"`
	Poi   SfPoi   `xml:"Poi"`
	Poa   SfPoa   `xml:"Poa"`
	LData SfLData `xml:"LData"`
	Pht   string  `xml:"Pht"`
}

type SfPoi struct {
	Dob    string `xml:"dob,attr"`
	Gender string `xml:"gender,attr"`
	Name   string `xml:"name,attr"`
}

type SfPoa struct {
	CareOf  string `xml:"co,attr"`
	Country string `xml:"country,attr"`
	Dist    string `xml:"dist,attr"`
	Loc     string `xml:"loc,attr"`
	Pincode string `xml:"pc,attr"`
	State   string `xml:"state,attr"`
	Vtc     string `xml:"vtc,attr"`
}

type SfLData struct {
	CareOf   string `xml:"co,attr"`
	Country  string `xml:"country,attr"`
	Dist     string `xml:"dist,attr"`
	Language string `xml:"lang,attr"`
	Loc      string `xml:"loc,attr"`
	Name     string `xml:"name,attr"`
	Pincode  string `xml:"pc,attr"`
	State    string `xml:"state,attr"`
	Vtc      string `xml:"vtc,attr"`
}

type SfSignature struct {
	XMLName        string       `xml:"Signature"`
	SignedInfo     SfSignedInfo `xml:"SignedInfo"`
	SignatureValue string       `xml:"SignatureValue"`
	KeyInfo        SfKeyInfo    `xml:"KeyInfo"`
}

type SfSignedInfo struct {
	CanonicalizationMethod SfCanonicalizationMethod `xml:"CanonicalizationMethod"`
	SignatureMethod        SfSignatureMethod        `xml:"SignatureMethod"`
	Reference              SfReference              `xml:"Reference"`
}

type SfCanonicalizationMethod struct {
	Algorithm string `xml:"Algorithm,attr"`
}

type SfSignatureMethod struct {
	Algorithm string `xml:"Algorithm,attr"`
}

type SfReference struct {
	URI          string         `xml:"URI,attr"`
	Transforms   SfTransforms   `xml:"Transforms"`
	DigestMethod SfDigestMethod `xml:"DigestMethod"`
	DigestValue  string         `xml:"DigestValue"`
}

type SfTransforms struct {
	Transform SfTransform `xml:"Transform"`
}

type SfTransform struct {
	Algorithm string `xml:"Algorithm,attr"`
}

type SfDigestMethod struct {
	Algorithm string `xml:"Algorithm,attr"`
}

type SfKeyInfo struct {
	X509Data SfX509Data `xml:"X509Data"`
}

type SfX509Data struct {
	X509SubjectNames []string `xml:"X509SubjectName"`
	X509Certificates []string `xml:"X509Certificate"`
}

type OkAadhaarKycData struct {
	XMLName     xml.Name    `xml:"OfflinePaperlessKyc"`
	ReferenceID string      `xml:"referenceId,attr"`
	UidData     OkUidData   `xml:"UidData"`
	Signature   OkSignature `xml:"Signature"`
}

type OkUidData struct {
	Poi OkPoi  `xml:"Poi"`
	Poa OkPoa  `xml:"Poa"`
	Pht string `xml:"Pht"`
}

type OkPoi struct {
	Dob    string `xml:"dob,attr"`
	Gender string `xml:"gender,attr"`
	Name   string `xml:"name,attr"`
}

type OkPoa struct {
	CareOf     string `xml:"careof,attr"`
	Country    string `xml:"country,attr"`
	Dist       string `xml:"dist,attr"`
	House      string `xml:"house,attr"`
	Landmark   string `xml:"landmark,attr"`
	Location   string `xml:"loc,attr"`
	Pincode    string `xml:"pc,attr"`
	PostOffice string `xml:"po,attr"`
	State      string `xml:"state,attr"`
	Street     string `xml:"street,attr"`
	SubDist    string `xml:"subdist,attr"`
	Vtc        string `xml:"vtc,attr"`
}

type OkSignature struct {
	SignedInfo     OkSignedInfo `xml:"SignedInfo"`
	SignatureValue string       `xml:"SignatureValue"`
	KeyInfo        OkKeyInfo    `xml:"KeyInfo"`
}

type OkSignedInfo struct {
	CanonicalizationMethod rmnkmranonicalizationMethod `xml:"CanonicalizationMethod"`
	SignatureMethod        OkSignatureMethod        `xml:"SignatureMethod"`
	Reference              OkReference              `xml:"Reference"`
}

type rmnkmranonicalizationMethod struct {
	Algorithm string `xml:"Algorithm,attr"`
}

type OkSignatureMethod struct {
	Algorithm string `xml:"Algorithm,attr"`
}

type OkReference struct {
	URI          string         `xml:"URI,attr"`
	Transforms   OkTransforms   `xml:"Transforms"`
	DigestMethod OkDigestMethod `xml:"DigestMethod"`
	DigestValue  string         `xml:"DigestValue"`
}

type OkTransforms struct {
	Transform OkTransform `xml:"Transform"`
}

type OkTransform struct {
	Algorithm string `xml:"Algorithm,attr"`
}

type OkDigestMethod struct {
	Algorithm string `xml:"Algorithm,attr"`
}

type OkKeyInfo struct {
	X509Data OkX509Data `xml:"X509Data"`
}

type OkX509Data struct {
	X509SubjectName []string `xml:"X509SubjectName"`
	X509Certificate []string `xml:"X509Certificate"`
}
