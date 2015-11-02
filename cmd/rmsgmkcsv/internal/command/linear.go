package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"math"
	"math/big"
	"strconv"
)

var (
	linearCmd = &cobra.Command{
		Use:   "linear",
		Short: "linear - create CSV records with linear values.",
		Long: "Create CSV records  with linear values with " +
			"the 'linear' command. " +
			"The 'linear' command starts at the 'min'(imum) decimal " +
			" bearing/azimuth, and works in 'step' decimal degrees up to the " +
			" 'max'(imum) decimal bearing/azimuth degree, with both the " +
			"'min',  and 'max' decimal degrees included, but not exceeded / " +
			"not out of bounds.",
		RunE: linearRunE,
	}
)

func linearRunE(cmd *cobra.Command, args []string) error {
	var (
		records   []*Record
		err       error
		maxInt64B = big.NewInt(math.MaxInt64)
		mcuB      = big.NewInt(flags.Create.MCU.Start)
		mcuStepB  = big.NewInt(flags.Create.MCU.Step)
	)
	if records, err = createDegreeValues(&flags.Create.Degrees); err != nil {
		return err
	}

	for _, record := range records {
		//Overflow error handling for the mcu value (right-hand side):
		//=============================================================
		if mcuB.Cmp(maxInt64B) > 0 {
			return fmt.Errorf("%s", "%s (%s) %s %s: '%s'.",
				"The micro processor value", mcuB.String(),
				"had become too positive. The largest positive microcontroller",
				"value can be: ", strconv.FormatInt(math.MaxInt64, 10))
		}
		//end of mcu value error/overflow handling

		record.McuValue = mcuB.Int64()
		if err := data.Writer.Write(record.Strings()); err != nil {
			return err
		}
		_ = mcuB.Add(mcuB, mcuStepB)
	}
	return nil
}
