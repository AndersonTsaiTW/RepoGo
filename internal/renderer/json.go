// Package renderer provides output rendering functionality for different formats.
package renderer

import (
	"encoding/json"
	"io"

	"github.com/AndersonTsaiTW/RepoGo/internal/models"
)

// RenderJSON renders the output document in JSON format with indentation.
func RenderJSON(w io.Writer, doc models.OutputDoc) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(doc)
}
