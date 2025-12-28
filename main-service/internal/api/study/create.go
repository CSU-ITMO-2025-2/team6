package study

import (
	"io"
	"net/http"

	"main-service/internal/api/response"
	"main-service/internal/model"

	"github.com/google/uuid"
)

// @Summary Создает новое исследование
// @Description Отправляет изображение сельхоз культуры на анализ
// @Tags study
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Image file"
// @Success 200 {object} CreateStudyResponse
// @Failure 400 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /study [post]

func (i *Implementation) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := r.ParseMultipartForm(32 << 20) // 32 MB
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_FORM", "invalid form-data")
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "MISSING_IMAGE", "image is required")
		return
	}
	defer file.Close()

	imageBytes, err := io.ReadAll(file)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_IMAGE", "failed to read image")
		return
	}

	userID := uuid.New() // FIXME брать из контекста, когда добавиться авторизация

	st := &model.Study{
		OwnerID: userID,
		Image: &model.Image{
			ChunkData: imageBytes,
			FileName:  "",
			MimeType:  header.Header.Get("Content-Type"),
		},
	}

	id, err := i.studyService.Create(ctx, st)
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.OK(w, CreateStudyResponse{ID: id.String()})
}
