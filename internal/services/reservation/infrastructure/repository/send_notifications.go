package repository

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
	"github.com/goplaceapp/goplace-user/pkg/userservice/domain"
)

func (r *ReservationRepository) SendSpecialOccasionMessage(messageBody string, reservationId int32) error {
	var queueUrl string

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		return err
	}

	sqsClient := sqs.New(sess)
	users := []*domain.User{}
	r.GetSharedDB().Find(&users)
	for _, u := range users {
		queueName := strings.ReplaceAll(strings.ReplaceAll(u.Email, "@", "-"), ".", "-")
		if queueUrl, err = getQueueURL(sess, queueName); err != nil {
			if _, err := createQueue(sess, queueName); err != nil {
				return err
			}

			queueUrl, err = getQueueURL(sess, queueName)
			if err != nil {
				return err
			}
		}

		id := uuid.New().String()
		fullMessageBody := fmt.Sprintf(`{"id": "%s", "from": "Support", "body": "%s", "type": "broadcast", "seen": false, "reservation": %v, "createdAt": "%v"}`, id, messageBody, reservationId, time.Now().UTC())

		_, err = sqsClient.SendMessage(&sqs.SendMessageInput{
			QueueUrl:    &queueUrl,
			MessageBody: aws.String(fullMessageBody),
		})
	}

	return err
}

func createQueue(sess *session.Session, queueName string) (*sqs.CreateQueueOutput, error) {
	sqsClient := sqs.New(sess)
	result, err := sqsClient.CreateQueue(&sqs.CreateQueueInput{
		QueueName: &queueName,
		Attributes: map[string]*string{
			"DelaySeconds": aws.String("0"),
		},
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func getQueueURL(sess *session.Session, queue string) (string, error) {
	sqsClient := sqs.New(sess)

	result, err := sqsClient.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queue,
	})

	if err != nil {
		return "", err
	}

	return *result.QueueUrl, nil
}
