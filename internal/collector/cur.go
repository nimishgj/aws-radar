package collector

import (
	"archive/zip"
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

// curManifest represents the CUR manifest JSON structure. Legacy CUR uses
// reportKeys; CUR 2.0 (Data Exports) manifests have used both dataFileS3Keys
// (plain strings) and dataFiles/files (objects), so all variants are accepted.
type curManifest struct {
	ReportKeys     []string           `json:"reportKeys"`
	DataFileS3Keys []string           `json:"dataFileS3Keys"`
	DataFiles      []curManifestEntry `json:"dataFiles"`
	Files          []curManifestEntry `json:"files"`
}

// curManifestEntry is either a plain string key/URI or an object with a
// key/s3Key field, depending on the export tooling that wrote the manifest.
type curManifestEntry struct {
	Key   string `json:"key"`
	S3Key string `json:"s3Key"`
}

func (e *curManifestEntry) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == '"' {
		return json.Unmarshal(data, &e.Key)
	}
	type entry curManifestEntry
	return json.Unmarshal(data, (*entry)(e))
}

func (m *curManifest) keys() []string {
	var keys []string
	keys = append(keys, m.ReportKeys...)
	keys = append(keys, m.DataFileS3Keys...)
	for _, e := range append(m.DataFiles, m.Files...) {
		if e.S3Key != "" {
			keys = append(keys, e.S3Key)
		} else if e.Key != "" {
			keys = append(keys, e.Key)
		}
	}
	for i, k := range keys {
		// CUR 2.0 keys may be full s3:// URIs; reduce to the object key.
		if strings.HasPrefix(k, "s3://") {
			if idx := strings.IndexByte(strings.TrimPrefix(k, "s3://"), '/'); idx >= 0 {
				keys[i] = strings.TrimPrefix(k, "s3://")[idx+1:]
			}
		}
	}
	return keys
}

// curAggregation holds aggregated cost data from CUR files. Costs are tracked in
// both the unblended and amortized views (see unblended/amortizedContribution).
type curAggregation struct {
	totalUnblended float64
	totalAmortized float64
	byService      map[svcKey]*costPair    // {service, region} -> costs
	byResource     map[string]*resourceAgg // key: "service|resource_id|region"
	byUsageType    map[string]float64      // key: "service|usage_type" (unblended)
	byTag          map[string]float64      // key: "tag_key|tag_value" (unblended)
	period         string
}

type svcKey struct{ service, region string }

type costPair struct{ unblended, amortized float64 }

type resourceAgg struct {
	service, resourceID, region, az, specKey, specVal string
	unblended, amortized, usage                       float64
}

// curRow is a normalized CUR line item (extracted from parquet or CSV) that the
// format-independent aggregateRow consumes.
type curRow struct {
	service, resourceID, usageType, lineItemType string
	region, az, specKey, specVal                 string
	unblended, effSP, effRI, usageAmount         float64
	tags                                         map[string]string
}

func newCurAggregation(period string) *curAggregation {
	return &curAggregation{
		byService:   make(map[svcKey]*costPair),
		byResource:  make(map[string]*resourceAgg),
		byUsageType: make(map[string]float64),
		byTag:       make(map[string]float64),
		period:      period,
	}
}

