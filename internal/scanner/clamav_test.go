package scanner

import (
	"aegis/internal/objectstore"
	"aegis/mocks"
	"fmt"
	"testing"

	"go.uber.org/zap"
	// "github.com/stretchr/testify/assert"
)

type CommonTestItems struct {
	Sugar *zap.SugaredLogger
}

func ProvideCommonTestItems(t *testing.T) *CommonTestItems {
	t.Helper()
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Printf("Error creating logger: %s", err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()
	return &CommonTestItems{Sugar: sugar}
}

func TestDispatcher(t *testing.T) {
	commonTestItems := ProvideCommonTestItems(t)

	mockObjectStore := new(mocks.ObjectStore)
}
