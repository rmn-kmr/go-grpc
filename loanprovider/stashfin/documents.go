package stashfin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"github.com/rmnkmr/go-common/errors"
	"github.com/rmnkmr/lsp"
	api "github.com/rmnkmr/lsp/proto"
	"github.com/rmnkmr/lsp/utils"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"

	"github.com/rmnkmr/lsp/log"
)

const (
	upload_documents = "/v3/upload-documents"
	xml_file_name    = "okyc_aadhaar.xml"
)

type stashfinDocumentType string

const (
	stashfinDocumentType_SELFIE    stashfinDocumentType = "1"
	stashfinDocumentType_PAN       stashfinDocumentType = "2"
	stashfinDocumentType_ADHAR_ZIP stashfinDocumentType = "7"
	stashfinDocumentType_AGREEMENT stashfinDocumentType = "8"
	stashfinDocumentType_KFS       stashfinDocumentType = "8"
)

var stashfinDocumentTypeMapping = map[api.DocumentType]stashfinDocumentType{
	api.DocumentType_SELFIE:                stashfinDocumentType_SELFIE,
	api.DocumentType_AADHAAR_ZIP:           stashfinDocumentType_ADHAR_ZIP,
	api.DocumentType_SIGNED_LOAN_AGREEMENT: stashfinDocumentType_AGREEMENT,
	api.DocumentType_KFS:                   stashfinDocumentType_KFS,
	api.DocumentType_PAN:                   stashfinDocumentType_PAN,
}

type uploadDocumentResponse struct {
	Status  bool        `json:"status"`
	Results interface{} `json:"results"`
	Errors  interface{} `json:"errors"`
}

// upload KYC documents
func uploadKYCDocuments(ctx context.Context, lp ApiClient, referenceId string,
	userKycDetails []*api.UserKycDetails, panMetadata string) (*api.UploadDocumentResponse, error) {
	for _, document := range userKycDetails {
		var err error
		// Get file from GCS bucket
		file, err := lsp.GetGCSObject(ctx, lp.GCSClient, document.Url, document.BucketName)
		if err != nil {
			log.Error(ctx, err, "uploadDocuments:unable to get file from GCS", "document", document)
			return nil, err
		}
		inputFileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			log.Error(ctx, err, "UploadDocuments: unable to read file")
			return nil, err
		}
		if document.Type == utils.PHOTOGRAPH {
			docWithFile := &api.Document{
				Type: api.DocumentType_SELFIE,
				Name: document.DocumentName,
				File: inputFileBytes,
			}
			_, err = uploadDocument(ctx, lp, referenceId, stashfinDocumentTypeMapping[api.DocumentType_SELFIE], docWithFile, "")
		} else if document.Type == utils.EKYC_ADHAR_ZIP {

			byteDta, err := ConvertAadhaarFileForStashfin(ctx, inputFileBytes, document.DocumentName, document.ShareCode)
			if err != nil {
				log.Error(ctx, err, "UploadDocuments: ConvertAadhaarFileForStashfin failed")
				return nil, err
			}
			docWithFile := &api.Document{
				Type: api.DocumentType_AADHAAR_ZIP,
				Name: xml_file_name,
				File: byteDta,
			}
			log.Info(ctx, "UploadDocuments: uploading aadhaar zip file",
				"bytes", byteDta, "shareCode", document.ShareCode)
			_, err = uploadDocument(ctx, lp, referenceId, stashfinDocumentTypeMapping[api.DocumentType_AADHAAR_ZIP], docWithFile, document.ShareCode)
		}
		if err != nil {
			log.Error(ctx, err, "UploadDocuments: upload document failed", "err", err, "referenceId", referenceId)
			return nil, err
		}
		log.Info(ctx, "UploadDocuments: upload success", "documentType", document.Type)
	}

	_, err := uploadPanDocument(ctx, lp, referenceId, stashfinDocumentTypeMapping[api.DocumentType_PAN], panMetadata)
	if err != nil {
		log.Error(ctx, err, "uploadKYCDocuments:uploadPanDocument err", "document", panMetadata)
		return nil, err
	}
	log.Info(ctx, "UploadDocuments: upload success", "documentType", "PAN")

	return nil, nil
}

// upload Loan documents
func uploadLoanDocuments(ctx context.Context, lp ApiClient, referenceId string,
	documents []*api.Document) (*api.UploadDocumentResponse, error) {
	acceptedDocuments := []string{api.DocumentType_SIGNED_LOAN_AGREEMENT.String(), api.DocumentType_KFS.String()}
	for _, document := range documents {
		if !utils.StringInSlice(document.Type.String(), acceptedDocuments) {
			continue
		}
		var err error
		// Get file from GCS bucket
		file, err := lsp.GetGCSObject(ctx, lp.GCSClient, document.Url, document.Bucket)
		if err != nil {
			log.Error(ctx, err, "uploadDocuments:unable to get file from GCS", "document", document)
			return nil, err
		}
		inputFileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			log.Error(ctx, err, "UploadDocuments: unable to read file")
			return nil, err
		}
		docWithFile := &api.Document{
			Name: document.Name,
			File: inputFileBytes,
		}
		_, err = uploadDocument(ctx, lp, referenceId, stashfinDocumentTypeMapping[document.Type], docWithFile, "")
		if err != nil {
			log.Error(ctx, err, "UploadDocuments: upload document failed", "err", err, "referenceId", referenceId)
			return nil, err
		}
	}

	return nil, nil
}

