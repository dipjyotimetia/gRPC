package main

import (
	"fmt"
	"github.com/bojand/ghz/printer"
	"github.com/bojand/ghz/runner"
	blogpb "github.com/dipjyotimetia/gogrpc/blog/blogPb"
	"github.com/google/uuid"
	"log"
	"os"
	"strconv"
	"time"
)

type Latency struct {
	Percentile int
	TotalTime  time.Duration
}

func main() {
	report, err := runner.Run(
		"blog.BlogService.CreateBlog",
		"localhost:50051",
		// runner.WithProtoFile("blog/blogPb/blog.proto", []string{}), //TODO: gRPC Reflection api is configured
		runner.WithData(&blogpb.CreateBlogRequest{
			Blog: &blogpb.Blog{
				Id:       strconv.Itoa(int(uuid.New().ID())),
				AuthorId: "f6a4c808-c33d-4c15-9722-fc62a3533e70",
				Title:    "Hello",
				Content:  "Hello",
			}}),
		runner.WithConcurrency(50),
		runner.WithTotalRequests(300),
		runner.WithInsecure(true),
	)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	printer := printer.ReportPrinter{
		Out:    os.Stdout,
		Report: report,
	}

	ts := make(map[int]int64)
	ts[50] = 10
	ts[75] = 12
	ts[90] = 15
	ts[95] = 16
	ts[99] = 25

	ValidateLatency(ts, printer.Report)
	printer.Print("influx-details")
}

// ValidateLatency validate latency
func ValidateLatency(latencyTime map[int]int64, report *runner.Report) {
	for _, details := range report.LatencyDistribution {
		switch details.Percentage {
		case 50:
			if details.Latency.Milliseconds() > latencyTime[50] {
				log.Fatalf("P90: latency differene :%d", latencyTime[50]-details.Latency.Milliseconds())
			}
		case 75:
			if details.Latency.Milliseconds() > latencyTime[75] {
				log.Fatalf("P75: latency differene :%d", latencyTime[75]-details.Latency.Milliseconds())
			}
		case 90:
			if details.Latency.Milliseconds() > latencyTime[90] {
				log.Fatalf("P90: latency differene :%d", latencyTime[90]-details.Latency.Milliseconds())
			}
		case 95:
			if details.Latency.Milliseconds() > latencyTime[95] {
				log.Fatalf("P95: latency differene :%d", latencyTime[95]-details.Latency.Milliseconds())
			}
		case 99:
			if details.Latency.Milliseconds() > latencyTime[99] {
				log.Fatalf("P99: latency differene :%d", latencyTime[99]-details.Latency.Milliseconds())
			}
		}
	}
}
