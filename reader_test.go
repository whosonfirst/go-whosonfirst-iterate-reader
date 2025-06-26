package reader

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"path/filepath"
	"testing"

	"github.com/whosonfirst/go-whosonfirst-iterate/v3"
)

func TestReaderIterator(t *testing.T) {

	ctx := context.Background()

	rel_path := "fixtures/data"
	abs_path, err := filepath.Abs(rel_path)

	if err != nil {
		t.Fatalf("Failed to derive absolute path for '%s', %v", rel_path, err)
	}

	reader_uri := fmt.Sprintf("fs://%s", abs_path)

	q := url.Values{}
	q.Set("reader", reader_uri)

	tests := map[string]int{
		fmt.Sprintf("reader://?%s", q.Encode()):                                                 1,
		fmt.Sprintf("reader://?include=properties.sfomuseum:placetype=locality&%s", q.Encode()): 0,
	}

	for iter_uri, expected_count := range tests {

		count := 0

		iter, err := iterate.NewIterator(ctx, iter_uri)

		if err != nil {
			t.Fatalf("Failed to create new iterator, %v", err)
		}

		for rec, err := range iter.Iterate(ctx, "102/527/513/102527513.geojson") {

			if err != nil {
				t.Fatalf("Failed to walk file, %v", err)
				break
			}

			defer rec.Body.Close()
			_, err = io.ReadAll(rec.Body)

			if err != nil {
				t.Fatalf("Failed to read body for %s, %v", rec.Path, err)
			}

			_, err = rec.Body.Seek(0, 0)

			if err != nil {
				t.Fatalf("Failed to rewind body for %s, %v", rec.Path, err)
			}

			_, err = io.ReadAll(rec.Body)

			if err != nil {
				t.Fatalf("Failed second read body for %s, %v", rec.Path, err)
			}

			count += 1
		}

		if count != expected_count {
			t.Fatalf("Unexpected count (%s): %d", iter_uri, count)
		}

		err = iter.Close()

		if err != nil {
			t.Fatalf("Failed to close iterator, %v", err)
		}
	}
}
