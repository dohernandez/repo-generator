package errors_test

import (
	"context"
	"fmt"

	"github.com/consensys-vertical-apps/platform-data-pipeline-toolkit/errors"
)

func ExampleNew() {
	err := errors.New("foo")

	fmt.Println(err)

	// Output:
	// foo
}

func ExampleNewf() {
	err := errors.Newf("foo: %s", "bar")

	fmt.Println(err)

	// Output:
	// foo: bar
}

func ExampleWrap() {
	err := errors.New("foo")

	err = errors.Wrap(err, "bar")

	fmt.Println(err)

	// Output:
	// bar: foo
}

func ExampleWrapf() {
	err := errors.New("foo")

	err = errors.Wrapf(err, "bar: %s", "baz")

	fmt.Println(err)

	// Output:
	// bar: baz: foo
}

func ExampleWrapWithError() {
	err := errors.New("foo")
	err2 := errors.New("bar")
	err3 := errors.New("baz")

	err = errors.WrapWithError(err, err2)
	err = errors.WrapWithError(err, err3)

	fmt.Println(err)

	// Output:
	// baz: bar: foo
}

func Example_withMessage_UnWrap() {
	err := errors.New("foo")

	err = errors.Wrap(err, "bar")

	for {
		fmt.Println(err)

		err = errors.Unwrap(err)
		if err == nil {
			return
		}
	}

	//nolint:dupword
	// Output:
	// bar: foo
	// foo
}

func Example_withError_UnWrap_New() {
	err := errors.New("foo")
	err2 := errors.New("bar")

	err = errors.WrapWithError(err, err2)

	for {
		fmt.Println(err)

		err = errors.Unwrap(err)
		if err == nil {
			return
		}
	}

	// Output:
	// bar: foo
	// bar
}

func Example_withError_UnWrap_Wrap() {
	err := errors.New("foo")
	err = errors.Wrapf(err, "bar")

	err2 := errors.New("baz")
	err2 = errors.Wrap(err2, "qux")

	err = errors.WrapWithError(err, err2)

	for {
		fmt.Println(err)

		err = errors.Unwrap(err)
		if err == nil {
			return
		}
	}

	//nolint:dupword
	// Output:
	// qux: baz: bar: foo
	// qux: baz
	// baz
}

func Example_withError_Cause_New() {
	err := errors.New("foo")
	err2 := errors.New("bar")

	err = errors.WrapWithError(err, err2)

	for {
		fmt.Println(err)

		err = errors.Cause(err)
		if err == nil {
			return
		}
	}

	//nolint:dupword
	// Output:
	// bar: foo
	// foo
}

func Example_withError_Cause_Wrap() {
	err := errors.New("foo")
	err = errors.Wrap(err, "bar")

	err2 := errors.New("baz")
	err2 = errors.Wrap(err2, "qux")

	err = errors.WrapWithError(err, err2)

	for {
		fmt.Println(err)

		err = errors.Cause(err)
		if err == nil {
			return
		}
	}

	// Output:
	// qux: baz: bar: foo
	// bar: foo
}

func Example_withError_Cause_Breadcrumb() {
	err := errors.New("foo")
	err = errors.Wrap(err, "bar")

	err2 := errors.New("baz")
	err2 = errors.Wrap(err2, "qux")

	err = errors.WrapWithError(err, err2)

	fmt.Println(err)

	wErr := errors.Unwrap(err)

	for {
		fmt.Println(wErr)

		wErr = errors.Unwrap(wErr)
		if wErr == nil {
			break
		}
	}

	cErr := errors.Cause(err)

	for {
		fmt.Println(cErr)

		cErr = errors.Unwrap(cErr)
		if cErr == nil {
			break
		}
	}

	//nolint:dupword
	// Output:
	// qux: baz: bar: foo
	// qux: baz
	// baz
	// bar: foo
	// foo
}

func Example_withError_Is() {
	err1 := errors.New("foo")
	err1 = errors.Wrap(err1, "bar")

	err2 := errors.New("baz")
	err2 = errors.Wrap(err2, "qux")

	err := errors.WrapWithError(err1, err2)

	fmt.Println(err)

	fmt.Println(errors.Is(err, err1))
	fmt.Println(errors.Is(err, err2))
	fmt.Println(errors.Is(err, context.DeadlineExceeded))

	//nolint:dupword
	// Output:
	// qux: baz: bar: foo
	// true
	// true
	// false
}
