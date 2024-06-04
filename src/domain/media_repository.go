package domain

type MediaRepository interface {
	Create(file MediaFile) (MediaFile, error)
	Find(id string) (MediaFile, error)
	FindByUser(u string) ([]MediaFile, error)
	Delete(id, u string) error
}
