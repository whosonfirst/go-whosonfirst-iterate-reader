package reader

import (
	"context"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/iterator"
	"io"
	"net/url"
	"path/filepath"
	"sync/atomic"
	"testing"
)

func TestReaderEmitter(t *testing.T) {

	ctx := context.Background()

	rel_path := "fixtures/data"
	abs_path, err := filepath.Abs(rel_path)

	if err != nil {
		t.Fatalf("Failed to derive absolute path for '%s', %v", rel_path, err)
	}

	reader_uri := fmt.Sprintf("fs://%s", abs_path)

	q := url.Values{}
	q.Set("reader", reader_uri)

	tests := map[string]int32{
		fmt.Sprintf("reader://?%s", q.Encode()):                                                 1,
		fmt.Sprintf("reader://?include=properties.sfomuseum:placetype=locality&%s", q.Encode()): 0,
	}

	for iter_uri, expected_count := range tests {

		count := int32(0)

		iter_cb := func(ctx context.Context, path string, r io.ReadSeeker, args ...interface{}) error {
			atomic.AddInt32(&count, 1)
			return nil
		}

		iter, err := iterator.NewIterator(ctx, iter_uri, iter_cb)

		if err != nil {
			t.Fatalf("Failed to create new iterator, %v", err)
		}

		err = iter.IterateURIs(ctx, "102/527/513/102527513.geojson")

		if err != nil {
			t.Fatalf("Failed to iterate %s, %v", abs_path, err)
		}

		if count != expected_count {
			t.Fatalf("Unexpected count: %d", count)
		}
	}
}
