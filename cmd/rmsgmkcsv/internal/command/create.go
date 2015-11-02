package command

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"math/big"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

type Record struct {
	Degree   int64
	McuValue int64
}

func (r *Record) Strings() []string {
	var strs = make([]string, 2, 2)
	strs[0] = strconv.FormatInt(r.Degree, 10)
	strs[1] = strconv.FormatInt(r.McuValue, 10)
	return strs
}

var (
	createCmd = &cobra.Command{
		Use:   "create",
		Short: "create CSV records.",
		Long: "Either create CSV records with linear values with " +
			"the 'linear' sub-command or CSV records with values from " +
			"a call to a program or script interpreter with the 'exec' " +
			" sub-command.",
		PersistentPreRunE:  persistentPreRunECreate,
		PersistentPostRunE: persistentPostRunECreate,
	}
)

func createDegreeValues(degrees *Degrees) (records []*Record, err error) {
	var (
		minB         = big.NewInt(degrees.Min)
		maxB         = big.NewInt(degrees.Max)
		stepB        = big.NewInt(int64(degrees.Step))
		maxInt32B    = big.NewInt(math.MaxInt32)
		steps        int
		stepsB       *big.Int
		maxMinusMinB *big.Int
		record       *Record
		degree       int64
	)
	maxMinusMinB = big.NewInt(int64(0))
	_ = maxMinusMinB.Sub(maxB, minB) //maxMinusMinB is modified, has the result
	if maxMinusMinB.Sign() == -1 {
		records = nil
		err = ErrMaxMinusMinIsNegative
		return
	}
	if maxMinusMinB.Sign() == 0 {
		// max - min == 0 is true:
		// There is only 1 degree (which is equal to both min and max):
		records = make([]*Record, 1, 1)
		records[0] = new(Record)
		records[0].Degree = degrees.Min
		return
	}
	if stepB.Sign() == 0 {
		//step == 0 is true:
		records = nil
		err = ErrStepDivisionByZero
		return
	}
	stepsB = big.NewInt(int64(0))
	_ = stepsB.Div(maxMinusMinB, stepB) //stepsB is modified, has the result
	if stepsB.Cmp(maxInt32B) > 0 {
		records = nil
		err = ErrStepsLargerThanMaxInt32
		return
	}
	steps = int(stepsB.Int64())
	records = make([]*Record, 0, steps)
	for degree = degrees.Min; degree <= degrees.Max; degree += int64(degrees.Step) {
		record = new(Record)
		record.Degree = degree
		records = append(records, record)
	}
	if len(records) > 0 {
		if records[(len(records)-1)].Degree < degrees.Max {
			record = new(Record)
			record.Degree = degrees.Max
			records = append(records, record)
		}
	}
	return records, nil
}

func persistentPreRunECreate(_ *cobra.Command, _ []string) error {
	var err error

	//min <= max check:
	if flags.Create.Degrees.Min > flags.Create.Degrees.Max {
		return ErrMaxLessThanMin
	}
	//step > 0 (zero) check:
	if flags.Create.Degrees.Step < 0 {
		return ErrStepLessThanZero
	}

	//Use a file(name), or standard output (STDOUT)?
	if len(flags.Create.FileName) > 0 {
		//create and truncate, or open and append?
		if flags.Create.TruncateFile == true {
			//create and truncate the file
			data.File, err = os.Create(flags.Create.FileName)
			if err != nil {
				outErr := fmt.Errorf("ERROR: Cannot create file '%s', \"%s\"\n",
					flags.Create.FileName, err.Error())
				return outErr
			}
		} else {
			//open and append to the file
			data.File, err = os.OpenFile(flags.Create.FileName,
				os.O_WRONLY|os.O_APPEND,
				os.FileMode(0666)) //filemode is before umask happens
			if err != nil {
				outErr := fmt.Errorf("%s: '%s' %s, ERROR: \"%s\"\n",
					"Could not open file", flags.Create.FileName,
					"for write-only and append access", err.Error())
				return outErr
			}
		}
		//Use the file ...
		data.Output = io.Writer(data.File)
	} else {
		//Use STDOUT ...
		data.Output = os.Stdout
	}

	//Create the CSV writer with some sane defaults ...
	data.Writer = csv.NewWriter(data.Output)
	data.Writer.Comma = ','
	data.Writer.UseCRLF = true

	return nil
}

func persistentPostRunECreate(_ *cobra.Command, _ []string) error {
	var err error

	if data.Writer != nil {
		//Flush buffer content
		data.Writer.Flush()
		if err = data.Writer.Error(); err != nil {
			return err
		}
		data.Writer = nil
	}

	data.Output = nil
	if data.File != nil {
		if err = data.File.Close(); err != nil {
			return err
		}
		data.File = nil
	}
	return nil
}

func initCreate() {
	var (
		persistentFlagSet = createCmd.PersistentFlags()
	)
	createCmd.AddCommand(linearCmd)
	createCmd.AddCommand(execCmd)

	//add persistent flags to the 'create' command, which are also visible
	//to the 'linear' and 'exec' sub-commands ...
	persistentFlagSet.Int64VarP(&flags.Create.Degrees.Min, "min", "i", 0,
		"Min. milli degrees that should be created, fx. 0=0.0 degrees.")
	persistentFlagSet.Int64VarP(&flags.Create.Degrees.Max, "max", "a",
		milliDegressInACircle,
		"Max milli degrees that should be created, fx. 180000=180.0 degrees.")
	persistentFlagSet.Int32VarP(&flags.Create.Degrees.Step, "step", "s", 100,
		"A 'step' is the amount of milli degrees that are added between 2 "+
			"following CSV records right hand side (degrees). Fx. if 'step' is"+
			" 100, and the 3rd record is 1.0 degrees, then the 4th record"+
			" will be 1.0+(100/1000)=1.1 degrees.")
	persistentFlagSet.StringVarP(&flags.Create.FileName, "output", "o", "",
		"Output the CSV content to this file. "+
			"The 't' or 'truncate' flag must be used, if the file does not exists. "+
			"If this flag is not given, then all records will be written to "+
			"standard output (STDOUT).")
	persistentFlagSet.BoolVarP(&flags.Create.TruncateFile, "truncate", "t", false,
		"Create and truncate the file if true, otherwise append to the "+
			"existing file (default)."+
			"The 't' or 'truncate' flag must be used, if the file does not exists. "+
			"This option has no effect then standard output (STDOUT) is used, "+
			"that is, when the 'o' or 'output' flag is not given.")
	persistentFlagSet.Int64VarP(&flags.Create.MCU.Step, "mcustep", "u", 2,
		`mcustep - Microprocessor step value: How much is the microprocessor value 
		changed between 2 follwing degrees? Negative, zero, and positive values 
		are accepted, but very large values will result in an overflow error, if
		a microprocessor value become less than -9,223,372,036,854,775,808 or 
		greater than (+)9,223,372,036,854,775,807`)
	persistentFlagSet.Int64VarP(&flags.Create.MCU.Start, "mcu", "m", 0,
		`mcu - is the start value - used as is - in the first CSV record. 
		The 'mcu' and 'mcustep' values are added together, then the 
		microprocessor value for the next CSV record is created.`)
}
