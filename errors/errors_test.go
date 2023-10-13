package errors_test

import (
	"context"
	"testing"

	"github.com/consensys-vertical-apps/platform-data-pipeline-toolkit/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	err := errors.New("failed")
	require.Error(t, err, "it is not an error")

	expected := "failed"
	assert.EqualError(t, err, expected, "error message mismatch, got %s want %s", err, expected)
}

func TestNewf(t *testing.T) {
	t.Parallel()

	err := errors.Newf("oops: %v", "failed")
	require.Error(t, err, "it is not an error")

	expected := "oops: failed"
	assert.EqualError(t, err, expected, "error message mismatch, got %s want %s", err, expected)
}

func TestWrap(t *testing.T) {
	t.Parallel()

	t.Run("Wrap with error", func(t *testing.T) {
		t.Parallel()

		err := errors.New("failed")

		errWrap := errors.Wrap(err, "oops")
		require.Error(t, errWrap, "it is not an error")

		expected := "oops: failed"
		assert.EqualError(t, errWrap, expected, "error message mismatch, got %s want %s", errWrap, expected)
	})

	t.Run("Wrap with nil", func(t *testing.T) {
		t.Parallel()

		errWrap := errors.Wrap(nil, "oops")
		require.NoError(t, errWrap, "error should be nil")
	})
}

func TestWrapf(t *testing.T) {
	t.Parallel()

	t.Run("Wrapf with error", func(t *testing.T) {
		t.Parallel()

		err := errors.New("failed")

		errWrap := errors.Wrapf(err, "oops id %d", 5)
		require.Error(t, errWrap, "it is not an error")

		expected := "oops id 5: failed"
		assert.EqualError(t, errWrap, expected, "error message mismatch, got %s want %s", errWrap, expected)
	})

	t.Run("Wrapf with nil", func(t *testing.T) {
		t.Parallel()

		errWrap := errors.Wrapf(nil, "oops id %d", 5)
		require.NoError(t, errWrap, "error should be nil")
	})
}

func TestWrapError(t *testing.T) {
	t.Parallel()

	t.Run("WrapWithError for errors", func(t *testing.T) {
		t.Parallel()

		err := errors.New("failed")
		sErr := errors.New("oops")

		errWrap := errors.WrapWithError(err, sErr)
		require.Error(t, errWrap, "it is not an error")

		expected := "oops: failed"
		assert.EqualError(t, errWrap, expected, "error message mismatch, got %s want %s", errWrap, expected)
	})

	t.Run("WrapWithError with cause nil", func(t *testing.T) {
		t.Parallel()

		sErr := errors.New("oops")

		errWrap := errors.WrapWithError(nil, sErr)
		require.Error(t, errWrap, "it is not an error")

		expected := "oops"
		assert.EqualError(t, errWrap, expected, "error message mismatch, got %s want %s", errWrap, expected)

		assert.Equal(t, sErr, errWrap)
	})

	t.Run("WrapWithError with supplied nil", func(t *testing.T) {
		t.Parallel()

		err := errors.New("oops")

		errWrap := errors.WrapWithError(err, nil)
		require.Error(t, errWrap, "it is not an error")

		expected := "oops"
		assert.EqualError(t, errWrap, expected, "error message mismatch, got %s want %s", errWrap, expected)

		assert.Equal(t, err, errWrap)
	})
}

func Test_Unwrap(t *testing.T) {
	t.Parallel()

	t.Run("Unwrap for errors.Wrap", func(t *testing.T) {
		t.Parallel()

		err := errors.New("failed")

		errWrap := errors.Wrap(err, "oops")
		require.Error(t, errWrap, "it is not an error")

		expected := "oops: failed"
		assert.EqualError(t, errWrap, expected, "error message mismatch, got %s want %s", errWrap, expected)

		uErr := errors.Unwrap(errWrap)
		require.NotNil(t, uErr, "err does not implement Unwrap interface")

		expected = "failed"
		assert.EqualError(t, uErr, expected, "error message mismatch, got %s want %s", uErr, expected)
	})

	t.Run("Unwrap for errors.Wrapf", func(t *testing.T) {
		t.Parallel()

		err := errors.New("failed")

		errWrap := errors.Wrapf(err, "oops id %d", 5)
		require.Error(t, errWrap, "it is not an error")

		expected := "oops id 5: failed"
		assert.EqualError(t, errWrap, expected, "error message mismatch, got %s want %s", errWrap, expected)

		uErr := errors.Unwrap(errWrap)
		require.NotNil(t, uErr, "err does not implement Unwrap interface")

		expected = "failed"
		assert.EqualError(t, uErr, expected, "error message mismatch, got %s want %s", uErr, expected)
	})

	t.Run("Unwrap for errors.WrapWithError", func(t *testing.T) {
		t.Parallel()

		err := errors.New("failed")
		sErr := errors.New("oops")

		errWrap := errors.WrapWithError(err, sErr)
		require.Error(t, errWrap, "it is not an error")

		uErr := errors.Unwrap(errWrap)
		require.NotNil(t, uErr, "err does not implement Unwrap interface")

		expected := "oops"
		assert.EqualError(t, uErr, expected, "error message mismatch, got %s want %s", uErr, expected)
	})
}

