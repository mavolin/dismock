// Package dismock creates mocks for the Discord API.
// The names of the mocks correspond to arikawa's API wrapper names, but as
// this are http mocks, any discord library can be mocked.
//
// Mocking Requests for Metadata
//
// Besides the regular API calls, dismock also features mocks for fetching
// an entities meta data, e.g. an icon or a splash.
// In order to mock requests for an entity's meta data, you need to make sure
// that those requests are made with Mocker.Client, so that the requests are
// correctly redirected to the mock server.
//
// Mocking Errors
//
// To send a discord error, use the Mocker.Error method with the path of the
// endpoint that should return an error.
//
//
// Important Notes
//
// BUG(mavolin): Due to an inconvenient behavior of json.Unmarshal where on
// JSON null the the UnmarshalJSON method doesn't get called, there is no way
// to differentiate between option.NullX and omitted, and therefore both will
// be seen as equal.
//
// BUG(mavolin): Due to the way Snowflakes are serialized, all fields that
// don't have the omitempty flag and are set to 0, will be sent as JSON null.
package dismock
