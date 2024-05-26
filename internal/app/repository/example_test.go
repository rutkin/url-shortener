package repository

func ExampleinMemoryRepository_CreateURL() {
	repository := NewInMemoryRepository()
	repository.CreateURL(URLRecord{ID: "1", URL: "http://example.com", UserID: "123"})
}
