package auth_repos_memory

import "context"

func NewPermissionMemoryRepository() *PermissionMemoryRepository {
	return &PermissionMemoryRepository{}
}

type PermissionMemoryRepository struct {
	data []string
}

func (x *PermissionMemoryRepository) Get(_ context.Context) []string {
	return []string{
		"user.permission.get",
		"user.role.get",
		"user.role.insert",
		"user.role.update",
		"user.role.delete",
		"user.user.get",
		"user.user.insert",
		"user.user.update",
		"user.user.update.password",
		"user.user.delete",

		"keep.transaksi.get",
		"keep.transaksi.insert",
		"keep.transaksi.update",
		"keep.transaksi.delete",
		"keep.transaksi.trash",
		"keep.lokasi.get",
		"keep.barang.get",
		"keep.pos.get",
		"keep.pos.insert",
		"keep.pos.update",
		"keep.pos.delete",
		"keep.pos.trash",
		"keep.kantong.get",
		"keep.kantong.insert",
		"keep.kantong.update",
		"keep.kantong.delete",
		"keep.kantong.trash",
		"keep.kantong.history.get",
		"keep.kantong.history.insert",
		"keep.kantong.history.update",
		"keep.kantong.history.delete",

		"changelog.get",
		"changelog.insert",
		"changelog.update",
		"changelog.delete",
	}
}
