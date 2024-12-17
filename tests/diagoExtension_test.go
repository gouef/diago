package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDiagoExtension(t *testing.T) {
	t.Run("Test GetPanelHtml", func(t *testing.T) {
		mock := &MockDiagoExtension{PanelHtml: "<div>Panel</div>"}

		result := mock.GetPanelHtml(nil)
		assert.Equal(t, "<div>Panel</div>", result)
	})

	t.Run("Test GetHtml", func(t *testing.T) {
		mock := &MockDiagoExtension{Html: "<div>Content</div>"}

		result := mock.GetHtml(nil)
		assert.Equal(t, "<div>Content</div>", result)
	})

	t.Run("Test GetJSHtml", func(t *testing.T) {
		mock := &MockDiagoExtension{JSHtml: "<script>console.log('test');</script>"}

		result := mock.GetJSHtml(nil)
		assert.Equal(t, "<script>console.log('test');</script>", result)
	})

	t.Run("Test BeforeNext", func(t *testing.T) {
		mock := &MockDiagoExtension{}

		mock.BeforeNext(nil)

		assert.True(t, mock.BeforeCalled)
	})

	t.Run("Test AfterNext", func(t *testing.T) {
		mock := &MockDiagoExtension{}

		mock.AfterNext(nil)

		assert.True(t, mock.AfterCalled)
	})
}
