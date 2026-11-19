package services

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// UploadService guarda las fotos de las capturas. El directorio sale de
// UPLOAD_DIR (en el contenedor es /app/uploads, montado como volumen en el
// compose) y se sirve por /uploads.
type UploadService struct {
	dir string
}

func NewUploadService() *UploadService {
	if err := os.MkdirAll(UploadDir(), 0o755); err != nil {
		// No se rompe el arranque: el upload va a fallar con un 500 explícito.
		fmt.Printf("no se pudo crear el directorio de fotos (%s): %v\n", UploadDir(), err)
	}
	return &UploadService{dir: UploadDir()}
}

// UploadDir es el directorio de destino; queda exported para poder testearlo.
func UploadDir() string {
	if dir := os.Getenv("UPLOAD_DIR"); dir != "" {
		return dir
	}
	return "uploads"
}

const maxUploadBytes = 8 << 20 // 8 MB

// allowedImageTypes mapea el tipo detectado en el contenido a la extensión.
var allowedImageTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
	"image/heic": ".heic",
}

// UploadPhoto recibe un multipart con el campo "photo" y devuelve la URL pública.
func (s *UploadService) UploadPhoto(c *fiber.Ctx) error {
	file, err := c.FormFile("photo")
	if err != nil {
		return c.Status(400).JSON(Error("Falta el archivo (campo 'photo')"))
	}
	if file.Size > maxUploadBytes {
		return c.Status(413).JSON(Error("La foto supera los 8 MB"))
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(400).JSON(Error("No se pudo leer la foto"))
	}
	defer src.Close()

	// Se inspecciona el contenido, no el header: el cliente puede mentir.
	head := make([]byte, 512)
	n, _ := src.Read(head)
	ext, ok := allowedImageTypes[http.DetectContentType(head[:n])]
	if !ok {
		return c.Status(415).JSON(Error("Formato no permitido (JPEG, PNG, WebP o HEIC)"))
	}

	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return c.Status(500).JSON(Error("No se pudo preparar la carpeta de fotos"))
	}

	name := uuid.NewString() + ext
	out, err := os.Create(filepath.Join(s.dir, name))
	if err != nil {
		return c.Status(500).JSON(Error("No se pudo guardar la foto"))
	}
	defer out.Close()

	if _, err := out.Write(head[:n]); err != nil {
		return c.Status(500).JSON(Error("No se pudo guardar la foto"))
	}
	if _, err := io.Copy(out, src); err != nil {
		return c.Status(500).JSON(Error("No se pudo guardar la foto"))
	}

	return c.Status(201).JSON(fiber.Map{
		"url":      "/uploads/" + name,
		"filename": name,
	})
}
