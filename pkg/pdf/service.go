package pdf

import pkgbucket "maicare_go/pkg/bucket"

type pdfService struct {
	bucketClient *pkgbucket.ObjectStorageClient
}

func NewPdfService(bucketClient *pkgbucket.ObjectStorageClient) *pdfService {
	return &pdfService{
		bucketClient: bucketClient,
	}
}
