package repository

import (
	"bytes"
	"context"
	"image"
	"image/png"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/kasragay/backend/internal/ports"
	"github.com/kasragay/backend/internal/utils"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const s3Caller = packageCaller + ".S3"

type S3 struct {
	logger        *utils.Logger
	client        *minio.Client
	avatarsBucket string
}

func NewS3Repo(logger *utils.Logger) ports.S3Repo {
	port := os.Getenv("MINIO_PORT")
	if port == "" {
		logger.Fatal(context.Background(), "MINIO_PORT is not set")
	}
	host := os.Getenv("MINIO_HOST")
	if host == "" {
		logger.Fatal(context.Background(), "MINIO_HOST is not set")
	}
	endpoint := host + ":" + port

	accessKey := os.Getenv("MINIO_USERNAME")
	if accessKey == "" {
		logger.Fatal(context.Background(), "MINIO_USERNAME is not set")
	}
	secretKey := os.Getenv("MINIO_PASSWORD")
	if secretKey == "" {
		logger.Fatal(context.Background(), "MINIO_PASSWORD is not set")
	}
	useSSL_ := os.Getenv("MINIO_USE_SSL")
	if useSSL_ == "" {
		logger.Fatal(context.Background(), "MINIO_USE_SSL is not set")
	}
	useSSL, err := strconv.ParseBool(useSSL_)
	if err != nil {
		logger.Fatal(context.Background(), "MINIO_USE_SSL is not valid bool")
	}
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		logger.Fatalf(context.Background(), "Failed to connect to MinIO: %v", err)
	}
	avatarsBucket := os.Getenv("MINIO_AVATARS_BUCKET")
	if avatarsBucket == "" {
		logger.Fatal(context.Background(), "MINIO_AVATAR_BUCKET is not set")
	}
	return &S3{
		logger:        logger,
		client:        minioClient,
		avatarsBucket: avatarsBucket,
	}
}

func (r *S3) UploadAvatar(ctx context.Context, userId uuid.UUID, userType ports.UserType, img *image.Image) (err error) {
	defer func() { err = utils.FuncPipe(s3Caller+".UploadAvatar", err) }()
	return r.uploadImage(ctx, r.avatarsBucket, string(userType)+"/"+userId.String()+".png", img)
}

func (r *S3) DeleteAvatar(ctx context.Context, userId uuid.UUID, userType ports.UserType) (err error) {
	defer func() { err = utils.FuncPipe(s3Caller+".DeleteAvatar", err) }()
	if err = r.client.RemoveObject(ctx, r.avatarsBucket, string(userType)+"/"+userId.String()+".png", minio.RemoveObjectOptions{}); err != nil {
		return err
	}
	r.logger.Infof(ctx, "deleted %s from bucket %s", string(userType)+"/"+userId.String()+".png", r.avatarsBucket)
	return nil
}

func (r *S3) uploadImage(ctx context.Context, bucketName, objectName string, img *image.Image) (err error) {
	defer func() { err = utils.FuncPipe(s3Caller+".uploadImage", err) }()
	var buf bytes.Buffer
	encoder := png.Encoder{
		CompressionLevel: png.BestCompression,
	}
	if err = encoder.Encode(&buf, *img); err != nil {
		return utils.BadRequestResponse.Clone().
			WithReason("error", "error encoding image")
	}
	size := int64(buf.Len())
	reader := bytes.NewReader(buf.Bytes())
	contentType := "image/png"

	_, err = r.client.PutObject(
		ctx,
		bucketName,
		objectName,
		reader,
		size,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return err
	}
	r.logger.Infof(ctx, "uploaded %s to bucket %s", objectName, bucketName)
	return nil
}
func (r *S3) GetAvatar(ctx context.Context, userId uuid.UUID, userType ports.UserType) (avatar *minio.Object, err error) {
	defer func() { err = utils.FuncPipe(s3Caller+".GetAvatar", err) }()
	avatar, err = r.client.GetObject(context.Background(), r.avatarsBucket, string(userType)+"/"+userId.String()+".png", minio.GetObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return nil, nil
		}
		return nil, err
	}
	return
}
