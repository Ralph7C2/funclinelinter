package testdata

import (
	"context"
)

type thing struct {
	longlyNamedFunctionField      func(ctx context.Context, andOfCourseALonglyNamedParam int, xyzzyx int, yzzxxyasdfas int) error
	wronglyFormattedFunctionField func(ctx context.Context,
		andOfCourseALonglyNamedParam int, xyzzyx int, yzzxxyasfs int) error
	goodFunctionField func(
		ctx context.Context, stuff int, morestuff string,
	) error
}

func (thing) longlyNamedMethod(ctx context.Context, andOfCourseALonglyNamedParam int, xyzzyx int, yzzxxyasdfas int) error {
	return nil
}

func (thing) wronglyFormattedMethod(
	ctx context.Context, andOfCourseALonglyNamedParam int, xyzzyx int, yzzxxyasdfas int) error {
	return nil
}

func (thing) goodMethod(
	ctx context.Context, andOfCourseALonglyNamedParam int, xyzzyx int, yzzxxyasdfas int,
) error {
	return nil
}

type longlyNamedFunctionType func(ctx context.Context, alsasjdfhasjdkflaskdfjaskdfhaskdjfahsi int, asdkjfahsdkjfhaskdjfhas string) error
type wronglyFormattedFunctionType func(
	context.Context) error
type goodFunctionType func(
	ctx context.Context, stuff int, morestuff string,
) error

var longlyNamedFunctionLiteral = func(ctx context.Context, alsasjdfhasjdkflaskdfjaskdfhaskdjfahsi int, asdkjfahsdkjfhaskdjfhas string) error {
	return nil
}
var wronglyFormattedFunctionLiteral = func(ctx context.Context,
	alsasjdfhasjdkflaskdfjaskdfhaskdjfahsi int, asdkjfahsdkjfhaskdjfhas string) error {
	return nil
}
var goodFunctionLiteral = func(ctx context.Context,
	alsasjdfhasjdkflaskdfjaskdfhaskdjfahsi int, asdkjfahsdkjfhaskdjfhas string,
) error {
	return nil
}

func longlyNamedFunction(someParamWithAReallyLongName thing, more bool, and string, andyet int) (int, bool, string, thing) {
	return 0, false, "", thing{}
}
func wronglyFormattedFunction(
	someParamWithAReallyLongName thing, more bool, and string, andyet int) (int, bool, string, thing) {
	return 0, false, "", thing{}
}

func goodFunction(
	someParamWithAReallyLongName thing, more bool, and string, andyet int,
) (int, bool, string, thing) {
	return 0, false, "", thing{}
}
