// This file is part of MinIO Orchestrator
// Copyright (c) 2020 MinIO, Inc.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package restapi

import (
	"context"
	"fmt"

	mc "github.com/minio/mc/cmd"
	"github.com/minio/mc/pkg/probe"
	"github.com/minio/minio-go/v6"
)

func init() {
	// All minio-go API operations shall be performed only once,
	// another way to look at this is we are turning off retries.
	minio.MaxRetry = 1
}

// Define MinioClient interface with all functions to be implemented
// by mock when testing, it should include all MinioClient respective api calls
// that are used within this project.
type MinioClient interface {
	listBucketsWithContext(ctx context.Context) ([]minio.BucketInfo, error)
	makeBucketWithContext(ctx context.Context, bucketName, location string) error
	setBucketPolicyWithContext(ctx context.Context, bucketName, policy string) error
	removeBucket(bucketName string) error
	getBucketNotification(bucketName string) (bucketNotification minio.BucketNotification, err error)
	getBucketPolicy(bucketName string) (string, error)
}

// Interface implementation
//
// Define the structure of a minIO Client and define the functions that are actually used
// from minIO api.
type minioClient struct {
	client *minio.Client
}

// implements minio.ListBucketsWithContext(ctx)
func (c minioClient) listBucketsWithContext(ctx context.Context) ([]minio.BucketInfo, error) {
	return c.client.ListBucketsWithContext(ctx)
}

// implements minio.MakeBucketWithContext(ctx, bucketName, location)
func (c minioClient) makeBucketWithContext(ctx context.Context, bucketName, location string) error {
	return c.client.MakeBucketWithContext(ctx, bucketName, location)
}

// implements minio.SetBucketPolicyWithContext(ctx, bucketName, policy)
func (c minioClient) setBucketPolicyWithContext(ctx context.Context, bucketName, policy string) error {
	return c.client.SetBucketPolicyWithContext(ctx, bucketName, policy)
}

// implements minio.RemoveBucket(bucketName)
func (c minioClient) removeBucket(bucketName string) error {
	return c.client.RemoveBucket(bucketName)
}

// implements minio.GetBucketNotification(bucketName)
func (c minioClient) getBucketNotification(bucketName string) (bucketNotification minio.BucketNotification, err error) {
	return c.client.GetBucketNotification(bucketName)
}

// implements minio.GetBucketPolicy(bucketName)
func (c minioClient) getBucketPolicy(bucketName string) (string, error) {
	return c.client.GetBucketPolicy(bucketName)
}

// Define MCS3Client interface with all functions to be implemented
// by mock when testing, it should include all mc/S3Client respective api calls
// that are used within this project.
type MCS3Client interface {
	addNotificationConfig(arn string, events []string, prefix, suffix string, ignoreExisting bool) *probe.Error
	removeNotificationConfig(arn string) *probe.Error
}

// Interface implementation
//
// Define the structure of a mc S3Client and define the functions that are actually used
// from mcS3client api.
type mcS3Client struct {
	client *mc.S3Client
}

// implements S3Client.AddNotificationConfig()
func (c mcS3Client) addNotificationConfig(arn string, events []string, prefix, suffix string, ignoreExisting bool) *probe.Error {
	return c.client.AddNotificationConfig(arn, events, prefix, suffix, ignoreExisting)
}

// implements S3Client.RemoveNotificationConfig()
func (c mcS3Client) removeNotificationConfig(arn string) *probe.Error {
	return c.client.RemoveNotificationConfig(arn)
}

// newMinioClient creates a new MinIO client to talk to the server
func newMinioClient() (*minio.Client, error) {
	endpoint := getMinIOEndpoint()
	accessKeyID := getAccessKey()
	secretAccessKey := getSecretKey()
	useSSL := getMinIOEndpointIsSecure()

	// Initialize minio client object.
	minioClient, err := minio.NewV4(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		return nil, err
	}

	return minioClient, nil
}

// newS3BucketClient creates a new mc S3Client to talk to the server based on a bucket
func newS3BucketClient(bucketName *string) (*mc.S3Client, error) {
	endpoint := getMinIOServer()
	accessKeyID := getAccessKey()
	secretAccessKey := getSecretKey()
	useSSL := getMinIOEndpointIsSecure()

	if bucketName != nil {
		endpoint += fmt.Sprintf("/%s", *bucketName)
	}
	s3Config := newS3Config(endpoint, accessKeyID, secretAccessKey, !useSSL)
	client, err := mc.S3New(s3Config)
	if err != nil {
		return nil, err.Cause
	}
	s3Client, ok := client.(*mc.S3Client)
	if !ok {
		return nil, fmt.Errorf("the provided url doesn't point to a S3 server")
	}

	return s3Client, nil
}

// newS3Config simply creates a new Config struct using the passed
// parameters.
func newS3Config(endpoint, accessKey, secretKey string, isSecure bool) *mc.Config {
	// We have a valid alias and hostConfig. We populate the
	// credentials from the match found in the config file.
	s3Config := new(mc.Config)

	s3Config.AppName = "mcs" // TODO: make this a constant
	s3Config.AppVersion = "" // TODO: get this from constant or build
	s3Config.AppComments = []string{}
	s3Config.Debug = false
	s3Config.Insecure = isSecure

	s3Config.HostURL = endpoint
	s3Config.AccessKey = accessKey
	s3Config.SecretKey = secretKey
	s3Config.Signature = "S3v4"
	return s3Config
}
