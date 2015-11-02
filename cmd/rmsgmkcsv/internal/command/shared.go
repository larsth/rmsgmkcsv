package command

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
)

type Degrees struct {
	Min  int64
	Max  int64
	Step int32
}

type CreateMCUFlags struct {
	Start int64
	Step  int64
}

type CreateFlags struct {
	Degrees      Degrees
	FileName     string
	TruncateFile bool
	MCU          CreateMCUFlags
}

type TopLevelCommands struct {
	Create CreateFlags
}

type Data struct {
	File   *os.File
	Output io.Writer
	Writer *csv.Writer
}

const (
	CommandName           = "rmsgmkcsv"
	formatFloatPrecision  = 5
	milliDegressInACircle = 360000
	milliDegreesPerDegree = float64(1000)
)

var (
	ErrMaxLessThanMin = errors.New(
		"Flag 'max' degrees flag is less than flag 'min' degrees flag.")
	ErrStepLessThanZero = errors.New(
		"The 'step' degrees flag is less than zero.")
	ErrNoExec = errors.New(
		"No executable - see the example from running this command: '" +
			CommandName + " help create exec' (w/o the quotes)")
	ErrStepsLargerThanMaxInt32 = errors.New(
		"The number of steps is larger than 2147483647 ((2³¹)-1)")
	ErrMaxMinusMinIsNegative = errors.New(
		"Flag 'max' minus flag 'min' has a negative sign (is < 0).")
	ErrStepDivisionByZero = errors.New(
		"Flag 'step' must not be zero (0).")
	flags *TopLevelCommands
	data  Data
)