// aggregateRow folds one normalized CUR line into the aggregation, tracking the
// unblended and amortized cost views plus usage quantity.
func aggregateRow(agg *curAggregation, r curRow) {
	u := unblendedContribution(r.lineItemType, r.unblended)
	a := amortizedContribution(r.lineItemType, r.unblended, r.effSP, r.effRI)
	if u == 0 && a == 0 && r.usageAmount == 0 {
		return
	}
	agg.totalUnblended += u
	agg.totalAmortized += a

	if r.service != "" {
		sk := svcKey{r.service, r.region}
		cp := agg.byService[sk]
		if cp == nil {
			cp = &costPair{}
			agg.byService[sk] = cp
		}
		cp.unblended += u
		cp.amortized += a
	}

	if r.resourceID != "" && r.service != "" {
		key := r.service + "|" + r.resourceID + "|" + r.region
		ra := agg.byResource[key]
		if ra == nil {
			ra = &resourceAgg{service: r.service, resourceID: r.resourceID, region: r.region,
				az: r.az, specKey: r.specKey, specVal: r.specVal}
			agg.byResource[key] = ra
		}
		ra.unblended += u
		ra.amortized += a
		ra.usage += r.usageAmount
	}

	if r.usageType != "" && r.service != "" {
		agg.byUsageType[r.service+"|"+r.usageType] += u
	}
	for k, v := range r.tags {
		if v != "" {
			agg.byTag[k+"|"+v] += u
		}
	}
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

	metrics.CURTotalCost.WithLabelValues(account, accountName, "unblended", agg.period).Set(agg.totalUnblended)
	metrics.CURTotalCost.WithLabelValues(account, accountName, "amortized", agg.period).Set(agg.totalAmortized)

	for sk, cp := range agg.byService {
		metrics.CURCostByService.WithLabelValues(account, accountName, sk.service, sk.region, "unblended", agg.period).Set(cp.unblended)
		metrics.CURCostByService.WithLabelValues(account, accountName, sk.service, sk.region, "amortized", agg.period).Set(cp.amortized)
	}

	// Emit top N resources by cost (unblended + amortized), plus usage quantity.
	topResources := topNResources(agg.byResource, c.cfg.MaxResources)
	for _, rc := range topResources {
		metrics.CURCostByResource.WithLabelValues(account, accountName, rc.service, rc.resourceID, rc.region, rc.az, rc.specKey, rc.specVal, "unblended", agg.period).Set(rc.unblended)
		metrics.CURCostByResource.WithLabelValues(account, accountName, rc.service, rc.resourceID, rc.region, rc.az, rc.specKey, rc.specVal, "amortized", agg.period).Set(rc.amortized)
		metrics.CURUsageAmount.WithLabelValues(account, accountName, rc.service, rc.resourceID, rc.region, rc.az, rc.specKey, rc.specVal, agg.period).Set(rc.usage)
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

	periodStart, _ := billingPeriod(now)

	var manifest *curManifest
	var err error
	candidates := c.manifestKeyCandidates(now)
	for _, manifestKey := range candidates {
		manifest, err = c.downloadManifest(ctx, client, manifestKey)
		if err == nil {
			break
		}
		log.Debug().Err(err).Str("key", manifestKey).Msg("CUR manifest candidate not found")
	}
	if manifest == nil {
		return fmt.Errorf("download CUR manifest (tried %v): %w", candidates, err)
	}

	reportKeys := manifest.keys()
	if len(reportKeys) == 0 {
		log.Warn().Msg("CUR manifest has no report keys")
		return nil
	}

	agg := newCurAggregation(periodStart)

	format := c.detectFormat(reportKeys[0])

	for _, key := range reportKeys {
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

// manifestKeyCandidates returns manifest paths to try, in order. Legacy CUR:
// {prefix}/{name}/{YYYYMMDD-YYYYMMDD}/{name}-Manifest.json. CUR 2.0 (Data
// Exports): {prefix}/{name}/metadata/BILLING_PERIOD=YYYY-MM/ containing either
// {name}-Manifest.json or Manifest.json.
func (c *CURCollector) manifestKeyCandidates(now time.Time) []string {
	if strings.EqualFold(c.cfg.ReportType, "cur2") {
		partition := "BILLING_PERIOD=" + now.UTC().Format("2006-01")
		base := fmt.Sprintf("%s/%s/metadata/%s", c.cfg.Prefix, c.cfg.ReportName, partition)
		return []string{
			fmt.Sprintf("%s/%s-Manifest.json", base, c.cfg.ReportName),
			base + "/Manifest.json",
		}
	}
	periodStart, periodEnd := billingPeriod(now)
	return []string{fmt.Sprintf("%s/%s/%s-%s/%s-Manifest.json",
		c.cfg.Prefix, c.cfg.ReportName, periodStart, periodEnd, c.cfg.ReportName)}
}

func (c *CURCollector) downloadManifest(ctx context.Context, client *s3.Client, key string) (*curManifest, error) {
	out, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.cfg.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	defer func() { _ = out.Body.Close() }()

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
	if strings.HasSuffix(reportKey, ".zip") {
		return "zip"
	}
	return "csv"
}

func (c *CURCollector) processFile(ctx context.Context, client *s3.Client, key, format string, agg *curAggregation) error {
	switch format {
	case "parquet":
		return c.processParquet(ctx, client, key, agg)
	case "zip":
		return c.processZIP(ctx, client, key, agg)
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
	defer func() { _ = out.Body.Close() }()

	// Parquet reader needs a file; download to temp.
	tmp, err := os.CreateTemp("", "cur-*.parquet")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	defer func() { _ = os.Remove(tmp.Name()) }()
	defer func() { _ = tmp.Close() }()

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
				_ = rows.Close()
				return fmt.Errorf("read parquet rows: %w", err)
			}
		}
		_ = rows.Close()
	}

	return nil
}

// costLineItemTypes are the line_item_line_item_type values that represent
// actual unblended spend, matching the AWS Cost Explorer "unblended cost" view.
//
// Credits and refunds are excluded. Savings Plans are counted via their
// RecurringFee/UpfrontFee only: SavingsPlanCoveredUsage carries the public
// on-demand price of the covered usage and SavingsPlanNegation carries its
// exact offset, so the two net to zero in unblended cost. Including
// CoveredUsage while excluding Negation would double-count the covered compute
// (e.g. inflating EC2), so both are left out here. DiscountedUsage (RI-covered
// usage) carries ~$0 unblended cost and is harmless to include.
var costLineItemTypes = map[string]struct{}{
	"Usage":                   {},
	"DiscountedUsage":         {},
	"Tax":                     {},
	"Fee":                     {},
	"RIFee":                   {},
	"SavingsPlanRecurringFee": {},
	"SavingsPlanUpfrontFee":   {},
}

// productAttrs pulls the region, availability zone, and a single adaptive
// "spec" (instance type / volume type / db engine) out of the CUR product map.
func productAttrs(product map[string]string) (region, az, specKey, specVal string) {
	region = product["region"]
	az = product["availability_zone"]
	switch {
	case product["database_engine"] != "":
		specKey, specVal = "engine", product["database_engine"]
	case product["volume_api_name"] != "":
		specKey, specVal = "volume_type", product["volume_api_name"]
	case product["instance_type_family"] != "":
		specKey, specVal = "instance_type", product["instance_type_family"]
	}
	return
}

// unblendedContribution returns the line's cost for the AWS "unblended" view:
// the raw unblended cost for real-charge line-item types, 0 otherwise
// (credits/refunds/negations and SP-covered usage are excluded — SP is counted
// via its recurring/upfront fee).
func unblendedContribution(lit string, unblended float64) float64 {
	if lit == "" { // no line-item-type column: treat as a plain cost line
		return unblended
	}
	if _, ok := costLineItemTypes[lit]; ok {
		return unblended
	}
	return 0
}

// amortizedContribution returns the line's cost for the "amortized" view, where
// Savings Plan / RI commitment is attributed to the usage that consumed it:
// covered usage is valued at its effective cost and the commitment fees become 0.
func amortizedContribution(lit string, unblended, effSP, effRI float64) float64 {
	switch lit {
	case "SavingsPlanCoveredUsage":
		return effSP
	case "DiscountedUsage":
		return effRI
	case "SavingsPlanRecurringFee", "SavingsPlanUpfrontFee", "RIFee":
		return 0
	case "Usage", "Tax", "Fee", "":
		return unblended
	default:
		return 0
	}
}

type parquetColIndices struct {
	cost         int
	service      int
	resource     int
	usageType    int
	accountID    int
	lineItemType int
	usageAmount  int
	effSP        int            // savings_plan_savings_plan_effective_cost
	effRI        int            // reservation_effective_cost
	productKey   int            // product.key_value.key
	productVal   int            // product.key_value.value
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
		usageAmount:  -1,
		effSP:        -1,
		effRI:        -1,
		productKey:   -1,
		productVal:   -1,
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
		case "line_item_usage_amount":
			idx.usageAmount = i
		case "savings_plan_savings_plan_effective_cost":
			idx.effSP = i
		case "reservation_effective_cost":
			idx.effRI = i
		case "product.key_value.key":
			idx.productKey = i
		case "product.key_value.value":
			idx.productVal = i
		default:
			if strings.HasPrefix(path, "resource_tags_user_") {
				tagName := strings.TrimPrefix(path, "resource_tags_user_")
				idx.tagCols[i] = tagName
			}
		}
	}

	return idx
}

