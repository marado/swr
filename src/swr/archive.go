/*  Space Wars Rebellion Mud
 *  Copyright (C) 2022 @{See Authors}
 *
 *  This program is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with this program.  If not, see <https://www.gnu.org/licenses/>.
 *
 */
package swr

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type ArchiveService struct {
	T *time.Ticker
}

var _archiver *ArchiveService

func Archiver() *ArchiveService {
	if _archiver == nil {
		_archiver = new(ArchiveService)
		_archiver.T = time.NewTicker(time.Duration(1) * time.Hour)
	}
	return _archiver
}
func BackupStart() {
	ar := Archiver()
	for {
		t := <-ar.T.C
		DoBackup(t)
		DoBackupCleanup(t)
	}
}

func DoBackup(t time.Time) {
	log.Printf("***** BACKUP STARTED *****\r\n")
	// Lock the database to prevent disk activity while we tar a backup archive
	db := DB()
	db.Lock()
	defer db.Unlock()

	if runtime.GOOS == "windows" {
		_, err := exec.Command("tar", "-cJf", fmt.Sprintf("backup\\archive-%s.tar.xz", t.Format(time.RFC3339Nano)), "data\\*").Output()
		ErrorCheck(err)
	} else {
		_, err := exec.Command("tar", "-cJf", fmt.Sprintf("./backup/archive-%s.tar.xz", t.Format(time.RFC3339Nano)), "data/*").Output()
		ErrorCheck(err)
	}

	log.Printf("***** BACKUP COMPLETE *****\r\n")
}

func DoBackupCleanup(t time.Time) {
	archives, err := ioutil.ReadDir("backup")
	ErrorCheck(err)
	for _, archive := range archives {
		p := archive.Name()
		p = strings.ReplaceAll(p, ".tar.xz", "")
		p = strings.ReplaceAll(p, "archive-", "")
		ar_time, err := time.Parse(time.RFC3339Nano, p)
		ErrorCheck(err)
		cut_time := t.Add(-72 * time.Hour) // 3 days worth of backups are stored.
		if ar_time.Before(cut_time) {
			if runtime.GOOS == "windows" {
				err := os.Remove(fmt.Sprintf("backup\\%s", archive.Name()))
				ErrorCheck(err)
			} else {
				err := os.Remove(fmt.Sprintf("backup/%s", archive.Name()))
				ErrorCheck(err)
			}

		}
	}
}
