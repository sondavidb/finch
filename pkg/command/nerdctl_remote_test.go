// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build darwin || windows

// Package command_test is not named as command to avoid circular dependency (command <-> mocks).
package command_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/runfinch/finch/pkg/command"
	"github.com/runfinch/finch/pkg/mocks"
)

const (
	mockLimaHomePath = "/lima/home"
	mockLimactlPath  = "/lima/bin/limactl"
	mockQemuBinPath  = "/lima/bin"
	mockSystemPath   = "/usr/bin"
	finalPath        = mockQemuBinPath + command.EnvKeyPathJoiner + mockSystemPath
)

var mockArgs = []string{"shell", "finch"}

func TestNerdctlCmdCreator_Create(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		mockSvc func(*mocks.Logger, *mocks.CommandCreator, *mocks.Command, *mocks.NerdctlCmdCreatorSystemDeps)
		wantErr error
	}{
		{
			name:    "happy path",
			wantErr: nil,
			mockSvc: func(logger *mocks.Logger, cmdCreator *mocks.CommandCreator, cmd *mocks.Command, lcd *mocks.NerdctlCmdCreatorSystemDeps) {
				logger.EXPECT().Debugf("Creating limactl command: ARGUMENTS: %v, %s: %s", mockArgs, command.EnvKeyLimaHome, mockLimaHomePath)
				cmdCreator.EXPECT().Create(mockLimactlPath, mockArgs).Return(cmd)
				lcd.EXPECT().Environ().Return([]string{})
				lcd.EXPECT().Stdin().Return(nil)
				lcd.EXPECT().Stdout().Return(nil)
				lcd.EXPECT().Stderr().Return(nil)
				lcd.EXPECT().Env(command.EnvKeyPath).Return(mockSystemPath)
				cmd.EXPECT().SetEnv([]string{
					fmt.Sprintf("%s=%s", command.EnvKeyLimaHome, mockLimaHomePath),
					fmt.Sprintf("%s=%s", command.EnvKeyPath, finalPath),
				})
				cmd.EXPECT().SetStdin(nil)
				cmd.EXPECT().SetStdout(nil)
				cmd.EXPECT().SetStderr(nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			cmdCreator := mocks.NewCommandCreator(ctrl)
			cmd := mocks.NewCommand(ctrl)
			logger := mocks.NewLogger(ctrl)
			lcd := mocks.NewNerdctlCmdCreatorSystemDeps(ctrl)
			tc.mockSvc(logger, cmdCreator, cmd, lcd)
			command.NewNerdctlCmdCreator(cmdCreator, logger, mockLimaHomePath, mockLimactlPath, mockQemuBinPath, lcd).Create(mockArgs...)
		})
	}
}

func TestNerdctlCmdCreator_CreateWithoutStdio(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		mockSvc func(*mocks.Logger, *mocks.CommandCreator, *mocks.Command, *mocks.NerdctlCmdCreatorSystemDeps)
		wantErr error
	}{
		{
			name:    "happy path",
			wantErr: nil,
			mockSvc: func(logger *mocks.Logger, cmdCreator *mocks.CommandCreator, cmd *mocks.Command, lcd *mocks.NerdctlCmdCreatorSystemDeps) {
				logger.EXPECT().Debugf("Creating limactl command: ARGUMENTS: %v, %s: %s", mockArgs, command.EnvKeyLimaHome, mockLimaHomePath)
				cmdCreator.EXPECT().Create(mockLimactlPath, mockArgs).Return(cmd)
				lcd.EXPECT().Environ().Return([]string{})
				lcd.EXPECT().Env(command.EnvKeyPath).Return(mockSystemPath)
				cmd.EXPECT().SetEnv([]string{
					fmt.Sprintf("%s=%s", command.EnvKeyLimaHome, mockLimaHomePath),
					fmt.Sprintf("%s=%s", command.EnvKeyPath, finalPath),
				})
				cmd.EXPECT().SetStdin(nil)
				cmd.EXPECT().SetStdout(nil)
				cmd.EXPECT().SetStderr(nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			cmdCreator := mocks.NewCommandCreator(ctrl)
			cmd := mocks.NewCommand(ctrl)
			logger := mocks.NewLogger(ctrl)
			lcd := mocks.NewNerdctlCmdCreatorSystemDeps(ctrl)
			tc.mockSvc(logger, cmdCreator, cmd, lcd)
			command.NewNerdctlCmdCreator(cmdCreator, logger, mockLimaHomePath, mockLimactlPath, mockQemuBinPath, lcd).
				CreateWithoutStdio(mockArgs...)
		})
	}
}

