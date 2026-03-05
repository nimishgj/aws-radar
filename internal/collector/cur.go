package collector

import (
	"compress/gzip"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/nimishgj/aws-radar/internal/config"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/parquet-go/parquet-go"
	"github.com/rs/zerolog/log"
)

// curManifest represents the CUR manifest JSON structure.
type curManifest struct {
	ReportKeys []string `json:"reportKeys"`
}

// curAggregation holds aggregated cost data from CUR files.
type curAggregation struct {
	total       float64
	byService   map[string]float64
	byResource  map[string]resourceCost // key: "service|resource_id"
	byUsageType map[string]float64      // key: "service|usage_type"
	byTag       map[string]float64      // key: "tag_key|tag_value"
	period      string
}

type resourceCost struct {
	service    string
	resourceID string
	cost       float64
}

// CURCollector reads CUR report files from S3 and emits cost metrics.
type CURCollector struct {
	cfg config.CostCURConfig

	mu        sync.RWMutex
	lastFetch time.Time
	cached    *curAggregation
}

// NewCURCollector creates a new CUR collector.
func NewCURCollector(cfg config.CostCURConfig) *CURCollector {
	freq := strings.ToLower(strings.TrimSpace(cfg.Frequency))
	switch freq {
	case "hourly", "daily":
	case "":
		freq = "daily"
	default:
		freq = "daily"
	}
	cfg.Frequency = freq

	if cfg.Region == "" {
		cfg.Region = "us-east-1"
	}
	if cfg.MaxResources <= 0 {
		cfg.MaxResources = 100
	}

	return &CURCollector{cfg: cfg}
}

func (c *CURCollector) Name() string {
	return "cost_cur"
}

func (c *CURCollector) Collect(ctx context.Context, cfg aws.Config, account, accountName string) error {
	now := time.Now().UTC()
	if c.shouldRefresh(now) {
		if err := c.refresh(ctx, cfg, now); err != nil {
			return err
		}
	}

	c.mu.RLock()
	agg := c.cached
	c.mu.RUnlock()

	if agg == nil {
		return nil
	}

	metrics.CURTotalCost.WithLabelValues(account, accountName, agg.period).Set(agg.total)

	for svc, cost := range agg.byService {
		metrics.CURCostByService.WithLabelValues(account, accountName, svc, agg.period).Set(cost)
	}

	// Emit top N resources by cost.
	topResources := topNResources(agg.byResource, c.cfg.MaxResources)
	for _, rc := range topResources {
		metrics.CURCostByResource.WithLabelValues(account, accountName, rc.service, rc.resourceID, agg.period).Set(rc.cost)
	}

	for key, cost := range agg.byUsageType {
		parts := strings.SplitN(key, "|", 2)
		if len(parts) == 2 {
			metrics.CURCostByUsageType.WithLabelValues(account, accountName, parts[0], parts[1], agg.period).Set(cost)
		}
	}

	for key, cost := range agg.byTag {
		parts := strings.SplitN(key, "|", 2)
		if len(parts) == 2 {
			metrics.CURCostByTag.WithLabelValues(account, accountName, parts[0], parts[1], agg.period).Set(cost)
		}
	}

	metrics.CURLastProcessed.WithLabelValues(account, accountName).Set(float64(c.lastFetch.Unix()))

	return nil
}

func (c *CURCollector) shouldRefresh(now time.Time) bool {
	c.mu.RLock()
	lastFetch := c.lastFetch
	c.mu.RUnlock()

	if lastFetch.IsZero() {
		return true
	}

	switch c.cfg.Frequency {
	case "hourly":
		return now.Format("2006-01-02T15") != lastFetch.Format("2006-01-02T15")
	default:
		return now.Format("2006-01-02") != lastFetch.Format("2006-01-02")
	}
}

