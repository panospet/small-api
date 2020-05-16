package cache

type Cacher interface {
	SetProduct(id string, prodStr string) error
	GetProduct(id string) (string, error)
	DeleteProduct(id string) error
	SetCategory(id string, catStr string) error
	GetCategory(id string) (string, error)
	DeleteCategory(id string) error
	GetAllProducts() (map[string]string, error)
	GetAllCategories() (map[string]string, error)
	SetApiRequest(path string, serializedResponse string) error
	GetApiRequest(path string) (string, error)
}
