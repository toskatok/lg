/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 15-12-2018
 * |
 * | File Name:     config.go
 * +===============================================
 */

package models

// Config is a configuration structure for toskatok/lg that can be passed from HTTP or file
type Config struct {
	Generator struct {
		Name string
		Info interface{} // this structure is not used in config, it is passed to generators to configure them
	}
	Token    string
	Messages []map[string]interface{}
}
