// Package bilibili provides a Go 1.18 Bilibili client library.
//
// The package is organized around four layers:
//  1. Client: lifecycle, configuration, authentication, module accessors.
//  2. Transport: request execution, response decoding, anti-spider signing.
//  3. Endpoint: declarative API descriptors inspired by bilibili-api's Api model.
//  4. Module: typed domain services such as video, user, search, live, and login.
package bilibili
