package core_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
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
		ctx := context.Background()
		defer subPub.Close(ctx)

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
		ctx := context.Background()
		subPub.Close(ctx)

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
		ctx := context.Background()
		defer subPub.Close(ctx)

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
		ctx := context.Background()
		subpub.Close(ctx)

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
		ctx := context.Background()
		subpub.Close(ctx)
	})

	t.Run("close with context cancel", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Настраиваем моки
		mockUUIDGenerator := mocks.NewMockUUIDGenerator(ctrl)
		mockLogger := mocks.NewMockLogger(ctrl)
		setupMockLogger(mockLogger)

		// Создаем SubPub
		subpub := core.NewSubPub(mockUUIDGenerator, mockLogger)

		// Создаем отменённый контекст
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Немедленно отменяем

		// Вызываем Close с отменённым контекстом
		err := subpub.Close(ctx)

		// Проверяем, что вернулась ошибка отмены контекста
		if err != context.Canceled {
			t.Errorf("Ожидалась ошибка context.Canceled, получили: %v", err)
		}
	})

	t.Run("close with context cancel in progress with 10 subs and 1 topic", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Настраиваем моки
		mockUUIDGenerator := mocks.NewMockUUIDGenerator(ctrl)
		mockUUIDGenerator.EXPECT().NewString().DoAndReturn(func() string {
			return uuid.New().String()
		}).Times(10)
		mockLogger := mocks.NewMockLogger(ctrl)
		setupMockLogger(mockLogger)

		// Создаем SubPub с 10 подписками
		subpub := core.NewSubPub(mockUUIDGenerator, mockLogger)
		for i := 0; i < 10; i++ {
			_, err := subpub.Subscribe("test", func(msg interface{}) {})
			if err != nil {
				t.Fatalf("Ошибка подписки: %v", err)
			}
		}

		ctx, cancel := context.WithCancel(context.Background())
		done := make(chan struct{})

		go func() {
			err := subpub.Close(ctx)
			if err != context.Canceled {
				t.Errorf("Ожидалась ошибка context.Canceled, получили: %v", err)
			}
			close(done)
		}()

		// Даем время начать выполнение Close
		time.Sleep(10 * time.Millisecond)

		// Отменяем контекст
		cancel()

		select {
		case <-done:
		case <-time.After(100 * time.Millisecond):
			t.Error("Timeout waiting for Close to complete")
		}
	})

	t.Run("close with context cancel in progress with 5000 subs and 5000 topics", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Настраиваем моки
		mockUUIDGenerator := mocks.NewMockUUIDGenerator(ctrl)
		mockUUIDGenerator.EXPECT().NewString().DoAndReturn(func() string {
			return uuid.New().String()
		}).Times(5000)
		mockLogger := mocks.NewMockLogger(ctrl)
		setupMockLogger(mockLogger)

		// Создаем SubPub с 5000 подписками
		subpub := core.NewSubPub(mockUUIDGenerator, mockLogger)
		for i := 0; i < 5000; i++ {
			topicName := uuid.New().String()
			_, err := subpub.Subscribe(topicName, func(msg interface{}) {})
			if err != nil {
				t.Fatalf("Ошибка подписки: %v", err)
			}
		}

		ctx, cancel := context.WithCancel(context.Background())
		done := make(chan struct{})

		go func() {
			err := subpub.Close(ctx)
			if err != context.Canceled {
				t.Errorf("Ожидалась ошибка context.Canceled, получили: %v", err)
			}
			close(done)
		}()

		// Даем время начать выполнение Close
		time.Sleep(5 * time.Millisecond)

		// Отменяем контекст
		cancel()

		select {
		case <-done:
		case <-time.After(100 * time.Millisecond):
			t.Error("Timeout waiting for Close to complete")
		}
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
		ctx := context.Background()
		defer subpub.Close(ctx)

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
		ctx := context.Background()
		defer subpub.Close(ctx)

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
		ctx := context.Background()
		subpub.Close(ctx)

		// Попытка публикации
		subpub.Publish("test", "Hello, World!")
	})
}