func (c *CURCollector) refresh(ctx context.Context, cfg aws.Config, now time.Time) error {
	s3Cfg := cfg.Copy()
	s3Cfg.Region = c.cfg.Region

	client := s3.NewFromConfig(s3Cfg)

	// Compute billing period path.
	periodStart, periodEnd := billingPeriod(now)
	period := periodStart + "-" + periodEnd

	manifestKey := fmt.Sprintf("%s/%s/%s/%s-Manifest.json",
		c.cfg.Prefix, c.cfg.ReportName, period, c.cfg.ReportName)

	manifest, err := c.downloadManifest(ctx, client, manifestKey)
	if err != nil {
		return fmt.Errorf("download CUR manifest: %w", err)
	}

	if len(manifest.ReportKeys) == 0 {
		log.Warn().Str("manifest", manifestKey).Msg("CUR manifest has no report keys")
		return nil
	}

	agg := &curAggregation{
		byService:   make(map[string]float64),
		byResource:  make(map[string]resourceCost),
		byUsageType: make(map[string]float64),
		byTag:       make(map[string]float64),
		period:      periodStart,
	}

	format := c.detectFormat(manifest.ReportKeys[0])

	for _, key := range manifest.ReportKeys {
		if err := c.processFile(ctx, client, key, format, agg); err != nil {
			log.Error().Err(err).Str("key", key).Msg("Failed to process CUR file")
			continue
		}
	}

	c.mu.Lock()
	c.cached = agg
	c.lastFetch = now
	c.mu.Unlock()

	return nil
}

func (c *CURCollector) downloadManifest(ctx context.Context, client *s3.Client, key string) (*curManifest, error) {
	out, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.cfg.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	defer out.Body.Close()

	var manifest curManifest
	if err := json.NewDecoder(out.Body).Decode(&manifest); err != nil {
		return nil, fmt.Errorf("decode manifest JSON: %w", err)
	}
	return &manifest, nil
}

func (c *CURCollector) detectFormat(reportKey string) string {
	if c.cfg.Format != "" {
		return strings.ToLower(c.cfg.Format)
	}
	if strings.HasSuffix(reportKey, ".parquet") || strings.HasSuffix(reportKey, ".snappy.parquet") {
		return "parquet"
	}
	return "csv"
}

func (c *CURCollector) processFile(ctx context.Context, client *s3.Client, key, format string, agg *curAggregation) error {
	switch format {
	case "parquet":
		return c.processParquet(ctx, client, key, agg)
	default:
		return c.processCSV(ctx, client, key, agg)
	}
}

func (c *CURCollector) processParquet(ctx context.Context, client *s3.Client, key string, agg *curAggregation) error {
	out, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.cfg.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("get S3 object %s: %w", key, err)
	}
	defer out.Body.Close()

	// Parquet reader needs a file; download to temp.
	tmp, err := os.CreateTemp("", "cur-*.parquet")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	defer os.Remove(tmp.Name())
	defer tmp.Close()

	if _, err := io.Copy(tmp, out.Body); err != nil {
		return fmt.Errorf("download parquet to temp: %w", err)
	}

	if _, err := tmp.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("seek temp file: %w", err)
	}

	stat, err := tmp.Stat()
	if err != nil {
		return fmt.Errorf("stat temp file: %w", err)
	}

	pf, err := parquet.OpenFile(tmp, stat.Size())
	if err != nil {
		return fmt.Errorf("open parquet file: %w", err)
	}

	schema := pf.Schema()
	colIndices := parquetColumnIndices(schema)

	for _, rg := range pf.RowGroups() {
		rows := rg.Rows()
		rowBuf := make([]parquet.Row, 128)
		for {
			n, err := rows.ReadRows(rowBuf)
			for i := 0; i < n; i++ {
				aggregateParquetRow(rowBuf[i], colIndices, agg)
			}
			if err != nil {
				if err == io.EOF {
					break
				}
				rows.Close()
				return fmt.Errorf("read parquet rows: %w", err)
			}
		}
		rows.Close()
	}

	return nil
}

// costLineItemTypes are the line_item_line_item_type values that represent
// actual charges. Credits, refunds, and savings-plan negations are excluded
// so that the aggregated totals reflect real spend.
var costLineItemTypes = map[string]struct{}{
	"Usage":                   {},
	"DiscountedUsage":         {},
	"SavingsPlanCoveredUsage": {},
	"Tax":                     {},
	"Fee":                     {},
	"RIFee":                   {},
	"SavingsPlanRecurringFee": {},
	"SavingsPlanUpfrontFee":   {},
}

