package initializers

import (
	"archive/zip"
	"io"
	"log"
)

func GenerateArchive(resultStream io.Writer, files map[string][]byte) error {
	zipOutputStream := zip.NewWriter(resultStream)
	defer zipOutputStream.Close()

	for filename, fileContent := range files {
		fileHeader, err := zip.FileInfoHeader(filename)
		if err != nil {
			return err
		}

		writer, err := zipOutputStream.CreateHeader(fileHeader)
		if err != nil {
			return err
		}

		_, err = writer.Write(fileContent)
		if err != nil {
			return err
		}
	}

	return nil
}

func GenerateArchiveFromStream(resultStream io.Reader, output io.Writer) error {
	zipOutputStream := zip.NewWriter(output)
	defer zipOutputStream.Close()

	buf := make([]byte, 1024)
	for {
		n, err := resultStream.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			return err
		}

		_, err = zipOutputStream.Write(buf[:n])
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}