func uploadDocument(ctx context.Context, lp ApiClient, referenceID string, documentType stashfinDocumentType, document *api.Document, shareCode string) (*uploadDocumentResponse, error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	writer.WriteField("document_type", string(documentType))
	writer.WriteField("application_id", referenceID)
	writer.WriteField("document_name", document.Name)
	if shareCode != "" {
		writer.WriteField("share_code", shareCode)
	}

	part, err := writer.CreatePart(createPartHeaders("files", document.Name))
	if err != nil {
		log.Error(ctx, err, "uploadLendboxDocumentFromLocal: failed in creating part header")
		return nil, err
	}

	_, err = io.Copy(part, bytes.NewReader(document.GetFile()))
	if err != nil {
		log.Error(ctx, err, "uploadLendboxDocumentFromLocal: failed in reading the file")
		return nil, err
	}

	writer.Close()
	httpReq, err := http.NewRequest(http.MethodPost, lp.ProviderBaseUrl+upload_documents, payload)
	if err != nil {
		log.Error(ctx, err, "uploadDocument: ", "err", err)
		return nil, errors.Internal()
	}

	authToken, err := lp.GetTokenFromProvider(ctx)
	if err != nil {
		log.Error(ctx, err, "uploadDocument: GetTokenFromProvider failed", "err", err, "authToken", authToken)
		return nil, err
	}
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())
	httpReq.Header.Set("client-token", authToken.AuthToken)

	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Error(ctx, err, "uploadDocument: upload documents failed", "err", err, "req", httpReq)
		return nil, errors.Internal()
	}

	defer httpResp.Body.Close()

	resp := uploadDocumentResponse{}

	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		log.Error(ctx, err, "uploadDocument: failed to parse response body: %v", err)
		return nil, errors.Internal()
	}
	if !resp.Status {
		err = errors.From(400, "upload documents failed")
		log.Error(ctx, err, "uploadDocument: upload document failed",
			"resp", resp)
		return nil, err
	}

	return &resp, nil
}

func uploadPanDocument(ctx context.Context, lp ApiClient, referenceID string, documentType stashfinDocumentType, panMetadata string) (*uploadDocumentResponse, error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	documentName := fmt.Sprintf("%s_%s.pdf", "pan", referenceID)
	writer.WriteField("document_type", string(documentType))
	writer.WriteField("application_id", referenceID)
	writer.WriteField("document_name", documentName)
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.MultiCell(0, 10, panMetadata, "", " ", false)
	// Create a buffer to store the PDF content as bytes
	var buffer bytes.Buffer

	// Output the PDF content to the buffer
	err := pdf.Output(&buffer)
	if err != nil {
		log.Error(ctx, err, "uploadPanDocument: failed in pdf.Output")
		return nil, err
	}

	// You can now use 'buffer.Bytes()' to get the PDF content as bytes
	pdfBytes := buffer.Bytes()
	part, err := writer.CreatePart(createPartHeaders("files", documentName))
	if err != nil {
		log.Error(ctx, err, "uploadPanDocument: failed in creating part header")
		return nil, err
	}

	_, err = io.Copy(part, bytes.NewReader(pdfBytes))
	if err != nil {
		log.Error(ctx, err, "uploadPanDocument: failed in reading the file")
		return nil, err
	}

	writer.Close()
	httpReq, err := http.NewRequest(http.MethodPost, lp.ProviderBaseUrl+upload_documents, payload)
	if err != nil {
		log.Error(ctx, err, "uploadPanDocument: ", "err", err)
		return nil, errors.Internal()
	}

	authToken, err := lp.GetTokenFromProvider(ctx)
	if err != nil {
		log.Error(ctx, err, "uploadDocument: GetTokenFromProvider failed", "err", err, "authToken", authToken)
		return nil, err
	}
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())
	httpReq.Header.Set("client-token", authToken.AuthToken)

	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Error(ctx, err, "uploadPanDocument: upload documents failed", "err", err, "req", httpReq)
		return nil, errors.Internal()
	}

	defer httpResp.Body.Close()

	resp := uploadDocumentResponse{}

	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		log.Error(ctx, err, "uploadPanDocument: failed to parse response body: %v", err)
		return nil, errors.Internal()
	}
	if !resp.Status {
		err = errors.From(400, "upload documents failed")
		log.Error(ctx, err, "uploadPanDocument: upload document failed",
			"resp", resp)
		return nil, err
	}

	return &resp, nil
}

func createPartHeaders(parameter string, documentName string) textproto.MIMEHeader {
	partHeader := textproto.MIMEHeader{}

	disp := fmt.Sprintf("form-data; name=\"%s\"; filename=\"%s\"", parameter, documentName)
	partHeader.Add("Content-Disposition", disp)

	if strings.Contains(strings.ToLower(documentName), "png") {
		partHeader.Add("Content-Type", "image/png")
	} else if strings.Contains(strings.ToLower(documentName), "jpg") || strings.Contains(strings.ToLower(documentName), "jpeg") {
		partHeader.Add("Content-Type", "image/jpeg")
	} else if strings.Contains(strings.ToLower(documentName), "xml") {
		partHeader.Add("Content-Type", "application/xml")
	} else if strings.Contains(strings.ToLower(documentName), "pdf") {
		partHeader.Add("Content-Type", "application/pdf")
	}

	return partHeader
}