// productMapFromRow reconstructs the CUR `product` map from the repeated
// key_value columns in a parquet row (key then its value, per map entry).
func productMapFromRow(row parquet.Row, keyCol, valCol int) map[string]string {
	if keyCol < 0 || valCol < 0 {
		return nil
	}
	m := make(map[string]string)
	var lastKey string
	for _, v := range row {
		switch v.Column() {
		case keyCol:
			lastKey = v.String()
		case valCol:
			if lastKey != "" {
				m[lastKey] = v.String()
				lastKey = ""
			}
		}
	}
	return m
}

func aggregateParquetRow(row parquet.Row, idx parquetColIndices, agg *curAggregation) {
	if idx.cost < 0 {
		return
	}
	region, az, specKey, specVal := productAttrs(productMapFromRow(row, idx.productKey, idx.productVal))

	var tags map[string]string
	for colIdx, tagName := range idx.tagCols {
		if tv := parquetString(row, colIdx); tv != "" {
			if tags == nil {
				tags = make(map[string]string)
			}
			tags[tagName] = tv
		}
	}

	aggregateRow(agg, curRow{
		service:      parquetString(row, idx.service),
		resourceID:   parquetString(row, idx.resource),
		usageType:    parquetString(row, idx.usageType),
		lineItemType: parquetString(row, idx.lineItemType),
		region:       region,
		az:           az,
		specKey:      specKey,
		specVal:      specVal,
		unblended:    parquetFloat64(row, idx.cost),
		effSP:        parquetFloat64(row, idx.effSP),
		effRI:        parquetFloat64(row, idx.effRI),
		usageAmount:  parquetFloat64(row, idx.usageAmount),
		tags:         tags,
	})
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
	defer func() { _ = out.Body.Close() }()

	var reader io.Reader = out.Body
	if strings.HasSuffix(key, ".gz") {
		gz, err := gzip.NewReader(out.Body)
		if err != nil {
			return fmt.Errorf("open gzip reader: %w", err)
		}
		defer func() { _ = gz.Close() }()
		reader = gz
	}

	return parseCSVStream(reader, agg)
}

