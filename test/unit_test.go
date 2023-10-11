package test

import (
	"context"
	"fmt"
	"go-gin-restful-service/util"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func FormatString(s string) string {
	return fmt.Sprintf("Formatted: %s", s)
}
func TestFormatString(t *testing.T) {
	input := "Hello, World!"
	expectedOutput, err := os.ReadFile("expected_output.golden")
	if err != nil {
		t.Fatalf("Failed to read golden file: %v", err)
	}

	actualOutput := FormatString(input)

	if string(actualOutput) != string(expectedOutput) {
		t.Errorf("Unexpected output:\nExpected: %s\nActual: %s", expectedOutput, actualOutput)
	}
}

func TestForCmdFunc(t *testing.T) {
	results, err := util.ExecuteCommand("ping 127.0.0.1 -n 10")
	if err != nil {
		t.Fatalf("Failed to read golden file: %v", err)
	}
	fmt.Println(results)
}

type TestQuery struct {
	ID   int64
	Resp string
}

func TestEventBus(t *testing.T) {
	bus := util.ProvideBus()

	var invoked bool

	bus.AddEventListener(func(ctx context.Context, query *TestQuery) error {
		invoked = true
		return nil
	})

	err := bus.Publish(context.Background(), &TestQuery{})
	require.NoError(t, err, "unable to publish event")

	require.True(t, invoked)
}
