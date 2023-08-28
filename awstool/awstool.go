package awstool

import (
	"bufio"
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"go-gin-restful-service/config"
	"go-gin-restful-service/log"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"go.uber.org/zap"
)

type AwsTool struct {
	CFG      *config.Config
	S3Client *s3.Client
}

var ctx = context.TODO()

const defaultRegion = "us-east-1"

// 初始bucket和s3 连接
func (m *AwsTool) InitBucket(bucketName *string, cfg *config.Config) error {
	// 创建OSSClient实例。
	hostAddress := cfg.AwsConfig.Endpoint
	resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...any) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:       "aws",
			SigningRegion:     defaultRegion,
			URL:               hostAddress,
			HostnameImmutable: true,
		}, nil
	})

	awsConf := aws.Config{
		Region:                      defaultRegion,
		EndpointResolverWithOptions: resolver,
		Credentials:                 credentials.NewStaticCredentialsProvider(cfg.AwsConfig.AccessKey, cfg.AwsConfig.SecretKey, ""),
	}

	s3Client := s3.NewFromConfig(awsConf)

	defaultBecket := cfg.AwsConfig.BucketName
	if bucketName == nil || len(*bucketName) == 0 {
		bucketName = &defaultBecket
	}
	found, err := BucketExists(s3Client, *bucketName)
	if err != nil || !found {
		log.Logger.Error(err)
		CreateBucket(s3Client, *bucketName)
	}
	m.S3Client = s3Client
	m.CFG = cfg
	return nil
}

func CreateBucket(s3Client *s3.Client, bucketName string) {
	bucket, err := s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(defaultRegion),
		},
	})
	if err != nil {
		log.Logger.Error(err)
		return
	}
	log.Logger.Infof("Successfully created mybucket %v", bucket)
}

// 判定bucket是否存在
func BucketExists(s3Client *s3.Client, bucketName string) (bool, error) {
	_, err := s3Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	exists := true
	if err != nil {
		var apiError smithy.APIError
		if errors.As(err, &apiError) {
			switch apiError.(type) {
			case *types.NotFound:
				log.Logger.Errorf("Bucket %v is available.\n", bucketName)
				exists = false
				err = nil
			default:
				log.Logger.Errorf("Either you don't have access to bucket %v or another error occurred. "+
					"Here's what happened: %v\n", bucketName, err)
			}
		}
	} else {
		log.Logger.Infof("Bucket %v exists and you already own it.", bucketName)
	}

	return exists, err
}

// 上传外部传入文件
func (m *AwsTool) UploadDataToFile(dataStr string, objectKey string, contentType string) (string, error) {
	s3Client := m.S3Client
	// 上传文件流。
	uploadInfo, err := s3Client.PutObject(ctx,
		&s3.PutObjectInput{
			Bucket:      aws.String(m.CFG.AwsConfig.BucketName),
			Key:         aws.String(objectKey),
			Body:        strings.NewReader(dataStr),
			ContentType: &contentType,
		})
	if err != nil {
		log.Logger.Error(err)
		return "", err
	}
	log.Logger.Info("Successfully uploaded bytes: ", uploadInfo)
	return m.CFG.AwsConfig.BucketUrl + "/" + objectKey, nil
}

// 上传外部传入文件
func (m *AwsTool) UploadFile(file multipart.File, header *multipart.FileHeader, key string) (string, string, error) {
	s3Client := m.S3Client
	// 上传阿里云路径 文件名格式 自己可以改 建议保证唯一性
	// yunFileTmpPath := filepath.Join("uploads", time.Now().Format("2006-01-02")) + "/" + file.Filename
	yunFileTmpPath := "uploads" + "/" + time.Now().Format("2006-01-02") + "/" + key + "/" + header.Filename

	// 上传文件流。
	uploadInfo, err := s3Client.PutObject(ctx,
		&s3.PutObjectInput{
			Bucket: aws.String(m.CFG.AwsConfig.BucketName),
			Key:    aws.String(yunFileTmpPath),
			Body:   file,
		})
	if err != nil {
		log.Logger.Error(err)
	}
	defer file.Close()
	log.Logger.Info("Successfully uploaded bytes: ", uploadInfo)
	if err != nil {
		log.Logger.Error("function formUploader.Put() Failed", zap.Any("err", err.Error()))
		return "", "", errors.New("function formUploader.Put() Failed, err:" + err.Error())
	}

	return m.CFG.AwsConfig.BucketUrl + "/" + yunFileTmpPath, yunFileTmpPath, nil
}