type parquetColIndices struct {
	cost         int
	service      int
	resource     int
	usageType    int
	accountID    int
	lineItemType int
	tagCols      map[int]string // column index -> tag name
}

func parquetColumnIndices(schema *parquet.Schema) parquetColIndices {
	idx := parquetColIndices{
		cost:         -1,
		service:      -1,
		resource:     -1,
		usageType:    -1,
		accountID:    -1,
		lineItemType: -1,
		tagCols:      make(map[int]string),
	}

	for i, col := range schema.Columns() {
		path := strings.Join(col, ".")
		switch path {
		case "line_item_unblended_cost":
			idx.cost = i
		case "line_item_product_code":
			idx.service = i
		case "line_item_resource_id":
			idx.resource = i
		case "line_item_usage_type":
			idx.usageType = i
		case "line_item_usage_account_id":
			idx.accountID = i
		case "line_item_line_item_type":
			idx.lineItemType = i
		default:
			if strings.HasPrefix(path, "resource_tags_user_") {
				tagName := strings.TrimPrefix(path, "resource_tags_user_")
				idx.tagCols[i] = tagName
			}
		}
	}

	return idx
}

func aggregateParquetRow(row parquet.Row, idx parquetColIndices, agg *curAggregation) {
	if idx.cost < 0 {
		return
	}

	// Skip non-cost line item types (credits, refunds, negations).
	if idx.lineItemType >= 0 {
		lit := parquetString(row, idx.lineItemType)
		if _, ok := costLineItemTypes[lit]; !ok {
			return
		}
	}

	cost := parquetFloat64(row, idx.cost)
	if cost == 0 {
		return
	}

	service := parquetString(row, idx.service)
	resourceID := parquetString(row, idx.resource)
	usageType := parquetString(row, idx.usageType)

	agg.total += cost

	if service != "" {
		agg.byService[service] += cost
	}

	if resourceID != "" && service != "" {
		key := service + "|" + resourceID
		rc := agg.byResource[key]
		rc.service = service
		rc.resourceID = resourceID
		rc.cost += cost
		agg.byResource[key] = rc
	}

	if usageType != "" && service != "" {
		agg.byUsageType[service+"|"+usageType] += cost
	}

	// Aggregate tag costs.
	for colIdx, tagName := range idx.tagCols {
		tagValue := parquetString(row, colIdx)
		if tagValue != "" {
			agg.byTag[tagName+"|"+tagValue] += cost
		}
	}
}

func parquetString(row parquet.Row, colIdx int) string {
	if colIdx < 0 || colIdx >= len(row) {
		return ""
	}
	return row[colIdx].String()
}

func parquetFloat64(row parquet.Row, colIdx int) float64 {
	if colIdx < 0 || colIdx >= len(row) {
		return 0
	}
	v := row[colIdx]
	switch v.Kind() {
	case parquet.Double:
		return v.Double()
	case parquet.Float:
		return float64(v.Float())
	case parquet.Int64:
		return float64(v.Int64())
	case parquet.Int32:
		return float64(v.Int32())
	default:
		// Try parsing from string representation.
		f, _ := strconv.ParseFloat(v.String(), 64)
		return f
	}
}

func (c *CURCollector) processCSV(ctx context.Context, client *s3.Client, key string, agg *curAggregation) error {
	out, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.cfg.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("get S3 object %s: %w", key, err)
	}
	defer out.Body.Close()

	var reader io.Reader = out.Body
	if strings.HasSuffix(key, ".gz") {
		gz, err := gzip.NewReader(out.Body)
		if err != nil {
			return fmt.Errorf("open gzip reader: %w", err)
		}
		defer gz.Close()
		reader = gz
	}

	csvReader := csv.NewReader(reader)
	csvReader.LazyQuotes = true
	csvReader.ReuseRecord = true

	header, err := csvReader.Read()
	if err != nil {
		return fmt.Errorf("read CSV header: %w", err)
	}

	colMap := csvColumnMap(header)

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read CSV row: %w", err)
		}
		aggregateCSVRow(record, colMap, agg)
	}

	return nil
}

