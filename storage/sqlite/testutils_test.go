package sqlite

import (
	"fmt"
	"time"

	"github.com/coline-carle/zhevra/storage/model"
)

func tearDown(s *Storage, table string, where string, params ...interface{}) {
	_, err := s.DB.Exec(
		fmt.Sprintf("DELETE FROM \"%s\" WHERE %s", table, where),
		params...)
	if err != nil {
		panic(fmt.Sprintf("Problem tearing down %s data: %v", table, err))
	}
}

func testUtilCreateRelease(ID int64, addonID int64) model.CurseRelease {
	return model.CurseRelease{
		ID:           ID,
		Filename:     fmt.Sprintf("Filename-%d.zip", ID),
		CreatedAt:    time.Date(2018, 9, 9, 0, 0, 0, 0, time.UTC),
		URL:          fmt.Sprintf("https://url/filename-%d", ID),
		GameVersions: []int{7*0x100*0x100 + 3*0X100 + 2},
		AddonID:      addonID,
		Directories: []string{
			fmt.Sprintf("directory-%d-1", ID),
			fmt.Sprintf("directory-%d-2", ID),
		},
	}
}

func testUtilCreateAddon(ID int64) model.CurseAddon {
	return model.CurseAddon{
		ID:            ID,
		Name:          fmt.Sprintf("Name-%d", ID),
		Summary:       fmt.Sprintf("Summary-%d", ID),
		URL:           fmt.Sprintf("https://url/addon-%d", ID),
		DownloadCount: 33 * ID,
		Releases:      []model.CurseRelease{},
	}
}
