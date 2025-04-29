package processor

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/writer"
)

type InputRecord struct {
	Version   string    `json:"version"`
	Name      string    `json:"name"`
	UUID      string    `json:"uuid"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Altitude  float64   `json:"altitude"`
	Course    float64   `json:"course"`
	Speed     float64   `json:"speed"`
	Timestamp time.Time `json:"timestamp"`
}

type ParquetRecord struct {
	Version    string  `parquet:"name=version, type=BYTE_ARRAY, convertedtype=UTF8"`
	Name       string  `parquet:"name=name, type=BYTE_ARRAY, convertedtype=UTF8"`
	UUID       string  `parquet:"name=uuid, type=BYTE_ARRAY, convertedtype=UTF8"`
	Latitude   float64 `parquet:"name=latitude, type=DOUBLE"`
	Longitude  float64 `parquet:"name=longitude, type=DOUBLE"`
	Altitude   float64 `parquet:"name=altitude, type=DOUBLE"`
	Course     float64 `parquet:"name=course, type=DOUBLE"`
	Speed      float64 `parquet:"name=speed, type=DOUBLE"`
	Timestamp  string  `parquet:"name=timestamp, type=BYTE_ARRAY, convertedtype=UTF8"`
	Year       int32   `parquet:"name=year, type=INT32"`
	Month      int32   `parquet:"name=month, type=INT32"`
	Day        int32   `parquet:"name=day, type=INT32"`
	ProcessedAt string  `parquet:"name=processed_at, type=BYTE_ARRAY, convertedtype=UTF8"`
}

func Process(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	
	// Map to store records by partition
	recordsByPartition := make(map[string][]ParquetRecord)
	
	for scanner.Scan() {
		line := scanner.Text()
		
		var record InputRecord
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			return fmt.Errorf("failed to parse JSON: %w", err)
		}
		
		// Extract time components for partitioning
		year := int32(record.Timestamp.Year())
		month := int32(record.Timestamp.Month())
		day := int32(record.Timestamp.Day())
		
		// Create the partition key (name/year/month/day)
		partitionKey := fmt.Sprintf("%s/%d/%02d/%02d", 
			record.Name, year, month, day)
		
		// Create a ParquetRecord
		parquetRecord := ParquetRecord{
			Version:    record.Version,
			Name:       record.Name,
			UUID:       record.UUID,
			Latitude:   record.Latitude,
			Longitude:  record.Longitude,
			Altitude:   record.Altitude,
			Course:     record.Course,
			Speed:      record.Speed,
			Timestamp:  record.Timestamp.Format(time.RFC3339),
			Year:       year,
			Month:      month,
			Day:        day,
			ProcessedAt: time.Now().Format(time.RFC3339),
		}
		
		// Add the record to the appropriate partition
		recordsByPartition[partitionKey] = append(recordsByPartition[partitionKey], parquetRecord)
	}
	
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}
	
	// Use the output directory from command line (main.go already has this flag)
	outputDir := "output"
	if envDir := os.Getenv("OUTPUT_DIR"); envDir != "" {
		outputDir = envDir
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	
	// Write each partition to a Parquet file
	for partition, records := range recordsByPartition {
		if len(records) == 0 {
			continue
		}
		
		// Create the partition directory
		partitionDir := filepath.Join(outputDir, partition)
		if err := os.MkdirAll(partitionDir, 0755); err != nil {
			return fmt.Errorf("failed to create partition directory %s: %w", partitionDir, err)
		}
		
		// Create a unique filename
		filename := filepath.Join(partitionDir, fmt.Sprintf("data_%s.parquet", uuid.New().String()))
		
		// Create a Parquet file
		fw, err := local.NewLocalFileWriter(filename)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", filename, err)
		}
		
		// Create a Parquet writer
		pw, err := writer.NewParquetWriter(fw, new(ParquetRecord), 4)
		if err != nil {
			fw.Close()
			return fmt.Errorf("failed to create Parquet writer: %w", err)
		}
		
		pw.RowGroupSize = 128 * 1024 * 1024 // 128MB
		pw.CompressionType = parquet.CompressionCodec_SNAPPY
		
		// Write the records
		for _, record := range records {
			if err := pw.Write(record); err != nil {
				pw.WriteStop()
				fw.Close()
				return fmt.Errorf("failed to write record: %w", err)
			}
		}
		
		// Close the Parquet writer and file
		if err := pw.WriteStop(); err != nil {
			fw.Close()
			return fmt.Errorf("failed to stop Parquet writer: %w", err)
		}
		if err := fw.Close(); err != nil {
			return fmt.Errorf("failed to close file %s: %w", filename, err)
		}
		
		fmt.Fprintf(w, "Created partition file: %s with %d records\n", filename, len(records))
	}
	
	return nil
}