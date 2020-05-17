package cache

type CacherMock struct {
	Products   map[string]string
	Categories map[string]string
	Responses  map[string]string
}

func NewCacherMock() *CacherMock {
	return &CacherMock{
		Products:   make(map[string]string),
		Categories: make(map[string]string),
		Responses:  make(map[string]string),
	}
}

func (c *CacherMock) SetProduct(id string, prodStr string) error {
	c.Products[id] = prodStr
	return nil
}

func (c *CacherMock) GetProduct(id string) (string, error) {
	if p, ok := c.Products[id]; ok {
		return p, nil
	}
	return "", nil
}

func (c *CacherMock) DeleteProduct(id string) error {
	delete(c.Products, id)
	return nil
}

func (c *CacherMock) SetCategory(id string, catStr string) error {
	c.Categories[id] = catStr
	return nil
}

func (c *CacherMock) GetCategory(id string) (string, error) {
	if c, ok := c.Categories[id]; ok {
		return c, nil
	}
	return "", nil
}

func (c *CacherMock) DeleteCategory(id string) error {
	delete(c.Categories, id)
	return nil
}

func (c *CacherMock) GetAllProducts() (map[string]string, error) {
	return c.Products, nil
}

func (c *CacherMock) GetAllCategories() (map[string]string, error) {
	return c.Categories, nil
}

func (c *CacherMock) SetApiRequest(path string, serializedResponse string) error {
	c.Responses[path] = serializedResponse
	return nil
}

func (c *CacherMock) GetApiRequest(path string) (string, error) {
	if response, ok := c.Responses[path]; ok {
		return response, nil
	}
	return "", nil
}
