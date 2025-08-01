// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/runfinch/finch/pkg/config (interfaces: LimaConfigApplier)
//
// Generated by this command:
//
//	mockgen -copyright_file=../../copyright_header -destination=../mocks/pkg_config_lima_config_applier.go -package=mocks -mock_names LimaConfigApplier=LimaConfigApplier . LimaConfigApplier
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// LimaConfigApplier is a mock of LimaConfigApplier interface.
type LimaConfigApplier struct {
	ctrl     *gomock.Controller
	recorder *LimaConfigApplierMockRecorder
	isgomock struct{}
}

// LimaConfigApplierMockRecorder is the mock recorder for LimaConfigApplier.
type LimaConfigApplierMockRecorder struct {
	mock *LimaConfigApplier
}

// NewLimaConfigApplier creates a new mock instance.
func NewLimaConfigApplier(ctrl *gomock.Controller) *LimaConfigApplier {
	mock := &LimaConfigApplier{ctrl: ctrl}
	mock.recorder = &LimaConfigApplierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *LimaConfigApplier) EXPECT() *LimaConfigApplierMockRecorder {
	return m.recorder
}

// ConfigureDefaultLimaYaml mocks base method.
func (m *LimaConfigApplier) ConfigureDefaultLimaYaml() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConfigureDefaultLimaYaml")
	ret0, _ := ret[0].(error)
	return ret0
}

// ConfigureDefaultLimaYaml indicates an expected call of ConfigureDefaultLimaYaml.
func (mr *LimaConfigApplierMockRecorder) ConfigureDefaultLimaYaml() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConfigureDefaultLimaYaml", reflect.TypeOf((*LimaConfigApplier)(nil).ConfigureDefaultLimaYaml))
}

// ConfigureOverrideLimaYaml mocks base method.
func (m *LimaConfigApplier) ConfigureOverrideLimaYaml() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConfigureOverrideLimaYaml")
	ret0, _ := ret[0].(error)
	return ret0
}

// ConfigureOverrideLimaYaml indicates an expected call of ConfigureOverrideLimaYaml.
func (mr *LimaConfigApplierMockRecorder) ConfigureOverrideLimaYaml() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConfigureOverrideLimaYaml", reflect.TypeOf((*LimaConfigApplier)(nil).ConfigureOverrideLimaYaml))
}

// GetFinchConfigPath mocks base method.
func (m *LimaConfigApplier) GetFinchConfigPath() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFinchConfigPath")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetFinchConfigPath indicates an expected call of GetFinchConfigPath.
func (mr *LimaConfigApplierMockRecorder) GetFinchConfigPath() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFinchConfigPath", reflect.TypeOf((*LimaConfigApplier)(nil).GetFinchConfigPath))
}
