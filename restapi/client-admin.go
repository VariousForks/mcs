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
	"io"
	"path/filepath"
	"runtime"

	mcCmd "github.com/minio/mc/cmd"
	"github.com/minio/mc/pkg/probe"
	"github.com/minio/minio/pkg/madmin"
)

const globalAppName = "mcs"

// NewAdminClient gives a new client interface
func NewAdminClient(url, accessKey, secretKey string) (*madmin.AdminClient, *probe.Error) {
	appName := filepath.Base(globalAppName)

	s3Client, err := s3AdminNew(&mcCmd.Config{
		HostURL:     url,
		AccessKey:   accessKey,
		SecretKey:   secretKey,
		AppName:     appName,
		AppVersion:  McsVersion,
		AppComments: []string{appName, runtime.GOOS, runtime.GOARCH},
	})
	if err != nil {
		return nil, err.Trace(url)
	}
	return s3Client, nil
}

// s3AdminNew returns an initialized minioAdmin structure. If debug is enabled,
// it also enables an internal trace transport.
var s3AdminNew = mcCmd.NewAdminFactory()

// Define MinioAdmin interface with all functions to be implemented
// by mock when testing, it should include all MinioAdmin respective api calls
// that are used within this project.
type MinioAdmin interface {
	listUsers(ctx context.Context) (map[string]madmin.UserInfo, error)
	addUser(ctx context.Context, acessKey, SecretKey string) error
	removeUser(ctx context.Context, accessKey string) error
	getUserInfo(ctx context.Context, accessKey string) (madmin.UserInfo, error)
	setUserStatus(ctx context.Context, accessKey string, status madmin.AccountStatus) error
	listGroups(ctx context.Context) ([]string, error)
	updateGroupMembers(ctx context.Context, greq madmin.GroupAddRemove) error
	getGroupDescription(ctx context.Context, group string) (*madmin.GroupDesc, error)
	setGroupStatus(ctx context.Context, group string, status madmin.GroupStatus) error
	listPolicies(ctx context.Context) (map[string][]byte, error)
	getPolicy(ctx context.Context, name string) ([]byte, error)
	removePolicy(ctx context.Context, name string) error
	addPolicy(ctx context.Context, name, policy string) error
	setPolicy(ctx context.Context, policyName, entityName string, isGroup bool) error
	getConfigKV(ctx context.Context, key string) (madmin.Targets, error)
	helpConfigKV(ctx context.Context, subSys, key string, envOnly bool) (madmin.Help, error)
	setConfigKV(ctx context.Context, kv string) (err error)
	serviceRestart(ctx context.Context) error
	serverInfo(ctx context.Context) (madmin.InfoMessage, error)
	startProfiling(ctx context.Context, profiler madmin.ProfilerType) ([]madmin.StartProfilingResult, error)
	stopProfiling(ctx context.Context) (io.ReadCloser, error)
}

// Interface implementation
//
// Define the structure of a minIO Client and define the functions that are actually used
// from minIO api.
type adminClient struct {
	client *madmin.AdminClient
}

// implements madmin.ListUsers()
func (ac adminClient) listUsers(ctx context.Context) (map[string]madmin.UserInfo, error) {
	return ac.client.ListUsers(ctx)
}

// implements madmin.AddUser()
func (ac adminClient) addUser(ctx context.Context, acessKey, secretKey string) error {
	return ac.client.AddUser(ctx, acessKey, secretKey)
}

// implements madmin.RemoveUser()
func (ac adminClient) removeUser(ctx context.Context, accessKey string) error {
	return ac.client.RemoveUser(ctx, accessKey)
}

//implements madmin.GetUserInfo()
func (ac adminClient) getUserInfo(ctx context.Context, accessKey string) (madmin.UserInfo, error) {
	return ac.client.GetUserInfo(ctx, accessKey)
}

