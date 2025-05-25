package file

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/tartoide/stori/stori-challenge/internal/domain"
	"github.com/tartoide/stori/stori-challenge/internal/services"
)

type csvFileProcessor struct {
	logger *slog.Logger
}

// NewCSVFileProcessor creates a new CSV file processor
func NewCSVFileProcessor(logger *slog.Logger) services.FileProcessor {
	return &csvFileProcessor{
		logger: logger,
	}
}

// ProcessFile processes a CSV file and returns transactions
func (p *csvFileProcessor) ProcessFile(ctx context.Context, filePath string) ([]domain.Transaction, error) {
	p.logger.Debug("opening file for processing", "file", filePath)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("opening file %s: %w", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 3 // ID, Date, Transaction

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("reading header from %s: %w", filePath, err)
	}

	p.logger.Debug("CSV header read", "header", header)

	// Validate header
	expectedHeaders := []string{"Id", "Date", "Transaction"}
	if !validateHeader(header, expectedHeaders) {
		return nil, fmt.Errorf("invalid CSV header in %s, expected %v, got %v", filePath, expectedHeaders, header)
	}

	var transactions []domain.Transaction
	lineNumber := 2 // content starts in line 2

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			p.logger.Warn("skipping invalid line", "line", lineNumber, "error", err, "file", filePath)
			lineNumber++
			continue
		}

		if len(record) != 3 {
			p.logger.Warn("skipping line with incorrect field count", "line", lineNumber, "fields", len(record), "file", filePath)
			lineNumber++
			continue
		}

		transaction, err := domain.NewTransaction(record[0], record[1], record[2])
		if err != nil {
			p.logger.Warn("skipping invalid transaction", "line", lineNumber, "error", err, "record", record, "file", filePath)
			lineNumber++
			continue
		}

		transactions = append(transactions, *transaction)
		lineNumber++
	}

	p.logger.Info("file processing completed", "file", filePath, "transactions", len(transactions), "lines_processed", lineNumber-1)
	return transactions, nil
}

// WatchDirectory watches a directory for new CSV files
func (p *csvFileProcessor) WatchDirectory(ctx context.Context, dirPath string) (<-chan string, error) {
	p.logger.Info("setting up directory watcher", "directory", dirPath)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("creating file watcher: %w", err)
	}

	err = watcher.Add(dirPath)
	if err != nil {
		watcher.Close()
		return nil, fmt.Errorf("adding directory to watcher %s: %w", dirPath, err)
	}

	fileChan := make(chan string, 10)

	go func() {
		defer watcher.Close()
		defer close(fileChan)

		for {
			select {
			case <-ctx.Done():
				p.logger.Info("stopping directory watcher due to context cancellation")
				return
			case event, ok := <-watcher.Events:
				if !ok {
					p.logger.Warn("watcher events channel closed")
					return
				}

				p.logger.Debug("file system event", "event", event.String())

				if (event.Op&fsnotify.Create == fsnotify.Create ||
					event.Op&fsnotify.Write == fsnotify.Write ||
					event.Op&fsnotify.Chmod == fsnotify.Chmod) &&
					strings.HasSuffix(strings.ToLower(event.Name), ".csv") {

					// Wait a bit to ensure file is fully written
					time.Sleep(100 * time.Millisecond)

					// Check if readable
					if _, err := os.Stat(event.Name); err == nil {
						p.logger.Info("new CSV file detected", "file", event.Name)
						select {
						case fileChan <- event.Name:
						case <-ctx.Done():
							return
						}
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					p.logger.Warn("watcher errors channel closed")
					return
				}
				p.logger.Error("file watcher error", "error", err)
			}
		}
	}()

	p.logger.Info("directory watcher started successfully", "directory", dirPath)
	return fileChan, nil
}

func validateHeader(header, expected []string) bool {
	if len(header) != len(expected) {
		return false
	}
	for i, h := range header {
		if strings.TrimSpace(h) != expected[i] {
			return false
		}
	}
	return true
}
