package rest

// This test is commented out because it requires a running server to work.
/* func TestRestClient(t *testing.T) {
	// Create a new rest client
	client := NewDefaultRestClient()
	// Create a new request
	res := client.Get("http://localhost:3002/entity").WithCache(1024 * time.Second).Do()
	to := make(map[string]interface{})
	err1 := res.MapTo(&to)
	assert.Nil(t, err1)
	// Execute the request
	res2 := client.Get("http://localhost:3002/entity").WithCache(1024 * time.Second).Do()
	to2 := make(map[string]interface{})
	err2 := res2.MapTo(&to2)
	assert.Nil(t, err2)
	// Check the status code
	assert.Equal(t, res, res2)
} */
