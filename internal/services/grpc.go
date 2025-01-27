package services

import (
	"context"
	"log"

	pb "screenresume/resume_screener"
)

type resumeScreenerServer struct {
	pb.UnimplementedResumeScreenerServer
}

func (s *resumeScreenerServer) ScreenResume(ctx context.Context, req *pb.ScreenResumeRequest) (*pb.ScreenResumeResponse, error) {
	// Implement your resume screening logic here.
	log.Printf("Received request to screen resume: %v", req)

	// Example response
	return &pb.ScreenResumeResponse{
		CriteriaDecisions: []*pb.CriteriaDecision{
			{
				Reasoning: "The candidate has relevant experience.",
				Decision:  true,
			},
		},
		OverallReasoning: "The candidate meets most of the criteria.",
		OverallDecision:  true,
		ResumeName:       req.Filename,
	}, nil
}

func (s *resumeScreenerServer) UploadResume(stream pb.ResumeScreener_UploadResumeServer) error {
	// Implement your resume upload logic here.
	log.Println("Starting resume upload")

	var fileInfo *pb.FileInfo
	var fileData []byte

	for {
		req, err := stream.Recv()
		if err != nil {
			log.Printf("Error receiving data: %v", err)
			return err
		}

		switch data := req.Data.(type) {
		case *pb.UploadResumeRequest_Info:
			fileInfo = data.Info
			log.Printf("Received file info: %v", fileInfo)
		case *pb.UploadResumeRequest_Chunk:
			fileData = append(fileData, data.Chunk...)
			log.Printf("Received file chunk of size: %d", len(data.Chunk))
		}

		// Check if the upload is complete (e.g., client closes the stream).
		if req.Data == nil {
			break
		}
	}

	// Save the file or process it further.
	log.Printf("Upload complete. File size: %d", len(fileData))

	// Send a response to the client.
	return stream.SendAndClose(&pb.UploadResumeResponse{
		ObjectName: fileInfo.Filename,
	})
}