type csvColMap struct {
	cost         int
	service      int
	resource     int
	usageType    int
	accountID    int
	lineItemType int
	tagCols      map[int]string // column index -> tag name
}

func csvColumnMap(header []string) csvColMap {
	m := csvColMap{
		cost:         -1,
		service:      -1,
		resource:     -1,
		usageType:    -1,
		accountID:    -1,
		lineItemType: -1,
		tagCols:      make(map[int]string),
	}

	for i, col := range header {
		switch col {
		case "lineItem/UnblendedCost":
			m.cost = i
		case "lineItem/ProductCode":
			m.service = i
		case "lineItem/ResourceId":
			m.resource = i
		case "lineItem/UsageType":
			m.usageType = i
		case "lineItem/UsageAccountId":
			m.accountID = i
		case "lineItem/LineItemType":
			m.lineItemType = i
		default:
			if strings.HasPrefix(col, "resourceTags/user:") {
				tagName := strings.TrimPrefix(col, "resourceTags/user:")
				m.tagCols[i] = tagName
			}
		}
	}

	return m
}

// AggregateCSVRow aggregates a single CSV row into the aggregation.
// Exported for testing.
func AggregateCSVRow(record []string, colMap csvColMap, agg *curAggregation) {
	aggregateCSVRow(record, colMap, agg)
}

func aggregateCSVRow(record []string, colMap csvColMap, agg *curAggregation) {
	if colMap.cost < 0 || colMap.cost >= len(record) {
		return
	}

	// Skip non-cost line item types (credits, refunds, negations).
	if colMap.lineItemType >= 0 {
		lit := csvField(record, colMap.lineItemType)
		if _, ok := costLineItemTypes[lit]; !ok {
			return
		}
	}

	cost, err := strconv.ParseFloat(record[colMap.cost], 64)
	if err != nil || cost == 0 {
		return
	}

	service := csvField(record, colMap.service)
	resourceID := csvField(record, colMap.resource)
	usageType := csvField(record, colMap.usageType)

	agg.total += cost

	if service != "" {
		agg.byService[service] += cost
	}

	if resourceID != "" && service != "" {
		key := service + "|" + resourceID
		rc := agg.byResource[key]
		rc.service = service
		rc.resourceID = resourceID
		rc.cost += cost
		agg.byResource[key] = rc
	}

	if usageType != "" && service != "" {
		agg.byUsageType[service+"|"+usageType] += cost
	}

	for colIdx, tagName := range colMap.tagCols {
		tagValue := csvField(record, colIdx)
		if tagValue != "" {
			agg.byTag[tagName+"|"+tagValue] += cost
		}
	}
}

func csvField(record []string, idx int) string {
	if idx < 0 || idx >= len(record) {
		return ""
	}
	return record[idx]
}

// billingPeriod returns the CUR billing period date strings (yyyymmdd format).
func billingPeriod(now time.Time) (string, string) {
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)
	return start.Format("20060102"), end.Format("20060102")
}

func topNResources(byResource map[string]resourceCost, n int) []resourceCost {
	resources := make([]resourceCost, 0, len(byResource))
	for _, rc := range byResource {
		resources = append(resources, rc)
	}

	sort.Slice(resources, func(i, j int) bool {
		return resources[i].cost > resources[j].cost
	})

	if len(resources) > n {
		resources = resources[:n]
	}
	return resources
}

// ParseManifest parses a CUR manifest JSON. Exported for testing.
func ParseManifest(data []byte) (*curManifest, error) {
	var m curManifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// CSVColumnMap builds a column mapping from a CSV header. Exported for testing.
func CSVColumnMap(header []string) csvColMap {
	return csvColumnMap(header)
}

// NewCURAggregation creates a new empty aggregation. Exported for testing.
func NewCURAggregation(period string) *curAggregation {
	return &curAggregation{
		byService:   make(map[string]float64),
		byResource:  make(map[string]resourceCost),
		byUsageType: make(map[string]float64),
		byTag:       make(map[string]float64),
		period:      period,
	}
}
