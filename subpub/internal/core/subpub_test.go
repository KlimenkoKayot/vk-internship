package core_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/klimenkokayot/vk-internship/subpub/internal/core"
	"github.com/klimenkokayot/vk-internship/subpub/testutils/mocks"
)

func setupMockLogger(mockLogger *mocks.MockLogger) {
	mockLogger.EXPECT().OK(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Warn(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()

	mockLogger.EXPECT().Fatal(gomock.Any(), gomock.Any()).Times(0)

	mockLogger.EXPECT().WithFields(gomock.Any()).Return(mockLogger).AnyTimes()
	mockLogger.EXPECT().WithLayer(gomock.Any()).Return(mockLogger).AnyTimes()
}

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUUIDGenerator := mocks.NewMockUUIDGenerator(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	setupMockLogger(mockLogger)

	core.NewSubPub(mockUUIDGenerator, mockLogger)
}

func TestSubscribe(t *testing.T) {
	t.Run("successful subscription", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUUIDGenerator := mocks.NewMockUUIDGenerator(ctrl)
		mockUUIDGenerator.EXPECT().NewString().Return("test-uuid")

		mockLogger := mocks.NewMockLogger(ctrl)
		setupMockLogger(mockLogger)

		subPub := core.NewSubPub(mockUUIDGenerator, mockLogger)

		// Используем простой обработчик без логики
		testMessageHandler := func(msg interface{}) {}

		sub, err := subPub.Subscribe("test", testMessageHandler)
		if err != nil {
			t.Fatalf("Subscribe failed: %v", err)
		}
		defer sub.Unsubscribe()

		// Проверяем публикацию (без проверки обработки в горутине)
		err = subPub.Publish("test", "test message")
		if err != nil {
			t.Errorf("Publish failed: %v", err)
		}
	})

	t.Run("closed subpub try subscription", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUUIDGenerator := mocks.NewMockUUIDGenerator(ctrl)
		mockLogger := mocks.NewMockLogger(ctrl)
		setupMockLogger(mockLogger)

		subPub := core.NewSubPub(mockUUIDGenerator, mockLogger)

		// Закрываем subpub
		subPub.Close(nil)

		// Используем простой обработчик без логики
		testMessageHandler := func(msg interface{}) {}

		_, err := subPub.Subscribe("test", testMessageHandler)
		if err != core.ErrSubPubClosed {
			t.Fatalf("Subscribe failed: %v", err)
		}
	})

	t.Run("subscription panic recover", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockLogger := mocks.NewMockLogger(ctrl)
		setupMockLogger(mockLogger)

		// Вызываем панику нерабочим UUIDGenerator`ом
		subPub := core.NewSubPub(nil, mockLogger)

		// Используем простой обработчик без логики
		testMessageHandler := func(msg interface{}) {}

		_, err := subPub.Subscribe("test", testMessageHandler)
		if err == nil {
			t.Fatalf("Subscribe failed: %v", err)
		}
	})
}

func TestClose(t *testing.T) {
	t.Run("default close", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUUIDGenerator := mocks.NewMockUUIDGenerator(ctrl)
		mockLogger := mocks.NewMockLogger(ctrl)
		setupMockLogger(mockLogger)

		subpub := core.NewSubPub(mockUUIDGenerator, mockLogger)
		subpub.Close(nil)
	})

	t.Run("close with subscribers", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUUIDGenerator := mocks.NewMockUUIDGenerator(ctrl)
		mockUUIDGenerator.EXPECT().NewString().Return("aboba").Times(3)
		mockLogger := mocks.NewMockLogger(ctrl)
		setupMockLogger(mockLogger)

		subpub := core.NewSubPub(mockUUIDGenerator, mockLogger)
		testCallback := func(msg interface{}) {
			return
		}
		_, err := subpub.Subscribe("test-1", testCallback)
		if err != nil {
			t.Fatalf("error subscribe")
		}
		_, err = subpub.Subscribe("test-2", testCallback)
		if err != nil {
			t.Fatalf("error subscribe")
		}
		_, err = subpub.Subscribe("test-3", testCallback)
		if err != nil {
			t.Fatalf("error subscribe")
		}
		subpub.Close(nil)
	})
}

func TestPublish(t *testing.T) {
	t.Run("publish with 1 subscriber", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUUIDGenerator := mocks.NewMockUUIDGenerator(ctrl)
		mockUUIDGenerator.EXPECT().NewString().Return("uuid").Times(1)
		mockLogger := mocks.NewMockLogger(ctrl)
		setupMockLogger(mockLogger)

		subpub := core.NewSubPub(mockUUIDGenerator, mockLogger)
		defer subpub.Close(nil)

		testCallback := func(msg interface{}) {
			return
		}
		_, err := subpub.Subscribe("test", testCallback)
		if err != nil {
			t.Fatalf("error subscribe")
		}

		subpub.Publish("test", "Hello, World!")
	})

	t.Run("publish with 0 subscriber", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUUIDGenerator := mocks.NewMockUUIDGenerator(ctrl)
		mockLogger := mocks.NewMockLogger(ctrl)
		setupMockLogger(mockLogger)

		subpub := core.NewSubPub(mockUUIDGenerator, mockLogger)
		defer subpub.Close(nil)

		subpub.Publish("test", "Hello, World!")
	})

	t.Run("publish in closed subpub", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUUIDGenerator := mocks.NewMockUUIDGenerator(ctrl)
		mockLogger := mocks.NewMockLogger(ctrl)
		setupMockLogger(mockLogger)

		subpub := core.NewSubPub(mockUUIDGenerator, mockLogger)
		// Закрываем шину
		subpub.Close(nil)

		// Попытка публикации
		subpub.Publish("test", "Hello, World!")
	})
}