// implements madmin.SetUserStatus()
func (ac adminClient) setUserStatus(ctx context.Context, accessKey string, status madmin.AccountStatus) error {
	return ac.client.SetUserStatus(ctx, accessKey, status)
}

// implements madmin.ListGroups()
func (ac adminClient) listGroups(ctx context.Context) ([]string, error) {
	return ac.client.ListGroups(ctx)
}

// implements madmin.UpdateGroupMembers()
func (ac adminClient) updateGroupMembers(ctx context.Context, greq madmin.GroupAddRemove) error {
	return ac.client.UpdateGroupMembers(ctx, greq)
}

// implements madmin.GetGroupDescription(group)
func (ac adminClient) getGroupDescription(ctx context.Context, group string) (*madmin.GroupDesc, error) {
	return ac.client.GetGroupDescription(ctx, group)
}

// implements madmin.SetGroupStatus(group, status)
func (ac adminClient) setGroupStatus(ctx context.Context, group string, status madmin.GroupStatus) error {
	return ac.client.SetGroupStatus(ctx, group, status)
}

// implements madmin.ListCannedPolicies()
func (ac adminClient) listPolicies(ctx context.Context) (map[string][]byte, error) {
	return ac.client.ListCannedPolicies(ctx)
}

// implements madmin.ListCannedPolicies()
func (ac adminClient) getPolicy(ctx context.Context, name string) ([]byte, error) {
	return ac.client.InfoCannedPolicy(ctx, name)
}

// implements madmin.RemoveCannedPolicy()
func (ac adminClient) removePolicy(ctx context.Context, name string) error {
	return ac.client.RemoveCannedPolicy(ctx, name)
}

// implements madmin.AddCannedPolicy()
func (ac adminClient) addPolicy(ctx context.Context, name, policy string) error {
	return ac.client.AddCannedPolicy(ctx, name, policy)
}

// implements madmin.SetPolicy()
func (ac adminClient) setPolicy(ctx context.Context, policyName, entityName string, isGroup bool) error {
	return ac.client.SetPolicy(ctx, policyName, entityName, isGroup)
}

// implements madmin.GetConfigKV()
func (ac adminClient) getConfigKV(ctx context.Context, key string) (madmin.Targets, error) {
	return ac.client.GetConfigKV(ctx, key)
}

// implements madmin.HelpConfigKV()
func (ac adminClient) helpConfigKV(ctx context.Context, subSys, key string, envOnly bool) (madmin.Help, error) {
	return ac.client.HelpConfigKV(ctx, subSys, key, envOnly)
}

// implements madmin.SetConfigKV()
func (ac adminClient) setConfigKV(ctx context.Context, kv string) (err error) {
	return ac.client.SetConfigKV(ctx, kv)
}

// implements madmin.ServiceRestart()
func (ac adminClient) serviceRestart(ctx context.Context) (err error) {
	return ac.client.ServiceRestart(ctx)
}

// implements madmin.ServerInfo()
func (ac adminClient) serverInfo(ctx context.Context) (madmin.InfoMessage, error) {
	return ac.client.ServerInfo(ctx)
}

// implements madmin.StartProfiling()
func (ac adminClient) startProfiling(ctx context.Context, profiler madmin.ProfilerType) ([]madmin.StartProfilingResult, error) {
	return ac.client.StartProfiling(ctx, profiler)
}

// implements madmin.DownloadProfilingData()
func (ac adminClient) stopProfiling(ctx context.Context) (io.ReadCloser, error) {
	return ac.client.DownloadProfilingData(ctx)
}

func newMAdminClient() (*madmin.AdminClient, error) {
	endpoint := getMinIOServer()
	accessKeyID := getAccessKey()
	secretAccessKey := getSecretKey()

	adminClient, pErr := NewAdminClient(endpoint, accessKeyID, secretAccessKey)
	if pErr != nil {
		return nil, pErr.Cause
	}
	return adminClient, nil
}
