package guest

import (
	"everlasting/src/domain/event"
	"everlasting/src/domain/sharedkernel/messagebroker"
	"everlasting/src/domain/sharedkernel/photo"
	"everlasting/src/domain/validator"
	"everlasting/src/infrastructure/http/middleware"
	"everlasting/src/infrastructure/http/routes"
	"everlasting/src/infrastructure/persistence"
	"everlasting/src/infrastructure/pkg/logger"
	minio "everlasting/src/infrastructure/pkg/storage"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sarulabs/di"
	"io"
	"net/http"
	"path/filepath"
)

const maxUploadBytes = 5 << 20

func wishingWallMessage(c echo.Context) (err error) {
	defer c.Request().Body.Close()

	// Load container
	var (
		container                = c.Get(string(middleware.MiddlewareValueContainer)).(di.Container)
		wishingMessageRepository = container.Get("persistence.wishing_wall_message").(*persistence.WishingWallMessagePersistence)
		storage                  = container.Get("pkg.storage.minio").(*minio.MinioStorage)
		broker                   = container.Get("pkg.messagebroker.amqp").(messagebroker.MessageBroker)
	)

	cc := c.Get(string(middleware.MiddlewareValueAppLoggerContext)).(*logger.AppLoggerContext)
	ctx := cc.GetContext()

	id := c.Param("id")
	if id == "" {
		return c.String(http.StatusInternalServerError, "ID should not be empty")
	}

	// Parse multipart form data
	err = c.Request().ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		return err
	}

	//handling photo
	photoFile := photo.Photo{}
	file, _ := c.FormFile("photo")
	if file != nil {
		if file.Size <= 0 || file.Size > maxUploadBytes {
			return routes.JsonResponse(c, nil, "Invalid file size",
				"image must be > 0 and <= 5MB", http.StatusBadRequest, nil)
		}

		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		// Read file content
		lr := &io.LimitedReader{R: src, N: maxUploadBytes + 1}
		photoBytes, err := io.ReadAll(lr)
		if err != nil {
			return err
		}
		if int64(len(photoBytes)) > maxUploadBytes {
			return routes.JsonResponse(c, nil, "File too large", "max 5MB", http.StatusRequestEntityTooLarge, nil)
		}

		// MIME sniffing (magic bytes), ignore filename/content-type headers
		detectedMIME, forcedExt, err := validator.DetectAndValidateImage(photoBytes)
		if err != nil {
			return routes.JsonResponse(c, nil, "Invalid image", err.Error(), http.StatusBadRequest, nil)
		}

		uniqueID := uuid.New().String()
		ext := filepath.Ext(file.Filename)
		uniqueFilename := "product/" + uniqueID + ext

		photoFile = photo.Photo{
			ContentType: detectedMIME,
			Byte:        photoBytes,
			PhotoUrl:    "",
			Filename:    uniqueFilename,
			Size:        int64(len(photoBytes)),
			FileType:    forcedExt,
		}
	}

	message := new(event.WishingWallInput)
	message.Name = c.FormValue("name")
	message.Message = c.FormValue("message")
	message.EventID = id
	message.Photo = photoFile

	wishWallMessage, err := message.SaveMessage(ctx, wishingMessageRepository, storage, broker)
	if err != nil {
		return err
	}

	return routes.JsonResponse(c, wishWallMessage, "Ok", "ok", 201, nil)
}

func RegisterGuestRoutes(container di.Container, server *echo.Group) {
	wishingWall := server.Group("/wishing-wall")

	wishingWall.POST("/:id/message", wishingWallMessage)
}
