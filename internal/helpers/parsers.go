package helpers

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/table"
	filev1 "github.com/rtmelsov/adv-keeper/gen/go/proto/file/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func pbTSLocalStr(ts *timestamppb.Timestamp) string {
	if ts == nil || !ts.IsValid() {
		return "-" // not set / invalid
	}
	return ts.AsTime().In(time.Local).Format("2006-01-02 15:04")
}

func FilesToRows(fs *filev1.GetFilesResponse) []table.Row {
	rows := make([]table.Row, 0, len(fs.Files))
	for _, f := range fs.Files {
		rows = append(rows, table.Row{
			f.Filename,
			fmt.Sprintf("%d B", f.Size),
			pbTSLocalStr(f.CreatedAt),
			// f.Fileid,
		})
	}
	return rows
}
