// Package migrations expone los archivos SQL de migración embebidos en el
// binario, para que funcionen igual en local que en Docker.
package migrations

import "embed"

// FS contiene los archivos *.sql de este directorio, ordenados por nombre.
//
//go:embed *.sql
var FS embed.FS
