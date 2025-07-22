// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build darwin

package main

import (
	"bytes"
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/runfinch/finch/pkg/config"
	"github.com/runfinch/finch/pkg/mocks"
)

func TestNewSettingsMCommand(t *testing.T) {
	t.Parallel()

	cmd := newSettingsVMCommand(nil, nil, nil, nil)
	assert.Equal(t, cmd.Name(), "settings")
}

func TestSettingsVMAction_runAdapter(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		wantErr error
		command *cobra.Command
		args    []string
		mockSvc func(
			*mocks.LimaConfigApplier,
			afero.Fs,
		)
	}{
		{
			name:    "should configure the instance for valid CPU and memory values",
			wantErr: nil,
			command: &cobra.Command{
				Use: "settings",
			},
			args: []string{
				"--cpus=1",
				"--memory=2GiB",
			},
			mockSvc: func(
				lca *mocks.LimaConfigApplier,
				fs afero.Fs,
			) {
				finchConfigPath := "/config.yaml"
				data := "cpus: 2\nmemory: 6GiB"
				require.NoError(t, afero.WriteFile(fs, finchConfigPath, []byte(data), 0o600))

				lca.EXPECT().GetFinchConfigPath().Return(finchConfigPath)
			},
		},
		{
			name:    "should configure the instance for valid CPU value",
			wantErr: nil,
			command: &cobra.Command{
				Use: "settings",
			},
			args: []string{
				"--cpus=1",
			},
			mockSvc: func(
				lca *mocks.LimaConfigApplier,
				fs afero.Fs,
			) {
				finchConfigPath := "/config.yaml"
				data := "cpus: 2\nmemory: 6GiB"
				require.NoError(t, afero.WriteFile(fs, finchConfigPath, []byte(data), 0o600))

				lca.EXPECT().GetFinchConfigPath().Return(finchConfigPath)
			},
		},
		{
			name:    "should configure the instance for valid memory value",
			wantErr: nil,
			command: &cobra.Command{
				Use: "settings",
			},
			args: []string{
				"--memory=2GiB",
			},
			mockSvc: func(
				lca *mocks.LimaConfigApplier,
				fs afero.Fs,
			) {
				finchConfigPath := "/config.yaml"
				data := "cpus: 2\nmemory: 6GiB"
				require.NoError(t, afero.WriteFile(fs, finchConfigPath, []byte(data), 0o600))

				lca.EXPECT().GetFinchConfigPath().Return(finchConfigPath)
			},
		},
		{
			name:    "should show settings --help when no flags are provided",
			wantErr: nil,
			command: &cobra.Command{
				Use: "settings",
			},
			args: []string{},
			mockSvc: func(
				_ *mocks.LimaConfigApplier,
				_ afero.Fs,
			) {
				// no expectations since help should be shown
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			logger := mocks.NewLogger(ctrl)
			lca := mocks.NewLimaConfigApplier(ctrl)
			fs := afero.NewMemMapFs()
			stdout := bytes.Buffer{}

			tc.mockSvc(lca, fs)

			cmd := newSettingsVMCommand(logger, lca, fs, &stdout)
			cmd.SetArgs(tc.args)
			err := cmd.Execute()
			assert.Equal(t, err, tc.wantErr)
		})
	}
}

func TestSettingsVMAction_run(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		wantErr          string
		wantStatusOutput string
		wantWarnOutput   string
		mockSvc          func(
			*mocks.LimaConfigApplier,
			afero.Fs,
		)
		opts config.VMConfigOpts
	}{
		{
			name:             "should update vm settings",
			wantErr:          "",
			wantStatusOutput: "Configurations have been successfully updated.\n",
			wantWarnOutput:   "",
			mockSvc: func(
				lca *mocks.LimaConfigApplier,
				fs afero.Fs,
			) {
				finchConfigPath := "/config.yaml"
				data := "cpus: 2\nmemory: 6GiB"
				require.NoError(t, afero.WriteFile(fs, finchConfigPath, []byte(data), 0o600))

				lca.EXPECT().GetFinchConfigPath().Return(finchConfigPath)
			},
			opts: config.VMConfigOpts{
				CPUs:          1,
				CPUsChanged:   true,
				Memory:        "2GiB",
				MemoryChanged: true,
			},
		},
		{
			name:             "should return an error if the configuration of CPU is invalid",
			wantErr:          "failed to validate config file: specified number of CPUs (0) must be greater than 0",
			wantStatusOutput: "",
			wantWarnOutput:   "",
			mockSvc: func(
				lca *mocks.LimaConfigApplier,
				fs afero.Fs,
			) {
				finchConfigPath := "/config.yaml"
				data := "cpus: 2\nmemory: 6GiB"
				require.NoError(t, afero.WriteFile(fs, finchConfigPath, []byte(data), 0o600))

				lca.EXPECT().GetFinchConfigPath().Return(finchConfigPath)
			},
			opts: config.VMConfigOpts{
				CPUs:          0,
				CPUsChanged:   true,
				Memory:        "2GiB",
				MemoryChanged: true,
			},
		},
		{
			name:             "should return an error if the configuration of memory is invalid",
			wantErr:          "failed to validate config file: failed to parse memory to uint: invalid suffix: 'gi'",
			wantStatusOutput: "",
			wantWarnOutput:   "",
			mockSvc: func(
				lca *mocks.LimaConfigApplier,
				fs afero.Fs,
			) {
				finchConfigPath := "/config.yaml"
				data := "cpus: 2\nmemory: 6GiB"
				require.NoError(t, afero.WriteFile(fs, finchConfigPath, []byte(data), 0o600))

				lca.EXPECT().GetFinchConfigPath().Return(finchConfigPath)
			},
			opts: config.VMConfigOpts{
				CPUs:          2,
				CPUsChanged:   true,
				Memory:        "2gi",
				MemoryChanged: true,
			},
		},
		{
			name:             "should not return an error if the configuration of CPU and memory matches existing config",
			wantErr:          "",
			wantStatusOutput: "",
			wantWarnOutput:   "Provided flags match existing settings, no changes made.",
			mockSvc: func(
				lca *mocks.LimaConfigApplier,
				fs afero.Fs,
			) {
				finchConfigPath := "/config.yaml"
				data := "cpus: 2\nmemory: 6GiB"
				require.NoError(t, afero.WriteFile(fs, finchConfigPath, []byte(data), 0o600))

				lca.EXPECT().GetFinchConfigPath().Return(finchConfigPath)
			},
			opts: config.VMConfigOpts{
				CPUs:          2,
				CPUsChanged:   true,
				Memory:        "6GiB",
				MemoryChanged: true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			logger := mocks.NewLogger(ctrl)
			lca := mocks.NewLimaConfigApplier(ctrl)
			fs := afero.NewMemMapFs()
			stdout := bytes.Buffer{}

			tc.mockSvc(lca, fs)
			if tc.wantWarnOutput != "" {
				logger.EXPECT().Warnln(tc.wantWarnOutput)
			}

			err := newSettingsVMAction(logger, lca, fs, &stdout).run(tc.opts)
			// Use string comparison since we don't save the error types
			if err == nil {
				assert.Equal(t, tc.wantErr, "")
			} else {
				assert.Equal(t, tc.wantErr, err.Error())
			}
			assert.Equal(t, tc.wantStatusOutput, stdout.String())
		})
	}
}