// 上传本地文件
func (m *AwsTool) UploadLocalFile(localFile string, objectKey string) (string, error) {
	s3Client := m.S3Client
	// Upload the glb file
	contentType := "application/octet-stream"
	// Upload the glb file with FPutObject
	file, err := os.Open(localFile)
	if err != nil {
		log.Logger.Errorf("os.Open - filename: %v, err: %v", localFile, err)
	}
	defer file.Close()
	info, err := s3Client.PutObject(ctx,
		&s3.PutObjectInput{
			Bucket: aws.String(m.CFG.AwsConfig.BucketName),
			Key:    aws.String(objectKey),
			//ACL:                aws.String(S3_ACL),
			Body: file, // bytes.NewReader(buffer),
			// ContentDisposition: aws.String("attachment"),
			// ContentLength:      aws.Int64(int64(len(buffer))),
			ContentType: aws.String(contentType),
			// ServerSideEncryption: aws.String("AES256"),
		})
	if err != nil {
		log.Logger.Error(err)
	}

	log.Logger.Info("Successfully uploaded bytes: ", info)
	if err != nil {
		log.Logger.Error("function formUploader.Put() Failed", zap.Any("err", err.Error()))
		return "", errors.New("function formUploader.Put() Failed, err:" + err.Error())
	}

	return m.CFG.AwsConfig.BucketUrl + objectKey, nil
}
func (m *AwsTool) EnsureDir(fileName string) {
	dirName := filepath.Dir(fileName)
	if _, serr := os.Stat(dirName); serr != nil {
		merr := os.MkdirAll(dirName, os.ModePerm)
		if merr != nil {
			log.Logger.Error(merr)
		}
	}
}

// 下载大文件
func (m *AwsTool) DownloadLargeObject(bucketName string, objectKey string, jobId string) (fpath string, err error) {
	var partMiBs int64 = 10
	downloader := manager.NewDownloader(m.S3Client, func(d *manager.Downloader) {
		d.PartSize = partMiBs * 1024 * 1024
	})
	fileName := filepath.Base(objectKey)
	fpath = "/tmp/" + jobId + "/" + fileName
	m.EnsureDir(fpath)
	localFile, err := os.Create(fpath)
	if err != nil {
		log.Logger.Error(err)
		return
	}
	defer localFile.Close()
	_, err = downloader.Download(context.TODO(), localFile, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		log.Logger.Errorf("Couldn't download large object from %v:%v. Here's why: %v\n",
			bucketName, objectKey, err)
	}
	return
}

// 下载文件
func (m *AwsTool) DownloadFile(bucketName string, key string, jobId string) (fpath string, err error) {
	s3Client := m.S3Client
	fileName := filepath.Base(key)
	fpath = "/tmp/" + jobId + "/" + fileName
	results, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Logger.Error(err)
		return
	}
	defer results.Body.Close()
	m.EnsureDir(fpath)
	localFile, err := os.Create(fpath)
	if err != nil {
		log.Logger.Error(err)
		return
	}
	defer localFile.Close()
	if _, err = io.Copy(localFile, results.Body); err != nil {
		log.Logger.Error("function save oss file to local Filed", zap.Any("err", err.Error()))
		return
	}
	return
}

// 删除文件
func (m *AwsTool) DeleteFile(key string) error {
	s3Client := m.S3Client
	_, err := s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(m.CFG.AwsConfig.BucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Logger.Error("function bucketManager.Delete() Filed", zap.Any("err", err.Error()))
		return errors.New("function bucketManager.Delete() Filed, err:" + err.Error())
	}

	return nil
}

// 获取下载授权url
func (m *AwsTool) GetObjectUrl(bucket string, objectKey string) (string, http.Header, error) {
	s3Client := m.S3Client
	// 路径 统一目录 + tenant+userid+date+文件名
	presignClient := s3.NewPresignClient(s3Client)
	request, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectKey),
	}, s3.WithPresignExpires(time.Minute*60))

	if err != nil {
		log.Logger.Error("function GetObjectUrl() Failed", zap.Any("err", err.Error()))
		return "", nil, errors.New("function GetObjectUrl() Failed, err:" + err.Error())
	}
	return request.URL, request.SignedHeader, err
}

// 获取上传授权url
func (m *AwsTool) UploadUrl(bizPath string, fileName string) (string, http.Header, error) {
	if len(fileName) == 0 {
		return "", nil, errors.New("function AwsTool.UploadUrl() Failed, err: 参数不允许留空")
	}
	s3Client := m.S3Client
	// 路径 统一目录 + tenant+userid+date+文件名
	yunFileTmpPath := "uploadedbyuser"
	if len(bizPath) > 0 {
		yunFileTmpPath += "/" + bizPath
	}
	yunFileTmpPath += "/" + time.Now().Format("2006-01-02") + "/" + fileName
	presignClient := s3.NewPresignClient(s3Client)
	request, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(m.CFG.AwsConfig.BucketName),
		Key:    aws.String(yunFileTmpPath),
	}, s3.WithPresignExpires(time.Minute*60))

	if err != nil {
		log.Logger.Error("function AwsTool.UploadUrl() Failed", zap.Any("err", err.Error()))
		return "", nil, errors.New("function AwsTool.UploadUrl() Failed, err:" + err.Error())
	}
	return request.URL, request.SignedHeader, err
}

