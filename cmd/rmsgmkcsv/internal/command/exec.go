package command

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/larsth/linescanner"
)

func mkExecCmdLongTxt() string {
	var buf bytes.Buffer

	buf.WriteString("\tCreate CSV records by calling a program or a script via a script interpreter.\n")
	buf.WriteString("\tThe program to be called is the 1st argument mentioned after the '--' token.\n")
	buf.WriteString("\tThe '--' token must have a whitespace around itself, \n")
	buf.WriteString("\tand must come after the last argument to ")
	buf.WriteString(CommandName)
	buf.WriteString("\n\t")
	buf.WriteString(CommandName)
	buf.WriteString(" will call that program for translation of many decimal degree values.\n")
	buf.WriteString("\tThat is done by writing ASCII 8-bit floating point values to the program's \n")
	buf.WriteString("\tstandard input (STDIN).\n")
	buf.WriteString("\t\n")
	buf.WriteString("\tThe floating point comma is a '.' (dot) and 2 floating-point values are \n")
	buf.WriteString("\tseperated by a '\\n' token, so the program can read a line - translate\n")
	buf.WriteString("\tthe degree to the microcontroller whole number value, which is 8-bit ASCII\n")
	buf.WriteString("\tencoded text.\n")
	buf.WriteString("\t\n")
	buf.WriteString("\tThe program is expected to write the microcontroller value to standard output\n")
	buf.WriteString("\t(STDOUT), and then repeat the whole sequence of actions again.\n")
	buf.WriteString("\t\n")
	buf.WriteString("\tThe program will be killed by " + CommandName + ", when there is nothing more to do.\n")

	return buf.String()
}

var (
	execCmd = &cobra.Command{
		Use: "exec",
		Short: "exec - Create CSV records by calling an external program or " +
			"script interpreter.",
		Long:    mkExecCmdLongTxt(),
		RunE:    execRunE,
		Example: `rmsgmkcsv create exec -t -o csv.csv -- python3 ./example-script.py`,
	}
)

func startExternalCmd(cmd *exec.Cmd) (stdinPipe io.WriteCloser,
	stdoutPipe io.ReadCloser, err error) {

	if stdinPipe, err = cmd.StdinPipe(); err != nil {
		return nil, nil, err
	}
	if stdoutPipe, err = cmd.StdoutPipe(); err != nil {
		return stdinPipe, nil, err
	}
	//	start the external program
	if err = cmd.Start(); err != nil {
		return stdinPipe, stdoutPipe, err
	}
	return stdinPipe, stdoutPipe, nil
}

func writeToSTDIN(degree int64, stdinPipe io.WriteCloser) error {
	var (
		f        float64 = float64(degree) / milliDegreesPerDegree
		stdinBuf bytes.Buffer
		err      error
	)
	stdinBuf.WriteString(strconv.FormatFloat(f, 'f', -1, 64))
	stdinBuf.WriteString("\n")
	if _, err = stdinBuf.WriteTo(stdinPipe); err != nil {
		return err
	}
	return nil
}

func readFromSTDOUT(ls *linescanner.LineScanner, record *Record) error {
	var (
		s         string
		err       error
	)
	
	for ls.Scan() == false {
		if ls.Err() != nil {
			return err
		}
	}
	if ls.Err() != nil {
		return err
	}
	s = ls.Text()
	if record.McuValue, err = strconv.ParseInt(s, 10, 64); err != nil {
		return fmt.Errorf("ERROR - %s:'%s'- this error ocuured: %s",
			"while parsing the microcontroller value", s, err.Error())
	}
	return nil
}

func recordsIO(records []*Record, stdinPipe io.WriteCloser, stdoutPipe io.ReadCloser) error {
	var err error
	var ls *linescanner.LineScanner
	
	if ls, err = linescanner.New(io.Reader(stdoutPipe)); err != nil {
		return err
	}
	
	for _, record := range records {
		if err = writeToSTDIN(record.Degree, stdinPipe); err != nil {
			return err
		}
		//Read from the program's STDOUT ...
		if err = readFromSTDOUT(ls, record); err != nil {
			return err
		}
		//Write the CSV record ...
		if err = data.Writer.Write(record.Strings()); err != nil {
			return err
		}
	}
	return nil
}

func closePipesAndKillChildProcess(i io.WriteCloser, o io.ReadCloser, p *os.Process) error {
	var err error
	if err = i.Close(); err != nil && err != io.EOF {
		return err
	}
	if err = o.Close(); err != nil && err != io.EOF {
		return err
	}
	if err = p.Kill(); err != nil {
		log.Println("INFO, child process:", err.Error())
	}
	return nil
}

func execRunE(cobraCmd *cobra.Command, _ []string) error {
	var (
		records    []*Record
		err        error
		execArgs   = cobraCmd.Flags().Args()
		cmd        *exec.Cmd
		stdinPipe  io.WriteCloser //connected to the external program's STDIN
		stdoutPipe io.ReadCloser  //connected to the external program's STDOUT
	)
	if len(execArgs) == 0 {
		return ErrNoExec
	}
	//len(execArgs) > 0, min. 1 argument exists ...
	if records, err = createDegreeValues(&flags.Create.Degrees); err != nil {
		return err
	}
	//	create an *exec.cmd ...
	cmd = exec.Command(execArgs[0], execArgs[1:]...)
	// ..., and start it ...
	if stdinPipe, stdoutPipe, err = startExternalCmd(cmd); err != nil {
		return err
	}
	if err = recordsIO(records, stdinPipe, stdoutPipe); err != nil {
		return err
	}
	err = closePipesAndKillChildProcess(stdinPipe, stdoutPipe, cmd.Process)
	if err != nil {
		return err
	}
	return nil
}
