package controller

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"screenresume/internal/models"
	pb "screenresume/internal/repositories"
	"screenresume/internal/services"
	"screenresume/pkg/s3"
	"time"

	"github.com/go-fuego/fuego"
	"github.com/google/uuid"
	"github.com/minio/minio-go"
)

type ScreenResumeResources struct {
	ScreenResumeService services.ScreenResumeService
}

func (rs ScreenResumeResources) Routes(s *fuego.Server) {
	ScreenResumesGroup := fuego.Group(s, "/screen-resume")

	fuego.Post(ScreenResumesGroup, "/", rs.postScreenResume)
}

func (rs ScreenResumeResources) postScreenResume(c fuego.ContextWithBody[models.ScreenResumeCreate]) (models.ScreenResume, error) {
	const (
		endpoint   = "localhost:9000"
		bucketName = "resume"
		location   = "us-east-1"
	)

	minioClient, err := s3.MinioConnection(bucketName, location)
	if err != nil {
		return models.ScreenResume{}, err
	}

	// Extract file from request
	file, fileHeader, err := c.Request().FormFile("file")
	if err != nil {
		log.Printf("Error retrieving file from request: %v", err)
		return models.ScreenResume{}, err
	}
	defer file.Close()

	// Read file content into buffer
	fileBuffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(fileBuffer, file); err != nil {
		log.Printf("Error reading file: %v", err)
		return models.ScreenResume{}, err
	}

	// Calculate checksum
	hash := sha256.New()
	if _, err := hash.Write(fileBuffer.Bytes()); err != nil {
		return models.ScreenResume{}, fmt.Errorf("checksum calculation error: %w", err)
	}
	checksum := fmt.Sprintf("%x", hash.Sum(nil))

	// Generate unique filename
	fileName := fmt.Sprintf("%s-%d.pdf", uuid.New().String(), time.Now().Unix())

	// Upload file to MinIO
	_, err = minioClient.PutObject(
		bucketName,
		fileName,
		bytes.NewReader(fileBuffer.Bytes()),
		int64(fileBuffer.Len()),
		minio.PutObjectOptions{ContentType: "application/pdf"},
	)
	if err != nil {
		log.Printf("Error uploading file to MinIO: %v", err)
		return models.ScreenResume{}, err
	}

	log.Printf("Successfully uploaded file: %s", fileName)

	screenResumeRequest := pb.ScreenResumeRequest{
		JobDescription: c.Request().FormValue("job_description"),
		Criteria:       c.Request().Form["criteria"],
		FileUrl:        fmt.Sprintf("http://%s/%s/%s", endpoint, bucketName, fileName),
	}

	var criterias []string
	// Check if it's a nested slice ([][]string) and extract the first element
	if len(c.Request().Form["criteria"]) > 0 {
		firstElement := c.Request().Form["criteria"][0]

		// Try to parse it as JSON if it looks like a list in string form
		err := json.Unmarshal([]byte(firstElement), &criterias)
		if err != nil {
			log.Printf("Error parsing criteria JSON: %v", err)
			return models.ScreenResume{}, fmt.Errorf("invalid criteria format")
		}
	} else {
		return models.ScreenResume{}, fmt.Errorf("no criteria provided")
	}
	processScreeningCreate := models.ProcessScreeningCreate{
		ApplicationID: c.Request().FormValue("application_id"),
		Criteria:      criterias,
	}

	// Create file record
	fileCreate := models.FilesCreate{
		Path:     fmt.Sprintf("%s/%s", bucketName, fileName),
		FileType: fileHeader.Header.Get("Content-Type"),
		Checksum: checksum,
	}
	screenResumeResponse, err := rs.ScreenResumeService.ProcessScreening(&screenResumeRequest, processScreeningCreate, fileCreate)
	if err != nil {
		return models.ScreenResume{}, err
	}

	return screenResumeResponse, nil
}