// 获取指定bucket 上传授权url
func (m *AwsTool) UploadUrlWithBucket(bucketName string, bizPath string, fileName string) (string, http.Header, error) {
	if len(bizPath) == 0 && len(fileName) == 0 {
		return "", nil, errors.New("function AwsTool.UploadUrl() Failed, err: 参数不允许留空")
	}
	s3Client := m.S3Client
	defaultBecket := m.CFG.AwsConfig.BucketName
	if len(bucketName) == 0 {
		bucketName = defaultBecket
	}
	found, err := BucketExists(m.S3Client, bucketName)
	if err != nil && !found {
		log.Logger.Error(err)
		CreateBucket(m.S3Client, bucketName)
	}
	yunFileTmpPath := bizPath + "/" + time.Now().Format("2006-01-02") + "/" + fileName

	presignClient := s3.NewPresignClient(s3Client)
	request, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(yunFileTmpPath),
	}, s3.WithPresignExpires(time.Minute*60))

	if err != nil {
		log.Logger.Error("function AwsTool.UploadUrl() Failed", zap.Any("err", err.Error()))
		return "", nil, errors.New("function AwsTool.UploadUrl() Failed, err:" + err.Error())
	}
	return request.URL, request.SignedHeader, err
}

// 判定对象是否存在
func (m *AwsTool) DoesObjectExsist(bucket string, key string) (bool, error) {
	_, err := m.S3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		var responseError *awshttp.ResponseError
		if errors.As(err, &responseError) && responseError.ResponseError.HTTPStatusCode() == http.StatusNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// 生成每个part的上传url
func (m *AwsTool) GeneratePartPresignedUrl(bizPath string, fileName string, uploadId string, partNumber int) (string, error) {
	yunFileTmpPath := "uploadedbyuser"
	if len(bizPath) > 0 {
		yunFileTmpPath += "/" + bizPath
	}
	yunFileTmpPath += "/" + time.Now().Format("2006-01-02") + "/" + fileName
	presignClient := s3.NewPresignClient(m.S3Client)

	request, err := presignClient.PresignUploadPart(ctx, &s3.UploadPartInput{
		Bucket:     aws.String(m.CFG.AwsConfig.BucketName),
		Key:        aws.String(yunFileTmpPath),
		PartNumber: int32(partNumber),
		UploadId:   &uploadId,
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(1 * int64(time.Hour))
	})
	return request.URL, err
}

// 初始化分片上传分片
func (m *AwsTool) NewMultipartUpload(bucket string, object string) (uploadID string, err error) {
	result, err := m.S3Client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(object),
	})
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	uploadID = *result.UploadId
	return
}

// 获取已经上传成功的parts
func (m *AwsTool) ListObjectParts(bucket string, object string, uploadID string, maxParts int32) (result *s3.ListPartsOutput, err error) {
	result, err = m.S3Client.ListParts(ctx, &s3.ListPartsInput{
		Bucket:   aws.String(bucket),
		Key:      aws.String(object),
		UploadId: aws.String(uploadID),
		MaxParts: maxParts,
	})
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	return
}

// 分片没有上去完成，丢失分片直接放弃上传
func (m *AwsTool) AbortMultipartUpload(bucket string, key string, uploadId string) (result *s3.AbortMultipartUploadOutput, err error) {
	input := &s3.AbortMultipartUploadInput{
		Bucket:   aws.String(bucket),
		Key:      aws.String(key),
		UploadId: aws.String(uploadId),
	}
	result, err = m.S3Client.AbortMultipartUpload(ctx, input)
	if err != nil {
		log.Logger.Panic(err.Error())
	}
	return
}

// 完成 分片合并
func (m *AwsTool) CompleteMultipartUpload(bucket string, key string, uploadId string, maxParts int32) (err error) {
	listPartsOutput, _ := m.ListObjectParts(bucket, key, uploadId, maxParts)
	if listPartsOutput != nil && len(listPartsOutput.Parts) == int(maxParts) {
		completedParts := make([]types.CompletedPart, maxParts)
		for _, p := range listPartsOutput.Parts {
			completedParts = append(completedParts, types.CompletedPart{
				ETag:       p.ETag,
				PartNumber: p.PartNumber,
			})
		}
		input := &s3.CompleteMultipartUploadInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
			MultipartUpload: &types.CompletedMultipartUpload{
				Parts: completedParts,
			},
			UploadId: aws.String(uploadId),
		}
		_, err := m.S3Client.CompleteMultipartUpload(ctx, input)
		if err != nil {
			log.Logger.Panic(err.Error())
		}
	}
	return
}

const bufferSize = 65536

// MD5sum returns MD5 checksum of filename
func (m *AwsTool) MD5sum(filename string) (string, error) {
	if info, err := os.Stat(filename); err != nil {
		return "", err
	} else if info.IsDir() {
		return "", nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	for buf, reader := make([]byte, bufferSize), bufio.NewReader(file); ; {
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		hash.Write(buf[:n])
	}

	checksum := fmt.Sprintf("%x", hash.Sum(nil))
	return checksum, nil
}
