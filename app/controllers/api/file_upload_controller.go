package api

import (
	"GIG-SDK/libraries"
	"GIG-SDK/models"
	"GIG/app/controllers"
	"GIG/app/storages"
	"github.com/revel/revel"
	"log"
	"net/url"
)

type FileUploadController struct {
	*revel.Controller
}

// swagger:operation POST /upload File upload
//
// Upload File
//
// This API allows to upload a file to the server from a remote source url
//
// ---
// produces:
// - application/json
//
// parameters:
//
// - name: upload
//   in: body
//   description: upload object
//   required: true
//   schema:
//       "$ref": "#/definitions/Upload"
//
// security:
//   - Bearer: []
//   - ApiKey: []
//
// responses:
//   '200':
//     description: file uploaded
//     schema:
//         type: string
//   '403':
//     description: input validation error
//     schema:
////       "$ref": "#/definitions/Response"
//   '500':
//     description: server error
//     schema:
//       "$ref": "#/definitions/Response"
func (c FileUploadController) Upload() revel.Result {
	var (
		upload models.Upload
	)
	c.Response.Status = 400
	if err := c.Params.BindJSON(&upload); err != nil {
		return c.RenderJSON(err)
	}

	go func(uploadedFile models.Upload) {
		decodedFileName, err := url.QueryUnescape(libraries.ExtractFileName(uploadedFile.GetSource()))
		if err != nil {
			log.Println("decode filename error:", err)
		}

		tempDir := storages.FileStorageHandler{}.GetCacheDirectory() + uploadedFile.GetTitle() + "/"
		tempFile := tempDir + decodedFileName
		if err = libraries.EnsureDirectory(tempDir); err != nil {
			log.Println("directory create error:", err)
		}

		if err := libraries.DownloadFile(tempFile, uploadedFile.GetSource());
			err != nil {
			log.Println("file download error:", err)
		}

		if err = (storages.FileStorageHandler{}.UploadFile(uploadedFile.GetTitle(), tempFile)); err != nil {
			log.Println("file upload error:", err)
		}
	}(upload)

	return c.RenderJSON(controllers.BuildSuccessResponse("file upload queued.", 200))
}