func parseCSVStream(reader io.Reader, agg *curAggregation) error {
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

func (c *CURCollector) processZIP(ctx context.Context, client *s3.Client, key string, agg *curAggregation) error {
	out, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.cfg.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("get S3 object %s: %w", key, err)
	}
	defer func() { _ = out.Body.Close() }()

	// ZIP reader needs random access; download to temp.
	tmp, err := os.CreateTemp("", "cur-*.zip")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	defer func() { _ = os.Remove(tmp.Name()) }()
	defer func() { _ = tmp.Close() }()

	if _, err := io.Copy(tmp, out.Body); err != nil {
		return fmt.Errorf("download zip to temp: %w", err)
	}

	stat, err := tmp.Stat()
	if err != nil {
		return fmt.Errorf("stat temp file: %w", err)
	}

	return parseZIP(tmp, stat.Size(), agg)
}

func parseZIP(r io.ReaderAt, size int64, agg *curAggregation) error {
	zr, err := zip.NewReader(r, size)
	if err != nil {
		return fmt.Errorf("open zip reader: %w", err)
	}

	for _, f := range zr.File {
		if !strings.HasSuffix(f.Name, ".csv") || strings.HasPrefix(f.Name, "__MACOSX/") {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("open zip entry %s: %w", f.Name, err)
		}
		parseErr := parseCSVStream(rc, agg)
		_ = rc.Close()
		if parseErr != nil {
			return fmt.Errorf("parse zip entry %s: %w", f.Name, parseErr)
		}
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
	usageAmount  int
	effSP        int
	effRI        int
	region       int
	az           int
	instanceType int
	volumeType   int
	engine       int
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
		usageAmount:  -1,
		effSP:        -1,
		effRI:        -1,
		region:       -1,
		az:           -1,
		instanceType: -1,
		volumeType:   -1,
		engine:       -1,
		tagCols:      make(map[int]string),
	}

	for i, col := range header {
		switch col {
		case "lineItem/UnblendedCost", "line_item_unblended_cost":
			m.cost = i
		case "lineItem/ProductCode", "line_item_product_code":
			m.service = i
		case "lineItem/ResourceId", "line_item_resource_id":
			m.resource = i
		case "lineItem/UsageType", "line_item_usage_type":
			m.usageType = i
		case "lineItem/UsageAccountId", "line_item_usage_account_id":
			m.accountID = i
		case "lineItem/LineItemType", "line_item_line_item_type":
			m.lineItemType = i
		case "lineItem/UsageAmount", "line_item_usage_amount":
			m.usageAmount = i
		case "savingsPlan/SavingsPlanEffectiveCost", "savings_plan_savings_plan_effective_cost":
			m.effSP = i
		case "reservation/EffectiveCost", "reservation_effective_cost":
			m.effRI = i
		case "product/region", "product_region", "product/regionCode", "product_region_code":
			m.region = i
		case "lineItem/AvailabilityZone", "line_item_availability_zone", "product/availabilityZone", "product_availability_zone":
			m.az = i
		case "product/instanceTypeFamily", "product_instance_type_family":
			m.instanceType = i
		case "product/volumeApiName", "product_volume_api_name":
			m.volumeType = i
		case "product/databaseEngine", "product_database_engine":
			m.engine = i
		default:
			switch {
			case strings.HasPrefix(col, "resourceTags/user:"):
				m.tagCols[i] = strings.TrimPrefix(col, "resourceTags/user:")
			case strings.HasPrefix(col, "resource_tags_user_"):
				m.tagCols[i] = strings.TrimPrefix(col, "resource_tags_user_")
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
	cost, _ := strconv.ParseFloat(csvField(record, colMap.cost), 64)

	product := map[string]string{}
	if v := csvField(record, colMap.region); v != "" {
		product["region"] = v
	}
	if v := csvField(record, colMap.az); v != "" {
		product["availability_zone"] = v
	}
	if v := csvField(record, colMap.instanceType); v != "" {
		product["instance_type_family"] = v
	}
	if v := csvField(record, colMap.volumeType); v != "" {
		product["volume_api_name"] = v
	}
	if v := csvField(record, colMap.engine); v != "" {
		product["database_engine"] = v
	}
	region, az, specKey, specVal := productAttrs(product)

	var tags map[string]string
	for colIdx, tagName := range colMap.tagCols {
		if tv := csvField(record, colIdx); tv != "" {
			if tags == nil {
				tags = make(map[string]string)
			}
			tags[tagName] = tv
		}
	}

	effSP, _ := strconv.ParseFloat(csvField(record, colMap.effSP), 64)
	effRI, _ := strconv.ParseFloat(csvField(record, colMap.effRI), 64)
	usageAmount, _ := strconv.ParseFloat(csvField(record, colMap.usageAmount), 64)

	aggregateRow(agg, curRow{
		service:      csvField(record, colMap.service),
		resourceID:   csvField(record, colMap.resource),
		usageType:    csvField(record, colMap.usageType),
		lineItemType: csvField(record, colMap.lineItemType),
		region:       region,
		az:           az,
		specKey:      specKey,
		specVal:      specVal,
		unblended:    cost,
		effSP:        effSP,
		effRI:        effRI,
		usageAmount:  usageAmount,
		tags:         tags,
	})
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

// topNResources returns the n most expensive resources by unblended cost.
func topNResources(byResource map[string]*resourceAgg, n int) []*resourceAgg {
	resources := make([]*resourceAgg, 0, len(byResource))
	for _, rc := range byResource {
		resources = append(resources, rc)
	}

	sort.Slice(resources, func(i, j int) bool {
		return resources[i].unblended > resources[j].unblended
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
	return newCurAggregation(period)
}
