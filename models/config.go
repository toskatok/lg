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

// Config is a configuration structure for I1820/lg
type Config struct {
	Generator struct {
		Name string
		Info interface{} // this structure is not used in config, it is passed to generators
	}
	Token    string
	Messages []map[string]interface{}
}
