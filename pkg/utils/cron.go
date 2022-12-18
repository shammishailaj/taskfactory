package utils

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/robfig/cron"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

type TaskFactory struct {
	Tasks *cron.Cron
	U     *Utils
}

func (c *TaskFactory) executeCommand(cmd string, args []string) ([]byte, error) {
	var out bytes.Buffer
	// Execute the command
	command := exec.Command(cmd, args...)
	c.U.Log.Infof("cmd: %s, args: %s", cmd, strings.Join(args, " "))
	command.Stdout = &out
	err := command.Run()
	if err != nil {
		c.U.Log.Errorf("Error running command: %s", err.Error())
		return out.Bytes(), err
	}
	// TODO need to figure-out a way to log this output
	fmt.Println(out.String())
	return out.Bytes(), nil
}

func (c *TaskFactory) executeCommand2(name string, args []string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	//cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		c.U.Log.Error(err)
		return out.Bytes(), err
	}
	fmt.Printf("in all caps: %q\n", out.String())
	return out.Bytes(), nil
}

func (c *TaskFactory) loadAvg() (float64, error) {
	sysinfo := syscall.Sysinfo_t{}
	err := syscall.Sysinfo(&sysinfo)
	if err != nil {
		return 0, err
	}
	c.U.Log.Println("sysinfo = ", sysinfo)
	return float64(sysinfo.Loads[0]), nil // https://www.includehelp.com/golang/get-the-system-information-using-syscall.aspx
}

func (c *TaskFactory) ScheduleCommand(schedule string, cmd string, args []string) error {
	// Parse the crontab schedule and add the task to the cron job
	addFuncErr := c.Tasks.AddFunc(schedule, func() {
		// Collect telemetry data
		var rusage syscall.Rusage
		rusageErr := syscall.Getrusage(syscall.RUSAGE_SELF, &rusage)
		if rusageErr != nil {
			c.U.Log.Errorf("Error getting telemetry data: CPU Utilization: %s", rusageErr.Error())
		}
		memStats := &runtime.MemStats{}
		runtime.ReadMemStats(memStats)
		lavg, lavgErr := c.loadAvg()
		if lavgErr != nil {
			c.U.Log.Errorf("Error getting telemetry data: LoadAverage1: %s", lavgErr.Error())
		}
		// Log the telemetry data
		c.U.Log.Printf("CPU usage: %d\n", rusage.Utime.Nano())
		c.U.Log.Printf("Memory usage: %d\n", memStats.Alloc)
		c.U.Log.Printf("Load average: %.2f\n", lavg)

		execCommandOutput, execCommandErr := c.executeCommand(cmd, args)
		if execCommandErr != nil {
			c.U.Log.Errorf("Error running command %s %s. %s", cmd, strings.Join(args, " "), execCommandErr.Error())
		}

		s3LogOutput, s3LogOutputErr := c.storeLogsOnS3(c.U.AwsRegion, execCommandOutput)
		if s3LogOutputErr != nil {
			c.U.Log.Errorf("Error Uploading Logs to S3. %s", s3LogOutputErr.Error())
		} else {
			c.U.Log.Infof("Successfully uploaded logs to S3. %v", s3LogOutput)
		}
	})
	if addFuncErr != nil {
		c.U.Log.Errorf("Error scheduling task: %s", addFuncErr.Error())
		return addFuncErr
	}

	return nil
}

func (c *TaskFactory) GetCronJobs() []*cron.Entry {
	// Retrieve the list of cron jobs
	return c.Tasks.Entries()
}

func (c *TaskFactory) storeLogsOnS3(awsRegion string, logs []byte) (*s3.PutObjectOutput, error) {
	// Create an AWS session
	awsConfig, awsConfigErr := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if awsConfigErr != nil {
		return nil, awsConfigErr
	}
	// Create an S3 client
	s3Client := s3.NewFromConfig(awsConfig)
	// Set the parameters for the PutObject operation
	params := &s3.PutObjectInput{
		Bucket: aws.String("BUCKET_NAME"),
		Key:    aws.String("FILE_NAME"),
		Body:   bytes.NewReader(logs),
	}
	// Store the logs on S3
	putObjectOutput, putObjectErr := s3Client.PutObject(context.TODO(), params)
	if putObjectErr != nil {
		return nil, putObjectErr
	}
	return putObjectOutput, nil
}
