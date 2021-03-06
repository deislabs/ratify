/*
Copyright The Ratify Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package utils

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/deislabs/ratify/pkg/homedir"
)

func TestHomepathReplacement(t *testing.T) {

	sampleFoldername := "test"
	testPath := homedir.GetShortcutString() + string(os.PathSeparator) + sampleFoldername

	result := ReplaceHomeShortcut(testPath)
	home, err := os.UserHomeDir()
	expectedPath := home + "/" + sampleFoldername

	if result != expectedPath {
		t.Fatalf("sample input %v ,result expected to be %v, actual result %v, error %v", testPath, expectedPath, result, err)
	}
}

func TestReadCertificatesFromPath_InvalidPath(t *testing.T) {

	files, err := GetCertificatesFromPath("/invalid/path")

	expectedFileCount := 0
	if len(files) != expectedFileCount || err != nil {
		t.Fatalf("response length expected to be %v, actual %v, error %v", expectedFileCount, len(files), err)
	}
}

func TestReadCertificatesFromPath_NestedDirectory(t *testing.T) {
	// Setup to create nested directory structure
	testDir := "TestDirectory"
	nestedDir := ".nestedFolder"
	testFile1 := testDir + string(os.PathSeparator) + "file1.txt"
	testFile2 := testDir + string(os.PathSeparator) + nestedDir + string(os.PathSeparator) + ".file2.crt"
	testFile3 := testDir + string(os.PathSeparator) + "file3.crt"
	testFile4 := testDir + string(os.PathSeparator) + "file4.crt"

	setupDirectoryForTesting(t, testDir)
	setupDirectoryForTesting(t, testDir+string(os.PathSeparator)+nestedDir)

	createFile(t, testFile1)
	createCertFile(t, testFile2)
	createCertFile(t, testFile3)
	createCertFile(t, testFile4)

	// Invoke method to test
	files, err := GetCertificatesFromPath(testDir)

	// Tear down
	os.RemoveAll(testDir)

	// Validate
	expectedFileCount := 3
	if len(files) != expectedFileCount || err != nil {
		t.Fatalf("response length expected to be %v, actual %v, error %v", expectedFileCount, len(files), err)
	}
}

func TestReadFilesFromPath_SymbolicLink(t *testing.T) {

	// Setup
	currPath, _ := os.Getwd()
	testDirName := "TestDirectory"
	absTestDirPath := currPath + string(os.PathSeparator) + testDirName + string(os.PathSeparator)
	testFile1 := absTestDirPath + "file1.Crt"

	setupDirectoryForTesting(t, testDirName)
	createCertFile(t, testFile1)

	symlink := absTestDirPath + "symlink"
	os.Symlink(testFile1, symlink)
	files, err := GetCertificatesFromPath(symlink)

	// Teardown
	os.RemoveAll(absTestDirPath)

	// Validate
	if len(files) != 1 || err != nil {
		t.Fatalf("response length expected to be 1, actual %v, error %v", len(files), err)
	}
}

func TestReadFilesFromPath_MultilevelSymbolicLink(t *testing.T) {

	// Setup
	currPath, _ := os.Getwd()
	testDirName := "TestDirectory"
	absTestDirPath := currPath + string(os.PathSeparator) + testDirName + string(os.PathSeparator)
	testFile1 := absTestDirPath + "file1.Crt"

	setupDirectoryForTesting(t, testDirName)
	createCertFile(t, testFile1)

	symlink := absTestDirPath + "symlink"
	twoLevelSymlink := absTestDirPath + "symlink2"
	os.Symlink(testFile1, symlink)
	os.Symlink(symlink, twoLevelSymlink)
	files, err := GetCertificatesFromPath(twoLevelSymlink)

	// Teardown
	os.RemoveAll(absTestDirPath)

	// validate
	if len(files) != 1 || err != nil {
		t.Fatalf("response length expected to be 1, actual %v, error %v", len(files), err)
	}
}

func TestReadFilesFromPath_SingleFile(t *testing.T) {
	// Setup
	testDir := "TestDirectory"
	testFile1 := testDir + string(os.PathSeparator) + "file1.Crt"

	setupDirectoryForTesting(t, testDir)
	createCertFile(t, testFile1)

	// Invoke method to test
	files, err := GetCertificatesFromPath(testDir)

	// Teardown
	os.RemoveAll(testDir)

	// Validate
	if len(files) != 1 || err != nil {
		t.Fatalf("response length expected to be 1, actual %v, error %v", len(files), err)
	}

}

func createFile(t *testing.T, path string) {
	_, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("write file '%s' failed with error '%v'", path, err)
	}
}

func createCertFile(t *testing.T, path string) {

	// open cert file
	content, err := ioutil.ReadFile("testCert1.crt")

	if err != nil {
		t.Fatalf("open cert file '%s' failed with error '%v'", path, err)
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		t.Fatalf("creating new file '%s' failed with error '%v'", path, err)

	}

	_, err = file.Write(content)
	if err != nil {
		t.Fatalf("write file '%s' failed with error '%v'", path, err)
	}

}

func setupDirectoryForTesting(t *testing.T, path string) {
	err := os.Mkdir(path, 0755)
	if err != nil {
		t.Fatalf("Creating directory '%s' failed with '%v'", path, err)
	}
}
