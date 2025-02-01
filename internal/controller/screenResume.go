package controller

import (
	"bytes"
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
	file, _, err := c.Request().FormFile("file")
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

	new, err := rs.ScreenResumeService.ScreenResume(&screenResumeRequest)
	if err != nil {
		return models.ScreenResume{}, err
	}

	criteriaDecisions := make([]*models.CriteriaDecision, len(new.CriteriaDecisions))
	for i, f := range new.CriteriaDecisions {
		criteriaDecisions[i] = &models.CriteriaDecision{
			Reasoning: f.Reasoning,
			Decision:  f.Decision,
		}
	}

	screenResumeResponse := models.ScreenResume{
		CriteriaDecisions: criteriaDecisions,
		OverallReasoning:  new.OverallReasoning,
		OverallDecision:   new.OverallDecision,
		ResumeName:        new.ResumeName,
	}

	return screenResumeResponse, nil
}
