package watcher

import "testing"

func TestAddRecursive(t *testing.T) {
	watcher := New(true)
	watcher.Start()
	watcher.AddRecursive("../../core")
	watcher.RemoveRecursive("../../core")
}
