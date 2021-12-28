// =================================================================
//
// Work of the U.S. Department of Defense, Defense Digital Service.
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package ziputil

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
)

func UnmarshalFile(r *zip.Reader, name string, v interface{}) error {
	fr, err := r.Open(name)
	if err != nil {
		return fmt.Errorf("error opening file %q in zip file: %w", name, err)
	}

	data, err := io.ReadAll(fr)
	if err != nil {
		return fmt.Errorf("error reading file %q in zip file: %w", name, err)
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		return fmt.Errorf("error unmarshaling file %q in zip file: %w", name, err)
	}

	return nil
}
