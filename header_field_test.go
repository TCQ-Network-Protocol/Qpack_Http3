//	Project: QPACK HTTP3
//	Author: Trần Nguyên Hiền (c)
//	Major: Electronic And Communication Engineering
//	Email: trannguyenhien29085@gmail.com
//	Date: 2/3/2026
//	GPL-3.0 Licence
//
// ----------------------------------------------------------------
package qpack

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHeaderFieldIsPseudo(t *testing.T) {
	t.Run("Pseudo headers", func(t *testing.T) {
		require.True(t, (HeaderField{Name: ":status"}).IsPseudo())
		require.True(t, (HeaderField{Name: ":authority"}).IsPseudo())
		require.True(t, (HeaderField{Name: ":foobar"}).IsPseudo())
	})

	t.Run("Non-pseudo headers", func(t *testing.T) {
		require.False(t, (HeaderField{Name: "status"}).IsPseudo())
		require.False(t, (HeaderField{Name: "foobar"}).IsPseudo())
	})
}
