package pb

import (
	"sort"
)

func (f *FileList) Sort() {
	sort.Slice(f.Files, func(i, j int) bool {
		return f.Files[i].GetEntry().GetName() < f.Files[j].GetEntry().GetName()
	})
}
