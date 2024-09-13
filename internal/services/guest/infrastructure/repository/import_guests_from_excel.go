package repository

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/goplaceapp/goplace-guest/utils"
	"github.com/xuri/excelize/v2"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (r *GuestRepository) ImportGuestsFromExcel(ctx context.Context, req *guestProto.ImportGuestsFromExcelRequest) (*guestProto.ImportGuestsFromExcelResponse, error) {
	var (
		successCount int
		failedCount  int
	)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})

	if err != nil {
		log.Fatal(err)
	}

	s3URL := req.GetFilePath()
	key, err := utils.ExtractKeyFromURL(s3URL)
	if err != nil {
		log.Fatal(err)
	}

	localFilePath, err := utils.DownloadFileFromS3(sess, os.Getenv("S3_BUCKET"), key)
	if err != nil {
		log.Fatal(err)
	}

	excelfile, err := excelize.OpenFile(localFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := excelfile.Close(); err != nil {
			fmt.Println(err)
		}

		if err := os.Remove(localFilePath); err != nil {
			fmt.Println(err)
		}
	}()

	rows, err := excelfile.GetRows("Sheet1")
	if err != nil {
		log.Fatal(err)
	}

	for i, row := range rows {
		if i == 0 {
			continue
		}

		guest := &guestProto.Guest{
			FirstName:   row[0],
			LastName:    row[1],
			PhoneNumber: row[3],
			Email:       row[4],
			Language:    row[6],
		}

		birthDateStr := row[2]
		if row[2] == "" {
			guest.BirthDate = nil
		} else {
			t, err := time.Parse(time.DateOnly, birthDateStr)
			if err != nil {
				t = time.Now().Truncate(time.Hour)
			}

			guest.BirthDate = timestamppb.New(t)
		}

		res, err := r.CreateGuest(ctx, &guestProto.CreateGuestRequest{
			Params: &guestProto.GuestParams{
				FirstName:   guest.GetFirstName(),
				LastName:    guest.GetLastName(),
				Email:       guest.GetEmail(),
				PhoneNumber: guest.GetPhoneNumber(),
				Language:    guest.GetLanguage(),
				BirthDate:   birthDateStr,
			},
		})
		if err != nil {
			failedCount++
			continue
		}

		if rows[i][5] != "" {
			_, err = r.AddGuestNote(ctx, &guestProto.AddGuestNoteRequest{
				Params: &guestProto.GuestNoteParams{
					GuestId:     res.GetResult().GetId(),
					Description: rows[i][5],
				},
			})
			if err != nil {
				continue
			}
		}

		successCount++
	}

	// delete the file from s3
	err = utils.DeleteFileFromS3(sess, os.Getenv("S3_BUCKET"), key)
	if err != nil {
		log.Fatal(err)
	}

	return &guestProto.ImportGuestsFromExcelResponse{
		Code:    http.StatusCreated,
		Message: fmt.Sprintf("Successfully imported %d guests and failed to import %d guests", successCount, failedCount),
	}, nil

}
