package controllers

import (
	"crud-api-local-storage/helper"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
)

// path buat file
const storageRoot = "D:/storage-local"
const maxUpload = 50 * 1024 // 5MB

func safeJoin(root, subDir string) (string, error) {
	path := filepath.Join(root, subDir)
	rel, err := filepath.Rel(root, path)
	if err != nil || strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("Path tidak bisa")
	}
	return path, nil
}

// UploadFile godoc
// @Summary Upload file
// @Description Upload file ke local storage
// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Success 200 {object} map[string]interface{} "Upload berhasil"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /files/upload [post]
func UploadFile(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUpload)
	if err := r.ParseMultipartForm(maxUpload); err != nil {
		helper.ResponseJSON(w, http.StatusBadRequest, "File upload terlalu besar", nil)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		helper.ResponseJSON(w, http.StatusBadRequest, "Gagal mengambil file", nil)
		return
	}
	defer file.Close()

	subDir := r.URL.Query().Get("dir")
	targetDir, err := safeJoin(storageRoot, subDir)
	if err != nil {
		helper.ResponseJSON(w, http.StatusBadRequest, "Path tidak valid", nil)
		return
	}

	os.MkdirAll(targetDir, os.ModePerm)

	destPath := filepath.Join(targetDir, header.Filename)
	dst, err := os.Create(destPath)
	if err != nil {
		helper.ResponseJSON(w, http.StatusInternalServerError, "Gagal membuat file di server", nil)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		helper.ResponseJSON(w, http.StatusInternalServerError, "Gagal menyimpan file", nil)
		return
	}

	helper.ResponseJSON(w, http.StatusOK, "File berhasil diupload", helper.UploadData{
		FilePath: destPath,
		FileName: header.Filename,
	})
}

// DownloadFile godoc
// @Summary Download/View file
// @Description Download atau view file dari local storage
// @Tags Files
// @Produce application/octet-stream
// @Param path path string true "File path"
// @Success 200 {file} binary "File content"
// @Failure 404 {object} map[string]string "File not found"
// @Router /files/view/{path} [get]
func DownloadFile(w http.ResponseWriter, r *http.Request) {
	fileName := chi.URLParam(r, "path")
	subDir := r.URL.Query().Get("dir")

	filePath := filepath.Join(storageRoot, subDir, fileName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		helper.ResponseJSON(w, http.StatusNotFound, "File tidak ditemukan", nil)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(filePath)))
	http.ServeFile(w, r, filePath)
}

// DeleteFiles godoc
// @Summary Delete file
// @Description Menghapus file dari local storage
// @Tags Files
// @Produce json
// @Param path path string true "File path to delete"
// @Success 200 {object} map[string]string "File deleted"
// @Failure 404 {object} map[string]string "File not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /files/{path} [delete]
func DeleteFiles(w http.ResponseWriter, r *http.Request) {
	fileName := chi.URLParam(r, "path")
	subDir := r.URL.Query().Get("dir")

	filePath := filepath.Join(storageRoot, subDir, fileName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		helper.ResponseJSON(w, http.StatusNotFound, "File tidak ditemukan di folder tersebut", nil)
		return
	}

	if err := os.Remove(filePath); err != nil {
		helper.ResponseJSON(w, http.StatusInternalServerError, "Gagal menghapus file", nil)
		return
	}

	helper.ResponseJSON(w, http.StatusOK, "File berhasil dihapus", helper.DeleteData{
		DeletedFile: fileName,
		Directory:   subDir,
	})
}

// ListFiles godoc
// @Summary List all files
// @Description Mendapatkan daftar semua file yang tersimpan
// @Tags Files
// @Produce json
// @Success 200 {array} map[string]interface{} "List of files"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /files [get]
func ListFiles(w http.ResponseWriter, r *http.Request) {
	subDir := r.URL.Query().Get("dir")
	targetDir, err := safeJoin(storageRoot, subDir)
	if err != nil {
		helper.ResponseJSON(w, http.StatusBadRequest, "Path tidak valid", nil)
		return
	}

	files, err := os.ReadDir(targetDir)
	if err != nil {
		helper.ResponseJSON(w, http.StatusNotFound, "Direktori tidak ditemukan", nil)
		return
	}

	var fileList []string
	for _, file := range files {
		if !file.IsDir() {
			fileList = append(fileList, file.Name())
		}
	}

	helper.ResponseJSON(w, http.StatusOK, "Berhasil mengambil daftar file", helper.FileListData{
		Directory: subDir,
		Files:     fileList,
	})
}
