package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/LitoleNINJA/json-parser/cmd/encoder"
	"github.com/LitoleNINJA/json-parser/cmd/parser"
)

type benchResult struct {
	custom       time.Duration
	stdlib       time.Duration
	testcaseSize int64
	percentDiff  float64
}

type benchJob struct {
	fileName string
	data     []byte
	fileSize int64
}

// Worker pool to process benchmark jobs
func benchmarkWorker(b *testing.B, jobs <-chan benchJob, results chan<- struct {
	fileName string
	result   benchResult
}, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		var customTime, stdlibTime time.Duration

		// Create separate sub-benchmarks for accurate measurements
		customDone := make(chan time.Duration)
		stdlibDone := make(chan time.Duration)

		// Run custom and stdlib benchmarks in parallel
		go func() {
			start := time.Now()
			for i := 0; i < b.N; i++ {
				var result any
				if err := parser.ParseJSON(job.data, &result, false); err != nil {
					b.Errorf("Custom parser failed for %s: %v", job.fileName, err)
					continue
				}
				if _, err := encoder.EncodeJSON(result, false); err != nil {
					b.Errorf("Custom encoder failed for %s: %v", job.fileName, err)
					continue
				}
			}
			customDone <- time.Since(start)
		}()

		go func() {
			start := time.Now()
			for i := 0; i < b.N; i++ {
				var result any
				if err := json.Unmarshal(job.data, &result); err != nil {
					b.Errorf("Stdlib unmarshal failed for %s: %v", job.fileName, err)
					continue
				}
				if _, err := json.Marshal(result); err != nil {
					b.Errorf("Stdlib marshal failed for %s: %v", job.fileName, err)
					continue
				}
			}
			stdlibDone <- time.Since(start)
		}()

		// Wait for both benchmarks to complete
		customTime = <-customDone
		stdlibTime = <-stdlibDone

		// Calculate results
		result := benchResult{
			custom:       customTime,
			stdlib:       stdlibTime,
			testcaseSize: job.fileSize,
			percentDiff:  (float64(customTime-stdlibTime) / float64(stdlibTime)) * 100,
		}

		results <- struct {
			fileName string
			result   benchResult
		}{job.fileName, result}
	}
}

func BenchmarkComparisonWithSummary(b *testing.B) {
	// Create results file
	resultsFile, err := os.Create("bench_results.txt")
	if err != nil {
		b.Fatalf("Failed to create results file: %v", err)
	}
	defer resultsFile.Close()

	// Write system information
	writeSystemInfo(resultsFile)

	// Note benchmark start time
	writeResultsToFile(resultsFile, "Benchmark Run: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	// Initialize channels for job distribution
	numWorkers := runtime.NumCPU()
	jobs := make(chan benchJob, numWorkers)
	results := make(chan struct {
		fileName string
		result   benchResult
	}, numWorkers)

	// Start worker pool
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go benchmarkWorker(b, jobs, results, &wg)
	}

	// Process test files
	go func() {
		testDir := "../test/JSONTestSuite/test_parsing"
		files, err := os.ReadDir(testDir)
		if err != nil {
			b.Fatalf("Failed to read test directory: %v", err)
			return
		}

		for _, file := range files {
			fileName := file.Name()
			// process only valid JSON file (with "y_" prefix)
			if !strings.HasPrefix(fileName, "y_") {
				continue
			}

			filePath := filepath.Join(testDir, fileName)
			data, err := os.ReadFile(filePath)
			if err != nil {
				b.Errorf("Failed to read file %s: %v", fileName, err)
				continue
			}

			fileInfo, err := os.Stat(filePath)
			if err != nil {
				b.Errorf("Failed to get file info %s: %v", fileName, err)
				continue
			}

			jobs <- benchJob{
				fileName: fileName,
				data:     data,
				fileSize: fileInfo.Size(),
			}
		}
		close(jobs)
	}()

	// Collect and process results
	go func() {
		wg.Wait()
		close(results)
	}()

	// Process results and generate summary
	benchmarkResults := make(map[string]benchResult)
	var totalCustomTime, totalStdlibTime time.Duration
	var totalBytes int64
	fileCount := 0

	for result := range results {
		benchmarkResults[result.fileName] = result.result
		totalCustomTime += result.result.custom
		totalStdlibTime += result.result.stdlib
		totalBytes += result.result.testcaseSize
		fileCount++
	}

	// Write detailed results
	writeDetailedResults(resultsFile, benchmarkResults)

	// Write summary statistics
	writeSummaryStats(resultsFile, benchmarkResults, fileCount, totalBytes, totalCustomTime, totalStdlibTime)
}

func writeSystemInfo(f *os.File) {
	writeResultsToFile(f, "System Information:\n")
	writeResultsToFile(f, "Go Version: %s\n", runtime.Version())
	writeResultsToFile(f, "OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	writeResultsToFile(f, "CPU Cores: %d\n\n\n", runtime.NumCPU())
}

func writeDetailedResults(f *os.File, results map[string]benchResult) {
	writeResultsToFile(f, "\n==================== Benchmark Summary ====================")
	writeResultsToFile(f, "\n-------------------------------------------------------------")

	for fileName, r := range results {
		performance := "slower"
		if r.percentDiff < 0 {
			performance = "faster"
		}

		writeResultsToFile(f, "\nFile: %s (%.2f KB)\n", fileName, float64(r.testcaseSize)/1024)
		writeResultsToFile(f, "Custom Implementation: %v\n", r.custom)
		writeResultsToFile(f, "Standard Library:     %v\n", r.stdlib)
		writeResultsToFile(f, "Performance:          %.2f%% %s than stdlib\n", math.Abs(r.percentDiff), performance)
		writeResultsToFile(f, "-------------------------------------------------------------")
	}
}

func writeSummaryStats(f *os.File, results map[string]benchResult, fileCount int, totalBytes int64, totalCustomTime, totalStdlibTime time.Duration) {
	var totalPercentDiff float64
	for _, r := range results {
		totalPercentDiff += r.percentDiff
	}

	avgPercentDiff := totalPercentDiff / float64(fileCount)
	avgPerformance := "slower"
	if avgPercentDiff < 0 {
		avgPerformance = "faster"
	}

	writeResultsToFile(f, "\n\n========================== Overall Statistics ==========================\n\n")
	writeResultsToFile(f, "Total files processed:     %d\n", fileCount)
	writeResultsToFile(f, "Total data processed:      %.2f KB\n", float64(totalBytes)/1024)
	writeResultsToFile(f, "Average performance:       %.2f%% %s than stdlib\n", math.Abs(avgPercentDiff), avgPerformance)
	writeResultsToFile(f, "Total time (Custom):       %v\n", totalCustomTime)
	writeResultsToFile(f, "Total time (Stdlib):       %v\n", totalStdlibTime)
	writeResultsToFile(f, "\n========================================================================")
	writeResultsToFile(f, "\nBenchmark End: %s\n", time.Now().Format("2006-01-02 15:04:05"))
}

func writeResultsToFile(f *os.File, format string, args ...interface{}) {
	fmt.Fprintf(f, format, args...)
}
