package handlers

import (
	"bytes"
	"encoding/json"
	"library_project/database"
	"library_project/internal/middleware"
	"library_project/internal/model"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// SongDetail - структура для хранения информации о песне
type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

// GetSongs - Получение списка песен с фильтрацией и пагинацией
// @Summary Получение списка песен
// @Description Получение списка песен с фильтрацией и пагинацией
// @Tags songs
// @Accept json
// @Produce json
// @Param group query string false "Группа"
// @Param song query string false "Название песни"
// @Param page query int false "Страница"
// @Param limit query int false "Лимит"
// @Success 200 {array} model.Song
// @Router /songs [get]
func GetSongs(c *gin.Context) {
	var songs []model.Song
	group := c.Query("group")
	song := c.Query("song")
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	offset := (page - 1) * limit
	query := database.DB.Model(&songs)

	if group != "" {
		query = query.Where("group = ?", group)
	}
	if song != "" {
		query = query.Where("song_title LIKE ?", "%"+song+"%")
	}

	query.Offset(offset).Limit(limit).Find(&songs)
	c.JSON(http.StatusOK, songs)
}

// GetSongText - Получение текста песни по ID
// @Summary Получение текста песни
// @Description Получение текста песни по ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Success 200 {array} string
// @Failure 404 {object} map[string]string
// @Router /songs/{id}/text [get]
func GetSongText(c *gin.Context) {
	id := c.Param("id")
	var song model.Song

	if err := database.DB.First(&song, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}

	verses := splitTextIntoVerses(song.Text)
	c.JSON(http.StatusOK, verses)
}

// splitTextIntoVerses - Вспомогательная функция для разбивки текста на куплеты
func splitTextIntoVerses(text string) []string {
	return strings.Split(text, "\n")
}

// DeleteSong - Удаление песни по ID
// @Summary Удаление песни
// @Description Удаление песни по ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Success 204 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /songs/{id} [delete]
func DeleteSong(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&model.Song{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// UpdateSong - Обновление данных существующей песни по ID
// @Summary Изменение данных песни
// @Description Обновление данных существующей песни по ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Param song body model.Song true "Обновление данных песни"
// @Success 200 {object} model.Song
// @Failure 404 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /songs/{id} [put]
func UpdateSong(c *gin.Context) {
	id := c.Param("id")
	var song model.Song

	if err := database.DB.First(&song, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}

	if err := c.ShouldBindJSON(&song); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Save(&song)
	c.JSON(http.StatusOK, song)
}

// AddSong - Добавление новой песни в библиотеку
// @Summary Добавление новой песни
// @Description Добавление новой песни в библиотеку
// @Tags songs
// @Accept json
// @Produce json
// @Param song body model.Song true "Создание новой песни"
// @Success 201 {object} model.Song
// @Failure 400 {object} map[string]string
// @Router /songs [post]
func AddSong(c *gin.Context) {
	var song model.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		middleware.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	middleware.Logger.Infof("Adding new song: %s by %s", song.SongTitle, song.Group)

	// Пример отправки запроса к внешнему API
	infoURL := "http://localhost:8080/info"
	reqBody, _ := json.Marshal(song)
	resp, err := http.Post(infoURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		middleware.Logger.Error("Error posting to external API")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error posting to external API"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		middleware.Logger.Error("External API returned non-200 status")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "External API returned non-200 status"})
		return
	}

	// Сохраняем песню в базе данных
	database.DB.Create(&song)
	c.JSON(http.StatusCreated, song)
}
