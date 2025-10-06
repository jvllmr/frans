package tasks

import (
	"testing"

	"github.com/jvllmr/frans/internal/services"
	"github.com/jvllmr/frans/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestFileLifecycleTaskPreserveFileData(t *testing.T) {
	cfg := testutil.SetupTestConfig()
	db := testutil.SetupTestDBClient(t)
	fs := services.NewFileService(cfg, db)
	testUser := testutil.SetupTestUser(t, db, nil)
	testFile := testutil.SetupTestFile(
		t,
		cfg,
		db,
		"test.txt",
		"Hello there!",
		testUser,
		"single",
		0,
		0,
		1,
	)
	testFile2 := testutil.SetupTestFile(
		t,
		cfg,
		db,
		"test2.txt",
		"Hello there!",
		testUser,
		"single",
		0,
		0,
		1,
	)

	filesDatas := db.FileData.Query().WithFiles().WithUsers().AllX(t.Context())
	assert.Equal(t, 1, len(filesDatas))
	assert.Equal(t, 2, len(filesDatas[0].Edges.Files))
	assert.Equal(t, 1, len(filesDatas[0].Edges.Users))

	db.File.UpdateOne(testFile).SetTimesDownloaded(1).ExecX(t.Context())
	FileLifecycleTask(db, fs)

	filesDatas = db.FileData.Query().WithFiles().WithUsers().AllX(t.Context())
	assert.Equal(t, 1, len(filesDatas))
	assert.Equal(t, 1, len(filesDatas[0].Edges.Files))
	assert.Equal(t, 1, len(filesDatas[0].Edges.Users))

	db.File.UpdateOne(testFile2).SetTimesDownloaded(1).ExecX(t.Context())
	FileLifecycleTask(db, fs)

	filesDatas = db.FileData.Query().WithFiles().WithUsers().AllX(t.Context())
	assert.Equal(t, 0, len(filesDatas))
}
