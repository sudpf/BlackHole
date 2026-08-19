package output

import (
	"BlackHole/internal/stash/config"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/olivere/elastic/v7"
	"github.com/rogpeppe/go-internal/semver"
	"github.com/zeromicro/go-zero/core/executors"
	"github.com/zeromicro/go-zero/core/logx"

	jsoniter "github.com/json-iterator/go"
)

const es8Version = "8.0.0"

type (
	EsWriter struct {
		ctx       context.Context
		docType   string
		esVersion string
		client    *elastic.Client
		inserter  *executors.ChunkExecutor
		indexer   *Index
	}

	valueWithIndex struct {
		index string
		val   string
	}
)

func NewElasticSearchWriter(ctx context.Context, c *config.ElasticSearchConf) (*EsWriter, error) {
	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(c.Hosts...),
		elastic.SetGzip(c.Compress),
		elastic.SetBasicAuth(c.Username, c.Password),
	)
	if err != nil {
		return nil, err
	}

	version, err := client.ElasticsearchVersion(c.Hosts[0])
	if err != nil {
		client.Stop()
		return nil, err
	}

	writer := EsWriter{
		ctx:       ctx,
		docType:   c.DocType,
		client:    client,
		esVersion: version,
	}
	writer.inserter = executors.NewChunkExecutor(writer.execute, executors.WithChunkBytes(c.MaxChunkBytes))

	var loc *time.Location
	if len(c.TimeZone) > 0 {
		loc, err = time.LoadLocation(c.TimeZone)
		if err != nil {
			client.Stop()
			return nil, fmt.Errorf("load Elasticsearch timezone %q: %w", c.TimeZone, err)
		}
	} else {
		loc = time.Local
	}
	writer.indexer = NewIndex(client, c.Index, loc)
	return &writer, nil
}

func (w *EsWriter) Write(ctx context.Context, val map[string]interface{}) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	index := w.indexer.GetIndex(ctx, val)

	bs, err := jsoniter.Marshal(val)
	if err != nil {
		return err
	}

	return w.inserter.Add(valueWithIndex{
		index: index,
		val:   string(bs),
	}, len(string(bs)))
}

func (w *EsWriter) Close(ctx context.Context) error {
	if w == nil {
		return nil
	}

	w.inserter.Flush()
	w.inserter.Wait()
	w.client.Stop()
	return nil
}

func (w *EsWriter) execute(vals []interface{}) {
	start := time.Now()
	if err := w.executeBatch(vals); err != nil {
		logx.Error(err)
		recordBatchWrite("elasticsearch", "failure", len(vals), time.Since(start))
		return
	}
	recordBatchWrite("elasticsearch", "success", len(vals), time.Since(start))
}

func (w *EsWriter) executeBatch(vals []interface{}) error {
	var bulk = w.client.Bulk()
	for _, val := range vals {
		pair := val.(valueWithIndex)
		req := elastic.NewBulkIndexRequest().Index(pair.index)
		if isSupportType(w.esVersion) && len(w.docType) > 0 {
			req = req.Type(w.docType)
		}
		req = req.Doc(pair.val)
		bulk.Add(req)
	}
	resp, err := bulk.Do(w.ctx)
	if err != nil {
		return fmt.Errorf("bulk index: %w", err)
	}

	// bulk error in docs will report in response items
	if !resp.Errors {
		return nil
	}

	var batchErr error
	for _, imap := range resp.Items {
		for _, item := range imap {
			if item.Error == nil {
				continue
			}

			logx.Error(item.Error)
			batchErr = errors.Join(batchErr, errorFromElasticDetails(item.Error))
		}
	}

	if batchErr != nil {
		return fmt.Errorf("bulk item errors: %w", batchErr)
	}
	return nil
}

func isSupportType(version string) bool {
	// es8.x not support type field
	return semver.Compare(version, es8Version) < 0
}

func errorFromElasticDetails(details *elastic.ErrorDetails) error {
	if details == nil {
		return nil
	}
	if details.Index != "" {
		return fmt.Errorf("index %s: %s: %s", details.Index, details.Type, details.Reason)
	}
	return fmt.Errorf("%s: %s", details.Type, details.Reason)
}