func Test_Cause(t *testing.T) {
	t.Parallel()

	t.Run("Cause for errors.WrapWithError", func(t *testing.T) {
		t.Parallel()

		err := errors.New("failed")
		sErr := errors.New("oops")

		errWrap := errors.WrapWithError(err, sErr)
		require.Error(t, errWrap, "it is not an error")

		cErr := errors.Cause(errWrap)
		require.NotNil(t, cErr, "err does not implement Cause interface")

		expected := "failed"
		assert.EqualError(t, cErr, expected, "error message mismatch, got %s want %s", cErr, expected)
	})
}

func Test_Is(t *testing.T) {
	t.Parallel()

	t.Run("Is for errors.New", func(t *testing.T) {
		t.Parallel()

		err := errors.New("failed")
		require.Error(t, err, "it is not an error")

		expected := errors.New("failed")
		assert.ErrorIs(t, err, expected)
	})

	t.Run("no Is for errors.New", func(t *testing.T) {
		t.Parallel()

		err := errors.New("failed")
		require.Error(t, err, "it is not an error")

		assert.NotErrorIs(t, err, context.Canceled)
	})

	t.Run("Is for errors.Newf", func(t *testing.T) {
		t.Parallel()

		err := errors.Newf("oops: %v", "failed")
		require.Error(t, err, "it is not an error")

		expected := errors.New("oops: failed")
		assert.ErrorIs(t, err, expected)
	})

	t.Run("Is for errors.Wrap", func(t *testing.T) {
		t.Parallel()

		err := errors.New("failed")

		errWrap := errors.Wrap(err, "oops")
		require.Error(t, errWrap, "it is not an error")

		expected := errors.New("failed")
		assert.ErrorIs(t, errWrap, expected)
	})

	t.Run("Is for errors.Wrapf", func(t *testing.T) {
		t.Parallel()

		err := errors.New("failed")

		errWrap := errors.Wrapf(err, "oops id %d", 5)
		require.Error(t, errWrap, "it is not an error")

		expected := errors.New("failed")
		assert.ErrorIs(t, errWrap, expected)
	})

	t.Run("Is for errors.WrapWithError", func(t *testing.T) {
		t.Parallel()

		err := errors.New("failed")
		sErr := errors.New("oops")

		errWrap := errors.WrapWithError(err, sErr)
		require.Error(t, errWrap, "it is not an error")

		assert.EqualError(t, errWrap, "oops: failed")

		assert.ErrorIs(t, errWrap, err)

		assert.ErrorIs(t, errWrap, sErr)
	})

	t.Run("no Is for errors.WrapWithError", func(t *testing.T) {
		t.Parallel()

		err := errors.New("failed")
		sErr := errors.New("oops")

		errWrap := errors.WrapWithError(err, sErr)
		require.Error(t, errWrap, "it is not an error")

		assert.NotErrorIs(t, errWrap, context.Canceled)
	})

	t.Run("Is for errors.WrapWithError two levels", func(t *testing.T) {
		t.Parallel()

		sErr1 := errors.New("failed")
		sErr2 := errors.New("oops")

		errWrap := errors.WrapWithError(context.Canceled, sErr1)
		require.Error(t, errWrap, "it is not an error")

		errWrap = errors.WrapWithError(errWrap, sErr2)
		require.Error(t, errWrap, "it is not an error")

		assert.ErrorIs(t, errWrap, sErr1)

		assert.ErrorIs(t, errWrap, sErr2)

		assert.ErrorIs(t, errWrap, context.Canceled)
	})
}
