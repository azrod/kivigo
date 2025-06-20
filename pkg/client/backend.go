package client

import "context"

// Get retrieves the value stored under the specified key and decodes it into dest.
// Returns an error if the key does not exist or decoding fails.
//
// Example:
//
//	var value string
//	err := client.Get(ctx, "myKey", &value)
func (c Client) Get(ctx context.Context, key string, value any) error {
	vV, err := c.GetRaw(ctx, key)
	if err != nil {
		return err
	}

	return c.opts.Encoder.Decode(vV, value)
}

// Set stores the given value under the specified key.
// Returns an error if the operation fails.
//
// Example:
//
//	err := client.Set(ctx, "myKey", "myValue")
func (c Client) Set(ctx context.Context, key string, value any) error {
	vV, err := c.opts.Encoder.Encode(value)
	if err != nil {
		return err
	}

	return c.SetRaw(ctx, key, vV)
}
