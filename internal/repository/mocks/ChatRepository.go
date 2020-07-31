// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	model "github.com/flaambe/avito/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// ChatRepository is an autogenerated mock type for the ChatRepository type
type ChatRepository struct {
	mock.Mock
}

// FindChatByID provides a mock function with given fields: id
func (_m *ChatRepository) FindChatByID(id string) (model.Chat, error) {
	ret := _m.Called(id)

	var r0 model.Chat
	if rf, ok := ret.Get(0).(func(string) model.Chat); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(model.Chat)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindChats provides a mock function with given fields: user
func (_m *ChatRepository) FindChats(user model.User) ([]model.Chat, error) {
	ret := _m.Called(user)

	var r0 []model.Chat
	if rf, ok := ret.Get(0).(func(model.User) []model.Chat); ok {
		r0 = rf(user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Chat)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(model.User) error); ok {
		r1 = rf(user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertChat provides a mock function with given fields: name, users
func (_m *ChatRepository) InsertChat(name string, users []model.User) (string, error) {
	ret := _m.Called(name, users)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, []model.User) string); ok {
		r0 = rf(name, users)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, []model.User) error); ok {
		r1 = rf(name, users)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}