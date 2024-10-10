package bucket

import "container/list"

var ExportGetList = (*Bucket).getList

func (b *Bucket) getList() *list.List {
	return b.list
}