func TestNerdctlCmdCreator_RunWithReplacingStdout(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		mockSvc  func(*mocks.Logger, *mocks.CommandCreator, *mocks.NerdctlCmdCreatorSystemDeps, *gomock.Controller, string, *os.File)
		wantErr  error
		stdoutRs []command.Replacement
		inOut    string
		outOut   string
	}{
		{
			name:     "happy path",
			wantErr:  nil,
			stdoutRs: []command.Replacement{{Source: "s1", Target: "t1"}, {Source: "s3", Target: "t3"}, {Source: "s6", Target: "t6"}},
			inOut:    "s1 s2 ,s3 /s4 s1.s5",
			outOut:   "t1 s2 ,t3 /s4 t1.s5",
			mockSvc: func(logger *mocks.Logger, cmdCreator *mocks.CommandCreator,
				lcd *mocks.NerdctlCmdCreatorSystemDeps, ctrl *gomock.Controller, inOut string, f *os.File,
			) {
				logger.EXPECT().Debugf("Creating limactl command: ARGUMENTS: %v, %s: %s", mockArgs, command.EnvKeyLimaHome, mockLimaHomePath)
				cmd := mocks.NewCommand(ctrl)
				cmdCreator.EXPECT().Create(mockLimactlPath, mockArgs).Return(cmd)
				lcd.EXPECT().Environ().Return([]string{})
				lcd.EXPECT().Stdin().Return(nil)
				lcd.EXPECT().Stderr().Return(nil)
				lcd.EXPECT().Env(command.EnvKeyPath).Return(mockSystemPath)
				cmd.EXPECT().SetEnv([]string{
					fmt.Sprintf("%s=%s", command.EnvKeyLimaHome, mockLimaHomePath),
					fmt.Sprintf("%s=%s", command.EnvKeyPath, finalPath),
				})
				cmd.EXPECT().SetStdin(nil)
				var stdoutBuf *bytes.Buffer
				cmd.EXPECT().SetStdout(gomock.Any()).Do(func(buf *bytes.Buffer) {
					stdoutBuf = buf
				})
				cmd.EXPECT().SetStderr(nil)
				cmd.EXPECT().Run().Do(func() {
					stdoutBuf.Write([]byte(inOut))
				})
				lcd.EXPECT().Stdout().Return(f)
			},
		},
		{
			name:     "overlapped replacements",
			wantErr:  nil,
			stdoutRs: []command.Replacement{{Source: "s1", Target: "s2"}, {Source: "s2", Target: "s3"}},
			inOut:    "s1 s2 ,s3 /s4 s1.s5",
			outOut:   "s3 s3 ,s3 /s4 s3.s5",
			mockSvc: func(logger *mocks.Logger, cmdCreator *mocks.CommandCreator,
				lcd *mocks.NerdctlCmdCreatorSystemDeps, ctrl *gomock.Controller, inOut string, f *os.File,
			) {
				logger.EXPECT().Debugf("Creating limactl command: ARGUMENTS: %v, %s: %s",
					mockArgs, command.EnvKeyLimaHome, mockLimaHomePath)
				cmd := mocks.NewCommand(ctrl)
				cmdCreator.EXPECT().Create(mockLimactlPath, mockArgs).Return(cmd)
				lcd.EXPECT().Environ().Return([]string{})
				lcd.EXPECT().Stdin().Return(nil)
				lcd.EXPECT().Stderr().Return(nil)
				lcd.EXPECT().Env(command.EnvKeyPath).Return(mockSystemPath)
				cmd.EXPECT().SetEnv([]string{
					fmt.Sprintf("%s=%s", command.EnvKeyLimaHome, mockLimaHomePath),
					fmt.Sprintf("%s=%s", command.EnvKeyPath, finalPath),
				})
				cmd.EXPECT().SetStdin(nil)
				var stdoutBuf *bytes.Buffer
				cmd.EXPECT().SetStdout(gomock.Any()).Do(func(buf *bytes.Buffer) {
					stdoutBuf = buf
				})
				cmd.EXPECT().SetStderr(nil)
				cmd.EXPECT().Run().Do(func() {
					stdoutBuf.Write([]byte(inOut))
				})
				lcd.EXPECT().Stdout().Return(f)
			},
		},
		{
			name:     "empty replacements",
			wantErr:  nil,
			stdoutRs: []command.Replacement{},
			inOut:    "s1 s2 ,s3 /s4 .s5",
			outOut:   "s1 s2 ,s3 /s4 .s5",
			mockSvc: func(logger *mocks.Logger, cmdCreator *mocks.CommandCreator,
				lcd *mocks.NerdctlCmdCreatorSystemDeps, ctrl *gomock.Controller, inOut string, f *os.File,
			) {
				logger.EXPECT().Debugf("Creating limactl command: ARGUMENTS: %v, %s: %s", mockArgs, command.EnvKeyLimaHome, mockLimaHomePath)
				cmd := mocks.NewCommand(ctrl)
				cmdCreator.EXPECT().Create(mockLimactlPath, mockArgs).Return(cmd)
				lcd.EXPECT().Environ().Return([]string{})
				lcd.EXPECT().Stdin().Return(nil)
				lcd.EXPECT().Stderr().Return(nil)
				lcd.EXPECT().Env(command.EnvKeyPath).Return(mockSystemPath)
				cmd.EXPECT().SetEnv([]string{
					fmt.Sprintf("%s=%s", command.EnvKeyLimaHome, mockLimaHomePath),
					fmt.Sprintf("%s=%s", command.EnvKeyPath, finalPath),
				})
				cmd.EXPECT().SetStdin(nil)
				var stdoutBuf *bytes.Buffer
				cmd.EXPECT().SetStdout(gomock.Any()).Do(func(buf *bytes.Buffer) {
					stdoutBuf = buf
				})
				cmd.EXPECT().SetStderr(nil)
				cmd.EXPECT().Run().Do(func() {
					stdoutBuf.Write([]byte(inOut))
				})
				lcd.EXPECT().Stdout().Return(f)
			},
		},
		{
			name:     "running cmd returns error",
			wantErr:  errors.New("run cmd error"),
			stdoutRs: []command.Replacement{{Source: "source-out", Target: "target-out"}},
			inOut:    "source-out",
			outOut:   "",
			mockSvc: func(logger *mocks.Logger, cmdCreator *mocks.CommandCreator,
				lcd *mocks.NerdctlCmdCreatorSystemDeps, ctrl *gomock.Controller, _ string, _ *os.File,
			) {
				logger.EXPECT().Debugf("Creating limactl command: ARGUMENTS: %v, %s: %s", mockArgs, command.EnvKeyLimaHome, mockLimaHomePath)
				cmd := mocks.NewCommand(ctrl)
				cmdCreator.EXPECT().Create(mockLimactlPath, mockArgs).Return(cmd)
				lcd.EXPECT().Environ().Return([]string{})
				lcd.EXPECT().Stdin().Return(nil)
				lcd.EXPECT().Stderr().Return(nil)
				lcd.EXPECT().Env(command.EnvKeyPath).Return(mockSystemPath)
				cmd.EXPECT().SetEnv([]string{
					fmt.Sprintf("%s=%s", command.EnvKeyLimaHome, mockLimaHomePath),
					fmt.Sprintf("%s=%s", command.EnvKeyPath, finalPath),
				})
				cmd.EXPECT().SetStdin(nil)
				cmd.EXPECT().SetStdout(gomock.Any())
				cmd.EXPECT().SetStderr(nil)
				cmd.EXPECT().Run().Return(errors.New("run cmd error"))
			},
		},
		{
			name:     "writing to stdout returns error",
			wantErr:  fs.ErrInvalid,
			stdoutRs: []command.Replacement{{Source: "source-out", Target: "target-out"}},
			inOut:    "source-out",
			outOut:   "",
			mockSvc: func(logger *mocks.Logger, cmdCreator *mocks.CommandCreator,
				lcd *mocks.NerdctlCmdCreatorSystemDeps, ctrl *gomock.Controller, inOut string, _ *os.File,
			) {
				logger.EXPECT().Debugf("Creating limactl command: ARGUMENTS: %v, %s: %s", mockArgs, command.EnvKeyLimaHome, mockLimaHomePath)
				cmd := mocks.NewCommand(ctrl)
				cmdCreator.EXPECT().Create(mockLimactlPath, mockArgs).Return(cmd)
				lcd.EXPECT().Environ().Return([]string{})
				lcd.EXPECT().Stdin().Return(nil)
				lcd.EXPECT().Stderr().Return(nil)
				lcd.EXPECT().Env(command.EnvKeyPath).Return(mockSystemPath)
				cmd.EXPECT().SetEnv([]string{
					fmt.Sprintf("%s=%s", command.EnvKeyLimaHome, mockLimaHomePath),
					fmt.Sprintf("%s=%s", command.EnvKeyPath, finalPath),
				})
				cmd.EXPECT().SetStdin(nil)
				var stdoutBuf *bytes.Buffer
				cmd.EXPECT().SetStdout(gomock.Any()).Do(func(buf *bytes.Buffer) {
					stdoutBuf = buf
				})
				cmd.EXPECT().SetStderr(nil)
				cmd.EXPECT().Run().Do(func() {
					stdoutBuf.Write([]byte(inOut))
				})
				lcd.EXPECT().Stdout().Return(nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			cmdCreator := mocks.NewCommandCreator(ctrl)
			logger := mocks.NewLogger(ctrl)
			lcd := mocks.NewNerdctlCmdCreatorSystemDeps(ctrl)

			stdoutFilepath := filepath.Clean(filepath.Join(t.TempDir(), "test"))
			stdoutFile, err := os.Create(stdoutFilepath)
			require.NoError(t, err)

			tc.mockSvc(logger, cmdCreator, lcd, ctrl, tc.inOut, stdoutFile)
			assert.Equal(t, tc.wantErr,
				command.NewNerdctlCmdCreator(
					cmdCreator,
					logger,
					mockLimaHomePath,
					mockLimactlPath,
					mockQemuBinPath,
					lcd,
				).RunWithReplacingStdout(tc.stdoutRs, mockArgs...))

			stdout, err := os.ReadFile(stdoutFilepath)
			require.NoError(t, err)
			assert.Equal(t, tc.outOut, string(stdout))
			assert.NoError(t, stdoutFile.Close())
			assert.NoError(t, os.Remove(stdoutFilepath))
		})
	}
}
